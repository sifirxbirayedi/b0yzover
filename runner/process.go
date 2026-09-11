package runner

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/logrusorgru/aurora"
)

func Process(config *Config) error {
	fingerprints, err := Fingerprints()

	if err != nil {
		return fmt.Errorf("Process: %v", err)
	}

	config.initHTTPClient()
	config.loadFingerprints()
	subdomains := getSubdomains(config)

	fmt.Println(aurora.Cyan("[ * ]"), "Yüklenen hedef sayısı:", aurora.BrightYellow(len(subdomains)))
	fmt.Println(aurora.Cyan("[ * ]"), "Yüklenen parmak izi sayısı:", aurora.BrightYellow(len(fingerprints)))

	if config.Output != "" {
		fmt.Printf("[ * ] Çıktı dosyası: %s\n", config.Output)
		fmt.Println(isEnabled(config.OnlyVuln), "Yalnızca savunmasız alt alan adları kaydedilir (--vuln)")
	}

	fmt.Println(isEnabled(config.HTTPS), "Varsayılan olarak HTTPS kullan (--https)")
	fmt.Println("[", config.Concurrency, "]", "Eşzamanlı istek sayısı (--concurrency)")
	fmt.Println(isEnabled(config.VerifySSL), "Yalnızca SSL sertifikası geçerli hedefler denetlenir (--verify_ssl)")
	fmt.Println("[", config.Timeout, "]", "HTTP istek zaman aşımı (saniye) (--timeout)")
	fmt.Println(isEnabled(config.HideFails), "Yalnızca potansiyel olarak savunmasız alt alan adları gösterilir (--hide_fails)")

	const ExtraChannelCapacity = 5
	subdomainCh := make(chan string, config.Concurrency+ExtraChannelCapacity)
	resCh := make(chan *subdomainResult, config.Concurrency)

	var wg sync.WaitGroup
	wg.Add(config.Concurrency)

	var results []*subdomainResult
	var vulnerableResults []VulnerableResult

	var collectorWg sync.WaitGroup
	collectorWg.Add(1)
	go collectResults(resCh, &results, &vulnerableResults, config, &collectorWg)

	for i := 0; i < config.Concurrency; i++ {
		go processor(subdomainCh, resCh, config, &wg)
	}

	distributeSubdomains(subdomains, subdomainCh)
	wg.Wait()
	close(resCh)
	collectorWg.Wait()

	if len(vulnerableResults) > 0 {
		if err := saveVulnerableResults(vulnerableResults); err != nil {
			fmt.Printf("[ ! ] Savunmasız sonuçlar kaydedilemedi: %v\n", err)
		}
	}

	if config.Output != "" {
		if err := saveResults(config.Output, results); err != nil {
			return err
		}
	}

	return nil
}

func processor(subdomainCh <-chan string, resCh chan<- *subdomainResult, c *Config, wg *sync.WaitGroup) {
	defer wg.Done()
	for subdomain := range subdomainCh {
		result := c.checkSubdomain(subdomain)

		res := &subdomainResult{
			Subdomain:     subdomain,
			Status:        string(result.ResStatus),
			Engine:        result.Entry.Service,
			Documentation: result.Entry.Documentation,
			Discussion:    result.Entry.Discussion,
			CICDPass:      result.Entry.CICDPass,
			CName:         result.Entry.CName,
			Fingerprint:   result.Entry.Fingerprint,
			NXDomain:      result.Entry.NXDomain,
			Service:       result.Entry.Service,
			Vulnerable:    result.ResStatus == ResultVulnerable,
		}
		if result.HTTPStatus > 0 {
			status := result.HTTPStatus
			res.HTTPStatus = &status
		}

		if result.ResStatus == ResultVulnerable {
			fmt.Print("-----------------\n")
			fmt.Println("[", result.Status, "]", " - ", subdomain, " [", result.Entry.Service, "]")
			fmt.Println("[", aurora.Blue("TARTIŞMA"), "]", " - ", result.Entry.Discussion)
			fmt.Println("[", aurora.Blue("BELGELENDİRME"), "]", " - ", result.Entry.Documentation)
			fmt.Print("-----------------\n")
		} else {
			if !c.HideFails {
				fmt.Println("[", result.Status, "]", " - ", subdomain)
			}
		}

		resCh <- res
	}
}

func distributeSubdomains(subdomains []string, subdomainCh chan<- string) {
	for _, subdomain := range subdomains {
		subdomainCh <- subdomain
	}
	close(subdomainCh)
}

func collectResults(resCh <-chan *subdomainResult, results *[]*subdomainResult, vulnerableResults *[]VulnerableResult, config *Config, wg *sync.WaitGroup) {
	defer wg.Done()
	for r := range resCh {
		if config.Output != "" && (!config.OnlyVuln || r.Vulnerable) {
			*results = append(*results, r)
		}

		if r.Vulnerable {
			cname := GetCNAME(r.Subdomain)
			vulnResult := VulnerableResult{
				Subdomain: r.Subdomain,
				Service:   r.Engine,
				CNAME:     cname,
			}
			*vulnerableResults = append(*vulnerableResults, vulnResult)
		}
	}
}

func saveResults(filename string, results []*subdomainResult) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return err
	}

	fmt.Printf("[ * ] Sonuçlar %q dosyasına kaydedildi\n", filename)
	return nil
}

func saveVulnerableResults(results []VulnerableResult) error {
	filename := "vulnerable_subdomains.txt"

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	header := fmt.Sprintf("=== SAVUNMASIZ SUBDOMAIN'LER ===\nBulunma Tarihi: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	if _, err := f.WriteString(header); err != nil {
		return err
	}

	for i, r := range results {
		line := fmt.Sprintf("%d. %s\n   Servis: %s\n   CNAME: %s\n\n",
			i+1, r.Subdomain, r.Service, r.CNAME)
		if _, err := f.WriteString(line); err != nil {
			return err
		}
	}

	fmt.Printf("[ * ] Savunmasız subdomain'ler %q dosyasına kaydedildi (toplam: %d)\n",
		filename, len(results))
	return nil
}

func RunDiscoveryScan(domain string, scan *Config) error {
	fmt.Println(aurora.BrightCyan("[ 1 / 2 ]"), aurora.White("Subdomain sorgulanıyor:"), aurora.Yellow(domain))

	subdomains, err := RunSubfinderWithCrt(domain)
	if err != nil {
		return err
	}

	if len(subdomains) == 0 {
		return fmt.Errorf("hata: %s için alt alan adı bulunamadı", domain)
	}

	subdomains = CheckAndFilterWildcard(domain, subdomains)

	for _, s := range subdomains {
		fmt.Println("  ", aurora.BrightBlack(s))
	}

	listFile := SubdomainsFile(domain)
	if err := saveSubdomains(listFile, subdomains); err != nil {
		return fmt.Errorf("alt alan adı listesi %q dosyasına yazılamadı: %v", listFile, err)
	}
	fmt.Printf("          %s %d\n", aurora.Cyan("Toplanan alt alan adı sayısı:"), len(subdomains))
	fmt.Printf("          %s %s\n", aurora.Cyan("Alt alan adı listesi:"), aurora.Green(listFile))

	fmt.Println(aurora.BrightCyan("[ 2 / 2 ]"), aurora.White("Takeover taraması başlıyor."))
	scan.Targets = listFile
	return Process(scan)
}

func getSubdomains(c *Config) []string {
	if c.Target == "" {
		subdomains, err := readSubdomains(c.Targets)
		if err != nil {
			log.Fatalf("Alt alan adları okunurken hata oluştu: %s", err)
		}
		return subdomains
	}

	parts := strings.Split(c.Target, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func CheckAndFilterWildcard(domain string, subdomains []string) []string {
	rand.Seed(time.Now().UnixNano())
	randomSub := fmt.Sprintf("b0yz-wildcard-%d.%s", rand.Intn(100000), domain)

	addrs, err := net.LookupHost(randomSub)
	if err != nil || len(addrs) == 0 {
		return subdomains
	}

	wildcardIP := addrs[0]
	fmt.Printf(" Uyarı: %s üzerinde Wildcard DNS aktif (IP: %s). Eşleşenler filtreleniyor...\n", domain, wildcardIP)

	var filtered []string
	for _, sub := range subdomains {
		subAddrs, subErr := net.LookupHost(sub)
		if subErr == nil {
			isWildcardIP := false
			for _, ip := range subAddrs {
				if ip == wildcardIP {
					isWildcardIP = true
					break
				}
			}
			if !isWildcardIP {
				filtered = append(filtered, sub)
			}
		} else {
			filtered = append(filtered, sub)
		}
	}
	return filtered
}

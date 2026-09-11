package runner

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func RunSubfinder(domain string) ([]string, error) {
	if !IsSubfinderInstalled() {
		return nil, fmt.Errorf("subfinder kurulu değil. Kurulum: go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest")
	}

	fmt.Printf("[ * ] Subfinder ile subdomain keşfi başlatılıyor: %s\n", domain)

	cmd := exec.Command("subfinder", "-d", domain, "-silent")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("subfinder çalıştırılamadı: %v", err)
	}

	var subdomains []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	seen := make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || seen[line] {
			continue
		}
		if strings.HasSuffix(line, "."+domain) || line == domain {
			seen[line] = true
			subdomains = append(subdomains, line)
		}
	}

	if len(subdomains) == 0 {
		return nil, fmt.Errorf("subfinder %s için subdomain bulamadı", domain)
	}

	fmt.Printf("[ * ] Subfinder %d subdomain buldu\n", len(subdomains))
	return subdomains, nil
}

func IsSubfinderInstalled() bool {
	_, err := exec.LookPath("subfinder")
	return err == nil
}

func RunSubfinderWithCrt(domain string) ([]string, error) {
	var allSubdomains []string
	seen := make(map[string]bool)

	if IsSubfinderInstalled() {
		subfinderSubs, err := RunSubfinder(domain)
		if err == nil {
			for _, s := range subfinderSubs {
				if !seen[s] {
					seen[s] = true
					allSubdomains = append(allSubdomains, s)
				}
			}
		} else {
			fmt.Printf("[ ! ] Subfinder hatası: %v\n", err)
		}
	} else {
		fmt.Println("[ ! ] Subfinder kurulu değil, sadece crt.name kullanılacak")
	}

	crtSubs, err := FetchSubdomains(domain)
	if err == nil {
		for _, s := range crtSubs {
			if !seen[s] {
				seen[s] = true
				allSubdomains = append(allSubdomains, s)
			}
		}
	} else {
		fmt.Printf("[ ! ] crt.name hatası: %v\n", err)
	}

	if len(allSubdomains) == 0 {
		return nil, fmt.Errorf("hiçbir kaynaktan subdomain bulunamadı: %s", domain)
	}

	fmt.Printf("[ * ] Toplam %d benzersiz subdomain bulundu\n", len(allSubdomains))
	return allSubdomains, nil
}

func InstallSubfinder() error {
	fmt.Println("[ * ] Subfinder kuruluyor...")

	cmd := exec.Command("go", "install", "-v", "github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("subfinder kurulumu başarısız: %v - %s", err, string(output))
	}

	fmt.Println("[ * ] Subfinder başarıyla kuruldu")
	return nil
}

package runner

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const crtSearchEndpoint = "https://crt.name/v1/search?apex="

func SubdomainsFile(domain string) string {
	return domain + "_altalanlar.txt"
}

func FetchSubdomains(domain string) ([]string, error) {
	client := &http.Client{Timeout: 90 * time.Second}

	req, err := http.NewRequest("GET", crtSearchEndpoint+url.QueryEscape(domain), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("crt.name sorgusu başarısız: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("crt.name beklenmeyen durum kodu döndürdü: %d", resp.StatusCode)
	}

	seen := make(map[string]bool)
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimSuffix(line, ".")
		if line == "" || seen[line] {
			continue
		}
		if line == domain || strings.HasSuffix(line, "."+domain) {
			seen[line] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("crt.name yanıtı okunamadı: %v", err)
	}

	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out, nil
}

func SaveSubdomains(path string, subdomains []string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, s := range subdomains {
		if _, err := w.WriteString(s + "\n"); err != nil {
			return err
		}
	}
	return w.Flush()
}

func saveSubdomains(path string, subdomains []string) error {
	return SaveSubdomains(path, subdomains)
}

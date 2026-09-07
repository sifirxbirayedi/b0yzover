package runner

import (
	"net/url"
	"strings"
)

type subdomainResult struct {
	Subdomain     string   `json:"subdomain"`
	Status        string   `json:"status"`
	Engine        string   `json:"engine"`
	Documentation string   `json:"documentation"`
	Discussion    string   `json:"discussion"`
	CICDPass      bool     `json:"cicd_pass"`
	CName         []string `json:"cname"`
	Fingerprint   string   `json:"fingerprint"`
	HTTPStatus    *int     `json:"http_status"`
	NXDomain      bool     `json:"nxdomain"`
	Service       string   `json:"service"`
	Vulnerable    bool     `json:"vulnerable"`
}

func isEnabled(setting bool) string {
	if setting {
		return "[ Evet ]"
	}
	return "[ Hayır ]"
}

func isValidUrl(toTest string) bool {
	u, err := url.ParseRequestURI(strings.TrimSpace(toTest))
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func CleanAndExtractTarget(raw string) []string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "\ufeff")
	if raw == "" {
		return nil
	}

	fields := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	})

	var targets []string
	seen := make(map[string]bool)

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		testStr := field
		if !strings.HasPrefix(testStr, "http://") && !strings.HasPrefix(testStr, "https://") {
			testStr = "http://" + testStr
		}

		u, err := url.Parse(testStr)
		var host string
		if err == nil && u.Host != "" {
			host = u.Host
		} else {
			host = field
			if idx := strings.Index(host, "/"); idx != -1 {
				host = host[:idx]
			}
		}

		host = strings.TrimSuffix(host, ".")
		host = strings.TrimPrefix(host, "www.")

		if host != "" && !seen[host] {
			seen[host] = true
			targets = append(targets, host)
		}
	}

	return targets
}

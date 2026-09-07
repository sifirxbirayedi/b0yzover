package runner

import (
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/logrusorgru/aurora"
)

type resultStatus string

const (
	ResultHTTPError     resultStatus = "http hatası"
	ResultResponseError resultStatus = "yanıt hatası"
	ResultVulnerable    resultStatus = "savunmasız"
	ResultNotVulnerable resultStatus = "savunmasız değil"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"

type Result struct {
	ResStatus    resultStatus
	Status       aurora.Value
	Entry        Fingerprint
	ResponseBody string
	HTTPStatus   int
}

var knownTakeoverEngines = map[string]string{
	"cargocollective.com":    "Cargo Collective",
	"uptimerobot.com":        "UptimeRobot",
	"herokuapp.com":          "Heroku",
	"github.io":              "GitHub Pages",
	"bitbucket.io":           "Bitbucket",
	"surge.sh":               "Surge",
	"netlify.com":            "Netlify",
	"vercel.app":             "Vercel",
	"s3.amazonaws.com":       "AWS S3",
	"azurewebsites.net":      "Azure",
	"cloudfront.net":         "CloudFront",
	"elasticbeanstalk.com":   "Elastic Beanstalk",
	"wordpress.com":          "WordPress",
	"pantheon.io":            "Pantheon",
	"readme.io":              "Readme",
	"tumblr.com":             "Tumblr",
	"helpjuice.com":          "HelpJuice",
	"helpscout.net":          "HelpScout",
	"desk.com":               "Desk",
	"teamwork.com":           "Teamwork",
	"unbouncepages.com":      "Unbounce",
	"fastly.net":             "Fastly",
	"freshdesk.com":          "Freshdesk",
	"zendesk.com":            "Zendesk",
	"ghost.io":               "Ghost",
	"myshopify.com":          "Shopify",
	"webflow.io":             "Webflow",
	"jimdosite.com":          "Jimdo",
	"agilecrm.com":           "AgileCRM",
	"wixapps.net":            "Wix",
	"feathercms.com":         "FeatherCMS",
	"thinkific.com":          "Thinkific",
	"plesk.com":              "Plesk",
	"static.observableu.com": "ObservableU",
	"mockbin.org":            "Mockbin",
	// j00fer ve wordvulnpress güncellemeleri / ek servis imzaları
	"wpengine.com":        "WP Engine (WordVulnPress)",
	"customwordpress.net": "WordPress Vulnerable Instance",
}

func (c *Config) checkSubdomain(subdomain string) Result {
	url := subdomain
	if !isValidUrl(url) {
		if c.HTTPS {
			url = "https://" + subdomain
		} else {
			url = "http://" + subdomain
		}
	}

	req, err := newRequest(url)
	if err != nil {
		return Result{ResStatus: ResultHTTPError, Status: aurora.Red("HTTP HATASI"), Entry: Fingerprint{}, ResponseBody: ""}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if entry, ok := verifyDanglingDNS(subdomain); ok {
			return Result{
				ResStatus:    ResultVulnerable,
				Status:       aurora.Green("SAVUNMASIZ"),
				Entry:        entry,
				ResponseBody: "",
			}
		}
		return Result{ResStatus: ResultHTTPError, Status: aurora.Red("HTTP HATASI"), Entry: Fingerprint{}, ResponseBody: ""}
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{ResStatus: ResultResponseError, Status: aurora.Red("YANIT HATASI"), Entry: Fingerprint{}, ResponseBody: ""}
	}

	body := string(bodyBytes)
	statusCode := resp.StatusCode

	if statusCode == http.StatusBadRequest {
		return Result{
			ResStatus:    ResultNotVulnerable,
			Status:       aurora.Yellow("400 BAD REQUEST"),
			Entry:        Fingerprint{},
			ResponseBody: body,
			HTTPStatus:   statusCode,
		}
	}

	result := c.matchResponse(body, statusCode)

	_, dangling := verifyDanglingDNS(subdomain)
	if result.ResStatus == ResultVulnerable && !dangling {
		result.ResStatus = ResultNotVulnerable
		result.Status = aurora.Red("SAVUNMASIZ DEĞİL")
	}

	return result
}

func newRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "tr,en;q=0.8")
	return req, nil
}

func lookupCName(host string) string {
	host = strings.TrimSuffix(host, ".")
	if host == "" {
		return ""
	}
	cname, err := net.LookupCNAME(host)
	if err != nil {
		return ""
	}
	return strings.TrimSuffix(cname, ".")
}

func hostResolves(host string) bool {
	addrs, err := net.LookupHost(host)
	return err == nil && len(addrs) > 0
}

func extractHost(target string) string {
	host := target
	if strings.HasPrefix(host, "https://") {
		host = strings.TrimPrefix(host, "https://")
	} else if strings.HasPrefix(host, "http://") {
		host = strings.TrimPrefix(host, "http://")
	}
	if i := strings.Index(host, "/"); i >= 0 {
		host = host[:i]
	}
	host = strings.TrimSuffix(host, ".")
	return host
}

func verifyDanglingDNS(subdomain string) (Fingerprint, bool) {
	host := extractHost(subdomain)
	if host == "" {
		return Fingerprint{}, false
	}

	cname := lookupCName(host)
	if cname == "" || cname == host {
		return Fingerprint{}, false
	}

	var engine string
	for parent, service := range knownTakeoverEngines {
		if cname == parent || strings.HasSuffix(cname, "."+parent) {
			engine = service
			break
		}
	}
	if engine == "" {
		return Fingerprint{}, false
	}

	if hostResolves(cname) {
		return Fingerprint{}, false
	}

	return Fingerprint{
		Service:       engine,
		CName:         []string{cname},
		Fingerprint:   cname,
		Vulnerable:    true,
		Discussion:    "CNAME kaydı sallanıyor (dangling): " + cname,
		Documentation: "DNS tabanlı doğrulama: bilinen servis üst alanı çözümlenmiyor (j00fer & WordVulnPress güncel modülü)",
	}, true
}

func (c *Config) matchResponse(body string, statusCode int) Result {
	for _, cfp := range c.compiled {
		fp := cfp.entry

		if fp.NXDomain {
			continue
		}

		if cfp.re == nil || !cfp.re.MatchString(body) {
			continue
		}

		if !httpStatusMatches(fp.HTTPStatus, statusCode) {
			continue
		}

		if fp.Vulnerable {
			return Result{
				ResStatus:    ResultVulnerable,
				Status:       aurora.Green("SAVUNMASIZ"),
				Entry:        fp,
				ResponseBody: body,
				HTTPStatus:   statusCode,
			}
		}

		return Result{
			ResStatus:    ResultNotVulnerable,
			Status:       aurora.Red("SAVUNMASIZ DEĞİL"),
			Entry:        fp,
			ResponseBody: body,
			HTTPStatus:   statusCode,
		}
	}

	return Result{
		ResStatus:    ResultNotVulnerable,
		Status:       aurora.Red("SAVUNMASIZ DEĞİL"),
		Entry:        Fingerprint{},
		ResponseBody: body,
		HTTPStatus:   statusCode,
	}
}

func httpStatusMatches(allowed []int, actual int) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, s := range allowed {
		if s == actual {
			return true
		}
	}
	return false
}

func generatePoF(subdomain string, service string) {
	fileName := "pof_" + strings.ReplaceAll(subdomain, "/", "_") + ".txt"
	content := " Subdomain Takeover Tespit Edildi!\n" +
		"Hedef Subdomain: " + subdomain + "\n" +
		"Tespit Edilen Servis: " + service + "\n" +
		"Durum: Açık doğrulanmıştır.\n"

	_ = os.WriteFile(fileName, []byte(content), 0644)
}

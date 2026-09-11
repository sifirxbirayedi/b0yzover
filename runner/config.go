package runner

import (
	"crypto/tls"
	"net/http"
	"regexp"
	"time"
)

type Config struct {
	HTTPS        bool
	VerifySSL    bool
	Emoji        bool
	HideFails    bool
	OnlyVuln     bool
	Concurrency  int
	Timeout      int
	Targets      string
	Target       string
	Output       string
	client       *http.Client
	fingerprints []Fingerprint
	compiled     []compiledFingerprint
}

type compiledFingerprint struct {
	entry Fingerprint
	re    *regexp.Regexp
}

func (s *Config) initHTTPClient() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !s.VerifySSL},
		Proxy:           http.ProxyFromEnvironment,
	}

	timeout := time.Duration(s.Timeout) * time.Second
	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	s.client = client
}

func (c *Config) loadFingerprints() error {
	fingerprints, err := Fingerprints()
	if err != nil {
		return err
	}
	c.fingerprints = fingerprints
	c.compileFingerprints()
	return nil
}

func (c *Config) compileFingerprints() {
	c.compiled = make([]compiledFingerprint, 0, len(c.fingerprints))
	for _, fp := range c.fingerprints {
		entry := fp
		if entry.NXDomain || entry.Fingerprint == "" {
			c.compiled = append(c.compiled, compiledFingerprint{entry: entry})
			continue
		}
		re, err := regexp.Compile(entry.Fingerprint)
		if err != nil {
			quoted := regexp.QuoteMeta(entry.Fingerprint)
			if re2, qerr := regexp.Compile(quoted); qerr == nil {
				re = re2
			}
		}
		c.compiled = append(c.compiled, compiledFingerprint{entry: entry, re: re})
	}
}

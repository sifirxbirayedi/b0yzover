package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type FlexInts []int

func (f *FlexInts) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*f = nil
		return nil
	}

	var single int
	if err := json.Unmarshal(data, &single); err == nil {
		*f = FlexInts{single}
		return nil
	}

	var list []int
	if err := json.Unmarshal(data, &list); err == nil {
		*f = list
		return nil
	}

	var raws []json.RawMessage
	if err := json.Unmarshal(data, &raws); err == nil {
		for _, r := range raws {
			var n int
			if json.Unmarshal(r, &n) == nil {
				*f = append(*f, n)
			}
		}
		return nil
	}

	return fmt.Errorf("http_status: beklenmeyen biçim: %s", data)
}

type FlexBool bool

func (b *FlexBool) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if bytes.Equal(data, []byte("null")) {
		*b = false
		return nil
	}

	var v bool
	if err := json.Unmarshal(data, &v); err == nil {
		*b = FlexBool(v)
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch normalizeASCII(s) {
		case "vulnerable", "true", "evet":
			*b = true
		default:
			*b = false
		}
		return nil
	}

	return fmt.Errorf("vulnerable: beklenmeyen biçim: %s", data)
}

func normalizeASCII(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			r += 32
		}
		out = append(out, r)
	}
	return string(out)
}

type Fingerprint struct {
	CICDPass      bool     `json:"cicd_pass"`
	CName         []string `json:"cname"`
	Discussion    string   `json:"discussion"`
	Documentation string   `json:"documentation"`
	Fingerprint   string   `json:"fingerprint"`
	HTTPStatus    FlexInts `json:"http_status"`
	NXDomain      bool     `json:"nxdomain"`
	Service       string   `json:"service"`
	Status        string   `json:"status"`
	Vulnerable    FlexBool `json:"vulnerable"`
}

func Fingerprints() ([]Fingerprint, error) {
	var fingerprints []Fingerprint

	fingerPrintsPath, err := GetFingerprintPath()
	if err != nil {
		return nil, fmt.Errorf("Fingerprints: %v", err)
	}

	if _, err := os.Stat(fingerPrintsPath); os.IsNotExist(err) {
		fmt.Println("[ * ] fingerprints.json bulunamadı")
		if err := DownloadFingerprints(); err != nil {
			return nil, fmt.Errorf("Fingerprints otomatik indirilemedi: %v", err)
		}
		fmt.Println("[ * ] fingerprints.json indiirldi ")
	}

	file, err := os.ReadFile(fingerPrintsPath)
	if err != nil {
		return nil, fmt.Errorf("Fingerprints: %v", err)
	}

	err = json.Unmarshal(file, &fingerprints)
	if err != nil {
		return nil, fmt.Errorf("Fingerprints: %v", err)
	}

	return fingerprints, nil
}

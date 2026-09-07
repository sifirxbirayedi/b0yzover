package runner

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var (
	fingerprintURL = "https://raw.githubusercontent.com/EdOverflow/can-i-take-over-xyz/master/fingerprints.json"
)

func GetFingerprintPath() (string, error) {

	execDir, err := os.Getwd()
	if err != nil {
		return "fingerprints.json", nil
	}
	return filepath.Join(execDir, "fingerprints.json"), nil
}

func DownloadFingerprints() error {
	fingerprintsPath, err := GetFingerprintPath()
	if err != nil {
		return err
	}

	out, err := os.OpenFile(fingerprintsPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("DownloadFingerprints: %v", err)
	}
	defer out.Close()

	resp, err := http.Get(fingerprintURL)
	if err != nil {
		return fmt.Errorf("DownloadFingerprints: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DownloadFingerprints: parmak izi veritabanı indirilemedi; HTTP durumu %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("DownloadFingerprints: %v", err)
	}

	return nil
}

func CheckIntegrity() (bool, error) {
	resp, err := http.Get(fingerprintURL)
	if err != nil {
		return false, fmt.Errorf("CheckIntegrity: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("CheckIntegrity: üst kaynağa erişilemedi; HTTP durumu %d", resp.StatusCode)
	}

	outBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	fingerprintsLocal, err := GetFingerprintPath()
	if err != nil {
		return false, err
	}

	localBytes, err := os.ReadFile(fingerprintsLocal)
	if err != nil {
		return false, err
	}

	upstreamSum := md5.Sum(outBytes)
	localSum := md5.Sum(localBytes)

	return hex.EncodeToString(upstreamSum[:]) == hex.EncodeToString(localSum[:]), nil
}

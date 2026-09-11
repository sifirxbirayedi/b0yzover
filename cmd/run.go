package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"b0yzover/runner"
	"github.com/spf13/cobra"
)

var opts = runner.Config{}

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "Taramayı başlat",
	Aliases: []string{"r"},
	Long: `Taramayı başlatır.

İlk çalıştırmada parmak izi (fingerprint) veritabanı otomatik olarak
indirilir. Sonraki çalıştırmalarda yerel veritabanı üst kaynakla
(can-i-take-over-xyz) karşılaştırılarak güncelliği kontrol edilir.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fingerprintsPath, err := runner.GetFingerprintPath()
		if err != nil {
			return err
		}
		if _, err := os.Stat(fingerprintsPath); errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("[ * ] Parmak izi veritabanı bulunamadı; indiriliyor ve %q konumuna kaydediliyor\n",
				fingerprintsPath)
			if err := runner.DownloadFingerprints(); err != nil {
				return err
			}
		} else {
			fmt.Printf("[ * ] Parmak izi veritabanı bulundu; üst kaynakla bütünlük kontrolü yapılıyor\n")
			found, err := runner.CheckIntegrity()
			if err != nil {
				return err
			}
			if !found {
				fmt.Printf("[ * ] Yerel ve üst kaynak parmak izleri arasında fark saptandı; güncel sürüm indiriliyor\n")
				if err := runner.DownloadFingerprints(); err != nil {
					return err
				}
			}
		}

		if err := runner.Process(&opts); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	runCmd.Flags().StringVar(&opts.Target, "target", "", "Virgülle ayrılmış hedef alan adı listesi (örn: hedef1.com,hedef2.com)")
	runCmd.Flags().StringVar(&opts.Targets, "targets", "", "Alt alan adı listesi içeren dosya (her satıra bir hedef)")
	runCmd.Flags().StringVar(&opts.Output, "output", "", "JSON sonuç çıktısı için dosya adı")
	runCmd.Flags().BoolVar(&opts.HTTPS, "https", false, "Protokol belirtilmemiş hedeflerde HTTPS kullan (varsayılan: false, http)")
	runCmd.Flags().BoolVar(&opts.VerifySSL, "verify_ssl", false, "true ise geçersiz SSL sertifikalı siteler HTTP hatası sayılır ve atlanır")
	runCmd.Flags().BoolVar(&opts.HideFails, "hide_fails", false, "Başarısız sonuçları ekranda gösterme")
	runCmd.Flags().BoolVar(&opts.OnlyVuln, "vuln", false, "Yalnızca savunmasız (açık bulunan) alt alan adlarını çıktıya kaydet")
	runCmd.Flags().IntVar(&opts.Concurrency, "concurrency", 10, "Eşzamanlı istek sayısı")
	runCmd.Flags().IntVar(&opts.Timeout, "timeout", 10, "HTTP istek zaman aşımı (saniye)")
	rootCmd.AddCommand(runCmd)
}

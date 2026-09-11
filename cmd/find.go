package cmd

import (
	"fmt"
	"os"

	"b0yzover/runner"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:     "find",
	Short:   "crt.name ile alt alan adlarını bul ve taramayı başlat",
	Aliases: []string{"f"},
	Long: `Alt alan adı toplama ve ele geçirme (takeover) taramasını tek komutta yürütür.

Verilen alan adı https://crt.name/v1/search?apex= adresinden sorgulanır;
bulunan tüm alt alan adları ekranda gösterilir ve <alan>_altalanlar.txt
dosyasına kaydedilir. Tarama bu liste üzerinden yürütülür.

Kullanım:
  b0yzover find hedef.com
  b0yzover find hedef.com --https --concurrency 20 --timeout 15
  b0yzover find hedef.com --hide_fails --vuln --output sonuc.json`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("hata: tek bir alan adı belirtilmelidir (örn: b0yzover find hedef.com)")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]

		fingerprintsPath, err := runner.GetFingerprintPath()
		if err != nil {
			return err
		}
		if _, err := os.Stat(fingerprintsPath); os.IsNotExist(err) {
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

		return runner.RunDiscoveryScan(domain, &findOpts)
	},
}

var findOpts = runner.Config{}

func init() {
	findCmd.Flags().StringVar(&findOpts.Output, "output", "", "JSON sonuç çıktısı için dosya adı")
	findCmd.Flags().BoolVar(&findOpts.HTTPS, "https", true, "Protokol belirtilmemiş hedeflerde HTTPS kullan (varsayılan: true)")
	findCmd.Flags().BoolVar(&findOpts.VerifySSL, "verify_ssl", false, "true ise geçersiz SSL sertifikalı siteler HTTP hatası sayılır ve atlanır")
	findCmd.Flags().BoolVar(&findOpts.HideFails, "hide_fails", false, "Başarısız sonuçları ekranda gösterme")
	findCmd.Flags().BoolVar(&findOpts.OnlyVuln, "vuln", false, "Yalnızca savunmasız (açık bulunan) alt alan adlarını çıktıya kaydet")
	findCmd.Flags().IntVar(&findOpts.Concurrency, "concurrency", 15, "Eşzamanlı istek sayısı")
	findCmd.Flags().IntVar(&findOpts.Timeout, "timeout", 12, "HTTP istek zaman aşımı (saniye)")
	rootCmd.AddCommand(findCmd)
}

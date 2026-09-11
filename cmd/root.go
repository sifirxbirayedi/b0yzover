package cmd

import (
	"fmt"
	"os"

	"b0yzover/runner"

	"github.com/spf13/cobra"
)

var (
	target       string
	targetsFile  string
	concurrency  int
	timeout      int
	output       string
	verbose      bool
	useSubfinder bool
)

var rootCmd = &cobra.Command{
	Use:   "b0yzover",
	Short: "Advanced Subdomain Takeover & Vulnerability Scanner",
	Long:  `b0yzover, bulut servisleri üzerindeki dangling CNAME ve takeover açıklarını tarayan gelişmiş bir araçtır.`,
	Run: func(cmd *cobra.Command, args []string) {
		if target == "" && targetsFile == "" {
			RunInteractive()
			return
		}

		runner.InitLogger(verbose)

		config := &runner.Config{
			HTTPS:       true,
			VerifySSL:   false,
			Concurrency: concurrency,
			Timeout:     timeout,
			Output:      output,
		}

		if targetsFile != "" {
			config.Targets = targetsFile
			if err := runner.Process(config); err != nil {
				fmt.Printf("[!] Hata: %v\n", err)
				os.Exit(1)
			}
		} else if target != "" {
			if useSubfinder {
				if !runner.IsSubfinderInstalled() {
					fmt.Println("[!] Subfinder kurulu değil. Kurulum için: go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest")
					fmt.Println("[!] Veya subfinder olmadan devam etmek için --subfinder flag'ini kaldırın")
					os.Exit(1)
				}

				subdomains, err := runner.RunSubfinder(target)
				if err != nil {
					fmt.Printf("[!] Subfinder hatası: %v\n", err)
					os.Exit(1)
				}

				listFile := runner.SubdomainsFile(target)
				if err := runner.SaveSubdomains(listFile, subdomains); err != nil {
					fmt.Printf("[!] Dosya kaydedilemedi: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("[ * ] %d subdomain %s dosyasına kaydedildi\n", len(subdomains), listFile)
				config.Targets = listFile

				if err := runner.Process(config); err != nil {
					fmt.Printf("[!] Hata: %v\n", err)
					os.Exit(1)
				}
			} else {
				if err := runner.RunDiscoveryScan(target, config); err != nil {
					fmt.Printf("[!] Hata: %v\n", err)
					os.Exit(1)
				}
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&target, "target", "d", "", "Tekil hedef domain (crt.name keşfi için)")
	rootCmd.PersistentFlags().StringVarP(&targetsFile, "list", "l", "", "Taranacak hedefleri içeren TXT dosyası")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "threads", "c", 30, "Eşzamanlı thread/çekirdek sayısı")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", 10, "HTTP istek zaman aşımı (saniye)")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Sonuçların kaydedileceği JSON dosya adı")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Detaylı log modu (DEBUG seviyesi)")
	rootCmd.PersistentFlags().BoolVarP(&useSubfinder, "subfinder", "s", false, "Subfinder kullanarak subdomain keşfi yap (sadece --target ile çalışır)")
}

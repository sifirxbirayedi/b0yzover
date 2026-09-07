package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"b0yzover/runner"

	"github.com/logrusorgru/aurora"
)

func RunInteractive() {
	banner := `
░▒▓███████▓▒░░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓███████▓▒░  
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░    ░▒▓██▓▒░░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓██████▓▒░   ░▒▓██▓▒░  ░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒▒▓█▓▒░░▒▓██████▓▒░ ░▒▓███████▓▒░  
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░    ░▒▓██▓▒░    ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▓█▓▒░ ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░  ░▒▓█▓▒░   ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ ░▒▓█▓▓█▓▒░ ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓███████▓▒░░▒▓████████▓▒░  ░▒▓█▓▒░   ░▒▓████████▓▒░░▒▓██████▓▒░   ░▒▓██▓▒░  ░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░
`

	fmt.Println(aurora.Cyan(banner))
	fmt.Println(aurora.BrightMagenta("b0yner/0x1A7\n"))

	reader := bufio.NewReader(os.Stdin)

	fmt.Println(aurora.Yellow(" Taranacak TXT dosya yolunu veya ana domaini girin:").Bold())
	fmt.Print(aurora.Cyan("Hedef > ").Bold())

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println(aurora.Red("Hedef boş olamaz. Program kapatılıyor.").Bold())
		return
	}

	config := &runner.Config{
		HTTPS:       true,
		VerifySSL:   false,
		Concurrency: 50,
		Timeout:     10,
	}

	start := time.Now()

	if _, err := os.Stat(input); err == nil {
		config.Targets = input
		fmt.Println(aurora.Cyan(fmt.Sprintf(" TXT dosyası algılandı (%s), tarama başlatılıyor...", input)).Bold())
		if err := runner.Process(config); err != nil {
			fmt.Println(aurora.Red(fmt.Sprintf(" Hata oluştu: %v", err)).Bold())
		}
	} else {
		fmt.Println(aurora.Cyan(fmt.Sprintf(" Domain algılandı (%s), crt.name üzerinden keşif ve tarama başlatılıyor...", input)).Bold())
		if err := runner.RunDiscoveryScan(input, config); err != nil {
			fmt.Println(aurora.Red(fmt.Sprintf("Keşif hatası: %v", err)).Bold())
		}
	}

	fmt.Println(aurora.Green(fmt.Sprintf("\n İşlem tamamlandı. Toplam Süre: %v", time.Since(start))).Bold())
}

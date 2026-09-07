package main

import (
	"bufio"
	"fmt"
	"os"

	"b0yzover/cmd"
)

func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
		holdConsole()
		return
	}
	cmd.RunInteractive()
	holdConsole()
}

func holdConsole() {
	fmt.Println("\n[ * ] İşlem tamamlandı. Çıkmak için Enter tuşuna basın...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

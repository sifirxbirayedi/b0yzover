package runner

import (
	"bufio"
	"os"
	"strings"
)

func readSubdomains(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimPrefix(line, "\ufeff")
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}

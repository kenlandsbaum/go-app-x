package env

import (
	"os"
	"strings"
)

func Load(envFile string) {
	bts, err := os.ReadFile(envFile)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(bts), "\n")
	for _, line := range lines {
		parts := strings.Split(line, "=")
		os.Setenv(parts[0], strings.TrimSpace(parts[1]))
	}
}

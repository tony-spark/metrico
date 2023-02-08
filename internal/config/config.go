package config

import (
	"fmt"
	"os"
	"strings"
)

func ParseConfigFileParameter() (string, error) {
	configFile := ""
	idx := -1
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "-c=") || strings.HasPrefix(arg, "-config=") {
			ss := strings.Split(arg, "=")
			if len(ss) < 2 {
				return "", fmt.Errorf("parameter value not set: %s", arg)
			}
			configFile = ss[1]
			break
		}
		if arg == "-c" || arg == "-config" || arg == "--c" || arg == "--config" {
			idx = i
			break
		}
	}
	if idx+1 >= len(os.Args) {
		return "", fmt.Errorf("config parameter value missed")
	}
	if idx != -1 {
		configFile = os.Args[idx+1]
	}

	if len(configFile) == 0 {
		configFile = os.Getenv("CONFIG")
	}

	return configFile, nil
}

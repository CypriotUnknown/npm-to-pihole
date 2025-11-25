package app

import (
	"log/slog"
	"os"
	"regexp"
	"strings"
)

func (app *AppType) getNginxProxyHostServerNames() []string {

	confFiles, err := os.ReadDir(app.mginxProxyDir)
	if err != nil {
		panic(err)
	}

	proxy_hosts := make([]string, 0)

	for _, file := range confFiles {
		fileData, err := os.ReadFile(app.mginxProxyDir + "/" + file.Name())
		if err != nil {
			slog.Error("error reading file", "error", err, "file", file.Name())
			continue
		}

		serverNameRegex := regexp.MustCompile(`server_name (.*);`)
		matches := serverNameRegex.FindStringSubmatch(string(fileData))
		if len(matches) > 1 {
			proxy_hosts = append(proxy_hosts, strings.TrimSpace(matches[1]))
		}
	}

	return proxy_hosts
}

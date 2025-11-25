package utils

import (
	"fmt"
	"strings"
)

func GetBaseDomainFromHostString(host string) (string, error) {
	parts := strings.Split(host, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid host in nginx proxy conf file: %s", host)
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1], nil
}

package utility

import (
	"fmt"
	"jethur/config"
	"strings"
)

var (
	externalAllowedIPSet map[string]bool
	ipsInitialized       bool
)

func initAllowedIPs() {
	externalAllowedIPSet = make(map[string]bool)

	allowedIPs := config.GetEnv("EXTERNAL_ALLOWED_IPS", "")
	if allowedIPs != "" {
		for _, ip := range strings.Split(allowedIPs, ",") {
			ip = strings.TrimSpace(ip)
			if ip != "" {
				externalAllowedIPSet[ip] = true
			}
		}
	}

	ipsInitialized = true
}

func ValidateExternalClientIP(clientIP string) error {
	if !ipsInitialized {
		initAllowedIPs()
	}

	if len(externalAllowedIPSet) == 0 {
		return nil // IP allowlisting disabled
	}

	if externalAllowedIPSet[clientIP] {
		return nil
	}

	return fmt.Errorf("access denied from this network")
}
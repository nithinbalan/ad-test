package utility

import (
	"fmt"
	"jethur/config"
	"strings"
)

func ValidateExternalClientIP(clientIP string) error {
	allowedIPs := config.GetEnv("EXTERNAL_ALLOWED_IPS", "")
	if allowedIPs == "" {
		return nil // IP allowlisting disabled
	}

	for _, ip := range strings.Split(allowedIPs, ",") {
		if strings.TrimSpace(ip) == clientIP {
			return nil
		}
	}

	return fmt.Errorf("Access denied from this network.")
}
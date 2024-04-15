package utility

import (
	"net/url"
	"strings"
)

func IsValidURL(toTest string) (bool, *url.URL) {
	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false, nil
	}
	return true, u
}

func IsValidPlatformURL(toTest string, allowedHosts []string) bool {
	ok, u := IsValidURL(toTest)
	if !ok {
		return false
	}
	host := strings.ToLower(u.Host)
	for _, allowedHost := range allowedHosts {
		if strings.HasSuffix(host, allowedHost) || host == allowedHost {
			return true
		}
	}
	return false
}

func IsValidLinkedInURL(toTest string) bool {
	return IsValidPlatformURL(toTest, []string{"www.linkedin.com", "linkedin.com"})
}

func IsValidInstagramURL(toTest string) bool {
	return IsValidPlatformURL(toTest, []string{"www.instagram.com", "instagram.com"})
}

func IsValidYouTubeURL(toTest string) bool {
	return IsValidPlatformURL(toTest, []string{"www.youtube.com", "youtube.com", "youtu.be"})
}

func IsValidFacebookURL(toTest string) bool {
	return IsValidPlatformURL(toTest, []string{"www.facebook.com", "facebook.com"})
}

func IsValidTwitterURL(toTest string) bool {
	return IsValidPlatformURL(toTest, []string{"www.twitter.com", "twitter.com", "m.twitter.com"})
}

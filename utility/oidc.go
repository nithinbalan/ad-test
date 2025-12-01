package utility

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
)

// Global cache so JWKS & issuer metadata are reused automatically.
// go-oidc handles caching internally once provider is created.
var oidcProviderCache = map[string]*oidc.Provider{}

// VerifyMicrosoftIDToken verifies Microsoft's ID token signature & claims using go-oidc.
// tenantID: your Microsoft tenant (GUID)
// clientID: your Azure AD Application (client ID)
func VerifyMicrosoftIDToken(ctx context.Context, idToken string, tenantID string, clientID string) (map[string]interface{}, error) {

	issuer := fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", tenantID)

	// Reuse provider per tenant (go-oidc caches JWKS keys inside)
	provider, ok := oidcProviderCache[tenantID]
	var err error

	if !ok {
		provider, err = oidc.NewProvider(ctx, issuer)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize OIDC provider: %w", err)
		}
		oidcProviderCache[tenantID] = provider
	}

	// Build token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	// Verify token (signature + claims)
	verifiedToken, err := verifier.Verify(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("id token verification failed: %w", err)
	}

	// Extract claims
	var claims map[string]interface{}
	if err := verifiedToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract id token claims: %w", err)
	}

	return claims, nil
}

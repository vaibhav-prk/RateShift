package auth

import "github.com/vaibhav-prk/RateShift/internal/config"

type Authenticator struct {
	apiKeyMap map[string]string
}

func New(cfg *config.Config) *Authenticator {
	m := make(map[string]string, len(cfg.Tenants))
	for _, t := range cfg.Tenants {
		m[t.APIKey] = t.TenantID
	}
	return &Authenticator{apiKeyMap: m}
}

func (a *Authenticator) Authenticate(apiKey string) (tenantID string, ok bool) {
	tenantID, ok = a.apiKeyMap[apiKey]
	return
}

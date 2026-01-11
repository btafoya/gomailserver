package checks

import (
	"sync"

	"github.com/btafoya/gomailserver/internal/testing/checks/config"
	"github.com/btafoya/gomailserver/internal/testing/checks/mailflow"
	"github.com/btafoya/gomailserver/internal/testing/checks/security"
	"github.com/btafoya/gomailserver/internal/testing/types"
)

type Registry struct {
	checks []Check
	mu     sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		checks: make([]Check, 0),
	}
}

func (r *Registry) Register(check Check) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.checks = append(r.checks, check)
}

func (r *Registry) RegisterAll() {
	r.Register(&config.ConfigSyntaxCheck{})
	r.Register(&config.TLSCertificateCheck{})
	r.Register(&config.PortAvailabilityCheck{})
	r.Register(&config.DatabaseCheck{})
	r.Register(&config.DomainConfigurationCheck{})

	r.Register(&mailflow.SMTPConnectivityCheck{})
	r.Register(&mailflow.SMTPAuthenticationCheck{})
	r.Register(&mailflow.IMAPConnectivityCheck{})
	r.Register(&mailflow.IMAPAuthenticationCheck{})
	r.Register(&mailflow.MailFlowEndToEndCheck{})
	r.Register(&mailflow.MessageIntegrityCheck{})

	r.Register(&security.DKIMConfigAudit{})
	r.Register(&security.DKIMSignatureTest{})
	r.Register(&security.SPFPolicyCheck{})
	r.Register(&security.DMARCPolicyCheck{})
	r.Register(&security.SecurityChainCheck{})
}

func (r *Registry) GetByCategory(category types.Category) []Check {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Check
	for _, check := range r.checks {
		if check.Category() == category {
			result = append(result, check)
		}
	}
	return result
}

func (r *Registry) GetAll() []Check {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Check{}, r.checks...)
}

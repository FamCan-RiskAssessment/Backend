package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
)

type SensitiveFieldHasher struct {
	key []byte
}

func NewSensitiveFieldHasher(security *bootstrap.Security) *SensitiveFieldHasher {
	return &SensitiveFieldHasher{
		key: []byte(security.EncryptionKey),
	}
}

func (h *SensitiveFieldHasher) HashSSN(ssn string) string {
	normalized := normalizeSSN(ssn)
	if normalized == "" {
		return ""
	}

	mac := hmac.New(sha256.New, h.key)
	mac.Write([]byte(normalized))
	return hex.EncodeToString(mac.Sum(nil))
}

func normalizeSSN(ssn string) string {
	return strings.TrimSpace(strings.ReplaceAll(ssn, "-", ""))
}

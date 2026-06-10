package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"
)

type RefreshTokenIssuer interface {
	IssueRefreshToken() (IssuedRefreshToken, error)
	HashRefreshToken(rawToken string) string
}

type IssuedRefreshToken struct {
	Value     string
	Hash      string
	ExpiresAt time.Time
	ExpiresIn int
}

type SecureRefreshTokenIssuer struct {
	ttl time.Duration
	now func() time.Time
}

func NewSecureRefreshTokenIssuer(ttl time.Duration) SecureRefreshTokenIssuer {
	if ttl <= 0 {
		ttl = defaultRefreshTokenTTL()
	}

	return SecureRefreshTokenIssuer{
		ttl: ttl,
		now: time.Now,
	}
}

func (i SecureRefreshTokenIssuer) IssueRefreshToken() (IssuedRefreshToken, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return IssuedRefreshToken{}, err
	}

	value := base64.RawURLEncoding.EncodeToString(bytes)
	return IssuedRefreshToken{
		Value:     value,
		Hash:      i.HashRefreshToken(value),
		ExpiresAt: i.now().UTC().Add(i.ttl),
		ExpiresIn: int(i.ttl.Seconds()),
	}, nil
}

func (i SecureRefreshTokenIssuer) HashRefreshToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func defaultRefreshTokenTTL() time.Duration {
	return 30 * 24 * time.Hour
}

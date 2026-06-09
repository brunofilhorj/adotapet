package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrTokenSecretMissing = errors.New("jwt access secret is required")

type HMACJWTIssuer struct {
	issuer string
	secret string
	ttl    time.Duration
	now    func() time.Time
}

func NewHMACJWTIssuer(issuer string, secret string) HMACJWTIssuer {
	return HMACJWTIssuer{
		issuer: issuer,
		secret: secret,
		ttl:    defaultAccessTokenTTL(),
		now:    time.Now,
	}
}

func (i HMACJWTIssuer) IssueAccessToken(subject TokenSubject) (IssuedToken, error) {
	if strings.TrimSpace(i.secret) == "" {
		return IssuedToken{}, ErrTokenSecretMissing
	}

	now := i.now().UTC()
	expiresAt := now.Add(i.ttl)

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	claims := map[string]any{
		"iss":    i.issuer,
		"sub":    subject.UserID,
		"email":  subject.Email,
		"role":   subject.Role,
		"status": subject.Status,
		"iat":    now.Unix(),
		"exp":    expiresAt.Unix(),
	}

	headerPart, err := encodeJWTPart(header)
	if err != nil {
		return IssuedToken{}, err
	}
	claimsPart, err := encodeJWTPart(claims)
	if err != nil {
		return IssuedToken{}, err
	}

	unsigned := headerPart + "." + claimsPart
	signature := signHS256(unsigned, i.secret)

	return IssuedToken{
		Value:     unsigned + "." + signature,
		ExpiresIn: int(i.ttl.Seconds()),
	}, nil
}

func encodeJWTPart(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func signHS256(value string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func VerifyHMACJWT(token string, secret string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("invalid jwt format")
	}
	expected := signHS256(parts[0]+"."+parts[1], secret)
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return fmt.Errorf("invalid jwt signature")
	}
	return nil
}

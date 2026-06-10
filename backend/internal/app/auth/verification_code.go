package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"adotapet/internal/domain/user"
)

type VerificationCodeIssuer interface {
	IssueCode(userID string, channel user.VerificationChannel, destination string) (IssuedVerificationCode, error)
	HashCode(userID string, channel user.VerificationChannel, destination string, code string) string
}

type IssuedVerificationCode struct {
	Value     string
	Hash      string
	ExpiresAt time.Time
}

type NumericVerificationCodeIssuer struct {
	ttl time.Duration
	now func() time.Time
}

func NewNumericVerificationCodeIssuer(ttl time.Duration) NumericVerificationCodeIssuer {
	if ttl <= 0 {
		ttl = 15 * time.Minute
	}

	return NumericVerificationCodeIssuer{
		ttl: ttl,
		now: time.Now,
	}
}

func (i NumericVerificationCodeIssuer) IssueCode(userID string, channel user.VerificationChannel, destination string) (IssuedVerificationCode, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return IssuedVerificationCode{}, err
	}

	code := fmt.Sprintf("%06d", n.Int64())
	return IssuedVerificationCode{
		Value:     code,
		Hash:      i.HashCode(userID, channel, destination, code),
		ExpiresAt: i.now().UTC().Add(i.ttl),
	}, nil
}

func (i NumericVerificationCodeIssuer) HashCode(userID string, channel user.VerificationChannel, destination string, code string) string {
	payload := strings.Join([]string{
		userID,
		string(channel),
		strings.ToLower(strings.TrimSpace(destination)),
		strings.TrimSpace(code),
	}, "|")
	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func normalizeVerificationChannel(value string) (user.VerificationChannel, error) {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "", string(user.VerificationChannelEmail):
		return user.VerificationChannelEmail, nil
	case string(user.VerificationChannelSMS):
		return user.VerificationChannelSMS, nil
	case string(user.VerificationChannelWhatsApp):
		return user.VerificationChannelWhatsApp, nil
	case string(user.VerificationChannelPush):
		return user.VerificationChannelPush, nil
	default:
		return "", ErrInvalidVerificationChannel
	}
}

func verificationDestination(channel user.VerificationChannel, email string, destination string) (string, error) {
	switch channel {
	case user.VerificationChannelEmail:
		if strings.TrimSpace(destination) == "" {
			return strings.ToLower(strings.TrimSpace(email)), nil
		}
		return strings.ToLower(strings.TrimSpace(destination)), nil
	case user.VerificationChannelSMS, user.VerificationChannelWhatsApp:
		destination := strings.TrimSpace(destination)
		if destination == "" {
			return "", ErrVerificationDestinationMissing
		}
		return destination, nil
	case user.VerificationChannelPush:
		if strings.TrimSpace(destination) == "" {
			return strings.ToLower(strings.TrimSpace(email)), nil
		}
		return strings.TrimSpace(destination), nil
	default:
		return "", ErrInvalidVerificationChannel
	}
}

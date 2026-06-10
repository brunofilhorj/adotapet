package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	inport "adotapet/internal/app/port/in"
	outport "adotapet/internal/app/port/out"
	"adotapet/internal/domain/user"
)

var (
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired or revoked")
)

type RefreshTokenService struct {
	users         outport.UserRepository
	refreshTokens outport.RefreshTokenRepository
	accessTokens  AccessTokenIssuer
	refreshIssuer RefreshTokenIssuer
	now           func() time.Time
}

func NewRefreshTokenService(
	users outport.UserRepository,
	refreshTokens outport.RefreshTokenRepository,
	accessTokens AccessTokenIssuer,
	refreshIssuer RefreshTokenIssuer,
) RefreshTokenService {
	return RefreshTokenService{
		users:         users,
		refreshTokens: refreshTokens,
		accessTokens:  accessTokens,
		refreshIssuer: refreshIssuer,
		now:           time.Now,
	}
}

func (s RefreshTokenService) Refresh(ctx context.Context, cmd inport.RefreshTokenCommand) (inport.AuthTokens, error) {
	rawRefreshToken := strings.TrimSpace(cmd.RefreshToken)
	if rawRefreshToken == "" {
		return inport.AuthTokens{}, ErrInvalidRefreshToken
	}

	tokenHash := s.refreshIssuer.HashRefreshToken(rawRefreshToken)
	current, err := s.refreshTokens.FindByHash(ctx, tokenHash)
	if err != nil {
		return inport.AuthTokens{}, err
	}
	if current == nil {
		return inport.AuthTokens{}, ErrInvalidRefreshToken
	}
	if !current.IsUsable(s.now().UTC()) {
		return inport.AuthTokens{}, ErrRefreshTokenExpired
	}

	found, err := s.users.FindByID(ctx, current.UserID)
	if err != nil {
		return inport.AuthTokens{}, err
	}
	if found == nil {
		return inport.AuthTokens{}, ErrInvalidRefreshToken
	}
	if found.Status != user.StatusActive {
		return inport.AuthTokens{}, ErrAccountNotActive
	}

	if err := s.refreshTokens.Revoke(ctx, current.ID); err != nil {
		return inport.AuthTokens{}, err
	}

	return issueTokenPair(ctx, found, s.accessTokens, s.refreshIssuer, s.refreshTokens)
}

package bootstrap

import (
	"database/sql"
	"log/slog"

	httpauth "adotapet/internal/adapters/inbound/http/auth"
	notificationadapter "adotapet/internal/adapters/outbound/notification"
	postgresrepo "adotapet/internal/adapters/outbound/postgres/repository"
	authapp "adotapet/internal/app/auth"
)

func newAccessTokens(cfg Config) authapp.HMACJWTIssuer {
	return authapp.NewHMACJWTIssuer(cfg.JWTIssuer, cfg.JWTAccessSecret, cfg.JWTAccessTTL)
}

func newAuthService(db *sql.DB, cfg Config, log *slog.Logger, accessTokens authapp.AccessTokenIssuer) httpauth.Service {
	userRepository := postgresrepo.NewUserRepository(db)
	refreshTokenRepository := postgresrepo.NewRefreshTokenRepository(db)
	verificationCodeRepository := postgresrepo.NewVerificationCodeRepository(db)
	verificationSender := notificationadapter.NewLogVerificationSender(log)
	passwords := authapp.BcryptPasswordHasher{}
	refreshTokens := authapp.NewSecureRefreshTokenIssuer(cfg.JWTRefreshTTL)
	verificationCodes := authapp.NewNumericVerificationCodeIssuer(cfg.VerificationTTL)

	registerUsers := authapp.NewRegisterUserService(
		userRepository,
		verificationCodeRepository,
		passwords,
		verificationCodes,
		verificationSender,
	)
	loginUsers := authapp.NewLoginService(
		userRepository,
		refreshTokenRepository,
		passwords,
		accessTokens,
		refreshTokens,
	)
	refreshUserTokens := authapp.NewRefreshTokenService(
		userRepository,
		refreshTokenRepository,
		accessTokens,
		refreshTokens,
	)
	verifyAccounts := authapp.NewVerifyAccountService(
		userRepository,
		verificationCodeRepository,
		verificationCodes,
	)
	resendVerification := authapp.NewResendVerificationService(
		userRepository,
		verificationCodeRepository,
		verificationCodes,
		verificationSender,
	)

	return httpauth.NewService(registerUsers, loginUsers, refreshUserTokens, verifyAccounts, resendVerification)
}

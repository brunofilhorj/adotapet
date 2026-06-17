package bootstrap

import (
	"database/sql"
	"log/slog"

	"adotapet/internal/adapters/inbound/http/middleware"
	notificationadapter "adotapet/internal/adapters/outbound/notification"
	postgresrepo "adotapet/internal/adapters/outbound/postgres/repository"
	authapp "adotapet/internal/app/auth"
	inport "adotapet/internal/app/port/in"
)

type authServices struct {
	registerUsers      inport.RegisterUserInputPort
	loginUsers         inport.LoginInputPort
	refreshUserTokens  inport.RefreshTokenInputPort
	verifyAccounts     inport.VerifyAccountInputPort
	resendVerification inport.ResendVerificationInputPort
	accessTokens       middleware.AccessTokenVerifier
}

func newAuthServices(db *sql.DB, cfg Config, log *slog.Logger) authServices {
	userRepository := postgresrepo.NewUserRepository(db)
	refreshTokenRepository := postgresrepo.NewRefreshTokenRepository(db)
	verificationCodeRepository := postgresrepo.NewVerificationCodeRepository(db)
	verificationSender := notificationadapter.NewLogVerificationSender(log)
	passwords := authapp.BcryptPasswordHasher{}
	accessTokens := authapp.NewHMACJWTIssuer(cfg.JWTIssuer, cfg.JWTAccessSecret, cfg.JWTAccessTTL)
	refreshTokens := authapp.NewSecureRefreshTokenIssuer(cfg.JWTRefreshTTL)
	verificationCodes := authapp.NewNumericVerificationCodeIssuer(cfg.VerificationTTL)

	return authServices{
		registerUsers: authapp.NewRegisterUserService(
			userRepository,
			verificationCodeRepository,
			passwords,
			verificationCodes,
			verificationSender,
		),
		loginUsers: authapp.NewLoginService(
			userRepository,
			refreshTokenRepository,
			passwords,
			accessTokens,
			refreshTokens,
		),
		refreshUserTokens: authapp.NewRefreshTokenService(
			userRepository,
			refreshTokenRepository,
			accessTokens,
			refreshTokens,
		),
		verifyAccounts: authapp.NewVerifyAccountService(
			userRepository,
			verificationCodeRepository,
			verificationCodes,
		),
		resendVerification: authapp.NewResendVerificationService(
			userRepository,
			verificationCodeRepository,
			verificationCodes,
			verificationSender,
		),
		accessTokens: accessTokens,
	}
}

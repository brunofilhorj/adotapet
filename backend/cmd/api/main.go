package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpadapter "adotapet/internal/adapters/inbound/http"
	notificationadapter "adotapet/internal/adapters/outbound/notification"
	postgresadapter "adotapet/internal/adapters/outbound/postgres"
	postgresrepo "adotapet/internal/adapters/outbound/postgres/repository"
	authapp "adotapet/internal/app/auth"
	"adotapet/internal/config"
	"adotapet/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	db, err := postgresadapter.Open(context.Background(), cfg)
	if err != nil {
		log.Error("database connection failed", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	userRepository := postgresrepo.NewUserRepository(db)
	refreshTokenRepository := postgresrepo.NewRefreshTokenRepository(db)
	verificationCodeRepository := postgresrepo.NewVerificationCodeRepository(db)
	verificationSender := notificationadapter.NewLogVerificationSender(log)
	passwords := authapp.BcryptPasswordHasher{}
	accessTokens := authapp.NewHMACJWTIssuer(cfg.JWTIssuer, cfg.JWTAccessSecret, cfg.JWTAccessTTL)
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

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpadapter.NewRouter(cfg, log, registerUsers, loginUsers, refreshUserTokens, verifyAccounts, resendVerification, accessTokens),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("starting api", slog.String("addr", cfg.HTTPAddr), slog.String("env", cfg.AppEnv))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("api stopped unexpectedly", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("api shutdown failed", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("api stopped")
}

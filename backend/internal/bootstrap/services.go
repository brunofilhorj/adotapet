package bootstrap

import (
	"database/sql"
	"log/slog"

	httpconversations "adotapet/internal/adapters/inbound/http/conversations"
	httpfavorites "adotapet/internal/adapters/inbound/http/favorites"
	httpmedia "adotapet/internal/adapters/inbound/http/media"
	"adotapet/internal/adapters/inbound/http/middleware"
	httppuppies "adotapet/internal/adapters/inbound/http/puppies"
	httpusers "adotapet/internal/adapters/inbound/http/users"
	"adotapet/internal/adapters/inbound/http/webserver"
	"adotapet/internal/adapters/inbound/ws/chat"
	postgresrepo "adotapet/internal/adapters/outbound/postgres/repository"
	puppiesapp "adotapet/internal/app/puppies"
	usersapp "adotapet/internal/app/users"
)

// Service exposes the HTTP routes owned by an application feature.
type Service = webserver.RouteProvider

func newServices(resources resources, cfg Config, log *slog.Logger) []Service {
	accessTokens := newAccessTokens(cfg)

	return []Service{
		newAuthService(resources.database, cfg, log, accessTokens),
		newUsersService(resources.database, accessTokens),
		newMediaService(accessTokens),
		newPuppiesService(resources.database, accessTokens),
		newFavoritesService(accessTokens),
		newConversationsService(accessTokens),
		newChatService(accessTokens),
	}
}

func newUsersService(db *sql.DB, accessTokens middleware.AccessTokenVerifier) Service {
	userRepository := postgresrepo.NewUserRepository(db)
	profileRepository := postgresrepo.NewProfileRepository(db)
	profiles := usersapp.NewProfileService(userRepository, profileRepository)
	return httpusers.NewService(accessTokens, profiles)
}

func newMediaService(accessTokens middleware.AccessTokenVerifier) Service {
	return httpmedia.NewService(accessTokens)
}

func newPuppiesService(db *sql.DB, accessTokens middleware.AccessTokenVerifier) Service {
	puppyRepository := postgresrepo.NewPuppyRepository(db)
	listings := puppiesapp.NewListingService(puppyRepository)
	return httppuppies.NewService(accessTokens, listings)
}

func newFavoritesService(accessTokens middleware.AccessTokenVerifier) Service {
	return httpfavorites.NewService(accessTokens)
}

func newConversationsService(accessTokens middleware.AccessTokenVerifier) Service {
	return httpconversations.NewService(accessTokens)
}

func newChatService(accessTokens middleware.AccessTokenVerifier) Service {
	return chat.NewService(accessTokens)
}

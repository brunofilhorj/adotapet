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
	usersapp "adotapet/internal/app/users"
)

// Service exposes the HTTP routes owned by an application feature.
type Service = webserver.RouteProvider

func newServices(resources resources, cfg Config, log *slog.Logger) []Service {
	accessTokens := newAccessTokens(cfg)

	return []Service{
		newAuthService(resources.database, cfg, log, accessTokens),
		newUsersService(resources.database, accessTokens),
		newMediaService(),
		newPuppiesService(),
		newFavoritesService(),
		newConversationsService(),
		newChatService(),
	}
}

func newUsersService(db *sql.DB, accessTokens middleware.AccessTokenVerifier) Service {
	userRepository := postgresrepo.NewUserRepository(db)
	profileRepository := postgresrepo.NewProfileRepository(db)
	profiles := usersapp.NewProfileService(userRepository, profileRepository)
	return httpusers.NewService(accessTokens, profiles)
}

func newMediaService() Service {
	return httpmedia.NewService()
}

func newPuppiesService() Service {
	return httppuppies.NewService()
}

func newFavoritesService() Service {
	return httpfavorites.NewService()
}

func newConversationsService() Service {
	return httpconversations.NewService()
}

func newChatService() Service {
	return chat.NewService()
}

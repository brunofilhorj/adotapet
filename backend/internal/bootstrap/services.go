package bootstrap

import (
	"log/slog"

	"adotapet/internal/adapters/inbound/http/conversations"
	"adotapet/internal/adapters/inbound/http/favorites"
	"adotapet/internal/adapters/inbound/http/media"
	"adotapet/internal/adapters/inbound/http/middleware"
	"adotapet/internal/adapters/inbound/http/puppies"
	"adotapet/internal/adapters/inbound/http/users"
	"adotapet/internal/adapters/inbound/http/webserver"
	"adotapet/internal/adapters/inbound/ws/chat"
)

// Service exposes the HTTP routes owned by an application feature.
type Service = webserver.RouteProvider

func newServices(resources resources, cfg Config, log *slog.Logger) []Service {
	accessTokens := newAccessTokens(cfg)

	return []Service{
		newAuthService(resources.database, cfg, log, accessTokens),
		newUsersService(accessTokens),
		newMediaService(),
		newPuppiesService(),
		newFavoritesService(),
		newConversationsService(),
		newChatService(),
	}
}

func newUsersService(accessTokens middleware.AccessTokenVerifier) Service {
	return users.NewService(accessTokens)
}

func newMediaService() Service {
	return media.NewService()
}

func newPuppiesService() Service {
	return puppies.NewService()
}

func newFavoritesService() Service {
	return favorites.NewService()
}

func newConversationsService() Service {
	return conversations.NewService()
}

func newChatService() Service {
	return chat.NewService()
}

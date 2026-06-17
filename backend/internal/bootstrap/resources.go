package bootstrap

import (
	"context"
	"database/sql"

	postgresadapter "adotapet/internal/adapters/outbound/postgres"
)

type resources struct {
	database *sql.DB
}

func prepareResources(ctx context.Context, cfg Config) (resources, error) {
	db, err := postgresadapter.Open(ctx, cfg)
	if err != nil {
		return resources{}, err
	}

	return resources{
		database: db,
	}, nil
}

func (r resources) all() []resource {
	return []resource{
		databaseResource{db: r.database},
	}
}

type databaseResource struct {
	db *sql.DB
}

func (r databaseResource) Shutdown(ctx context.Context) error {
	return r.db.Close()
}

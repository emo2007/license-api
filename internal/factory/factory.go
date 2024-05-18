package factory

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/emo2007/block-accounting/examples/license-api/internal/interface/rest"
	"github.com/emo2007/block-accounting/examples/license-api/internal/pkg/config"
	"github.com/emo2007/block-accounting/examples/license-api/internal/usecases/repository"

	_ "github.com/lib/pq"
)

func NewService(
	log *slog.Logger,
	conf config.Config,
) (*rest.Server, func(), error) {
	db, f, err := ProvideDatabaseConnection(conf)
	if err != nil {
		return nil, func() {}, err
	}

	return rest.NewServer(log, conf.Rest, repository.NewRepository(db)), f, nil
}

func ProvideDatabaseConnection(c config.Config) (*sql.DB, func(), error) {
	sslmode := "disable"
	if c.DB.EnableSSL {
		sslmode = "enable"
	}

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=%s",
		c.DB.User, c.DB.Secret, c.DB.Host, c.DB.Database, sslmode,
	)

	fmt.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, func() {}, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, func() {
		db.Close()
	}, nil
}

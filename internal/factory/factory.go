package factory

import (
	"log/slog"

	"github.com/emo2007/block-accounting/examples/license-api/internal/interface/rest"
	"github.com/emo2007/block-accounting/examples/license-api/internal/pkg/config"
)

func NewService(
	log *slog.Logger,
	conf config.RestConfig,
) *rest.Server {
	return rest.NewServer(log, conf)
}

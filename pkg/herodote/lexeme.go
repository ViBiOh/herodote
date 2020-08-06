package herodote

import (
	"context"
	"time"

	"github.com/ViBiOh/httputils/v3/pkg/logger"
)

const refreshMaterializedView = `REFRESH MATERIALIZED VIEW herodote.lexeme`

func (a app) refreshLexeme(_ time.Time) error {
	logger.Info("Refreshing lexeme")
	_, err := a.db.ExecContext(context.Background(), refreshMaterializedView)
	return err
}

package source

import (
	"context"
	"time"

	"github.com/havocked/leipzig-cli/internal/model"
)

type Source interface {
	ID() string
	Fetch(ctx context.Context, from, to time.Time) ([]model.Event, error)
}

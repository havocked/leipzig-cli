package engine

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/havocked/leipzig-cli/internal/model"
	"github.com/havocked/leipzig-cli/internal/source"
)

type Engine struct {
	sources []source.Source
}

func New(sources ...source.Source) *Engine {
	return &Engine{sources: sources}
}

func (e *Engine) Fetch(ctx context.Context, from, to time.Time) ([]model.Event, error) {
	var all []model.Event

	for _, src := range e.sources {
		events, err := src.Fetch(ctx, from, to)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: source %s failed: %v\n", src.ID(), err)
			continue
		}
		all = append(all, events...)
	}

	all = Dedup(all)

	// Populate map URLs
	for i := range all {
		all[i].MapURL = all[i].MapsURL()
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].StartTime.Before(all[j].StartTime)
	})

	return all, nil
}

func (e *Engine) Sources() []source.Source {
	return e.sources
}

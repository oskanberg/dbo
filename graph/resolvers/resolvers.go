package resolvers

import (
	"context"
	"errors"
	"time"

	"github.com/oskanberg/dbo/boltstore"
	"github.com/oskanberg/dbo/graph"
	"github.com/oskanberg/dbo/graph/model"
)

// NewRootResolver returns a new root resolver
func NewRootResolver(store *boltstore.Store) *Resolver {
	return &Resolver{
		Store: store,
	}
}

// Resolver is the top level resolver
type Resolver struct {
	Store *boltstore.Store
}

// Query returns the query resolver
func (r *Resolver) Query() graph.QueryResolver {
	return &QueryResolver{r}
}

func (r *Resolver) Film() graph.FilmResolver {
	return &QueryResolver{r}
}

// QueryResolver resolves films
type QueryResolver struct{ *Resolver }

// GetFilm resolves a film by id
func (r *QueryResolver) GetFilm(ctx context.Context, id string) (*model.Film, error) {
	f, err := r.Store.GetDetails(id)
	return &f, err
}

// GrossDaily gets daily gross between two dates
func (r *QueryResolver) GrossDaily(ctx context.Context, film *model.Film, from, to *time.Time) ([]model.DailyGross, error) {
	today := time.Now().Truncate(24 * time.Hour)

	if from == nil {
		from = &today
	}

	if to == nil {
		to = &today
	}

	if to.Before(*from) {
		return nil, errors.New("from must be before to")
	}

	return r.Store.GetGross(film.ID, *from, *to)
}

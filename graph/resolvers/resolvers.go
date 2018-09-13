package resolvers

import (
	"context"

	"github.com/oskanberg/dbo/boltstore"
	"github.com/oskanberg/dbo/graph"
	"github.com/oskanberg/dbo/model"
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

// QueryResolver resolves films
type QueryResolver struct{ *Resolver }

// GetFilm resolves a film by id
func (r *QueryResolver) GetFilm(ctx context.Context, id string) (*model.Film, error) {
	film, err := r.Store.GetDetails(id)
	return &film, err
}

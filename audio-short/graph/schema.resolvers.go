package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/kv9c12/audio-shorts/graph/generated"
	"github.com/kv9c12/audio-shorts/graph/model"
	"github.com/kv9c12/audio-shorts/middleware"
)

func (r *mutationResolver) UploadShort(ctx context.Context, input model.NewShort) (string, error) {
	id, err := middleware.CreateShort(input)
	if err != nil {
		return "0", err
	}

	return strconv.Itoa(id), nil
}

func (r *queryResolver) GetShortsByPage(ctx context.Context, page int) ([]*model.Short, error) {
	shorts, err := middleware.GetAllShorts(page)

	if err != nil {
		return nil, err
	}

	var graphqlShorts []*model.Short

	for i := 0; i < len(shorts); i++ {
		var graphqlShort model.Short
		graphqlShort.Title = shorts[i].Title
		graphqlShort.Description = shorts[i].Description
		graphqlShort.Category = shorts[i].Category
		graphqlShort.ID = shorts[i].ID
		json.Unmarshal(shorts[i].Creator, &graphqlShort.Creator)
		graphqlShorts = append(graphqlShorts, &graphqlShort)
	}

	return graphqlShorts, nil
}

func (r *queryResolver) GetShortByID(ctx context.Context, id int) (*model.Short, error) {
	short, err := middleware.GetShort(id)

	if err != nil {
		return nil, err
	}

	var graphqlShort model.Short

	graphqlShort.Title = short.Title
	graphqlShort.Description = short.Description
	graphqlShort.Category = short.Category
	graphqlShort.ID = short.ID
	json.Unmarshal(short.Creator, &graphqlShort.Creator)

	return &graphqlShort, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

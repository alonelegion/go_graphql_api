package graph

import (
	"context"
	"errors"
	"github.com/alonelegion/go_graphql_api/graph/model"
)

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	user, err := r.UserService.GetByID(uint(id))
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		Active:    user.Active,
	}, nil
}

func (r *queryResolver) UserProfile(ctx context.Context) (*model.User, error) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return nil, errors.New("Unauthorized: Token is invlaid")
	}

	user, err := r.UserService.GetByID(userID.(uint))
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:        int(user.ID),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		Active:    user.Active,
	}, nil
}

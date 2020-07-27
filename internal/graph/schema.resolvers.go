package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"github.com/alonelegion/go_graphql_api/internal/graph/generated"
)

//func (r *mutationResolver) Register(ctx context.Context, input model.RegisterLogin) (*model.RegisterLoginOutput, error) {
//	panic(fmt.Errorf("not implemented"))
//}
//
//func (r *mutationResolver) Login(ctx context.Context, input model.RegisterLogin) (*model.RegisterLoginOutput, error) {
//	panic(fmt.Errorf("not implemented"))
//}
//
//func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUser) (*model.User, error) {
//	panic(fmt.Errorf("not implemented"))
//}
//
//func (r *mutationResolver) ForgotPassword(ctx context.Context, email string) (bool, error) {
//	panic(fmt.Errorf("not implemented"))
//}
//
//func (r *mutationResolver) ResetPassword(ctx context.Context, resetToken string, password string) (*model.RegisterLoginOutput, error) {
//	panic(fmt.Errorf("not implemented"))
//}
//
//func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
//	panic(fmt.Errorf("not implemented"))
//}
//
//func (r *queryResolver) UserProfile(ctx context.Context) (*model.User, error) {
//	panic(fmt.Errorf("not implemented"))
//}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

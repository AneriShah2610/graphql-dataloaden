//go:generate gorunpkg github.com/99designs/gqlgen

package handler

import (
	context "context"
	"log"
	"net/http"

	"github.com/aneri/graphql-dataloaden/dal"
	"github.com/aneri/graphql-dataloaden/dataloader"
	graph "github.com/aneri/graphql-dataloaden/graph"
	model "github.com/aneri/graphql-dataloaden/model"
)

type Resolver struct{}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var ctx context.Context

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func DbMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		crConn, err := dal.DbConnect()
		if err != nil {
			log.Fatal(err)
		}
		ctx = context.WithValue(request.Context(), "crConn", crConn)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
func (r *mutationResolver) NewUser(ctx context.Context, input model.CreateUser) (model.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) NewAppliation(ctx context.Context, input model.CreateApplication) (model.Application, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) (model.User, error) {
	panic("not implemented")
}
func (r *queryResolver) UserByID(ctx context.Context, userId int) (model.User, error) {
	user, err := dataloader.CtxLoaders(ctx).UserByID.Load(userId)
	if err != nil {
		log.Println("Error to fetch user data")
	}
	return *user, nil
}
func (r *queryResolver) UserByIds(ctx context.Context, userIds []int) ([]model.User, error) {
	panic("not implemented")
}
func (r *queryResolver) Applications(ctx context.Context) (model.Application, error) {
	panic("not implemented")
}
func (r *queryResolver) ApplicationByUserID(ctx context.Context, userId int) (model.Application, error) {
	panic("not implemented")
}
func (r *queryResolver) ApplicationByUserIds(ctx context.Context, userIds []int) (model.Application, error) {
	panic("not implemented")
}

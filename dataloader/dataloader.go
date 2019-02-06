package dataloader

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aneri/graphql-dataloaden/dal"
	"github.com/aneri/graphql-dataloaden/model"
	"github.com/jmoiron/sqlx"
)

type ctxKeyType struct{ name string }

var ctxKey = ctxKeyType{"userCtx"}

type loaders struct {
	UserByID        *UserLoader
	ApplicationByID *ApplicationLoader
}

func LoadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ldrs := loaders{}
		ctx := r.Context()
		crConn := ctx.Value("crConn").(*dal.DbConnection)
		// Set waiting time
		wait := 350 * time.Microsecond

		// User Loader
		ldrs.UserByID = &UserLoader{
			wait:     wait,
			maxBatch: 200,
			fetch: func(ids []int) ([]*model.User, []error) {
				var sqlQuery string
				if len(ids) == 1 {
					sqlQuery = "SELECT id, name, email, contact, createat FROM user_mst WHERE id = ?"
				} else {
					sqlQuery = "SELECT id, name, email, contact, createat FROM user_mst WHERE id IN (?)"
				}
				sqlQuery, arguments, err := sqlx.In(sqlQuery, ids)
				if err != nil {
					log.Println("Error to fetch user query")
				}
				sqlQuery = sqlx.Rebind(sqlx.DOLLAR, sqlQuery)
				rows, err := crConn.Db.Query(sqlQuery, arguments...)
				if err != nil {
					log.Println("Error to read user data")
				}
				userById := map[int]*model.User{}
				defer func() {
					if err := rows.Close(); err != nil {
						log.Println("Error at defer function of user")
					}
				}()
				for rows.Next() {
					var user model.User
					err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Contact, &user.CreateAt)
					if err != nil {
						log.Println("Error to scan user data")
					}
					// Store user.ID in key & model of user in value
					userById[user.ID] = &user
				}
				users := make([]*model.User, len(ids))
				for i, id := range ids {
					users[i] = userById[id]
					i++
				}
				return users, nil
			},
		}

		ctx = context.WithValue(ctx, ctxKey, ldrs)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func CtxLoaders(ctx context.Context) loaders {
	return ctx.Value(ctxKey).(loaders)
}

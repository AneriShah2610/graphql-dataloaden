package model

import (
	"fmt"
	"io"
	strconv "strconv"
	time "time"

	graphql "github.com/99designs/gqlgen/graphql"
)

type Application struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	AppliedBy   int        `json:"appliedBy"`
	AppliedAt   time.Time  `json:"appliedAt"`
	VerifiedBy  *int       `json:"verifiedBy"`
	VerifiedAt  *time.Time `json:"verifiedAt"`
}

type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Contact  string    `json:"contact"`
	CreateAt time.Time `json:"createAt"`
}

type CreateApplication struct {
	Description string `json:"description"`
	AppliedBy   int    `json:"appliedBy"`
	VerifiedBy  *int   `json:"verifiedBy"`
}

type CreateUser struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
}

// Lets redefine the base ID type to use an id from an external library
func MarshalID(id int) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", id)))
	})
}

// And the same for the unmarshaler
func UnmarshalID(v interface{}) (int, error) {
	id, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("ids must be strings")
	}
	return strconv.Atoi(id)
}

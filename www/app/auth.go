package app

import (
	"errors"
	"net/http"
)

type User struct {
	id string
}

func (u User) ID() string {
	return u.id
}

//pacis:authentication
func AuthHandler(r *http.Request) (*User, error) {
	// return &User{id: "test-user"}
	return nil, errors.New("test error")
}

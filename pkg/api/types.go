package api

import (
	"net/http"

	"github.com/lsowen/hoover-dam/pkg/api/service"
)

type GroupCreation service.GroupCreation

func (g *GroupCreation) Bind(r *http.Request) error {
	return nil
}

type Policy service.Policy

func (p *Policy) Bind(r *http.Request) error {
	return nil
}

type UserCreation service.UserCreation

func (u *UserCreation) Bind(r *http.Request) error {
	return nil
}

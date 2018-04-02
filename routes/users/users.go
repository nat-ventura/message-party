package users

import (
	"github.com/go-kit/kit/log"
	_ "github.com/julienschmidt/httprouter"
	_ "net/http"
)

type User struct {
	logger log.Logger
}

func NewServer(logger log.Logger) *User {
	return &User{
		logger: logger,
	}
}

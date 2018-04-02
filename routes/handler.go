package routes

import (
	"github.com/go-kit/kit/log"

	"github.com/nat/socket-party/routes/users"
)

type Handler struct {
	Logger log.Logger
	User   *users.User
}

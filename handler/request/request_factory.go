package request

import (
	"github.com/go-playground/validator"
	appSrv "github.com/mamadeusia/RequestSrv/service/application"
)

type RequestHandler struct {
	service  *appSrv.ApplicationService
	Validate *validator.Validate
}

func NewRequestHandler(srv *appSrv.ApplicationService, validate *validator.Validate) *RequestHandler {
	return &RequestHandler{
		service:  srv,
		Validate: validate,
	}
}

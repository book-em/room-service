package internal

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) Handler {
	return Handler{s}
}

type Route struct {
	handler Handler
}

func NewRoute(handler Handler) *Route {
	return &Route{handler}
}

func (r *Route) Route(rg *gin.RouterGroup) {

}

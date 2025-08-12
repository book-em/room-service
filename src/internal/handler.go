package internal

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/util"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Route struct{ handler Handler }

func NewRoute(handler Handler) *Route { return &Route{handler} }

func (r *Route) Route(rg *gin.RouterGroup) {
	rg.POST("/new", r.handler.createRoom)
	rg.GET("/:id", r.handler.findRoomById)
	rg.GET("/host/:id", r.handler.findRoomByHostId)
}

type Handler struct{ service Service }

func NewHandler(s Service) Handler { return Handler{s} }

func (h *Handler) createRoom(ctx *gin.Context) {
	jwt, err := util.GetJwt(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	if jwt.Role != userclient.Host {
		ctx.Error(fmt.Errorf("Unauthorized"))
		return
	}

	var dto CreateRoomDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	room, err := h.service.Create(dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, NewRoomDTO(room))
}

func (h *Handler) findRoomById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID: %s", err.Error())
		ctx.Error(err)
		return
	}

	log.Printf("Find room by id %d", id)

	room, err := h.service.FindById(uint(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomDTO(room))
}

func (h *Handler) findRoomByHostId(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID: %s", err.Error())
		ctx.Error(err)
		return
	}

	log.Printf("Find room by id %d", id)

	rooms, err := h.service.FindByHost(uint(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	result := make([]RoomDTO, 0)
	for _, room := range rooms {
		result = append(result, NewRoomDTO(&room))
	}

	ctx.JSON(http.StatusOK, result)
}

package internal

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/util"
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
	rg.GET("/host/:id", r.handler.findRoomsByHostId)
	rg.GET("/all", r.handler.findAvailableRooms)

	rg.GET("/available/room/:id", r.handler.findCurrentAvailabilityListOfRoom)
	rg.GET("/available/room/all/:id", r.handler.findAvailabilityListsByRoomId)
	rg.GET("/available/:id", r.handler.findAvailabilityListById)
	rg.POST("/available", r.handler.updateAvailability)

	rg.GET("/price/room/:id", r.handler.findCurrentPriceListOfRoom)
	rg.GET("/price/room/all/:id", r.handler.findPriceListsByRoomId)
	rg.GET("/price/:id", r.handler.findPriceListById)
	rg.POST("/price", r.handler.updatePriceList)
}

type Handler struct{ service Service }

func NewHandler(s Service) Handler { return Handler{s} }

func (h *Handler) createRoom(ctx *gin.Context) {
	jwt, err := util.GetJwt(ctx)
	if err != nil {
		AbortError(ctx, ErrUnauthenticated)
		return
	}

	if jwt.Role != userclient.Host {
		AbortError(ctx, ErrUnauthorized)
		return
	}

	var dto CreateRoomDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		AbortError(ctx, err)
		return
	}

	room, err := h.service.Create(jwt.ID, dto)
	if err != nil {
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, NewRoomDTO(room))
}

func (h *Handler) findRoomById(ctx *gin.Context) {
	log.Printf("findRoomById called")
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID %s: %s", ctx.Param("id"), err.Error())
		AbortError(ctx, ErrBadRequest)
		return
	}

	log.Printf("Find room by id %d", id)

	room, err := h.service.FindById(uint(id))
	if err != nil {
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomDTO(room))
}

func (h *Handler) findRoomsByHostId(ctx *gin.Context) {
	log.Printf("findRoomsByHostId called")

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID %s: %s", ctx.Param("id"), err.Error())
		AbortError(ctx, ErrBadRequest)
		return
	}

	rooms, err := h.service.FindByHost(uint(id))
	if err != nil {
		AbortError(ctx, err)
		return
	}

	result := make([]RoomDTO, 0)
	for _, room := range rooms {
		result = append(result, NewRoomDTO(&room))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findCurrentAvailabilityListOfRoom(ctx *gin.Context) {
	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID %s: %s", ctx.Param("id"), err.Error())
		AbortError(ctx, ErrBadRequest)
		return
	}

	list, err := h.service.FindCurrentAvailabilityListOfRoom(uint(roomId))
	if err != nil {
		AbortError(ctx, err)
		return
	}

	result := NewRoomAvailabilityListDTO(list)
	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findAvailabilityListsByRoomId(ctx *gin.Context) {
	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID %s: %s", ctx.Param("id"), err.Error())
		AbortError(ctx, ErrBadRequest)
		return
	}

	lists, err := h.service.FindAvailabilityListsByRoomId(uint(roomId))
	if err != nil {
		AbortError(ctx, err)
		return
	}

	result := make([]RoomAvailabilityListDTO, 0)
	for _, list := range lists {
		result = append(result, NewRoomAvailabilityListDTO(&list))
	}
	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findAvailabilityListById(ctx *gin.Context) {
	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse ID %s: %s", ctx.Param("id"), err.Error())
		AbortError(ctx, ErrBadRequest)
		return
	}

	list, err := h.service.FindAvailabilityListById(uint(listId))
	if err != nil {
		AbortError(ctx, err)
		return
	}

	result := NewRoomAvailabilityListDTO(list)
	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) updateAvailability(ctx *gin.Context) {
	jwt, err := util.GetJwt(ctx)
	if err != nil {
		AbortError(ctx, ErrUnauthenticated)
		return
	}

	if jwt.Role != userclient.Host {
		AbortError(ctx, ErrUnauthorized)
		return
	}

	var dto CreateRoomAvailabilityListDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		AbortError(ctx, err)
		return
	}

	list, err := h.service.UpdateAvailability(jwt.ID, dto)
	if err != nil {
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, NewRoomAvailabilityListDTO(list))
}

func (h *Handler) findCurrentPriceListOfRoom(ctx *gin.Context) {
	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse room ID %s: %s", ctx.Param("id"), err.Error())
		AbortError(ctx, ErrBadRequest)
		return
	}

	list, err := h.service.FindCurrentPriceListOfRoom(uint(roomId))
	if err != nil {
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomPriceListDTO(list))
}

func (h *Handler) findPriceListsByRoomId(ctx *gin.Context) {
	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse room ID %s: %s", ctx.Param("id"), err.Error())
		AbortError(ctx, ErrBadRequest)
		return
	}

	lists, err := h.service.FindPriceListsByRoomId(uint(roomId))
	if err != nil {
		AbortError(ctx, err)
		return
	}

	result := make([]RoomPriceListDTO, 0)
	for _, list := range lists {
		result = append(result, NewRoomPriceListDTO(&list))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findPriceListById(ctx *gin.Context) {
	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Printf("Could not parse price list ID %s: %s", ctx.Param("id"), err.Error())
		AbortError(ctx, ErrBadRequest)
		return
	}

	list, err := h.service.FindPriceListById(uint(listId))
	if err != nil {
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomPriceListDTO(list))
}

func (h *Handler) updatePriceList(ctx *gin.Context) {
	jwt, err := util.GetJwt(ctx)
	if err != nil {
		AbortError(ctx, ErrUnauthenticated)
		return
	}

	if jwt.Role != userclient.Host {
		AbortError(ctx, ErrUnauthorized)
		return
	}

	var dto CreateRoomPriceListDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		AbortError(ctx, err)
		return
	}

	list, err := h.service.UpdatePriceList(jwt.ID, dto)
	if err != nil {
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, NewRoomPriceListDTO(list))
}

func (h *Handler) findAvailableRooms(ctx *gin.Context) {

	var dto RoomsQueryDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		AbortError(ctx, err)
		return
	}

	rooms, resultInfo, err := h.service.FindAvailableRooms(dto)
	if err != nil {
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomsResultDTO(rooms, *resultInfo))
}

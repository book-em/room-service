package internal

import (
	"bookem-room-service/util"
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

	rg.POST("/reservation/query", r.handler.queryForReservation)
}

type Handler struct{ service Service }

func NewHandler(s Service) Handler { return Handler{s} }

func (h *Handler) createRoom(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "create-room-api")
	defer util.TEL.Pop()

	jwt, err := util.GetJwt(ctx)
	if err != nil {
		util.TEL.Event("failed fetching JWT", err)
		AbortError(ctx, ErrUnauthenticated)
		return
	}

	if jwt.Role != util.Host {
		util.TEL.Event("user is not host", nil)
		AbortError(ctx, ErrUnauthorized)
		return
	}

	var dto CreateRoomDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		util.TEL.Event("failed binding JSON", err)
		AbortError(ctx, err)
		return
	}

	room, err := h.service.Create(util.TEL.Ctx(), jwt.ID, dto)
	if err != nil {
		util.TEL.Event("failed creating room", err)
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, NewRoomDTO(room))
}

func (h *Handler) findRoomById(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-room-by-id-api")
	defer util.TEL.Pop()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.TEL.Eventf("could not parse ID %s", err, ctx.Param("id"))
		AbortError(ctx, ErrBadRequest)
		return
	}

	room, err := h.service.FindById(util.TEL.Ctx(), uint(id))
	if err != nil {
		util.TEL.Eventf("failed finding room by id %d", err, id)
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomDTO(room))
}

func (h *Handler) findRoomsByHostId(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-rooms-by-host-id-api")
	defer util.TEL.Pop()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.TEL.Eventf("could not parse ID %s", err, ctx.Param("id"))
		AbortError(ctx, ErrBadRequest)
		return
	}

	rooms, err := h.service.FindByHost(util.TEL.Ctx(), uint(id))
	if err != nil {
		util.TEL.Eventf("could not find rooms by host with ID %d", err, id)
		AbortError(ctx, err)
		return
	}

	util.TEL.Eventf("creating json output with %d rooms", nil, len(rooms))
	result := make([]RoomDTO, 0)
	for _, room := range rooms {
		result = append(result, NewRoomDTO(&room))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findCurrentAvailabilityListOfRoom(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-current-availability-list-of-room-api")
	defer util.TEL.Pop()

	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.TEL.Eventf("could not parse ID %s", err, ctx.Param("id"))
		AbortError(ctx, ErrBadRequest)
		return
	}

	list, err := h.service.FindCurrentAvailabilityListOfRoom(util.TEL.Ctx(), uint(roomId))
	if err != nil {
		util.TEL.Eventf("could not get current availability list of room with ID %d", err, roomId)
		AbortError(ctx, err)
		return
	}

	result := NewRoomAvailabilityListDTO(list)
	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findAvailabilityListsByRoomId(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-availability-lists-by-room-api")
	defer util.TEL.Pop()

	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.TEL.Eventf("could not parse ID %s", err, ctx.Param("id"))
		AbortError(ctx, ErrBadRequest)
		return
	}

	lists, err := h.service.FindAvailabilityListsByRoomId(util.TEL.Ctx(), uint(roomId))
	if err != nil {
		util.TEL.Eventf("could not find availability lists of rooms with ID %d", err, roomId)
		AbortError(ctx, err)
		return
	}

	util.TEL.Eventf("creating json output with %d lists", nil, len(lists))
	result := make([]RoomAvailabilityListDTO, 0)
	for _, list := range lists {
		result = append(result, NewRoomAvailabilityListDTO(&list))
	}
	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findAvailabilityListById(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-availability-list-by-id-api")
	defer util.TEL.Pop()

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.TEL.Eventf("could not parse ID %s", err, ctx.Param("id"))
		AbortError(ctx, ErrBadRequest)
		return
	}

	list, err := h.service.FindAvailabilityListById(util.TEL.Ctx(), uint(listId))
	if err != nil {
		util.TEL.Eventf("could not find availability list with ID %d", err, listId)
		AbortError(ctx, err)
		return
	}

	result := NewRoomAvailabilityListDTO(list)
	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) updateAvailability(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "update-room-availability-api")
	defer util.TEL.Pop()

	jwt, err := util.GetJwt(ctx)
	if err != nil {
		util.TEL.Event("could not get JWT", err)
		AbortError(ctx, ErrUnauthenticated)
		return
	}

	if jwt.Role != util.Host {
		util.TEL.Event("user is not host", nil)
		AbortError(ctx, ErrUnauthorized)
		return
	}

	var dto CreateRoomAvailabilityListDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		util.TEL.Event("failed to bind JSON", err)
		AbortError(ctx, err)
		return
	}

	list, err := h.service.UpdateAvailability(util.TEL.Ctx(), jwt.ID, dto)
	if err != nil {
		util.TEL.Event("could not update room availability", err)
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, NewRoomAvailabilityListDTO(list))
}

func (h *Handler) queryForReservation(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "query-room-for-reservation-api")
	defer util.TEL.Pop()

	jwt, err := util.GetJwt(ctx)
	if err != nil {
		util.TEL.Event("could not get JWT", err)
		AbortError(ctx, ErrUnauthenticated)
		return
	}

	if jwt.Role != util.Guest {
		util.TEL.Eventf("user is not guest (role = %s)", nil, jwt.Role)
		AbortError(ctx, ErrUnauthorized)
		return
	}

	var dto RoomReservationQueryDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		util.TEL.Event("failed to bind JSON", err)
		AbortError(ctx, err)
		return
	}

	result, err := h.service.QueryForReservation(util.TEL.Ctx(), jwt.ID, dto)
	if err != nil {
		util.TEL.Event("could not query room for availability", err)
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findCurrentPriceListOfRoom(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-current-price-list-of-room-api")
	defer util.TEL.Pop()

	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.TEL.Eventf("could not parse ID %s", err, ctx.Param("id"))
		AbortError(ctx, ErrBadRequest)
		return
	}

	list, err := h.service.FindCurrentPriceListOfRoom(util.TEL.Ctx(), uint(roomId))
	if err != nil {
		util.TEL.Eventf("could not get current price list of room with ID %d", err, roomId)
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomPriceListDTO(list))
}

func (h *Handler) findPriceListsByRoomId(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-price-lists-by-room-api")
	defer util.TEL.Pop()

	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.TEL.Eventf("could not parse ID %s", err, ctx.Param("id"))
		AbortError(ctx, ErrBadRequest)
		return
	}

	lists, err := h.service.FindPriceListsByRoomId(util.TEL.Ctx(), uint(roomId))
	if err != nil {
		util.TEL.Eventf("could not find price lists of rooms with ID %d", err, roomId)
		AbortError(ctx, err)
		return
	}

	util.TEL.Eventf("creating json output with %d lists", nil, len(lists))
	result := make([]RoomPriceListDTO, 0)
	for _, list := range lists {
		result = append(result, NewRoomPriceListDTO(&list))
	}

	ctx.JSON(http.StatusOK, result)
}

func (h *Handler) findPriceListById(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-price-list-by-id-api")
	defer util.TEL.Pop()

	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.TEL.Eventf("could not parse ID %s", err, ctx.Param("id"))
		AbortError(ctx, ErrBadRequest)
		return
	}

	list, err := h.service.FindPriceListById(util.TEL.Ctx(), uint(listId))
	if err != nil {
		util.TEL.Eventf("could not find price list with ID %d", err, listId)
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomPriceListDTO(list))
}

func (h *Handler) updatePriceList(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "update-room-pricelist-api")
	defer util.TEL.Pop()

	jwt, err := util.GetJwt(ctx)
	if err != nil {
		util.TEL.Event("could not get JWT", err)
		AbortError(ctx, ErrUnauthenticated)
		return
	}

	if jwt.Role != util.Host {
		util.TEL.Event("user is not host", nil)
		AbortError(ctx, ErrUnauthorized)
		return
	}

	var dto CreateRoomPriceListDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		util.TEL.Event("failed to bind JSON", err)
		AbortError(ctx, err)
		return
	}

	list, err := h.service.UpdatePriceList(util.TEL.Ctx(), jwt.ID, dto)
	if err != nil {
		util.TEL.Event("could not update room pricelist", err)
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, NewRoomPriceListDTO(list))
}

func (h *Handler) findAvailableRooms(ctx *gin.Context) {
	util.TEL.Push(ctx.Request.Context(), "find-available-rooms-api")
	defer util.TEL.Pop()

	var dto RoomsQueryDTO
	if err := ctx.ShouldBindQuery(&dto); err != nil {
		util.TEL.Event("failed to bind query", err)
		AbortError(ctx, err)
		return
	}

	rooms, resultInfo, err := h.service.FindAvailableRooms(util.TEL.Ctx(), dto)
	if err != nil {
		util.TEL.Event("failed to find available rooms", err)
		AbortError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, NewRoomsResultDTO(rooms, *resultInfo))
}

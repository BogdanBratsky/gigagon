package handler

import (
	"net/http"

	"github.com/BogdanBratsky/proverb/internal/dto"
	"github.com/BogdanBratsky/proverb/internal/model"
	"github.com/BogdanBratsky/proverb/internal/service"
	"github.com/gin-gonic/gin"
)

type roomHandler struct {
	service *service.RoomService
}

func NewRoomHandler(s *service.RoomService) *roomHandler {
	return &roomHandler{service: s}
}

func ok(c *gin.Context, room *model.Room, playerID string, status int) {
	c.JSON(status, dto.Resp{
		Message:  "ok",
		Room:     dto.BuildRoomResponse(room),
		PlayerID: playerID,
	})
}

func fail(c *gin.Context, err error, status int) {
	c.JSON(status, gin.H{
		"message": err.Error(),
	})
}

func (h *roomHandler) GetRoom(c *gin.Context) {
	id := c.Param("id")
	room, err := h.service.GetRoom(id)
	if err != nil {
		fail(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Message: "ok",
		Room:    dto.BuildRoomResponse(room),
	})
}

func (h *roomHandler) CreateRoom(c *gin.Context) {
	var req dto.CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	room, err := h.service.CreateRoom(req.Name)
	if err != nil {
		fail(c, err, http.StatusInternalServerError)
		return
	}

	ok(c, room, room.HostID, http.StatusCreated)
}

func (h *roomHandler) JoinRoom(c *gin.Context) {
	var req dto.JoinRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	player, room, err := h.service.JoinRoom(req.RoomID, req.Name)
	if err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	ok(c, room, player.ID, http.StatusOK)
}

func (h *roomHandler) StartGame(c *gin.Context) {
	var req dto.StartGameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	room, err := h.service.StartGame(req.RoomID, req.HostID)
	if err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	ok(c, room, "", http.StatusOK)
}

func (h *roomHandler) SubmitAnswer(c *gin.Context) {
	var req dto.SubmitAnswerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	err := h.service.SubmitAnswer(req.RoomID, req.PlayerID, req.Text)
	if err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *roomHandler) Vote(c *gin.Context) {
	var req dto.VoteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	err := h.service.Vote(req.RoomID, req.PlayerID, req.AnswerID)
	if err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *roomHandler) NextRound(c *gin.Context) {
	var req dto.NextRoundReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	room, err := h.service.NextRound(req.RoomID, req.HostID)
	if err != nil {
		fail(c, err, http.StatusBadRequest)
		return
	}

	ok(c, room, "", http.StatusOK)
}

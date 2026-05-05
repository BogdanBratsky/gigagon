package app

import (
	"net/http"

	"github.com/BogdanBratsky/proverb/internal/config"
	"github.com/BogdanBratsky/proverb/internal/handler"
	"github.com/BogdanBratsky/proverb/internal/service"
	"github.com/BogdanBratsky/proverb/internal/store"
	"github.com/BogdanBratsky/proverb/internal/ws"
	"github.com/gin-gonic/gin"
)

func NewServer(appCfg *config.AppConfig) (*http.Server, error) {
	// store
	roomStore := store.NewMemoryRoomStore()

	// hub (ВАЖНО — нужен для WS)
	hub := ws.NewHub()
	go hub.Run()

	// services
	proverbService := service.NewProverbService()
	roomService := service.NewRoomService(roomStore, proverbService, hub)

	// handler
	roomHandler := handler.NewRoomHandler(roomService, hub)

	r := gin.Default()

	// =========================
	// WS
	// =========================
	r.GET("/ws/rooms/:id", roomHandler.WS)

	// =========================
	// HTTP
	// =========================
	r.GET("/rooms/:id", roomHandler.GetRoom)

	r.POST("/rooms", roomHandler.CreateRoom)
	r.POST("/rooms/join", roomHandler.JoinRoom)
	r.POST("/rooms/start", roomHandler.StartGame)

	r.POST("/rooms/answer", roomHandler.SubmitAnswer)
	r.POST("/rooms/vote", roomHandler.Vote)
	r.POST("/rooms/next-round", roomHandler.NextRound)

	return &http.Server{
		Addr:    appCfg.Port,
		Handler: r,
	}, nil
}

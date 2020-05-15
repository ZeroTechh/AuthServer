package handler

import (
	"github.com/ZeroTechh/AuthServer/core/code"
	"github.com/ZeroTechh/AuthServer/core/jwt"
	"github.com/ZeroTechh/AuthServer/core/user"
	"github.com/ZeroTechh/VelocityCore/logger"
	"github.com/ZeroTechh/hades"
)

var (
	config = hades.GetConfig(
		"main.yaml",
		[]string{"config", "../config", "../../config"},
	)
	messages = config.Map("messages")
	log      = logger.GetLogger(
		config.Map("service").Str("logFile"),
		config.Map("service").Bool("debug"),
	)
)

// New returns new handler.
func New() *Handler {
	h := Handler{}
	h.init()
	return &h
}

// Handler handles all http function.
type Handler struct {
	code *code.Code
	user *user.User
	jwt  *jwt.JWT
}

// init initializes.
func (h *Handler) init() {
	h.code = code.New(log)
	h.user = user.New(log)
	h.jwt = jwt.New(log)
}

// Disconnect Disconnects.
func (h *Handler) Disconnect() {
	h.code.Disconnect()
	h.user.Disconnect()
	h.jwt.Disconnect()
}

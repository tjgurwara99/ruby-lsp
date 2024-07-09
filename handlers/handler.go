package handlers

import "log"

type Handler struct {
	logger *log.Logger
}

func New(l *log.Logger) *Handler {
	return &Handler{
		logger: l,
	}
}

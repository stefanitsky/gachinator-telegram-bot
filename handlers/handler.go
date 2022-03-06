package handlers

import (
	tele "gopkg.in/telebot.v3"
)

type SuccessHandleLog struct {
	Message     string
	EventAction string
}

type ExtraContext map[string]interface{}

type Handler struct {
	Endpoint         interface{}
	HandlerFunc      tele.HandlerFunc
	SuccessHandleLog SuccessHandleLog
	ExtraContext     ExtraContext
}

// Handle represents a custom handle to set an extra context before actual handle
func (h *Handler) Handle(c tele.Context) error {
	for k, v := range h.ExtraContext {
		c.Set(k, v)
	}

	return h.HandlerFunc(c)
}

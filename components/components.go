package components

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type D map[string]any

func (d D) Render(w io.Writer) error {
	enc, err := json.Marshal(d)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(strings.ReplaceAll(string(enc), "\"", "'")))
	return err
}

func (d D) Key() string {
	return "x-data"
}

type EventHandler struct {
	event     string
	handler   string
	modifiers []string
}

func (e *EventHandler) Render(w io.Writer) error {
	_, err := w.Write([]byte(e.handler))
	return err
}

func (e *EventHandler) Key() string {
	return fmt.Sprintf("@%s", e.event)
}

func On(event string, handler string) *EventHandler {
	return &EventHandler{event: event, handler: handler}
}

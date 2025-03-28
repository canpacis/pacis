package components

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	r "github.com/canpacis/pacis/renderer"
)

type D map[string]any

func (d D) Render(ctx context.Context, w io.Writer) error {
	enc, err := json.Marshal(d)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(strings.ReplaceAll(string(enc), "\"", "'")))
	return err
}

func (D) GetKey() string {
	return "x-data"
}

func (d D) GetValue() any {
	return d
}

func On(event string, handler string) r.Attribute {
	return r.Attr(fmt.Sprintf("x-on:%s", event), handler)
}

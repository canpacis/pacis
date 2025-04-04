package i18n

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"net/http"
	p "path"

	c "github.com/canpacis/pacis/ui/components"
	h "github.com/canpacis/pacis/ui/html"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Setup(dir fs.FS, defaultlang language.Tag) error {
	paths := []string{}
	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if p.Ext(d.Name()) != ".json" {
			return nil
		}

		paths = append(paths, path)
		return nil
	})
	if err != nil {
		return err
	}

	bundle = i18n.NewBundle(defaultlang)
	for _, path := range paths {
		_, err = bundle.LoadMessageFileFS(dir, path)
		if err != nil {
			return err
		}
	}

	return nil
}

type Message struct {
	key  string
	data any
}

func (m Message) Render(ctx context.Context, w io.Writer) error {
	var localizer *i18n.Localizer
	rctx, ok := ctx.(interface{ Request() *http.Request })
	if ok {
		langs := []string{}
		req := rctx.Request()
		locale, err := req.Cookie("app-locale")

		if err == nil {
			langs = append(langs, locale.Value)
		}
		langs = append(langs, req.FormValue("lang"), req.Header.Get("Accept-Language"))
		localizer = i18n.NewLocalizer(bundle, langs...)
	} else {
		localizer = i18n.NewLocalizer(bundle, language.English.String())
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    m.key,
		TemplateData: m.data,
	})
	if err != nil {
		return c.ErrorText(err).Render(ctx, w)
	}
	return h.Text(message).Render(ctx, w)
}

func (Message) NodeType() h.NodeType {
	return h.NodeText
}

func (m Message) String(ctx context.Context) string {
	buf := bytes.NewBuffer([]byte{})
	if err := m.Render(ctx, buf); err != nil {
		buf.Reset()
		c.ErrorText(err).Render(ctx, buf)
	}
	return buf.String()
}

func T(key string, data ...any) Message {
	var d any
	if len(data) > 0 {
		d = data[0]
	}

	return Message{key: key, data: d}
}

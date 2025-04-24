package i18n

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	p "path"

	"github.com/canpacis/pacis/pages/internal"
	c "github.com/canpacis/pacis/ui/components"
	h "github.com/canpacis/pacis/ui/html"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func Setup(dir fs.FS, defaultlang language.Tag) (*i18n.Bundle, error) {
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
		return nil, err
	}

	bundle := i18n.NewBundle(defaultlang)
	for _, path := range paths {
		_, err = bundle.LoadMessageFileFS(dir, path)
		if err != nil {
			return nil, err
		}
	}

	return bundle, nil
}

type Message struct {
	key  string
	data any
}

func (m Message) Render(ctx context.Context, w io.Writer) error {
	localizer := internal.Get[*i18n.Localizer](ctx, "localizer")
	if localizer == nil {
		return c.ErrorText(fmt.Errorf("no localizer in the context, have you registered the i18n middleware correctly?")).Render(ctx, w)
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

func Text(key string, data ...any) Message {
	var d any
	if len(data) > 0 {
		d = data[0]
	}

	return Message{key: key, data: d}
}

func Locale(ctx context.Context) (*language.Tag, error) {
	locale := internal.Get[*language.Tag](ctx, "locale")
	if locale == nil {
		return nil, fmt.Errorf("no localizer in the context, have you registered the i18n middleware correctly?")
	}

	return locale, nil
}

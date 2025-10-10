package i18n

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	p "path"

	"github.com/canpacis/pacis/html"
	"github.com/canpacis/pacis/server/middleware"
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

func (Message) Item() {}

func (m Message) Render(ctx context.Context, w io.Writer) error {
	localizer := middleware.GetLocalizer(ctx)
	if localizer == nil {
		return fmt.Errorf("no localizer in the context, have you registered the i18n middleware correctly?")
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    m.key,
		TemplateData: m.data,
	})
	if err != nil {
		return err
	}
	return html.Text(message).Render(ctx, w)
}

func (m Message) String(ctx context.Context) string {
	localizer := middleware.GetLocalizer(ctx)
	if localizer == nil {
		log.Fatal("no localizer in the context, have you registered the i18n middleware correctly?")
	}

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    m.key,
		TemplateData: m.data,
	})
	if err != nil {
		log.Fatal(err)
	}

	return message
}

func Text(key string, data ...any) Message {
	var d any
	if len(data) > 0 {
		d = data[0]
	}

	return Message{key: key, data: d}
}

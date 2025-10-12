// Package i18n provides internationalization (i18n) support for the application,
// enabling the loading and rendering of localized messages using JSON translation files.
// It integrates with the go-i18n library and provides utilities for loading translation
// bundles, rendering localized messages in HTML templates, and retrieving localized
// strings from the request context.
//
// The package exposes the following key features:
//   - Setup: Initializes an i18n.Bundle by loading all JSON translation files from a given file system.
//   - Message: Represents a localized message, implementing html.Item and html.Node interfaces for seamless
//     integration with HTML rendering.
//   - Text: Helper function to create a Message for a given key and optional template data.
//
// The package expects the i18n middleware to be registered so that a localizer is available in the request context.
package i18n

import (
	"context"
	"io/fs"
	"log"
	p "path"

	"github.com/canpacis/pacis/server/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

/*
Setup initializes an i18n.Bundle by loading all JSON translation files from the provided
file system (dir). It walks through the directory, finds all files with a ".json" extension,
and loads them into the bundle using the specified default language tag. Returns the
initialized i18n.Bundle or an error if any file fails to load.
*/
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

// Message represents a localized message with a key and associated data.
// The key identifies the message, and data holds any additional information
// required for formatting or interpolation.
type Message struct {
	key  string
	data any
}

// Implements the html.Item interface
func (Message) Item() {}

// TODO: Fix
// // Implements the html.Node interface
// func (m Message) Render(ctx context.Context, w io.Writer) error {
// 	localizer := middleware.GetLocalizer(ctx)
// 	if localizer == nil {
// 		return fmt.Errorf("no localizer in the context, have you registered the i18n middleware correctly?")
// 	}

// 	message, err := localizer.Localize(&i18n.LocalizeConfig{
// 		MessageID:    m.key,
// 		TemplateData: m.data,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	return html.Text(message).Render(ctx, w)
// }

// Returns the string value for the message based on the context.
// Used for rendering the key in places like element attributes.
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

// Text returns a Message struct that mimics the html.Text type
// and is meant to be a direct replacment for it for localized
// content.
func Text(key string, data ...any) Message {
	var d any
	if len(data) > 0 {
		d = data[0]
	}

	return Message{key: key, data: d}
}

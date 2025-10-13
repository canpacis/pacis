package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"text/template"
	"unicode"
)

func toPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})

	for i, word := range words {
		words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
	}

	return strings.Join(words, "")
}

var filetempl = `package {{.Package}}

import html "{{.HTMLPackage}}"

type r = html.RawUnsafe

{{ range .Icons }}
func {{.Name}}(items ...html.Item) html.Node {
	return Icon(join(items, r([]byte{ {{.Content}} }))...)
}
{{ end }}
`

type Icon struct {
	Name    string
	Content string
}

type TemplateData struct {
	Package     string
	HTMLPackage string
	Icons       []Icon
}

func main() {
	data := &TemplateData{
		Package:     "lucide",
		HTMLPackage: "github.com/canpacis/pacis/html",
		Icons:       []Icon{},
	}

	dir := "lucide/module/icons"
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	type svg struct {
		Content []byte `xml:",innerxml"`
	}

	for _, entry := range entries {
		name := entry.Name()
		ext := path.Ext(name)
		if ext != ".svg" {
			continue
		}

		var svg svg

		file, err := os.OpenFile(path.Join(dir, name), os.O_RDONLY, 0o644)
		if err != nil {
			log.Fatal(err)
		}
		decoder := xml.NewDecoder(file)
		if err := decoder.Decode(&svg); err != nil {
			log.Fatal(err)
		}
		svg.Content = slices.DeleteFunc(svg.Content, func(b byte) bool {
			return b == 10
		})

		icon := Icon{Name: toPascalCase(strings.TrimSuffix(name, ext))}

		c := []string{}
		for _, b := range svg.Content {
			c = append(c, fmt.Sprintf("%d", b))
		}
		icon.Content = strings.TrimSpace(strings.Join(c, ","))

		data.Icons = append(data.Icons, icon)
	}

	templ, err := template.New("file").Parse(filetempl)
	if err != nil {
		log.Fatal(err)
	}

	os.Remove("lucide/icons.go")
	file, err := os.OpenFile("lucide/icons.go", os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := templ.Execute(file, data); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"encoding/xml"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"slices"
	"strings"
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

func template(name string, content []byte) string {
	c := []string{}

	for _, b := range content {
		c = append(c, fmt.Sprintf("%d", b))
	}

	return fmt.Sprintf(`func %s(props ...h.I) h.Node { return Icon(join(props, r([]byte{%s}))...) }`, toPascalCase(name), strings.Join(c, ", "))
}

type SvgIcon struct {
	Content []byte `xml:",innerxml"`
}

func main() {
	os.Remove("icons.go")
	file, err := os.OpenFile("icons.go", os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	dir := "../lucide/icons"
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	program := "package icons\n\nimport h \"github.com/canpacis/pacis/ui/html\"\n\ntype r = h.RawUnsafe\n\n"

	for _, entry := range entries {
		name := entry.Name()
		ext := path.Ext(name)
		if ext != ".svg" {
			continue
		}

		var icon SvgIcon
		file, err := os.OpenFile(path.Join(dir, name), os.O_RDONLY, 0o644)
		if err != nil {
			log.Fatal(err)
		}
		decoder := xml.NewDecoder(file)
		if err := decoder.Decode(&icon); err != nil {
			log.Fatal(err)
		}
		icon.Content = slices.DeleteFunc(icon.Content, func(b byte) bool {
			return b == 10
		})

		program += template(strings.TrimSuffix(name, ext), icon.Content)
		program += "\n"
	}

	formatted, err := format.Source([]byte(program))
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(formatted)
	if err != nil {
		log.Fatal(err)
	}
}

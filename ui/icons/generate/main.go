package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path"
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

	return fmt.Sprintf(`func %s(props ...h.I) h.Node {
	props = join(props, h.RawUnsafe([]byte{%s}))
	return Icon(props...) 
}`, toPascalCase(name), strings.Join(c, ", "))
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

	program := "package icons\n\nimport h \"github.com/canpacis/pacis/ui/html\"\n\n"

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

		program += template(strings.TrimSuffix(name, ext), icon.Content)
		program += "\n"
	}

	_, err = file.Write([]byte(program))
	if err != nil {
		log.Fatal(err)
	}
}

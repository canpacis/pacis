package main

import (
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

func template(name string) string {
	return fmt.Sprintf("func %s(items ...r.I) r.Node { return Icon(\"%s\", items...) }", toPascalCase(name), name)
}

func main() {
	file, err := os.OpenFile("icons.go", os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	entries, err := os.ReadDir("../lucide/icons")
	if err != nil {
		log.Fatal(err)
	}

	program := "package icons\n\nimport r \"github.com/canpacis/pacis/renderer\"\n\n"

	for _, entry := range entries {
		name := entry.Name()
		ext := path.Ext(name)
		if ext != ".svg" {
			continue
		}

		program += template(strings.TrimSuffix(name, ext))
		program += "\n"
	}

	_, err = file.Write([]byte(program))
	if err != nil {
		log.Fatal(err)
	}
}

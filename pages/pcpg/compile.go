package main

import (
	"bytes"
	"embed"
	"encoding/hex"
	"errors"
	"hash/adler32"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/canpacis/pacis/ui/components"
	"github.com/evanw/esbuild/pkg/api"
)

// func snakecase(s string) string {
// 	var result strings.Builder
// 	for i, r := range s {
// 		if unicode.IsUpper(r) {
// 			if i > 0 && result.String()[result.Len()-1] != '_' {
// 				result.WriteRune('_')
// 			}
// 			result.WriteRune(unicode.ToLower(r))
// 		} else {
// 			result.WriteRune(r)
// 		}
// 	}
// 	return result.String()
// }

// func parseparams(text string) (map[string]string, error) {
// 	params := map[string]string{}

// 	for part := range strings.SplitSeq(strings.TrimSpace(text), " ") {
// 		param := strings.Split(part, "=")
// 		if len(param) != 2 {
// 			return nil, fmt.Errorf("malformed param list")
// 		}
// 		params[param[0]] = param[1]
// 	}

// 	return params, nil
// }

// type page struct {
// 	path       string
// 	layout     *layout
// 	middleware []middleware
// }

// type layout struct {
// 	name string
// }

// type middleware struct {
// 	name string
// }

// var layouts = []layout{}

// var pages = []page{}

// var middlewares = []page{}

// func main() {
// wd, err := os.Getwd()
// if err != nil {
// 	log.Fatal(err)
// }

// if len(os.Args) < 2 {
// 	log.Fatal("dir required")
// }

// root := path.Join(wd, os.Args[1])
// app := path.Join(root, "app")

// fset := token.NewFileSet()

// pkgs, err := parser.ParseDir(fset, app, nil, parser.ParseComments)
// if err != nil {
// 	log.Fatal(err)
// }
// apppkg := pkgs["app"]
// files := apppkg.Files

// type def struct {
// 	text string
// 	decl *ast.FuncDecl
// }
// pagedefs := []def{}
// layoutdefs := []def{}
// middlewaredefs := []def{}

// for _, file := range files {
// 	for _, decl := range file.Decls {
// 		fn, ok := decl.(*ast.FuncDecl)
// 		if !ok {
// 			continue
// 		}

// 		if len(fn.Doc.List) == 0 {
// 			continue
// 		}
// 		text := fn.Doc.List[0].Text
// 		if !strings.HasPrefix(text, "//pacis:") {
// 			continue
// 		}

// 		switch {
// 		case strings.HasPrefix(text, "//pacis:page"):
// 			pagedefs = append(pagedefs, def{text, fn})
// 		case strings.HasPrefix(text, "//pacis:layout"):
// 			layoutdefs = append(layoutdefs, def{text, fn})
// 		case strings.HasPrefix(text, "//pacis:middleware"):
// 			middlewaredefs = append(middlewaredefs, def{text, fn})
// 		default:
// 			continue
// 		}
// 	}

// 	for _, def := range layoutdefs {
// 		layouts = append(layouts, layout{snakecase(def.decl.Name.String())})
// 	}

// 	for _, def := range pagedefs {
// 		paramstext, _ := strings.CutPrefix(def.text, "//pacis:page")
// 		params, err := parseparams(paramstext)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		layoutname, ok := params["layout"]
// 		var l *layout
// 		if ok {
// 			idx := slices.IndexFunc(layouts, func(l layout) bool {
// 				return l.name == layoutname
// 			})
// 			if idx >= 0 {
// 				l = &layouts[idx]
// 			}
// 		}
// 		pages = append(pages, page{params["path"], l, nil})
// 	}
// 	fmt.Println(pages)
// }
// }

func hash(src []byte, prefix, suffix string) string {
	hasher := adler32.New()
	hasher.Write(src)
	return prefix + hex.EncodeToString(hasher.Sum(nil)) + suffix
}

type dirconfig struct {
	root     string
	app      string
	messages string
	assets   string
	static   string
}

func setupdir(target string) (*dirconfig, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	root := path.Join(wd, target)
	app := path.Join(root, "app")
	cfg := &dirconfig{
		root:     root,
		app:      app,
		messages: path.Join(app, "messages"),
		assets:   path.Join(app, "assets"),
		static:   path.Join(app, "static"),
	}
	return cfg, nil
}

func compile(target string) error {
	dircfg, err := setupdir(target)
	if err != nil {
		return err
	}

	os.RemoveAll(dircfg.static)
	os.Mkdir(dircfg.static, 0o755)

	assets, err := createassets(dircfg)
	if err != nil {
		return err
	}
	for _, asset := range assets {
		file, err := os.OpenFile(asset.name, os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			return err
		}
		defer file.Close()
		if _, err := file.Write(asset.content); err != nil {
			return err
		}
	}
	return nil
}

type asset struct {
	base    string
	name    string
	content []byte
}

//go:embed static
var static embed.FS

func createassets(dircfg *dirconfig) ([]asset, error) {
	assets := []asset{}

	var name string

	script := components.AppScript()
	name = hash(script, "main_", ".js")
	assets = append(assets, asset{name, path.Join(dircfg.static, name), script})

	style := components.AppStyle()

	entries, err := os.ReadDir(dircfg.assets)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		ext := path.Ext(name)
		base, _ := strings.CutSuffix(name, ext)

		switch ext {
		case ".ts":
			if name != "main.ts" {
				continue
			}
			result := api.Build(api.BuildOptions{
				EntryPoints: []string{path.Join(dircfg.assets, name)},
				Bundle:      true,
				Write:       false,
			})

			if len(result.Errors) != 0 {
				return nil, errors.New(result.Errors[0].Text)
			}
			raw := result.OutputFiles[0].Contents
			name = hash(raw, "main_", ".js")
			assets = append(assets, asset{name, path.Join(dircfg.static, name), raw})
		case ".css":
			if name != "main.css" {
				continue
			}
			raw, err := os.ReadFile(path.Join(dircfg.assets, name))
			if err != nil {
				return nil, err
			}
			generated, err := stdiotailwind(raw)
			if err != nil {
				return nil, err
			}
			style = append(style, 10)
			style = append(style, generated...)
		default:
			raw, err := os.ReadFile(path.Join(dircfg.assets, name))
			if err != nil {
				return nil, err
			}
			name = hash(raw, base+"_", ext)
			assets = append(assets, asset{base, path.Join(dircfg.static, name), raw})
		}
	}

	// TODO: find a solution for local static assets
	entries, err = static.ReadDir("static")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		ext := path.Ext(name)
		base, _ := strings.CutSuffix(name, ext)

		raw, err := static.ReadFile(path.Join("static", name))
		if err != nil {
			return nil, err
		}
		name = hash(raw, base+"_", ext)
		assets = append(assets, asset{base, path.Join(dircfg.static, name), raw})
	}

	name = hash(style, "main_", ".css")
	assets = append(assets, asset{name, path.Join(dircfg.static, name), style})

	return assets, nil
}

func stdiotailwind(src []byte) ([]byte, error) {
	inbuf := bytes.NewBuffer(src)

	infile, err := os.CreateTemp("", "input-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(infile.Name())
	defer infile.Close()

	if _, err := io.Copy(infile, inbuf); err != nil {
		return nil, err
	}
	infile.Sync()

	outfile, err := os.CreateTemp("", "output-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(outfile.Name())
	defer outfile.Close()

	cmd := exec.Command("pcpg_tw", "-i", infile.Name(), "-o", outfile.Name(), "-m")
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	if _, err := outfile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	outbuf := new(bytes.Buffer)
	if _, err := io.Copy(outbuf, outfile); err != nil {
		return nil, err
	}

	return outbuf.Bytes(), nil
}

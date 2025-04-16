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

func compile(target string) (map[string]string, error) {
	dircfg, err := setupdir(target)
	if err != nil {
		return nil, err
	}

	os.RemoveAll(dircfg.static)
	os.Mkdir(dircfg.static, 0o755)

	assets, assetmap, err := createassets(dircfg)
	if err != nil {
		return nil, err
	}
	for _, asset := range assets {
		file, err := os.OpenFile(asset.name, os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		if _, err := file.Write(asset.content); err != nil {
			return nil, err
		}
	}
	return assetmap, nil
}

type asset struct {
	base    string
	name    string
	content []byte
}

//go:embed static
var staticfs embed.FS

func createassets(dircfg *dirconfig) ([]asset, map[string]string, error) {
	assets := []asset{}
	assetmap := map[string]string{}

	static := func(str string) string {
		return "/static/" + str
	}

	var name string

	script := components.AppScript()
	name = hash(script, "app_", ".js")
	assets = append(assets, asset{name, path.Join(dircfg.static, name), script})
	assetmap["app.ts"] = static(name)

	style := components.AppStyle()

	entries, err := os.ReadDir(dircfg.assets)
	if err != nil {
		return nil, nil, err
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
				return nil, nil, errors.New(result.Errors[0].Text)
			}
			raw := result.OutputFiles[0].Contents
			old := name
			name = hash(raw, "main_", ".js")
			assets = append(assets, asset{name, path.Join(dircfg.static, name), raw})
			assetmap[old] = static(name)
		case ".css":
			if name != "main.css" {
				continue
			}
			raw, err := os.ReadFile(path.Join(dircfg.assets, name))
			if err != nil {
				return nil, nil, err
			}
			generated, err := stdiotailwind(raw)
			if err != nil {
				return nil, nil, err
			}
			style = append(style, 10)
			style = append(style, generated...)
		default:
			raw, err := os.ReadFile(path.Join(dircfg.assets, name))
			if err != nil {
				return nil, nil, err
			}
			old := name
			name = hash(raw, base+"_", ext)
			assets = append(assets, asset{base, path.Join(dircfg.static, name), raw})
			assetmap[old] = static(name)
		}
	}

	// TODO: find a solution for local static assets
	entries, err = staticfs.ReadDir("static")
	if err != nil {
		return nil, nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		ext := path.Ext(name)
		base, _ := strings.CutSuffix(name, ext)

		raw, err := staticfs.ReadFile(path.Join("static", name))
		if err != nil {
			return nil, nil, err
		}
		old := name
		name = hash(raw, base+"_", ext)
		assets = append(assets, asset{base, path.Join(dircfg.static, name), raw})
		assetmap[old] = static(name)
	}

	name = hash(style, "main_", ".css")
	assets = append(assets, asset{name, path.Join(dircfg.static, name), style})
	assetmap["main.css"] = static(name)

	return assets, assetmap, nil
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

	dir := getInstallPath()
	cmd := exec.Command(path.Join(dir, "pcpg_tw"), "-i", infile.Name(), "-o", outfile.Name(), "-m")
	cmd.Stderr = os.Stderr
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

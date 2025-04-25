package bundler

import (
	"encoding/hex"
	"hash/adler32"
	"os"
	"path"
	"strings"

	"github.com/gobwas/glob"
)

type AssetType int

const (
	InvalidAsset = AssetType(iota)
	TypeScript
	// Unused
	JavaScript
	// Unused
	RawCSS
	TailwindCSS
	SVG
	Folder
	Other
)

type Asset struct {
	Type       AssetType
	Path       string
	Name       string
	Extension  string
	OutputName string
	Content    []byte
	// For folders
	Children []*Asset
}

var ignore = []glob.Glob{
	glob.MustCompile(".DS_Store"),
	glob.MustCompile(".git"),
	glob.MustCompile(".dockerfile"),
	glob.MustCompile(".env"),
	glob.MustCompile("*.exe"),
	glob.MustCompile("*.exe~"),
	glob.MustCompile("*.dll"),
	glob.MustCompile("*.so"),
	glob.MustCompile("*.dylib"),
	glob.MustCompile("node_modules"),
}

func hash(src []byte, prefix, suffix string) string {
	hasher := adler32.New()
	hasher.Write(src)
	return prefix + hex.EncodeToString(hasher.Sum(nil)) + suffix
}

func ExtractAssets(target string) ([]*Asset, error) {
	assets := []*Asset{}
	entries, err := os.ReadDir(target)
	if err != nil {
		return nil, err
	}

loop:
	for _, entry := range entries {
		name := entry.Name()
		for _, pattern := range ignore {
			if pattern.Match(name) {
				continue loop
			}
		}

		pth := path.Join(target, name)
		ext := path.Ext(name)
		basename := strings.TrimSuffix(path.Base(name), ext)

		if entry.IsDir() {
			subassets, err := ExtractAssets(pth)
			if err != nil {
				return nil, err
			}
			if len(subassets) != 0 {
				assets = append(assets, &Asset{
					Type:       Folder,
					Path:       pth,
					Name:       basename,
					OutputName: basename,
					Children:   subassets,
				})
			}
		} else {
			switch ext {
			case ".ts":
			case ".js":
			case ".css":
			case ".svg":
			default:
				content, err := os.ReadFile(pth)
				if err != nil {
					return nil, err
				}
				sum := hash(content, basename+"_", ext)

				assets = append(assets, &Asset{
					Type:       Other,
					Path:       pth,
					Name:       basename,
					Extension:  ext,
					OutputName: sum,
					Content:    content,
				})
			}
		}
	}

	return assets, nil
}

func MoveAssets(target string, assets []*Asset) error {
	os.RemoveAll(target)
	os.Mkdir(target, 0o700)

	for _, asset := range assets {
		switch asset.Type {
		case Folder:
			pth := path.Join(target, asset.Name)
			if err := os.Mkdir(pth, 0o700); err != nil {
				return err
			}
			if err := MoveAssets(pth, asset.Children); err != nil {
				return err
			}
		default:
			pth := path.Join(target, asset.OutputName)
			if err := os.WriteFile(pth, asset.Content, 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

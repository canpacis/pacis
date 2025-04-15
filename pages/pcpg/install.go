package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func install() error {
	ros := runtime.GOOS
	arch := runtime.GOARCH

	switch ros {
	case "darwin":
		ros = "macos"
	}

	var dir string

	setup := func() error {
		if ros == "windows" {
			dir = filepath.Join(os.Getenv("PROGRAMFILES"), "pcpg")
		} else {
			dir = "/usr/local/bin"
		}
		// Create install directory if it doesn't exist
		os.MkdirAll(dir, 0755)

		// Remove old artifacts
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to remove artifacts: %w", err)
		}
		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, "pcpg_") {
				os.Remove(path.Join(dir, name))
			}
		}
		return nil
	}

	if err := setup(); err != nil {
		return err
	}

	savebin := func(url, name string) error {
		fmt.Printf("downloading binary %s from %s\n", name, url)

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to download binary: %w", err)
		}
		defer resp.Body.Close()
		fmt.Println("download complete")

		binpath := filepath.Join(dir, name)
		if ros == "windows" {
			binpath += ".exe"
		}

		out, err := os.Create(binpath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("failed to save binary: %w", err)
		}

		if ros != "windows" {
			err = os.Chmod(binpath, 0755)
			if err != nil {
				fmt.Println("Failed to make binary executable:", err)
				os.Exit(1)
			}
		}

		return nil
	}

	url := fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-%s-%s", ros, arch)
	if err := savebin(url, "pcpg_tw"); err != nil {
		return err
	}

	// potentially more binaries will come here
	fmt.Println("install complete")
	return nil
}

package main

import (
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
	"log"
	"os"
	"path/filepath"
)

func main() {
	m := minify.New()
	m.AddFunc("text/javascript", js.Minify)

	files := []string{
		"static/js/login.js",
		"static/js/monthly.js",
	}

	for _, file := range files {

		f, err := os.Open(file)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer f.Close()

		base := filepath.Base(file)
		name := base[:len(base)-len(filepath.Ext(base))]
		outFileName := filepath.Join(filepath.Dir(file), name+".min.js")

		outFile, err := os.Create(outFileName)
		if err != nil {
			log.Fatalf("failed to create minified file: %v", err)
		}
		defer outFile.Close()

		if err := m.Minify("text/javascript", outFile, f); err != nil {
			log.Fatalf("failed to minify file: %v", err)
		}
	}

	log.Println("JS Minification complete")
}

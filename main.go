package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	// cbf to write markdown parser. maybe another day
	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func main() {
	mux := http.NewServeMux()

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	mux.Handle("/static/", http.FileServer(http.Dir("./static/")))

	entries, err := os.ReadDir("./pages/")
	if err != nil {
		log.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		basename := strings.TrimSuffix(name, filepath.Ext(name))
		if strings.Contains(basename, ".") {
			continue
		}

		filepath := fmt.Sprintf("./pages/%s", name)

		log.Printf("adding path: /pages/%s\n", basename)
		mux.HandleFunc(fmt.Sprintf("/pages/%s", basename), func(w http.ResponseWriter, r *http.Request) {
			log.Printf(fmt.Sprintf("accessing page: %v\n", basename))

			mdbytes, err := os.ReadFile(filepath)
			if err != nil {
				log.Fatalf("failed to read file when adding handlers: %v\n", err)
			}

			htmlbytes := mdToHTML(mdbytes)

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, err = w.Write(htmlbytes)
		})
	}

	log.Printf("running server on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func mdToHTML(md []byte) []byte {
	// pasted thank you github user who wrote this
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank
	opts := mdhtml.RendererOptions{Flags: htmlFlags}
	renderer := mdhtml.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

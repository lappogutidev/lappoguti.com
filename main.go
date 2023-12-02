package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
  "path/filepath"
  "strings"
  "os"
  "fmt"
  "html/template"
)

type Page struct {
    Body  template.HTML
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		id = "index"
	}

	md, err := os.ReadFile(filepath.Join("pages", (id + ".md")))
	if err != nil {
		fmt.Fprint(w, "custom 404")
		return
	}
	html := mdToHTML(md)
	p := &Page{Body: template.HTML(html)}
	t, _ := template.ParseFiles("page.html")
	t.Execute(w, p)
}

func main() {
	port := flag.String("p", "8100", "port to serve on")
	assets := flag.String("d", ".", "assets directory")
	flag.Parse()

	http.Handle("/assets/", http.FileServer(http.Dir(*assets)))
	http.HandleFunc("/", pageHandler)

	log.Printf("Serving %s on HTTP port: %s\n", *assets, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
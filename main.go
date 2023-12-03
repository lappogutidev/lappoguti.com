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

func pageHandlerFactory(pages string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/")
		if id == "" {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		path := filepath.Join(pages, (id + ".md"))

		md, err := os.ReadFile(path)
		if err != nil {
			log.Println("Missing page: " + id)
			fmt.Fprint(w, "custom 404")
			return
		}
		html := mdToHTML(md)
		p := &Page{Body: template.HTML(html)}
		t, _ := template.ParseFiles("page.html")
		t.Execute(w, p)
	}
}

func main() {
	port := flag.String("p", "8100", "port to serve on")
	assets := flag.String("assets", "../assets/", "assets directory")
	pages := flag.String("pages", "../pages", "pages directory")
	flag.Parse()

	http.HandleFunc("/", pageHandlerFactory(*pages))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(*assets))))

	log.Printf("Serving %s on HTTP port: %s\n", *assets, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
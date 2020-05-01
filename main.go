package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const Index = `
<!DOCTYPE html>
<html>
	<head>{{ .Head }}</head>
	<body>{{ .Body }}</body>
</html>
`

type Page struct {
	*template.Template
	Head template.HTML
	Body template.HTML

	buf bytes.Buffer
}

func NewPageWithHeader(tmpl *template.Template, headData interface{}) (*Page, error) {
	var buf bytes.Buffer
	defer buf.Reset()

	if err := tmpl.ExecuteTemplate(&buf, "head", headData); err != nil {
		return nil, err
	}

	return &Page{
		Template: tmpl,
		Head:     template.HTML(buf.String()),
		buf:      buf,
	}, nil
}

func (p *Page) RenderToFile(path, name string, data interface{}) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to open/create file")
	}

	defer p.buf.Reset()

	if err := p.ExecuteTemplate(&p.buf, name, data); err != nil {
		return errors.Wrap(err, "Failed to render")
	}

	p.Body = template.HTML(p.buf.String())

	return errors.Wrap(p.Execute(f, p), "Failed to render")
}

func main() {
	var tmpl = parseTemplates()

	var output string
	flag.StringVar(&output, "o", "./public", "The folder to generate HTML files to.")
	var postDir string
	flag.StringVar(&postDir, "d", "./posts", "The folder to read posts from.")

	flag.Parse()

	// Get the absolute path
	o, err := filepath.Abs(output)
	if err != nil {
		log.Fatalln("Failed to get absolute path of output:", err)
	}
	output = o

	// Parse blog posts first.
	articles, err := ParseArticles(postDir)
	if err != nil {
		log.Fatalln("Failed to parse articles:", err)
	}

	blog := &Blog{
		Articles: articles,
	}

	if err := os.MkdirAll(output, 0755|os.ModeDir); err != nil {
		log.Fatalln("Failed to mkdir output:", err)
	}

	// Render the head.
	p, err := NewPageWithHeader(tmpl, blog)
	if err != nil {
		log.Fatalln("Failed to render header:", err)
	}

	// Render the homepage.
	var hppath = join(output, "index.html")
	log.Println("Rendering the homepage to", hppath)

	if err := p.RenderToFile(hppath, "index", blog); err != nil {
		log.Fatalln("Failed to render homepage:", err)
	}

	// Render all posts
	for _, post := range blog.Articles {
		var path = join(output, post.Path)
		log.Println("Rendering article titled", post.Title, "to", path)

		if err := p.RenderToFile(path, "article", post); err != nil {
			log.Fatalln("Failed to render", post.Slug+":", err)
		}
	}
}

func join(paths ...string) string {
	return filepath.Join(paths...)
}

func parseTemplates() (t *template.Template) {
	t = template.New("page")
	t = t.Funcs(Funcs)
	t = template.Must(t.ParseGlob("*.html"))
	t = template.Must(t.Parse(Index))
	return
}

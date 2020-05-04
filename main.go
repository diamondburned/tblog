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

const HTMLPrefix = `<!DOCTYPE html>
<html>
	`
const HTMLSuffix = `
</html>`

type Page struct {
	*template.Template
	buf bytes.Buffer
}

func NewPage(tmpl *template.Template) *Page {
	return &Page{
		Template: tmpl,
		buf:      bytes.Buffer{},
	}
}

func (p *Page) RenderToFile(path, name string, data interface{}) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to open/create file")
	}
	defer f.Close()

	// Write the prefix first
	if _, err := f.Write([]byte(HTMLPrefix)); err != nil {
		return errors.Wrap(err, "Failed to write prefix")
	}

	if err := p.ExecuteTemplate(f, name, data); err != nil {
		return errors.Wrap(err, "Failed to render")
	}

	// Write the suffix past.
	if _, err := f.Write([]byte(HTMLSuffix)); err != nil {
		return errors.Wrap(err, "Failed to write suffix")
	}

	return nil
}

func main() {
	var output string
	flag.StringVar(&output, "o", "./docs", "The folder to generate HTML files to.")
	var postDir string
	flag.StringVar(&postDir, "d", "./posts", "The folder to read posts from.")
	var tmplDir string
	flag.StringVar(&tmplDir, "t", "./templates", "The folder to read templates from.")

	flag.Parse()

	var tmpl = parseTemplates(tmplDir)

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

	// Make a new page constructor.
	p := NewPage(tmpl)

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

func parseTemplates(dir string) (t *template.Template) {
	t = template.New("page")
	t = t.Funcs(Funcs)
	t = template.Must(t.ParseGlob(filepath.Join(dir, "*")))
	return
}

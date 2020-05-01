package main

import (
	"bytes"
	"go/doc"
	"html/template"
	"regexp"
	"strings"
)

const (
	DateFmt = "_2 January 2006"
)

var (
	funcBuf bytes.Buffer

	index = template.Must(template.New("index").Parse(`
		<section id="{{ .slug }}">
			<h3>{{ .date }}</h3>
			<ul>
				{{ range .section }}
				<li><a id="{{ .Slug }}" href="{{ .Path }}">{{ .Title }}</a></li>
				{{ end }}
			</ul>
		</section>
	`))

	mdLink = regexp.MustCompile(`!<a href="(\S+)">\S+</a>`)
)

var Funcs = template.FuncMap{
	// toc generates a table of content
	"index": func(blog *Blog) template.HTML {
		defer funcBuf.Reset()

		var date = ""
		var last = ""
		var section = make([]Article, 0, 5)

		var endSection = func() {
			// Flush the articles into HTML.
			index.Execute(&funcBuf, map[string]interface{}{
				"slug":    strings.ReplaceAll(date, " ", ""),
				"date":    date,
				"section": section,
			})

			section = section[:0]
			date = last
		}

		for i, article := range blog.Articles {
			last = article.ModTime().Format(DateFmt)

			// If this is the first article, then the date variable isn't set.
			if i == 0 {
				date = last
			}

			// Add the article into the section.
			section = append(section, article)

			// If the date of this section is different from the last one.
			if last != date {
				endSection()
			}
		}

		// End the last section, if there's any.
		if len(section) > 0 {
			endSection()
		}

		return template.HTML(funcBuf.String())
	},

	// render converts string to HTML
	"render": func(markup string) template.HTML {

		funcBuf.Reset()
		doc.ToText(&funcBuf, markup, "", "\t", 100)
		markup = funcBuf.String()

		funcBuf.Reset()
		doc.ToHTML(&funcBuf, markup, nil)
		markup = funcBuf.String()

		funcBuf.Reset()
		markup = parseLinks(markup)

		return template.HTML(markup)
	},
}

func parseLinks(inputString string) string {
	return mdLink.ReplaceAllString(inputString, `<img src="$1" />`)
}

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type Blog struct {
	Articles []Article
}

type Article struct {
	os.FileInfo

	Slug string // filename
	Path string

	Title string            // first line always
	Meta  map[string]string // additional metadata
	Body  string
}

func ParseArticles(postDir string) ([]Article, error) {
	dir, err := ioutil.ReadDir(postDir)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read directory")
	}

	var articles = make([]Article, len(dir))

	for i, f := range dir {
		b, err := ioutil.ReadFile(join(postDir, f.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read post")
		}

		// Set the file info for use.
		articles[i].FileInfo = f

		text := string(b)

		// Grab the filename without the extension.
		articles[i].Slug = strings.SplitN(filepath.Base(f.Name()), ".", 2)[0]
		// Form a path based on the slug.
		articles[i].Path = articles[i].Slug + ".html"

		if newLine := strings.Index(text, "\n\n"); newLine > 1 {
			// Set title, parse meta, and get body.
			head := strings.Split(text[:newLine], "\n")
			articles[i].Title = head[0]
			// Metadata are "key: value" lines.
			articles[i].Meta = map[string]string{}
			for _, line := range head[1:] {
				sep := strings.IndexByte(line, ':')
				if sep >= 0 {
					articles[i].Meta[line[:sep]] = strings.TrimSpace(line[sep+1:])
				}
			}
			articles[i].Body = text[newLine:]

			log.Println("Parsed article", articles[i].Title)
			continue
		}

		return nil, errors.Wrap(err, "Failed to parse "+f.Name())
	}

	// Sort articles, latest first.
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].ModTime().After(articles[j].ModTime())
	})

	return articles, nil
}

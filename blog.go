package main

import (
	"bytes"
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

	Title string // first line always
	Body  string // without first line
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

		// Try and get the newline for a blog post.
		if newLine := bytes.IndexByte(b, '\n'); newLine > 1 {
			// Set the title and body.
			articles[i].Title = text[:newLine]
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

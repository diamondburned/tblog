# tblog

A simple blog generator made to be forked.

![Homepage](https://i.imgur.com/H9oBjxy.png)

![Article](https://i.imgur.com/DXknZGR.png)

## Quick Start

### Setting Up

1. Fork and clone the repository.
3. `cd tblog && go run .`
4. Go to docs/

#### Why "docs"?

Because GitHub Pages isn't that good.

### Adding Posts

Add a post by making a new file in `posts/` (or the directory that `-d` is set
to). A post can have any file extension, but the file name will be used as a
URL slug.

#### Markup

tblog uses the GoDoc format with some small extensions. The grammar could be
generalized as following:

```
Title
Date: Whenever
Other article metadata: the whole rest of the line

Paragraph. More sentences.

	Code block goes here
	More code goes here

This URL is automatically highlighted. https://google.com.

!https://images-start-with.an/exclam.jpeg
```

## Deploying

Use your preferred web service to host the `public/` directory of your forked
repository.

## Modding

This blog was made simple so you could easily modify and extend it.

### Requirements

All HTML files (globbed with `templates/*.html`) must declare the
following templates.

- `index` will be the homepage (defined in `templates/index.html`).
- `article` will be every article (defined in `templates/article.html`).

### Files

#### head.html

This file contains the head content. It uses the [Sakura framework](https://github.com/oxalorg/sakura) by default.

#### index.html

This is the main file. It contains the `index` template. As an example, it
declares an extra `head` template that `article.html` uses.

This file demonstrates the built-in `index` function declared in `funcs.go`.
This function generates an index of parsed articles.

#### article.html

This file declares the body for each article.

The very trivial example of this template demonstrates the `render` function in
`funcs.go`, which renders HTML from plain text. Refer to the Markup section for
its format.

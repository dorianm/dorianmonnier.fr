package main

import (
	"context"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"log/slog"
	"os"
)

func main() {
	content, err := os.ReadFile("content/index.md")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	f, err := os.Create("docs/index.html")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = page(string(mdToHTML(content))).Render(context.Background(), f)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

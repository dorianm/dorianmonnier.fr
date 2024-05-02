package main

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"log/slog"
	"os"
	"path/filepath"
	"text/template"
)

func main() {
	if err := os.RemoveAll("docs"); err != nil {
		slog.Error("Error while removing docs directory", err)
		os.Exit(1)
	}
	if err := filepath.Walk("content", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if err := os.Mkdir("docs"+path[len("content"):], 0755); err != nil {
				return err
			}
			return nil
		}
		if filepath.Ext(path) == ".md" {
			outputFile := "docs" + path[len("content"):len(path)-len(".md")] + ".html"
			if err := mdFileToHtmlFile(path, outputFile); err != nil {
				return err
			}
			return nil
		}
		return copyFile(path, "docs"+path[len("content"):])
	}); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func copyFile(src string, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(dst, data, 0644); err != nil {
		return err
	}

	return nil
}

func mdFileToHtmlFile(md string, html string) error {
	templ, err := getTemplate()
	if err != nil {
		return nil
	}
	f, err := os.Create(html)
	if err != nil {
		return err
	}
	content, err := os.ReadFile(md)
	if err != nil {
		return err
	}
	if err := templ.Execute(f, string(mdToHTML(content))); err != nil {
		return err
	}
	return nil
}

func getTemplate() (*template.Template, error) {
	fileContent, err := os.ReadFile("template.html")
	if err != nil {
		return nil, err
	}
	return template.New("template").Parse(string(fileContent))
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

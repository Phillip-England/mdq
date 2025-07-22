package mdq

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type MdFile struct {
	Path     string
	Text     string
	Theme    string
	Html     string
	Context  map[string]any
	MetaHtml string
	Name     string
}

func NewMdFileFromPath(path string, theme string) (MdFile, error) {
	var mdFile MdFile
	mdBytes, err := os.ReadFile(path)
	if err != nil {
		return mdFile, err
	}
	mdFile.Text = string(mdBytes)
	mdFile.Path = path
	mdFile.Name = filepath.Base(path)
	mdFile.Theme = theme
	md := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
			highlighting.NewHighlighting(
				highlighting.WithStyle(theme),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert(mdBytes, &buf, parser.WithContext(context)); err != nil {
		return mdFile, err
	}
	mdFile.Html = buf.String()
	mdFile.Context = meta.Get(context)
	for key, value := range mdFile.Context {
		mdFile.MetaHtml = mdFile.MetaHtml + fmt.Sprintf("<meta name='%s' content='%s'>\n", key, value)
	}
	return mdFile, nil
}

func NewMdFilesFromDir(path string, theme string) ([]MdFile, error) {
	var mds []MdFile
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		md, err := NewMdFileFromPath(path, theme)
		if err != nil {
			return err
		}
		mds = append(mds, md)
		return nil
	})

	if err != nil {
		panic(err)
	}
	return mds, nil
}

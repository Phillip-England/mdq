package mdq

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

type MdFileConfig map[string]string

func newMdFileConfigFromGoQuery(doc *goquery.Document) (MdFileConfig, error) {
	conf := make(MdFileConfig)
	doc.Find("set").Each(func(i int, s *goquery.Selection) {
		key, _ := s.Attr("key")
		value, _ := s.Attr("value")
		key = strings.ToLower(key)
		conf[key] = value
	})

	return conf, nil
}

type MdFile struct {
	Path     string
	Text     string
	Html     string
	Config   MdFileConfig
	Endpoint string
	Depth    int
}

func NewMdFileFromPath(path string, theme string) (MdFile, error) {
	relPath, err := filepath.Rel("./docs", path)
	if err != nil {
		relPath = path // fallback if relative fails
	}
	depth := len(strings.Split(filepath.Dir(relPath), string(os.PathSeparator)))

	md := &MdFile{
		Path:  path,
		Depth: depth,
	}

	f, err := os.ReadFile(md.Path)
	if err != nil {
		return *md, err
	}
	md.Text = string(f)
	gm := goldmark.New(
		goldmark.WithExtensions(
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
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
			goldmarkhtml.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	err = gm.Convert([]byte(md.Text), &buf)
	if err != nil {
		return *md, err
	}
	html := buf.String()
	lines := strings.Split(html, "\n")
	var configLines []string
	var mdLines []string
	foundConfigEnd := false

	for _, line := range lines {
		if strings.Contains(strings.Trim(line, " "), "</config>") {
			foundConfigEnd = true
			configLines = append(configLines, line)
			continue
		}
		if !foundConfigEnd {
			configLines = append(configLines, line)
			continue
		}
		mdLines = append(mdLines, line)
	}

	md.Html = strings.Join(mdLines, "\n")
	configHtml := strings.Join(configLines, "\n")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(configHtml))
	if err != nil {
		return *md, err
	}

	config, err := newMdFileConfigFromGoQuery(doc)
	if err != nil {
		return *md, nil
	}

	md.Config = config

	parts := strings.Split(relPath, string(os.PathSeparator))
	if len(parts) == 1 && parts[0] == "index.md" {
		md.Endpoint = "/"
	} else {
		endpoint := ""
		for _, part := range parts {
			part = strings.Replace(part, ".md", "", 1)
			endpoint += "/" + part
		}
		md.Endpoint = endpoint
	}

	return *md, nil
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

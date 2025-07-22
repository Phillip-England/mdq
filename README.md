# mdq

Quick `.md` to `.html` conversion with sensible defaults.

## Installation

```bash
go get github.com/Phillip-England/mdq
```

## Usage

### Convert a single Markdown file

```go
package main

import (
    "fmt"
    "github.com/Phillip-England/mdq"
)

func main() {
    mdFile, err := mdq.NewMdFileFromPath("example.md", "dracula")
    if err != nil {
        panic(err)
    }
    fmt.Println(mdFile.Html)
}
```

This will:

✅ Read `example.md`  
✅ Convert it to HTML with syntax highlighting (`dracula` theme)  
✅ Prepend meta tags (from the markdown context) to the HTML `<head>` section

---

### Convert all Markdown files in a directory

```go
package main

import (
    "fmt"
    "github.com/Phillip-England/mdq"
)

func main() {
    mdFiles, err := mdq.NewMdFilesFromDir("./markdowns", "dracula")
    if err != nil {
        panic(err)
    }
    for _, mdFile := range mdFiles {
        fmt.Println(mdFile.Path, mdFile.Html)
    }
}
```

---

### Access meta configuration in a Markdown file

`mdFile.Context` is a `map[string]any` that holds key-value pairs from the frontmatter (meta) of the Markdown file.

Example:

```markdown
---
title: Hello World
author: Phillip
---

# Welcome

This is a markdown file.
```

You can access:

```go
fmt.Println(mdFile.Context["title"])  // Hello World
fmt.Println(mdFile.Context["author"]) // Phillip
```

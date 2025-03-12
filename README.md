# mdq
Place config in your markdown files and parse both the config and html at once.

## Credit
All credit for this project is to [goldmark](https://github.com/yuin/goldmark) and [goquery](https://github.com/PuerkitoBio/goquery) for making this possible.

## Hello, World!
```go
package main

import "github.com/Phillip-England/mdq"

func main() {
	mdFile, err := mdq.NewMdFileFromPath("./README.md", "dracula")
	if err != nil {
		panic(err)
	}
}
```

## Config
Place config at the **TOP** of your markdown files like so:
```md
<config>
  <set name='title' value='A good title for config-friendly markdown content' />
</config>
```

Access the config value:
```go
mdFile, err := mdq.NewMdFileFromPath("./README.md", "dracula")
if err != nil {
  panic(err)
}
mdFile.Config["title"] // <== A good title for config-friendly markdown content
```
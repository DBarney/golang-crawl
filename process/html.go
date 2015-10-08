package process

import (
	"errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type (
	Add func(*Page, string)
)

var (
	notFound = errors.New("Not Found")
	srcAsset = map[string]Add{
		"src": addAsset,
	}
	hrefAsset = map[string]Add{
		"href": addAsset,
	}
	hrefLink = map[string]Add{
		"href": addLink,
	}
	mappings = map[atom.Atom]map[string]Add{
		atom.Img:    srcAsset,
		atom.Script: srcAsset,
		atom.Link:   hrefAsset,
		atom.A:      hrefLink,
	}
)

func addAsset(page *Page, asset string) {
	page.Assets = append(page.Assets, asset)
}

func addLink(page *Page, asset string) {
	page.Links = append(page.Links, asset)
}

func ParseXML(job interface{}) (interface{}, error) {
	page := job.(*Page)
	root, err := html.Parse(page.Res.Body)
	if err != nil {
		return nil, err
	}

	page.Root = root
	page.Links = make([]string, 0)
	page.Assets = make([]string, 0)

	return page, nil
}

func CompileNodeInfo(job interface{}) (interface{}, error) {
	page := job.(*Page)

	serachDoc(page.Root, page)
	page.Root = nil
	return page, nil
}

func serachDoc(node *html.Node, page *Page) {
	neededAttributes, found := mappings[node.DataAtom]
	if found {
		for key, add := range neededAttributes {
			if attr, err := findAttr(key, node.Attr); err == nil {
				add(page, attr.Val)
			}
		}
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		serachDoc(child, page)
	}
}

func findAttr(key string, attributes []html.Attribute) (html.Attribute, error) {
	for _, attribute := range attributes {
		if attribute.Key == key {
			return attribute, nil
		}
	}
	return html.Attribute{}, notFound
}

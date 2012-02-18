package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"libxml"
	"os"
)

var filename *string = flag.String("file", "", "File to parse")

func main() {
	flag.Parse()

	if *filename != "" {
		ParseFile(*filename)
	} else {
		ParseHTML("<html><body some_attr='b'><div id='boo'>hey<span class='boo'>some span text</span></div></body></html>")
		ParseHTML("<html><body some_attr='b'><div id='boo'>hey<span class='boo'>some span text</span></div><div>bah</div></body></html>")
		ParseHTML("")
	}
}

func ParseFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	ParseHTML(string(buf))
}

func ParseHTML(src string) bool {
	doc, err := libxml.ParseHTML(src)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer doc.Close()
	TraverseNode(doc.Root())

	result := doc.XPath("string(//div/*)")
	nodes := result.Nodes()

	for _, n := range nodes {
		fmt.Println("node text: ", n.Text())
	}

	fmt.Println(result.String())

	return true
}

func TraverseNode(node *libxml.XmlNode) {
	var curNode *libxml.XmlNode

	for curNode = node; curNode != nil; curNode = curNode.Next() {
		//Do something here...
		fmt.Println("NODE > ", curNode.Name(), " TYPE > ",
			curNode.Type(), " TEXT > ", curNode.Text(), " SOMEATTR > ", curNode.Attr("some_attr"), "ATTRS > ", curNode.Attrs())

		TraverseNode(curNode.Children())
	}
}

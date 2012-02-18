package main

import (
    "fmt"
    "libxml"
)

func main() {
    ParseHTML("<html><body some_attr='b'><div id='boo'>hey<span class='boo'>some span text</span></div></body></html>")
    ParseHTML("")
}

func ParseHTML(src string) bool {
    doc, err := libxml.ParseHTML(src)
    if err != nil {
        fmt.Println(err)
        return false
    }
    defer doc.Close()
    TraverseNode(doc.Root())

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

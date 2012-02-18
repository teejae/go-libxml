package main

import (
    "fmt"
    "libxml"
)

func main() {
    ParseHTML("<html><body><div id='boo'>hey</div></body></html>")
    ParseHTML("")
}

func ParseHTML(src string) bool {
    doc, err := libxml.ParseHTML(src)
    if err != nil {
        fmt.Println(err)
    }
    TraverseNode(doc.Root())

    doc.Close()

    return true
}

func TraverseNode(node *libxml.XmlNode) {
    var curNode *libxml.XmlNode

    for curNode = node; curNode != nil; curNode = curNode.Next() {
        //Do something here...
        fmt.Println("NODE > ", curNode.Name(), " TYPE > ",
curNode.Type())

        TraverseNode(curNode.Children())
    }
}

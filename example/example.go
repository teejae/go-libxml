package main

import (
    "fmt"
    "libxml"
)

func main() {
    ParseHTML("<html><body><div id='boo'>hey</div></body></html>")
}

func ParseHTML(src string) bool {
    d := libxml.HtmlReadDoc(src, "", "",
        libxml.HTML_PARSE_COMPACT | libxml.HTML_PARSE_NOBLANKS |
        libxml.HTML_PARSE_NOERROR | libxml.HTML_PARSE_NOWARNING)

    defer libxml.XmlFreeDoc(d) //free doc on exit

    //get root node
    root := libxml.XmlDocGetRootElement(d);
    if root == nil { return false } //no nodes

    //traverse tree
    var n libxml.XmlNode; n.Ptr = root
    NextNode(&n)

    return true
}

func NextNode(node *libxml.XmlNode) {
    var curNode libxml.XmlNode
    var childNode libxml.XmlNode

    for curNode.Ptr = node.Ptr; curNode.Ptr != nil; curNode.Ptr =
libxml.NodeNext(curNode.Ptr) {
        //Do something here...
        fmt.Println("NODE > ", curNode.Name(), " TYPE > ",
curNode.Type())

        childNode.Ptr = libxml.NodeChildren(curNode.Ptr)
        NextNode(&childNode)
    }
}

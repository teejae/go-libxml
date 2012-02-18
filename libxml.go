// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

// Original Source before edits:
// https://groups.google.com/group/golang-nuts/browse_thread/thread/d8b220af6fdb7075

// Package libxml provides access libxml2, an XML parsing library, found here:
// http://xmlsoft.org/
package libxml

/*
#cgo CFLAGS: -I/usr/include/libxml2
#cgo LDFLAGS: -L/usr/lib
#cgo LDFLAGS: -lxml2

#include <libxml/xmlversion.h>
#include <libxml/parser.h>
#include <libxml/HTMLparser.h>
#include <libxml/HTMLtree.h>
#include <libxml/xmlstring.h>

char* xmlChar2C(xmlChar* x) { return (char *) x; }
xmlNode * NodeNext(xmlNode *node) { return node->next; }
xmlNode * NodeChildren(xmlNode *node) { return node->children; }
int NodeType(xmlNode *node) { return (int)node->type; }
*/
import "C"

const (
	//parser option
	HTML_PARSE_RECOVER   = 1 << 0  //relaxed parsing
	HTML_PARSE_NOERROR   = 1 << 5  //suppress error reports
	HTML_PARSE_NOWARNING = 1 << 6  //suppress warning reports
	HTML_PARSE_PEDANTIC  = 1 << 7  //pedantic error reporting
	HTML_PARSE_NOBLANKS  = 1 << 8  //remove blank nodes
	HTML_PARSE_NONET     = 1 << 11 //forbid network access
	HTML_PARSE_COMPACT   = 1 << 16 //compact small text nodes

	//element type
	XML_ELEMENT_NODE       = 1
	XML_ATTRIBUTE_NODE     = 2
	XML_TEXT_NODE          = 3
	XML_CDATA_SECTION_NODE = 4
	XML_ENTITY_REF_NODE    = 5
	XML_ENTITY_NODE        = 6
	XML_PI_NODE            = 7
	XML_COMMENT_NODE       = 8
	XML_DOCUMENT_NODE      = 9
	XML_DOCUMENT_TYPE_NODE = 10
	XML_DOCUMENT_FRAG_NODE = 11
	XML_NOTATION_NODE      = 12
	XML_HTML_DOCUMENT_NODE = 13
	XML_DTD_NODE           = 14
	XML_ELEMENT_DECL       = 15
	XML_ATTRIBUTE_DECL     = 16
	XML_ENTITY_DECL        = 17
	XML_NAMESPACE_DECL     = 18
	XML_XINCLUDE_START     = 19
	XML_XINCLUDE_END       = 20
	XML_DOCB_DOCUMENT_NODE = 21
)

type XmlNode struct {
	Ptr *C.xmlNode
}

type XmlDoc struct {
	Ptr *C.xmlDoc
}

func XmlCheckVersion() int {
	var v C.int
	C.xmlCheckVersion(v)
	return int(v)
}

func XmlCleanUpParser() {
	C.xmlCleanupParser()
}

func XmlFreeDoc(d *C.xmlDoc) {
	C.xmlFreeDoc(d)
}

func HtmlReadFile(url string, encoding string, opts int) *C.xmlDoc {
	return C.htmlReadFile(C.CString(url), C.CString(encoding), C.int(opts))
}

func HtmlReadDoc(content string, url string, encoding string, opts int) *C.xmlDoc {
	c := C.xmlCharStrdup(C.CString(content))
	return C.htmlReadDoc(c, C.CString(url), C.CString(encoding), C.int(opts))
}

func HtmlGetMetaEncoding(d *C.xmlDoc) string {
	s := C.htmlGetMetaEncoding(d)
	return C.GoString(C.xmlChar2C(s))
}

func XmlDocGetRootElement(d *C.xmlDoc) *C.xmlNode {
	return C.xmlDocGetRootElement(d)
}

func XmlGetProp(n *C.xmlNode, name string) string {
	c := C.xmlCharStrdup(C.CString(name))
	s := C.xmlGetProp(n, c)
	return C.GoString(C.xmlChar2C(s))
}

func HtmlTagLookup(name string) *C.htmlElemDesc {
	c := C.xmlCharStrdup(C.CString(name))
	return C.htmlTagLookup(c)
}

func HtmlEntityLookup(name string) *C.htmlEntityDesc {
	c := C.xmlCharStrdup(C.CString(name))
	return C.htmlEntityLookup(c)
}

func HtmlEntityValueLookup(value uint) *C.htmlEntityDesc {
	return C.htmlEntityValueLookup(C.uint(value))
}

//Helpers
func NewDoc() (doc *C.xmlDoc)    { return }
func NewNode() (node *C.xmlNode) { return }
func NodeNext(node *C.xmlNode) *C.xmlNode {
	return C.NodeNext(node)
}
func NodeChildren(node *C.xmlNode) *C.xmlNode {
	return C.NodeChildren(node)
}
func NodeName(node *C.xmlNode) string {
	return C.GoString(C.xmlChar2C(node.name))
}
func NodeType(node *C.xmlNode) int {
	return int(C.NodeType(node))
}

// // Usage Example
// // -----------------------

// func ParseHTML(src string) bool {
//     d := libxml.HtmlReadDoc(src, "", "",
//         libxml.HTML_PARSE_COMPACT | libxml.HTML_PARSE_NOBLANKS |
//         libxml.HTML_PARSE_NOERROR | libxml.HTML_PARSE_NOWARNING)

//     defer libxml.XmlFreeDoc(d) //free doc on exit

//     //get root node
//     root := libxml.XmlDocGetRootElement(d);
//     if root == nil { return false } //no nodes

//     //traverse tree
//     var n libxml.XmlNode; n.Ptr = root
//     NextNode(&n)

//     return true
// }

// func NextNode(node *libxml.XmlNode) {
//     var curNode libxml.XmlNode
//     var childNode libxml.XmlNode

//     for curNode.Ptr = node.Ptr; curNode.Ptr != nil; curNode.Ptr =
// libxml.NodeNext(curNode.Ptr) {
//         //Do something here...
//         fmt.Println("NODE > ", libxml.NodeName(curNode.Ptr), " TYPE > ",
// libxml.NodeType(curNode.Ptr))

//         childNode.Ptr = libxml.NodeChildren(curNode.Ptr)
//         NextNode(&childNode)
//     }
// }

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
*/
import "C"

import (
	"errors"
	"unsafe"
)

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

func xmlCharToString(s *C.xmlChar) string {
	return C.GoString((*C.char)(unsafe.Pointer(s)))
}

type XmlNode struct {
	Ptr *C.xmlNode
}

func (n *XmlNode) Name() string {
	return xmlCharToString(n.Ptr.name)
}

func (n *XmlNode) Type() int {
	return int(n.Ptr._type)
}

func (n *XmlNode) Next() *XmlNode {
	next := (*C.xmlNode)(unsafe.Pointer(n.Ptr.next))
	if next == nil {
		return nil
	}
	return &XmlNode{Ptr: next}
}

func (n *XmlNode) Children() *XmlNode {
	children := (*C.xmlNode)(unsafe.Pointer(n.Ptr.children))
	if children == nil {
		return nil
	}
	return &XmlNode{Ptr: children}
}

func (n *XmlNode) IsText() bool {
	return C.xmlNodeIsText(n.Ptr) != 0
}

func (n *XmlNode) Text() string {
	if !n.IsText() {
		return ""
	}

	return xmlCharToString(C.xmlNodeListGetString(nil, n.Ptr, 0))
}

type XmlDoc struct {
	Ptr  *C.xmlDoc
	root *C.xmlNode
}

func (d *XmlDoc) Root() *XmlNode {
	return &XmlNode{Ptr: d.root}
}

func (d *XmlDoc) Close() error {
	XmlFreeDoc(d.Ptr)
	d.Ptr = nil
	d.root = nil

	return nil
}

const DEFAULT_HTML_PARSE_FLAGS = HTML_PARSE_COMPACT | HTML_PARSE_NOBLANKS |
	HTML_PARSE_NOERROR | HTML_PARSE_NOWARNING

func ParseHTML(src string) (*XmlDoc, error) {
	d := HtmlReadDoc(src, "", "", DEFAULT_HTML_PARSE_FLAGS)

	//get root node
	root := XmlDocGetRootElement(d)
	if root == nil {
		//no nodes
		XmlFreeDoc(d)
		return nil, errors.New("No nodes")
	}

	return &XmlDoc{Ptr: d, root: root}, nil
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
	return xmlCharToString(s)
}

func XmlDocGetRootElement(d *C.xmlDoc) *C.xmlNode {
	return C.xmlDocGetRootElement(d)
}

func XmlGetProp(n *C.xmlNode, name string) string {
	c := C.xmlCharStrdup(C.CString(name))
	s := C.xmlGetProp(n, c)
	return xmlCharToString(s)
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

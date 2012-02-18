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
#include <libxml/xpath.h>

xmlNodePtr nodeSetGetItem(xmlNodeSetPtr ns, int n) {
	return xmlXPathNodeSetItem(ns, n);
}
*/
import "C"

import (
	"errors"
	"fmt"
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

const (
	// enum xmlXPathObjectType:
	XPATH_UNDEFINED = iota
	XPATH_NODESET
	XPATH_BOOLEAN
	XPATH_NUMBER
	XPATH_STRING
	XPATH_POINT
	XPATH_RANGE
	XPATH_LOCATIONSET
	XPATH_USERS
	XPATH_XSLT_TREE
)

const (
	// enum xmlXPathError:
	XPATH_EXPRESSION_OK = iota
	XPATH_NUMBER_ERROR
	XPATH_UNFINISHED_LITERAL_ERROR
	XPATH_START_LITERAL_ERROR
	XPATH_VARIABLE_REF_ERROR
	XPATH_UNDEF_VARIABLE_ERROR
	XPATH_INVALID_PREDICATE_ERROR
	XPATH_EXPR_ERROR
	XPATH_UNCLOSED_ERROR
	XPATH_UNKNOWN_FUNC_ERROR
	XPATH_INVALID_OPERAND
	XPATH_INVALID_TYPE
	XPATH_INVALID_ARITY
	XPATH_INVALID_CTXT_SIZE
	XPATH_INVALID_CTXT_POSITION
	XPATH_MEMORY_ERROR
	XPTR_SYNTAX_ERROR
	XPTR_RESOURCE_ERROR
	XPTR_SUB_RESOURCE_ERROR
	XPATH_UNDEF_PREFIX_ERROR
	XPATH_ENCODING_ERROR
	XPATH_INVALID_CHAR_ERROR
	XPATH_INVALID_CTXT
)

func freeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

func xmlCharToString(s *C.xmlChar) string {
	cstr := unsafe.Pointer(s)
	str := C.GoString((*C.char)(cstr))
	return str
}

func stringToXmlChar(s string) *C.xmlChar {
	cstr := C.CString(s)
	defer freeCString(cstr)
	c := C.xmlCharStrdup(cstr)
	return c
}

type XmlNode struct {
	Ptr *C.xmlNode

	// cached attrs
	attrs map[string]string
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

func (n *XmlNode) Attr(name string) string {
	xname := stringToXmlChar(name)
	return xmlCharToString(C.xmlGetProp(n.Ptr, xname))
}

func (n *XmlNode) Attrs() map[string]string {
	if n.attrs != nil {
		return n.attrs
	}
	propList := n.Ptr.properties

	attrs := make(map[string]string)

	for prop := propList; prop != nil; prop = prop.next {
		name := xmlCharToString(prop.name)
		attrs[name] = n.Attr(name)
	}

	n.attrs = attrs
	return attrs
}

type XmlDoc struct {
	Ptr  *C.xmlDoc
	root *C.xmlNode
}

func (d *XmlDoc) Root() *XmlNode {
	return &XmlNode{Ptr: d.root}
}

func (d *XmlDoc) XPath(xpathExpr string) *XPathResult {
	context := C.xmlXPathNewContext(d.Ptr)
	defer C.xmlXPathFreeContext(context)

	result := C.xmlXPathEvalExpression(stringToXmlChar(xpathExpr), context)

	fmt.Println("XPath: ", xpathExpr, " Type: ", result._type, " Result: ", result)

	return &XPathResult{ptr: result}
}

func (d *XmlDoc) Close() error {
	XmlFreeDoc(d.Ptr)
	d.Ptr = nil
	d.root = nil

	return nil
}

type XPathResult struct {
	ptr C.xmlXPathObjectPtr
}

func (r *XPathResult) Type() uint {
	return uint(r.ptr._type)
}

func (r *XPathResult) Nodes() []*XmlNode {
	if r.Type() != XPATH_NODESET {
		return nil
	}

	nodeSet := r.ptr.nodesetval
	fmt.Println("nodeSet: ", nodeSet)

	nodeSetLen := int(nodeSet.nodeNr)
	fmt.Println("nodeSetLen: ", nodeSetLen)
	nodes := make([]*XmlNode, nodeSetLen)

	for i := 0; i < nodeSetLen; i++ {
		node := C.nodeSetGetItem(nodeSet, C.int(i))
		nodes[i] = &XmlNode{Ptr: node}
	}

	return nodes
}

func (r *XPathResult) String() string {
	if r.Type() != XPATH_STRING {
		panic("libxml: Not a string")
	}
	return xmlCharToString(C.xmlXPathCastToString(r.ptr))
}

func (r *XPathResult) Number() float64 {
	if r.Type() != XPATH_NUMBER {
		panic("libxml: Not a number")
	}
	return float64(C.xmlXPathCastToNumber(r.ptr))
}

func (r *XPathResult) Boolean() bool {
	if r.Type() != XPATH_BOOLEAN {
		panic("libxml: Not a boolean")
	}
	return C.xmlXPathCastToBoolean(r.ptr) != 0
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

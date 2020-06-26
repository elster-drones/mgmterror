// Copyright (c) 2017,2019-2020, AT&T Intellectual Property.
// All rights reserved.
//
// Copyright (c) 2016-2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package mgmterror

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	// RFC6241 Sect 3.1
	netconf_module    = "ietf-netconf"
	netconf_namespace = "urn:ietf:params:xml:ns:netconf:base:1.0"

	// RFC6020 Sect 5.3.1
	yang_module    = "ietf-yang" // There is no yang module
	yang_namespace = "urn:ietf:params:xml:ns:yang:1"

	// TODO: need to find a better place for this
	vyattaModule    = "vyatta-yang"
	VyattaNamespace = "urn:vyatta.com:mgmt:error:1"

	// Used to separate fields in error message strings
	error_msg_separator = ": "
)

// based on RFC6241 Sect 4.3

type errtype uint

const (
	// Secure Transport
	transport errtype = iota

	// Messages
	rpc

	// Operations
	protocol

	// Content
	application
)

var errtypemap = map[string]errtype{
	"transport":   transport,
	"rpc":         rpc,
	"protocol":    protocol,
	"application": application,
}

func (t *errtype) set(typ string) error {
	if v, ok := errtypemap[typ]; ok {
		*t = v
		return nil
	}
	return errors.New("Invalid error target")
}

func (t errtype) String() string {
	for s, v := range errtypemap {
		if t == v {
			return s
		}
	}
	return ""
}

// MgmtErrorInfoTag holds additional information for an error.
//
// Some NETCONF and YANG defined errors have mandatory InfoTag
// information. Implementations may include additional elements to
// provide extended and/or implementation- specific debugging
// information.
type MgmtErrorInfoTag struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

// RFC7951 section 4 dictates that module name instead of namespace is
// used to differentiate elements when JSON encoded.
//
// Since we only have 3 namespaces to worry about right now, we can do
// the mapping locally. In the future, yangd can get a method to
// lookup namespaces/modules.
func (i *MgmtErrorInfoTag) lookupNamespace(module string) string {
	modmap := map[string]string{
		netconf_module: netconf_namespace,
		yang_module:    yang_namespace,
		vyattaModule:   VyattaNamespace,
	}
	ns, ok := modmap[module]
	if !ok {
		return module
	}
	return ns
}

func (i *MgmtErrorInfoTag) lookupModule(ns string) string {
	nsmap := map[string]string{
		netconf_namespace: netconf_module,
		yang_namespace:    yang_module,
		VyattaNamespace:   vyattaModule,
	}
	module, ok := nsmap[ns]
	if !ok {
		return ns
	}
	return module
}

func (i *MgmtErrorInfoTag) UnmarshalJSON(value []byte) error {
	var obj map[string]string
	if err := json.Unmarshal(value, &obj); err != nil {
		return err
	}
	if len(obj) != 1 {
		return errors.New("malformed error-info tag")
	}
	for k, v := range obj {
		s := strings.Split(k, ":")
		if len(s) == 2 {
			i.XMLName.Space = i.lookupNamespace(s[0])
			i.XMLName.Local = s[1]
		} else {
			i.XMLName.Local = s[0]
		}
		i.Value = v
	}
	return nil
}

func (i *MgmtErrorInfoTag) MarshalJSON() ([]byte, error) {
	var tag string
	var out bytes.Buffer
	out.WriteString("{")
	if len(i.XMLName.Space) > 0 {
		tag = i.lookupModule(i.XMLName.Space) + ":" + i.XMLName.Local
	} else {
		tag = i.XMLName.Local
	}
	b, err := json.Marshal(tag)
	if err != nil {
		return []byte(""), err
	}
	out.Write(b)
	out.WriteString(":")
	b, err = json.Marshal(i.Value)
	if err != nil {
		return []byte(""), err
	}
	out.Write(b)
	out.WriteString("}")
	return out.Bytes(), nil
}

func NewMgmtErrorInfoTag(ns, name, value string) *MgmtErrorInfoTag {
	return &MgmtErrorInfoTag{
		XMLName: xml.Name{
			Space: ns,
			Local: name,
		},
		Value: value,
	}
}

type MgmtErrorInfo []MgmtErrorInfoTag

func (e *MgmtErrorInfo) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	error_info_name := xml.Name{
		Local: "error-info",
	}

	if err := enc.EncodeToken(xml.StartElement{Name: error_info_name}); err != nil {
		return err
	}
	for _, v := range *e {
		if err := enc.Encode(v); err != nil {
			return err
		}
	}
	if err := enc.EncodeToken(xml.EndElement{Name: error_info_name}); err != nil {
		return err
	}
	return nil
}

func (e *MgmtErrorInfo) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var infos []MgmtErrorInfoTag
Loop:
	for {
		tok, _ := dec.Token()
		if tok == nil {
			return nil
		}
		switch elem := tok.(type) {
		case xml.StartElement:
			var i MgmtErrorInfoTag
			if err := dec.DecodeElement(&i, &elem); err != nil {
				return err
			}
			infos = append(infos, i)
		case xml.EndElement:
			break Loop
		}
	}
	*e = infos
	return nil
}

func (e MgmtErrorInfo) FindMgmtErrorTag(ns, name string) string {
	for _, t := range e {
		if t.XMLName.Space == ns && t.XMLName.Local == name {
			return t.Value
		}
	}
	return ""
}

// Error structure based on RFC6241 Sect 4.3.
type MgmtError struct {
	XMLName xml.Name `json:"-"`

	// Typ is the type of error
	//
	// RFC6241 Sect 4.3 defines four types of errors:
	//  - transport (layer: Secure Transport)
	//  - rpc (layer: Messages)
	//  - protocol (layer: Operations)
	//  - application (layer: Content)
	Typ string `xml:"error-type" json:"error-type"`

	// Tag identifies the error condition
	//
	// Allowed values for NETCONF and YANG tags are defined by
	// RFC6241 Appendix A and RFC6020 Section 13 respectively.
	Tag string `xml:"error-tag" json:"error-tag"`

	// Severity indicates the error severity
	//
	// RFC6241 Sect 4.3 defines two of error severities:
	//  - error
	//  - warning
	Severity string `xml:"error-severity" json:"error-severity"`

	// AppTag identifies the data-model-specific or
	// implementation-specific error condition, if one exists.
	AppTag string `xml:"error-app-tag,omitempty" json:"error-app-tag,omitempty"`

	// Path contains the absolute XPath expression identifying the
	// element path to the node that is associated with the error
	// being reported if an appropriate payload element or
	// datastore node can be associated with a particular error
	// condition.
	//
	// See RFC6241 Section 4.3 for details on the context for the
	// XPath expression.
	Path string `xml:"error-path,omitempty" json:"error-path,omitempty"`

	// Message contains a string suitable for human display that
	// describes the error condition.
	Message string `xml:"error-message,omitempty" json:"error-message,omitempty"`

	// Info contains protocol- or data-model-specific error
	// content. This element will not be present if no such error
	// content is provided for a particular error condition. An
	// implementation MAY include additional elements to provide
	// extended and/or implementation- specific debugging
	// information.
	Info MgmtErrorInfo `xml:"error-info,omitempty" json:"error-info,omitempty"`
}

func newMgmtError() *MgmtError {
	e := &MgmtError{}
	e.setXMLName()
	return e
}

// MgmtErrorRef - interface that allows us to identify all types of MgmtError
// in a single check.  Use of private function (mgmtErrorRef) ensures no one
// else can create an object that meets the interface unless it explicitly
// includes a MgmtError object, in which case that's ok.
type MgmtErrorRef interface {
	mgmtErrorRef()
}

var _ MgmtErrorRef = (*MgmtError)(nil)

func (me *MgmtError) mgmtErrorRef() {}

// Formattable - interface provided by error types to allow formatting
// (NB: MgmtError is just one example of such a type.)
//
// Unlike plain Go errors, MgmtError types have useful properties which
// we may wish to print out differently depending on the end user / context.
// GetMessage() may well be customised for certain error types to include
// the contents of MgmtError.Info.
//
type Formattable interface {
	GetMessage() string
	GetPath() string
	GetSeverity() string
	GetTag() string
	GetAppTag() string
	GetType() string
	GetInfo() MgmtErrorInfo
}

// Ensure *MgmtError implements interface
var _ Formattable = (*MgmtError)(nil)

func (me *MgmtError) GetMessage() string     { return me.Message }
func (me *MgmtError) GetPath() string        { return me.Path }
func (me *MgmtError) GetSeverity() string    { return me.Severity }
func (me *MgmtError) GetTag() string         { return me.Tag }
func (me *MgmtError) GetAppTag() string      { return me.AppTag }
func (me *MgmtError) GetType() string        { return me.Typ }
func (me *MgmtError) GetInfo() MgmtErrorInfo { return me.Info }

func callCreate(fn interface{}, err *MgmtError) error {
	ty := reflect.TypeOf(fn)
	if ty.Kind() != reflect.Func ||
		ty.NumIn() != 1 ||
		ty.NumOut() != 1 {
		return err
	}
	vfn := reflect.ValueOf(fn)
	result := vfn.Call([]reflect.Value{reflect.ValueOf(err)})
	if len(result) != 1 {
		return err
	}
	e, ok := result[0].Interface().(error)
	if !ok {
		return err
	}
	return e
}

func (e *MgmtError) setXMLName() {
	e.XMLName = xml.Name{
		Space: netconf_namespace,
		Local: "rpc-error",
	}
}

func (e MgmtError) Error() string {
	var b bytes.Buffer

	b.WriteString(strings.Title(e.Severity))
	b.WriteString(error_msg_separator)

	if e.Path != "" {
		b.WriteString(e.Path)
		b.WriteString(error_msg_separator)
	}

	if e.Message != "" {
		b.WriteString(e.Message)
	}

	return b.String()
}

const errpfx = "com.vyatta.rpcerror."

// Encode error for DBus
func (e *MgmtError) DBusError() (string, []interface{}) {
	name := fmt.Sprintf("%s%s", errpfx, e.Typ)
	body := make([]interface{}, 1)
	body[0] = e
	return name, body
}

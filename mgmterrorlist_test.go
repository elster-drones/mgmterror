// Copyright (c) 2017,2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package mgmterror

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/danos/utils/pathutil"
	"github.com/kr/pretty"
	"html"
	"reflect"
	"strings"
	"testing"
)

type testMgmtError struct {
	*MgmtError
}

func (e *testMgmtError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *testMgmtError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func (e *testMgmtError) Error() string {
	return e.MgmtError.Error()
}

func genMgmtErrorPath(n uint) string {
	return fmt.Sprintf("/path/%d", n)
}

func genMgmtErrorMessage(n uint) string {
	return fmt.Sprintf("Message %d", n)
}

func genTestMgmtError(n uint) *testMgmtError {
	err := &testMgmtError{
		MgmtError: newMgmtError(),
	}
	err.Severity = "error"
	err.Typ = application.String()
	err.Tag = operation_failed.String()
	err.Path = genMgmtErrorPath(n)
	err.Message = genMgmtErrorMessage(n)
	return err
}

func genMgmtErrorListXml(n uint) string {
	genXML := func(i uint) string {
		return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + html.EscapeString(operation_failed.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-path>` + html.EscapeString(genMgmtErrorPath(i)) + `</error-path>
	<error-message>` + html.EscapeString(genMgmtErrorMessage(i)) + `</error-message>
</rpc-error>`
	}
	var i uint
	var b bytes.Buffer
	for i = 1; i <= n; i++ {
		if i > 1 {
			b.WriteByte('\n')
		}
		b.WriteString(genXML(i))
	}
	return b.String()
}

func TestMgmtErrorListJSON(t *testing.T) {
	var errs MgmtErrorList
	errs.MgmtErrorListAppend(NewOperationFailedApplicationError(),
		NewMustViolationError(),
		NewExecError([]string{"foo", "bar"}, "boom"),
		fmt.Errorf("This is not a MgmtError error"))

	marshal, err := json.MarshalIndent(errs, "", "\t")
	if err != nil {
		t.Errorf("Marshal MgmtErrorList error: %v\n", err)
		return
	}
	unmarshal := MgmtErrorList{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MgmtErrlrList error: %v\n", err)
		return
	}
	if !reflect.DeepEqual(errs, unmarshal) {
		t.Errorf("Failed JSON marshal/unmarshal")
		t.Logf("Expected: %# v", pretty.Formatter(errs))
		t.Logf("Result:   %# v", pretty.Formatter(unmarshal))
	}
}

func TestMgmtErrorListXML(t *testing.T) {
	var n uint
	gen := func() *testMgmtError {
		n++
		return genTestMgmtError(n)
	}
	var errs MgmtErrorList
	errs.MgmtErrorListAppend(gen(),
		gen(),
		gen(),
		gen(),
		fmt.Errorf("This is not a MgmtError error"))
	expected := genMgmtErrorListXml(n) + "\n"
	expected += `<rpc-error xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
	<error-type>application</error-type>
	<error-tag>operation-failed</error-tag>
	<error-severity>error</error-severity>
	<error-message>This is not a MgmtError error</error-message>
</rpc-error>`
	verifyXmlMarshal(t, errs, expected)
}

func ExampleMgmtErrorList() {
	var n uint
	gen := func() *testMgmtError {
		n++
		return genTestMgmtError(n)
	}
	var elist MgmtErrorList
	elist.MgmtErrorListAppend(gen(),
		gen(),
		gen(),
		gen(),
		fmt.Errorf("This is not a MgmtError error"))

	fmt.Println(elist.Error())

	// Output:
	// Error: /path/1: Message 1
	// Error: /path/2: Message 2
	// Error: /path/3: Message 3
	// Error: /path/4: Message 4
	// Error: This is not a MgmtError error
}

func formatCommitFailErrors(err error) string {
	var b bytes.Buffer

	if me, ok := err.(Formattable); ok {
		pathStr := strings.Join(pathutil.Makepath(me.GetPath()), " ")
		b.WriteString(fmt.Sprintf("[%s]\n\n", pathStr))
		b.WriteString(fmt.Sprintf("%s\n\n", me.GetMessage()))
		b.WriteString(fmt.Sprintf("[[%s]] failed.", pathStr))
	} else {
		b.WriteString(err.Error())
	}
	return b.String()
}

// All errors are mgmterror, though may not have path set.
// Here we check that the ones w/o path are processed correctly by
// custom format function.
// Also need to check override to GetMessage / GetPath works
// TBD
func ExampleMgmtErrorList_customFormat() {
	var n uint
	gen := func() *testMgmtError {
		n++
		return genTestMgmtError(n)
	}
	var elist MgmtErrorList
	elist.MgmtErrorListAppend(gen(),
		gen(),
		fmt.Errorf("This is not a MgmtError error"))

	fmt.Println(elist.CustomError(formatCommitFailErrors))

	// Output:
	// [path 1]
	//
	// Message 1
	//
	// [[path 1]] failed.
	// [path 2]
	//
	// Message 2
	//
	// [[path 2]] failed.
	// []
	//
	// This is not a MgmtError error
	//
	// [[]] failed.
}

// Copyright (c) 2017,2019, AT&T Intellectual Property.  All rights reserved.
//
// Copyright (c) 2016-2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package mgmterror

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"testing"
)

const error_type = "application"

func genNonUniqueXml(basepath string, paths []string) string {
	var error_info bytes.Buffer
	for _, path := range paths {
		error_info.WriteString("\n\t\t<")
		error_info.WriteString(non_unique_info.String())
		error_info.WriteString(" xmlns=\"")
		error_info.WriteString(yang_namespace)
		error_info.WriteString("\">")
		error_info.WriteString(html.EscapeString(path))
		error_info.WriteString("</")
		error_info.WriteString(non_unique_info.String())
		error_info.WriteString(">")
	}
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + html.EscapeString(operation_failed.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-app-tag>` + html.EscapeString(data_not_unique.String()) + `</error-app-tag>
	<error-path>` + html.EscapeString(basepath) + `</error-path>
	<error-message>` + html.EscapeString(msg_yang_operation_failed) + `</error-message>
	<error-info>` + error_info.String() + `
	</error-info>
</rpc-error>`
}

func TestNonUniqueError(t *testing.T) {
	basepath := "/testcontainer/testlist"
	paths := []string{
		basepath + "/foo/bar/baz",
		basepath + "/foo/biz/baz",
	}
	ncerr := NewNonUniqueError(paths)
	ncerr.Path = basepath
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal NonUniqueError error: %v\n", err)
		return
	}
	unmarshal := NonUniqueError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal NonUniqueError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genNonUniqueXml(basepath, paths))
}

func ExampleNonUniqueError() {
	basepath := "/testcontainer/testlist"
	paths := []string{
		basepath + "/name/dev1/attr/value",
		basepath + "/name/dev2/attr/value",
		basepath + "/name/dev3/attr/value",
	}
	err := NewNonUniqueError(paths)
	err.Path = basepath
	fmt.Println(err.Error())

	// Output:
	// Error: /testcontainer/testlist: Non-unique paths name/dev1/attr/value, name/dev2/attr/value, name/dev3/attr/value
}

func genTooManyElementsXml(path string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + html.EscapeString(operation_failed.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-app-tag>` + html.EscapeString(too_many_elements.String()) + `</error-app-tag>
	<error-path>` + html.EscapeString(path) + `</error-path>
	<error-message>` + html.EscapeString(msg_yang_operation_failed) + `</error-message>
</rpc-error>`
}

func TestTooManyElementsError(t *testing.T) {
	const path = "/foo/bar/baz"
	ncerr := NewTooManyElementsError(path)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal TooManyElementsError error: %v\n", err)
		return
	}
	unmarshal := TooManyElementsError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal TooManyElementsError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genTooManyElementsXml(path))
}

func ExampleTooManyElementsError() {
	err := NewTooManyElementsError("biz")
	err.Path = "/foo/bar"
	fmt.Println(err.Error())

	// Output:
	// Error: /foo/bar: The requested operation failed.
}

func genTooFewElementsXml(path string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + html.EscapeString(operation_failed.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-app-tag>` + html.EscapeString(too_few_elements.String()) + `</error-app-tag>
	<error-path>` + html.EscapeString(path) + `</error-path>
	<error-message>` + html.EscapeString(msg_yang_operation_failed) + `</error-message>
</rpc-error>`
}

func TestTooFewElementsError(t *testing.T) {
	const path = "/foo/bar/baz"
	ncerr := NewTooFewElementsError(path)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal TooFewElementsError error: %v\n", err)
		return
	}
	unmarshal := TooFewElementsError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal TooFewElementsError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genTooFewElementsXml(path))
}

func genMustViolationXml(path string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + html.EscapeString(operation_failed.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-app-tag>` + html.EscapeString(must_violation.String()) + `</error-app-tag>
	<error-path>` + html.EscapeString(path) + `</error-path>
	<error-message>` + html.EscapeString(msg_yang_operation_failed) + `</error-message>
</rpc-error>`
}

func TestMustViolationError(t *testing.T) {
	const path = "/foo/bar/baz"
	ncerr := NewMustViolationError()
	ncerr.Path = path
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal MustViolationError error: %v\n", err)
		return
	}
	unmarshal := MustViolationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MustViolationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genMustViolationXml(path))
}

func genInstanceRequiredXml(path string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + html.EscapeString(data_missing.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-app-tag>` + html.EscapeString(instance_required.String()) + `</error-app-tag>
	<error-path>` + html.EscapeString(path) + `</error-path>
	<error-message>` + html.EscapeString(msg_yang_data_missing) + `</error-message>
</rpc-error>`
}

func TestInstanceRequiredError(t *testing.T) {
	const path = "/foo/bar/baz"
	ncerr := NewInstanceRequiredError(path)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal InstanceRequiredError error: %v\n", err)
		return
	}
	unmarshal := InstanceRequiredError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal InstanceRequiredError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genInstanceRequiredXml(path))
}

func genLeafrefMismatchXml(path string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + html.EscapeString(data_missing.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-app-tag>` + html.EscapeString(instance_required.String()) + `</error-app-tag>
	<error-path>` + html.EscapeString(path) + `</error-path>
	<error-message>` + html.EscapeString(msg_yang_data_missing) + `</error-message>
</rpc-error>`
}

func TestLeafrefMismatchError(t *testing.T) {
	const path = "/foo/bar/baz"
	ncerr := NewLeafrefMismatchError(path, "foo")
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal LeafrefMismatchError error: %v\n", err)
		return
	}
	unmarshal := LeafrefMismatchError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal LeafrefMismatchError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genLeafrefMismatchXml(path))
}

func genMissingChoiceXml(path, name string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + (operation_failed.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-app-tag>` + html.EscapeString(missing_choice.String()) + `</error-app-tag>
	<error-path>` + html.EscapeString(path) + `</error-path>
	<error-message>` + html.EscapeString(msg_yang_operation_failed) + `</error-message>
	<error-info>
		<` + missing_choice_info.String() + ` xmlns="` + yang_namespace + `">` + html.EscapeString(name) + `</` + missing_choice_info.String() + `>
	</error-info>
</rpc-error>`
}

func TestMissingChoiceError(t *testing.T) {
	const (
		path = "/foo/bar"
		name = "baz"
	)
	ncerr := NewMissingChoiceError(path, name)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal MissingChoiceError error: %v\n", err)
		return
	}
	unmarshal := MissingChoiceError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MissingChoiceError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genMissingChoiceXml(path, name))
}

func genInsertFailedXml() string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(error_type) + `</error-type>
	<error-tag>` + html.EscapeString(bad_attribute.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(yang_severity_error.String()) + `</error-severity>
	<error-app-tag>` + html.EscapeString(missing_instance.String()) + `</error-app-tag>
	<error-message>` + html.EscapeString(msg_yang_bad_attribute) + `</error-message>
</rpc-error>`
}

func TestInsertFailedError(t *testing.T) {
	ncerr := NewInsertFailedError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal InsertFailedError error: %v\n", err)
		return
	}
	unmarshal := InsertFailedError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal InsertFailedError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genInsertFailedXml())
}

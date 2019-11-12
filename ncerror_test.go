// Copyright (c) 2017,2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2016-2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package mgmterror

import (
	"encoding/json"
	"fmt"
	"html"
	"testing"
)

const (
	bad_attr_value = "bad-attr-value"
	bad_elem_value = "bad-elem-value"
	bad_ns_value   = "urn:bad:namespace"
)

func genInUseXml(typ string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(in_use.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_in_use) + `</error-message>
</rpc-error>`
}

func TestInUseProtocolError(t *testing.T) {
	ncerr := NewInUseProtocolError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal JSON InUseProtocolError error: %v\n", err)
		return
	}
	unmarshal := InUseProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal JSON InUseProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genInUseXml(protocol.String()))
}

func TestInUseApplicationError(t *testing.T) {
	ncerr := NewInUseApplicationError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal InUseApplicationError error: %v\n", err)
		return
	}
	unmarshal := InUseApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal InUseApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genInUseXml(application.String()))
}

func genInvalidValueXml(typ string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(invalid_value.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_invalid_value) + `</error-message>
</rpc-error>`
}

func TestInvalidValueProtocolError(t *testing.T) {
	ncerr := NewInvalidValueProtocolError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal InvalidValueProtocolError error: %v\n", err)
		return
	}
	unmarshal := InvalidValueProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal InvalidValueProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genInvalidValueXml(protocol.String()))
}

func TestInvalidValueApplicationError(t *testing.T) {
	ncerr := NewInvalidValueApplicationError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal InvalidValueApplicationError error: %v\n", err)
		return
	}
	unmarshal := InvalidValueApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal InvalidValueApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genInvalidValueXml(application.String()))
}

func genTooBigXml(typ string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(too_big.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_too_big) + `</error-message>
</rpc-error>`
}

func TestTooBigTransportError(t *testing.T) {
	ncerr := NewTooBigTransportError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal TooBigTransportError error: %v\n", err)
		return
	}
	unmarshal := TooBigTransportError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal TooBigTransportError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genTooBigXml(transport.String()))
}

func TestTooBigRpcError(t *testing.T) {
	ncerr := NewTooBigRpcError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal TooBigRpcError error: %v\n", err)
		return
	}
	unmarshal := TooBigRpcError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal TooBigRpcError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genTooBigXml(rpc.String()))
}

func TestTooBigProtocolError(t *testing.T) {
	ncerr := NewTooBigProtocolError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal TooBigProtocolError error: %v\n", err)
		return
	}
	unmarshal := TooBigProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal TooBigProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genTooBigXml(protocol.String()))
}

func TestTooBigApplicationError(t *testing.T) {
	ncerr := NewTooBigApplicationError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal TooBigApplicationError error: %v\n", err)
		return
	}
	unmarshal := TooBigApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal TooBigApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genTooBigXml(application.String()))
}

func genMissingAttrXml(typ, bad_attr_value, bad_elem_value string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(missing_attribute.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_missing_attribute) + `</error-message>
	<error-info>
		<bad-attribute>` + html.EscapeString(bad_attr_value) + `</bad-attribute>
		<bad-element>` + html.EscapeString(bad_elem_value) + `</bad-element>
	</error-info>
</rpc-error>`
}

func TestMissingAttributeRpcError(t *testing.T) {
	ncerr := NewMissingAttrRpcError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal MissingAttributeRpcError error: %v\n", err)
		return
	}
	unmarshal := MissingAttrRpcError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MissingAttributeRpcError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genMissingAttrXml(rpc.String(), bad_attr_value, bad_elem_value))
}

func TestMissingAttributeProtocolError(t *testing.T) {
	ncerr := NewMissingAttrProtocolError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal MissingAttributeProtocolError error: %v\n", err)
		return
	}
	unmarshal := MissingAttrProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MissingAttributeProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genMissingAttrXml(protocol.String(), bad_attr_value, bad_elem_value))
}

func TestMissingAttributeApplicationError(t *testing.T) {
	ncerr := NewMissingAttrApplicationError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal MissingAttributeApplicationError error: %v\n", err)
		return
	}
	unmarshal := MissingAttrApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MissingAttributeApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genMissingAttrXml(application.String(), bad_attr_value, bad_elem_value))
}

func genBadAttrXml(typ, bad_attr_value, bad_elem_value string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(bad_attribute.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_bad_attribute) + `</error-message>
	<error-info>
		<bad-attribute>` + html.EscapeString(bad_attr_value) + `</bad-attribute>
		<bad-element>` + html.EscapeString(bad_elem_value) + `</bad-element>
	</error-info>
</rpc-error>`
}

func TestBadAttrRpcError(t *testing.T) {
	ncerr := NewBadAttrRpcError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal BadAttrApplicationError error: %v\n", err)
		return
	}
	unmarshal := BadAttrApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal BadAttrApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genBadAttrXml(rpc.String(), bad_attr_value, bad_elem_value))
}

func TestBadAttrProtocolError(t *testing.T) {
	ncerr := NewBadAttrProtocolError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal BadAttrApplicationError error: %v\n", err)
		return
	}
	unmarshal := BadAttrApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal BadAttrApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genBadAttrXml(protocol.String(), bad_attr_value, bad_elem_value))
}

func TestBadAttrApplicationError(t *testing.T) {
	ncerr := NewBadAttrApplicationError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal BadAttrApplicationError error: %v\n", err)
		return
	}
	unmarshal := BadAttrApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal BadAttrApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genBadAttrXml(application.String(), bad_attr_value, bad_elem_value))
}

func genUnknownAttrXml(typ, bad_attr_value, bad_elem_value string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(unknown_attribute.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_unknown_attribute) + `</error-message>
	<error-info>
		<bad-attribute>` + html.EscapeString(bad_attr_value) + `</bad-attribute>
		<bad-element>` + html.EscapeString(bad_elem_value) + `</bad-element>
	</error-info>
</rpc-error>`
}

func TestUnknownAttrRpcError(t *testing.T) {
	ncerr := NewUnknownAttrRpcError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal UnknownAttrRpcError error: %v\n", err)
		return
	}
	unmarshal := UnknownAttrApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal UnknownAttrRpcError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genUnknownAttrXml(rpc.String(), bad_attr_value, bad_elem_value))
}

func TestUnknownAttrProtocolError(t *testing.T) {
	ncerr := NewUnknownAttrProtocolError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal UnknownAttrProtocolError error: %v\n", err)
		return
	}
	unmarshal := UnknownAttrApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal UnknownAttrProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genUnknownAttrXml(protocol.String(), bad_attr_value, bad_elem_value))
}

func TestUnknownAttrApplicationError(t *testing.T) {
	ncerr := NewUnknownAttrApplicationError(bad_attr_value, bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal UnknownAttrApplicationError error: %v\n", err)
		return
	}
	unmarshal := UnknownAttrApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal UnknownAttrApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genUnknownAttrXml(application.String(), bad_attr_value, bad_elem_value))
}

func genMissingElementXml(typ, bad_elem_value string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(missing_element.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_missing_element) + `</error-message>
	<error-info>
		<bad-element>` + html.EscapeString(bad_elem_value) + `</bad-element>
	</error-info>
</rpc-error>`
}

func TestMissingElementProtocolError(t *testing.T) {
	ncerr := NewMissingElementProtocolError(bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal MissingElementProtocolError error: %v\n", err)
		return
	}
	unmarshal := UnknownAttrApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MissingElementProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genMissingElementXml(protocol.String(), bad_elem_value))
}

func TestMissingElementApplicationError(t *testing.T) {
	ncerr := NewMissingElementApplicationError(bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal MissingElementApplicationError error: %v\n", err)
		return
	}
	unmarshal := MissingElementApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MissingElementApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genMissingElementXml(application.String(), bad_elem_value))
}

func genBadElementXml(typ, bad_elem_value string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(bad_element.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_bad_element) + `</error-message>
	<error-info>
		<bad-element>` + html.EscapeString(bad_elem_value) + `</bad-element>
	</error-info>
</rpc-error>`
}

func TestBadElementProtocolError(t *testing.T) {
	ncerr := NewBadElementProtocolError(bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal BadElementProtocolError error: %v\n", err)
		return
	}
	unmarshal := BadElementApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal BadElementProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genBadElementXml(protocol.String(), bad_elem_value))
}

func TestBadElementApplicationError(t *testing.T) {
	ncerr := NewBadElementApplicationError(bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal BadElementApplicationError error: %v\n", err)
		return
	}
	unmarshal := BadElementApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal BadElementApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genBadElementXml(application.String(), bad_elem_value))
}

func genUnknownElementXml(typ, bad_elem_value string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(unknown_element.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_unknown_element) + `</error-message>
	<error-info>
		<bad-element>` + html.EscapeString(bad_elem_value) + `</bad-element>
	</error-info>
</rpc-error>`
}

func TestUnknownElementProtocolError(t *testing.T) {
	ncerr := NewUnknownElementProtocolError(bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal UnknownElementProtocolError error: %v\n", err)
		return
	}
	unmarshal := UnknownElementProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal UnknownElementProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genUnknownElementXml(protocol.String(), bad_elem_value))
}

func TestUnknownElementApplicationError(t *testing.T) {
	ncerr := NewUnknownElementApplicationError(bad_elem_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal UnknownElementApplicationError error: %v\n", err)
		return
	}
	unmarshal := UnknownElementApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal UnknownElementApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genUnknownElementXml(application.String(), bad_elem_value))
}

func ExampleUnknownElementApplicationError() {
	err := NewUnknownElementApplicationError("biz")
	err.Path = "/foo/bar"
	fmt.Println(err.Error())
}

func genUnknownNamespaceXml(typ, bad_elem_value, bad_ns_value string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(unknown_namespace.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_unknown_namespace) + `</error-message>
	<error-info>
		<bad-element>` + html.EscapeString(bad_elem_value) + `</bad-element>
		<bad-namespace>` + html.EscapeString(bad_ns_value) + `</bad-namespace>
	</error-info>
</rpc-error>`

	// Output:
	// Error: /foo/bar/biz: An unexpected element is present.
}

func TestUnknownNamespaceProtocolError(t *testing.T) {
	ncerr := NewUnknownNamespaceProtocolError(bad_elem_value, bad_ns_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal UnknownNamespaceProtocolError error: %v\n", err)
		return
	}
	unmarshal := UnknownNamespaceProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal UnknownNamespaceProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genUnknownNamespaceXml(protocol.String(), bad_elem_value, bad_ns_value))
}

func TestUnknownNamespaceApplicationError(t *testing.T) {
	ncerr := NewUnknownNamespaceApplicationError(bad_elem_value, bad_ns_value)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal UnknownNamespaceApplicationError error: %v\n", err)
		return
	}
	unmarshal := UnknownNamespaceApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal UnknownNamespaceApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genUnknownNamespaceXml(application.String(), bad_elem_value, bad_ns_value))
}

func genAccessDeniedXml(typ string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(access_denied.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_access_denied) + `</error-message>
</rpc-error>`
}

func TestAccessDeniedProtocolError(t *testing.T) {
	ncerr := NewAccessDeniedProtocolError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal AccessDeniedProtocolError error: %v\n", err)
		return
	}
	unmarshal := AccessDeniedProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal AccessDeniedProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genAccessDeniedXml(protocol.String()))
}

func TestAccessDeniedApplicationError(t *testing.T) {
	ncerr := NewAccessDeniedApplicationError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal AccessDeniedApplicationError error: %v\n", err)
		return
	}
	unmarshal := AccessDeniedApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal AccessDeniedApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genAccessDeniedXml(application.String()))
}

func genLockDeniedXml(sess string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(protocol.String()) + `</error-type>
	<error-tag>` + html.EscapeString(lock_denied.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_lock_denied) + `</error-message>
	<error-info>
		<session-id>` + html.EscapeString(sess) + `</session-id>
	</error-info>
</rpc-error>`
}

func TestLockDeniedError(t *testing.T) {
	const (
		sess = "1234"
	)
	ncerr := NewLockDeniedError(sess)
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal LockDeniedError error: %v\n", err)
		return
	}
	unmarshal := LockDeniedError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal LockDeniedError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genLockDeniedXml(sess))
}

func genResourceDeniedXml(typ string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(resource_denied.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_resource_denied) + `</error-message>
</rpc-error>`
}

func TestResourceDeniedTransportError(t *testing.T) {
	ncerr := NewResourceDeniedTransportError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal ResourceDeniedTransportError error: %v\n", err)
		return
	}
	unmarshal := ResourceDeniedTransportError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal ResourceDeniedTransportError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genResourceDeniedXml(transport.String()))
}

func TestResourceDeniedRpcError(t *testing.T) {
	ncerr := NewResourceDeniedRpcError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal ResourceDeniedRpcError error: %v\n", err)
		return
	}
	unmarshal := ResourceDeniedRpcError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal ResourceDeniedRpcError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genResourceDeniedXml(rpc.String()))
}

func TestResourceDeniedProtocolError(t *testing.T) {
	ncerr := NewResourceDeniedProtocolError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal ResourceDeniedProtocolError error: %v\n", err)
		return
	}
	unmarshal := ResourceDeniedProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal ResourceDeniedProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genResourceDeniedXml(protocol.String()))
}

func TestResourceDeniedApplicationError(t *testing.T) {
	ncerr := NewResourceDeniedApplicationError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal ResourceDeniedApplicationError error: %v\n", err)
		return
	}
	unmarshal := ResourceDeniedApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal ResourceDeniedApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genResourceDeniedXml(application.String()))
}

func genRollbackFailedError(typ string) *MgmtError {
	e := newMgmtError()
	e.Severity = nc_severity_error.String()
	e.Tag = rollback_failed.String()
	e.Typ = typ
	e.Message = msg_nc_rollback_failed
	return e
}

func genRollbackFailedXml(typ string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(rollback_failed.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_rollback_failed) + `</error-message>
</rpc-error>`
}

func TestRollbackFailedProtocolError(t *testing.T) {
	ncerr := NewRollbackFailedProtocolError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal RollbackFailedProtocolError error: %v\n", err)
		return
	}
	unmarshal := RollbackFailedProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal RollbackFailedProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genRollbackFailedXml(protocol.String()))
}

func TestRollbackFailedApplicationError(t *testing.T) {
	ncerr := NewRollbackFailedApplicationError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal RollbackFailedApplicationError error: %v\n", err)
		return
	}
	unmarshal := RollbackFailedApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal RollbackFailedApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genRollbackFailedXml(application.String()))
}

func genDataExistsXml() string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(application.String()) + `</error-type>
	<error-tag>` + html.EscapeString(data_exists.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_data_exists) + `</error-message>
</rpc-error>`
}

func TestDataExistsError(t *testing.T) {
	ncerr := NewDataExistsError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal DataExistsError error: %v\n", err)
		return
	}
	unmarshal := DataExistsError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal DataExistsError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genDataExistsXml())
}

func genDataMissingXml() string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(application.String()) + `</error-type>
	<error-tag>` + html.EscapeString(data_missing.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_data_missing) + `</error-message>
</rpc-error>`
}

func TestDataMissingError(t *testing.T) {
	ncerr := NewDataMissingError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal DataMissingError error: %v\n", err)
		return
	}
	unmarshal := DataMissingError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal DataMissingError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genDataMissingXml())
}

func genOperationNotSupportedXml(typ string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(operation_not_supported.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_operation_not_supported) + `</error-message>
</rpc-error>`
}

func TestOperationNotSupportedProtocolError(t *testing.T) {
	ncerr := NewOperationNotSupportedProtocolError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal OperationNotSupportedProtocolError error: %v\n", err)
		return
	}
	unmarshal := OperationNotSupportedProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal OperationNotSupportedProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genOperationNotSupportedXml(protocol.String()))
}

func TestOperationNotSupportedApplicationError(t *testing.T) {
	ncerr := NewOperationNotSupportedApplicationError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal OperationNotSupportedApplicationError error: %v\n", err)
		return
	}
	unmarshal := OperationNotSupportedApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal OperationNotSupportedApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genOperationNotSupportedXml(application.String()))
}

func genOperationFailedXml(typ string) string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(typ) + `</error-type>
	<error-tag>` + html.EscapeString(operation_failed.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_operation_failed) + `</error-message>
</rpc-error>`
}

func TestOperationFailedProtocolError(t *testing.T) {
	ncerr := NewOperationFailedProtocolError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal OperationFailedProtocolError error: %v\n", err)
		return
	}
	unmarshal := OperationFailedProtocolError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal OperationFailedProtocolError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genOperationFailedXml(protocol.String()))
}

func TestOperationFailedApplicationError(t *testing.T) {
	ncerr := NewOperationFailedApplicationError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal OperationFailedApplicationError error: %v\n", err)
		return
	}
	unmarshal := OperationFailedApplicationError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal OperationFailedApplicationError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genOperationFailedXml(application.String()))
}

func ExampleOperationFailedApplicationError() {
	err := NewOperationFailedApplicationError()
	err.Path = "/foo/bar"
	fmt.Println(err.Error())

	// Output:
	// Error: /foo/bar: Request could not be completed because the requested operation failed for some reason not covered by any other error condition.
}

func TestOperationFailedRpcError(t *testing.T) {
	ncerr := NewOperationFailedRpcError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal OperationFailedRpcError error: %v\n", err)
		return
	}
	unmarshal := OperationFailedRpcError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal OperationFailedRpcError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genOperationFailedXml(rpc.String()))
}

func genMalformedMessageXml() string {
	return `<rpc-error xmlns="` + netconf_namespace + `">
	<error-type>` + html.EscapeString(rpc.String()) + `</error-type>
	<error-tag>` + html.EscapeString(malformed_message.String()) + `</error-tag>
	<error-severity>` + html.EscapeString(nc_severity_error.String()) + `</error-severity>
	<error-message>` + html.EscapeString(msg_nc_malformed_message) + `</error-message>
</rpc-error>`
}

func TestMalformedMessageError(t *testing.T) {
	ncerr := NewMalformedMessageError()
	marshal, err := json.MarshalIndent(ncerr, "", "\t")
	if err != nil {
		t.Errorf("Marshal MalformedMessageError error: %v\n", err)
		return
	}
	unmarshal := MalformedMessageError{}
	if err := json.Unmarshal(marshal, &unmarshal); err != nil {
		t.Errorf("Unmarshal MalformedMessageError error: %v\n", err)
		return
	}
	cmpMgmtError(t, ncerr.MgmtError, unmarshal.MgmtError)

	verifyXmlMarshal(t, ncerr, genMalformedMessageXml())
}

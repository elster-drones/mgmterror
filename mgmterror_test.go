// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2016-2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package mgmterror

import (
	"encoding/json"
	"encoding/xml"
	"reflect"
	"testing"
)

func cmpMgmtError(t *testing.T, exp, unmarshal *MgmtError) {
	if !reflect.DeepEqual(exp, unmarshal) {
		t.Errorf("Failed JSON marshal/unmarshal")
		t.Logf("Expected: %#v", exp)
		t.Logf("Result:   %#v", unmarshal)
	}
}

func verifyMgmtErrorConstruction(t *testing.T, exp, rpcerr *MgmtError) {
	if !reflect.DeepEqual(exp, rpcerr) {
		t.Errorf("Unexpected %s/%s error constructed", exp.Tag, exp.Typ)
		t.Logf("Expected: %+v", exp)
		t.Logf("Constructed: %+v", rpcerr)
	}
}

func verifyXmlMarshal(t *testing.T, e interface{}, exp string) {
	marshal, err := xml.MarshalIndent(e, "", "\t")
	if err != nil {
		t.Errorf("XML Marshal error: %v\n", err)
		return
	}
	if exp != string(marshal) {
		t.Error("Unexpected XML marshal result")
		t.Logf("Expected: %s", exp)
		t.Logf("Marshal:  %s", marshal)
	}
}

func verifyMgmtErrorJson(t *testing.T, e *MgmtError) {
	marshal, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		t.Errorf("Marshal error: %v\n", err)
		return
	}

	unmarshal := newMgmtError()
	if err := json.Unmarshal(marshal, unmarshal); err != nil {
		t.Errorf("Unmarshal error: %v\n", err)
		return
	}
	if !reflect.DeepEqual(e, unmarshal) {
		t.Errorf("Failed %s/%s error JSON marshal/unmarshal", e.Tag, e.Typ)
		t.Logf("Expected: %#v", e)
		t.Logf("Result:   %#v", unmarshal)
	}
}

func TestMgmtErrorConstruction(t *testing.T) {
	exp := &MgmtError{
		XMLName: xml.Name{
			Space: "urn:ietf:params:xml:ns:netconf:base:1.0",
			Local: "rpc-error",
		},
	}
	verifyMgmtErrorConstruction(t, exp, newMgmtError())
}

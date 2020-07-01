// Copyright (c) 2017,2019-2020, AT&T Intellectual Property.
// All rights reserved.
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
)

type MgmtErrorList struct {
	errs []error
}

func (e MgmtErrorList) Errors() []error { return e.errs }

// Make sure the error has either a JSON or XML Marshaler.  If not,
// convert the "error" to a standard error.
func mkMgmtError(e error) error {
	switch e.(type) {
	case json.Marshaler, xml.Marshaler:
		return e
	case Formattable:
		me, _ := e.(Formattable)
		err := NewOperationFailedApplicationError()
		err.Message = me.GetMessage()
		err.Path = me.GetPath()
		return err
	default:
		err := NewOperationFailedApplicationError()
		err.Message = e.Error()
		return err
	}
}

func (e MgmtErrorList) MarshalJSON() ([]byte, error) {
	var out bytes.Buffer
	out.WriteString("{\"error-list\":[")
	for i, err := range e.errs {
		b, e := json.Marshal(err)
		if e != nil {
			return out.Bytes(), e
		}
		if i > 0 {
			out.WriteByte(',')
		}
		out.Write(b)
	}
	out.WriteString("]}")
	return out.Bytes(), nil
}

func (e *MgmtErrorList) UnmarshalJSON(value []byte) error {
	var errs struct {
		ErrorList []*MgmtError `json:"error-list"`
	}
	if err := json.Unmarshal(value, &errs); err != nil {
		return err
	}
	e.errs = []error{}
	for _, err := range errs.ErrorList {
		err.setXMLName()
		// NETCONF errors are the most generic (don't use
		// app-tag) so search them last.
		if vyerr := getVyattaError(err); vyerr != nil {
			e.MgmtErrorListAppend(vyerr)
		} else if yerr := getYangError(err); yerr != nil {
			e.MgmtErrorListAppend(yerr)
		} else if ncerr := getNetconfError(err); ncerr != nil {
			e.MgmtErrorListAppend(ncerr)
		} else {
			e.MgmtErrorListAppend(err)
		}
	}
	return nil
}

func (e MgmtErrorList) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	for _, err := range e.errs {
		if e := enc.Encode(err); e != nil {
			return e
		}
	}
	return nil
}

func (e *MgmtErrorList) MgmtErrorListAppend(errs ...error) {
	for _, err := range errs {
		e.errs = append(e.errs, mkMgmtError(err))
	}
}

func (e MgmtErrorList) Error() string {
	var b bytes.Buffer

	for i, err := range e.errs {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

type Formatter func(err error) string

func (e MgmtErrorList) CustomError(fmtFn Formatter) string {
	var b bytes.Buffer

	if fmtFn == nil {
		fmtFn = func(e error) string {
			return e.Error()
		}
	}

	for i, err := range e.errs {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(fmtFn(err))
	}

	return b.String()
}

// Encode error for DBus
func (e *MgmtErrorList) DBusError() (string, []interface{}) {
	body := make([]interface{}, 1)
	body[0] = e
	return "com.vyatta.mgmterror.list", body
}

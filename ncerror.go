// Copyright (c) 2017,2019, AT&T Intellectual Property. All rights reserved.
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
	"github.com/danos/utils/pathutil"
	"strings"
)

type ncerrseverity uint

const (
	nc_severity_error ncerrseverity = iota
)

var ncerrseveritymap = map[ncerrseverity]string{
	nc_severity_error: "error",
}

func (s ncerrseverity) String() string {
	return ncerrseveritymap[s]
}

type ncerrtag uint

// RFC6241 Apdx A
const (
	in_use ncerrtag = iota
	invalid_value
	too_big
	missing_attribute
	bad_attribute
	unknown_attribute
	missing_element
	bad_element
	unknown_element
	unknown_namespace
	access_denied
	lock_denied
	resource_denied
	rollback_failed
	data_exists
	data_missing
	operation_not_supported
	operation_failed
	//partial_operation /* obsolete */
	malformed_message
)

var ncerrtagmap = map[string]ncerrtag{
	"in-use":                  in_use,
	"invalid-value":           invalid_value,
	"too-big":                 too_big,
	"missing-attribute":       missing_attribute,
	"bad-attribute":           bad_attribute,
	"unknown-attribute":       unknown_attribute,
	"missing-element":         missing_element,
	"bad-element":             bad_element,
	"unknown-element":         unknown_element,
	"unknown-namespace":       unknown_namespace,
	"access-denied":           access_denied,
	"lock-denied":             lock_denied,
	"resource-denied":         resource_denied,
	"rollback-failed":         rollback_failed,
	"data-exists":             data_exists,
	"data-missing":            data_missing,
	"operation-not-supported": operation_not_supported,
	"operation-failed":        operation_failed,
	"malformed-message":       malformed_message,
}

func (t *ncerrtag) set(tag string) error {
	if v, ok := ncerrtagmap[tag]; ok {
		*t = v
		return nil
	}
	return errors.New("Invalid error tag")
}

func (t ncerrtag) String() string {
	for s, v := range ncerrtagmap {
		if t == v {
			return s
		}
	}
	// Can not happen
	return ""
}

const (
	msg_nc_in_use                  = `The request requires a resource that already is in use.`
	msg_nc_invalid_value           = `The request specifies an unacceptable value for one or more parameters.`
	msg_nc_too_big                 = `The request or response (that would be generated) is too large for the implementation to handle.`
	msg_nc_missing_attribute       = `An expected attribute is missing.`
	msg_nc_bad_attribute           = `An attribute value is not correct; e.g., wrong type, out of range, pattern mismatch.`
	msg_nc_unknown_attribute       = `An unexpected attribute is present.`
	msg_nc_missing_element         = `An expected element is missing.`
	msg_nc_bad_element             = `An element value is not correct; e.g., wrong type, out of range, pattern mismatch.`
	msg_nc_unknown_element         = `An unexpected element is present.`
	msg_nc_unknown_namespace       = `An unexpected namespace is present.`
	msg_nc_access_denied           = `Access to the requested protocol operation or data model is denied because authorization failed.`
	msg_nc_lock_denied             = `Access to the requested lock is denied because the lock is currently held by another entity.`
	msg_nc_resource_denied         = `Request could not be completed because of insufficient resources.`
	msg_nc_rollback_failed         = `Request to roll back some configuration change (via rollback-on-error or <discard-changes> operations) was not completed for some reason.`
	msg_nc_data_exists             = `Request could not be completed because the relevant data model content already exists.  For example, a "create" operation was attempted on data that already exists.`
	msg_nc_data_missing            = `Request could not be completed because the relevant data model content does not exist.  For example, a "delete" operation was attempted on data that does not exist.`
	msg_nc_operation_not_supported = `Request could not be completed because the requested operation is not supported by this implementation.`
	msg_nc_operation_failed        = `Request could not be completed because the requested operation failed for some reason not covered by any other error condition.`
	msg_nc_partial_operation       = `This error-tag is obsolete, and SHOULD NOT be sent by servers conforming to this document.`
	msg_nc_malformed_message       = `A message could not be handled because it failed to be parsed correctly.  For example, the message is not well-formed XML or it uses an invalid character set.`
)

type typeMap map[errtype]interface{}

type ncErrTag struct {
	severity ncerrseverity
	msg      string
	typ      typeMap
}

var ncErrTable map[ncerrtag]ncErrTag

func init() {
	// RFC6241 Apdx A
	ncErrTable = map[ncerrtag]ncErrTag{
		in_use: {
			severity: nc_severity_error,
			msg:      msg_nc_in_use,
			typ: typeMap{
				protocol:    createInUseProtocolError,
				application: createInUseApplicationError,
			},
		},
		invalid_value: {
			severity: nc_severity_error,
			msg:      msg_nc_invalid_value,
			typ: typeMap{
				protocol:    createInvalidValueProtocolError,
				application: createInvalidValueApplicationError,
			},
		},
		too_big: {
			severity: nc_severity_error,
			msg:      msg_nc_too_big,
			typ: typeMap{
				transport:   createTooBigTransportError,
				rpc:         createTooBigRpcError,
				protocol:    createTooBigProtocolError,
				application: createTooBigApplicationError,
			},
		},
		missing_attribute: {
			severity: nc_severity_error,
			msg:      msg_nc_missing_attribute,
			typ: typeMap{
				rpc:         createMissingAttrRpcError,
				protocol:    createMissingAttrProtocolError,
				application: createMissingAttrApplicationError,
			},
		},
		bad_attribute: {
			severity: nc_severity_error,
			msg:      msg_nc_bad_attribute,
			typ: typeMap{
				rpc:         createBadAttrRpcError,
				protocol:    createBadAttrProtocolError,
				application: createBadAttrApplicationError,
			},
		},
		unknown_attribute: {
			severity: nc_severity_error,
			msg:      msg_nc_unknown_attribute,
			typ: typeMap{
				rpc:         createUnknownAttrRpcError,
				protocol:    createUnknownAttrProtocolError,
				application: createUnknownAttrApplicationError,
			},
		},
		missing_element: {
			severity: nc_severity_error,
			msg:      msg_nc_missing_element,
			typ: typeMap{
				protocol:    createMissingElementProtocolError,
				application: createMissingElementApplicationError,
			},
		},
		bad_element: {
			severity: nc_severity_error,
			msg:      msg_nc_bad_element,
			typ: typeMap{
				protocol:    createBadElementProtocolError,
				application: createBadElementApplicationError,
			},
		},
		unknown_element: {
			severity: nc_severity_error,
			msg:      msg_nc_unknown_element,
			typ: typeMap{
				protocol:    createUnknownElementProtocolError,
				application: createUnknownElementApplicationError,
			},
		},
		unknown_namespace: {
			severity: nc_severity_error,
			msg:      msg_nc_unknown_namespace,
			typ: typeMap{
				protocol:    createUnknownNamespaceProtocolError,
				application: createUnknownNamespaceApplicationError,
			},
		},
		access_denied: {
			severity: nc_severity_error,
			msg:      msg_nc_access_denied,
			typ: typeMap{
				protocol:    createAccessDeniedProtocolError,
				application: createAccessDeniedApplicationError,
			},
		},
		lock_denied: {
			severity: nc_severity_error,
			msg:      msg_nc_lock_denied,
			typ: typeMap{
				protocol: createLockDeniedError,
			},
		},
		resource_denied: {
			severity: nc_severity_error,
			msg:      msg_nc_resource_denied,
			typ: typeMap{
				transport:   createResourceDeniedTransportError,
				rpc:         createResourceDeniedRpcError,
				protocol:    createResourceDeniedProtocolError,
				application: createResourceDeniedApplicationError,
			},
		},
		rollback_failed: {
			severity: nc_severity_error,
			msg:      msg_nc_rollback_failed,
			typ: typeMap{
				protocol:    createRollbackFailedProtocolError,
				application: createRollbackFailedApplicationError,
			},
		},
		data_exists: {
			severity: nc_severity_error,
			msg:      msg_nc_data_exists,
			typ: typeMap{
				application: createDataExistsError,
			},
		},
		data_missing: {
			severity: nc_severity_error,
			msg:      msg_nc_data_missing,
			typ: typeMap{
				application: createDataMissingError,
			},
		},
		operation_not_supported: {
			severity: nc_severity_error,
			msg:      msg_nc_operation_not_supported,
			typ: typeMap{
				protocol:    createOperationNotSupportedProtocolError,
				application: createOperationNotSupportedApplicationError,
			},
		},
		operation_failed: {
			severity: nc_severity_error,
			msg:      msg_nc_operation_failed,
			typ: typeMap{
				rpc:         createOperationFailedRpcError,
				protocol:    createOperationFailedProtocolError,
				application: createOperationFailedApplicationError,
			},
		},

		// Obsolete error
		// partial_operation: {
		// 	severity: nc_severity_error,
		// 	msg:      msg_nc_partial_operation
		// 	typ: typeMap{
		// 		application: struct{}{},
		// 	},
		// },

		// This error-tag is new in :base:1.1 and MUST NOT be sent to old clients.
		malformed_message: {
			severity: nc_severity_error,
			msg:      msg_nc_malformed_message,
			typ: typeMap{
				rpc: createMalformedMessageError,
			},
		},
	}
}

func getNetconfError(err *MgmtError) error {
	tag, ok := ncerrtagmap[err.Tag]
	if !ok {
		return nil
	}

	ncErrTag, ok := ncErrTable[tag]
	if !ok {
		return nil
	}

	errTypeId, ok := errtypemap[err.Typ]
	if !ok {
		return nil
	}

	if fn, ok := ncErrTag.typ[errTypeId]; ok {
		return callCreate(fn, err)
	}
	return nil
}

type ncErrInfoId uint

// RFC6421 Apdx A
const (
	bad_attribute_info ncErrInfoId = iota
	bad_element_info
	bad_namespace_info
	session_id_info
)

var ncErrInfoIdMap = map[ncErrInfoId]string{
	bad_attribute_info: "bad-attribute",
	bad_element_info:   "bad-element",
	bad_namespace_info: "bad-namespace",
	session_id_info:    "session-id",
}

func (i ncErrInfoId) String() string {
	if s, ok := ncErrInfoIdMap[i]; ok {
		return s
	}
	return ""
}

// Errors returned when trying to create a MgmtError
var invalid_error_tag = errors.New("invalid error tag")
var invalid_error_type = errors.New("invalid error type")
var invalid_error_tag_type = errors.New("invalid error type for tag")

func (e *MgmtError) setNcError(tag ncerrtag, typ, apptag, path string, info *MgmtErrorInfo) error {
	var errTypeId errtype
	ncErrTag, ok := ncErrTable[tag]
	if !ok {
		return invalid_error_tag
	}
	if errTypeId, ok = errtypemap[typ]; !ok {
		return invalid_error_type
	}
	if _, ok := ncErrTag.typ[errTypeId]; !ok {
		return invalid_error_tag_type
	}
	e.Tag = tag.String()
	e.Typ = typ
	e.Severity = ncErrTag.severity.String()
	e.Message = ncErrTag.msg
	e.AppTag = apptag
	e.Path = path
	if info != nil {
		e.Info = *info
	}
	return nil
}

func newNcError(tag ncerrtag, typ, apptag, path string, info *MgmtErrorInfo) *MgmtError {
	e := newMgmtError()
	if err := e.setNcError(tag, typ, apptag, path, info); err != nil {
		panic(err)
	}
	return e
}

func newAttrError(tag ncerrtag, typ, badAttr, badElem string) *MgmtError {
	info := MgmtErrorInfo{
		MgmtErrorInfoTag{
			XMLName: xml.Name{
				Local: bad_attribute_info.String(),
			},
			Value: badAttr,
		},
		MgmtErrorInfoTag{
			XMLName: xml.Name{
				Local: bad_element_info.String(),
			},
			Value: badElem,
		},
	}
	return newNcError(tag, typ, "", "", &info)
}

func newElemError(tag ncerrtag, typ, badElem string) *MgmtError {
	info := MgmtErrorInfo{
		MgmtErrorInfoTag{
			XMLName: xml.Name{
				Local: bad_element_info.String(),
			},
			Value: badElem,
		},
	}
	return newNcError(tag, typ, "", "", &info)
}

func newInUseError(typ string) *MgmtError {
	return newNcError(in_use, typ, "", "", nil)
}

type InUseProtocolError struct {
	*MgmtError
}

func (e *InUseProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *InUseProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createInUseProtocolError(err *MgmtError) *InUseProtocolError {
	return &InUseProtocolError{
		MgmtError: err,
	}
}

// Protocol error when a resource is already in use.
func NewInUseProtocolError() *InUseProtocolError {
	return createInUseProtocolError(newInUseError(protocol.String()))
}

type InUseApplicationError struct {
	*MgmtError
}

func (e *InUseApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *InUseApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createInUseApplicationError(err *MgmtError) *InUseApplicationError {
	return &InUseApplicationError{
		MgmtError: err,
	}
}

// Application error when a resource is already in use.
func NewInUseApplicationError() *InUseApplicationError {
	return createInUseApplicationError(newInUseError(application.String()))
}

func newInvalidValueError(typ string) *MgmtError {
	return newNcError(invalid_value, typ, "", "", nil)
}

type InvalidValueProtocolError struct {
	*MgmtError
}

func (e *InvalidValueProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *InvalidValueProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createInvalidValueProtocolError(err *MgmtError) *InvalidValueProtocolError {
	return &InvalidValueProtocolError{
		MgmtError: err,
	}
}

// Protocol error when a value for one or more parameters is invalid.
func NewInvalidValueProtocolError() *InvalidValueProtocolError {
	return createInvalidValueProtocolError(newInvalidValueError(protocol.String()))
}

type InvalidValueApplicationError struct {
	*MgmtError
}

func (e *InvalidValueApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *InvalidValueApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createInvalidValueApplicationError(err *MgmtError) *InvalidValueApplicationError {
	return &InvalidValueApplicationError{
		MgmtError: err,
	}
}

// Application error when a value for one or more parameters is invalid.
func NewInvalidValueApplicationError() *InvalidValueApplicationError {
	return createInvalidValueApplicationError(newInvalidValueError(application.String()))
}

func newTooBigError(typ string) *MgmtError {
	return newNcError(too_big, typ, "", "", nil)
}

type TooBigTransportError struct {
	*MgmtError
}

func (e *TooBigTransportError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *TooBigTransportError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createTooBigTransportError(err *MgmtError) *TooBigTransportError {
	return &TooBigTransportError{
		MgmtError: err,
	}
}

// Transport error when request or response (that would be generated)
// is too large for the implementation to handle.
func NewTooBigTransportError() *TooBigTransportError {
	return createTooBigTransportError(newTooBigError(transport.String()))
}

type TooBigRpcError struct {
	*MgmtError
}

func (e *TooBigRpcError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *TooBigRpcError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createTooBigRpcError(err *MgmtError) *TooBigRpcError {
	return &TooBigRpcError{
		err,
	}
}

// Rpc error when request or response (that would be generated)
// is too large for the implementation to handle.
func NewTooBigRpcError() *TooBigRpcError {
	return createTooBigRpcError(newTooBigError(rpc.String()))
}

type TooBigProtocolError struct {
	*MgmtError
}

func (e *TooBigProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *TooBigProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createTooBigProtocolError(err *MgmtError) *TooBigProtocolError {
	return &TooBigProtocolError{
		MgmtError: err,
	}
}

// Protocol error when request or response (that would be generated)
// is too large for the implementation to handle.
func NewTooBigProtocolError() *TooBigProtocolError {
	return createTooBigProtocolError(newTooBigError(protocol.String()))
}

type TooBigApplicationError struct {
	*MgmtError
}

func (e *TooBigApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *TooBigApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createTooBigApplicationError(err *MgmtError) *TooBigApplicationError {
	return &TooBigApplicationError{
		MgmtError: err,
	}
}

// Application error when request or response (that would be generated)
// is too large for the implementation to handle.
func NewTooBigApplicationError() *TooBigApplicationError {
	return createTooBigApplicationError(newTooBigError(application.String()))
}

func newMissingAttrError(typ, badAttr, badElem string) *MgmtError {
	return newAttrError(missing_attribute, typ, badAttr, badElem)
}

type MissingAttrRpcError struct {
	*MgmtError
}

func (e *MissingAttrRpcError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *MissingAttrRpcError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createMissingAttrRpcError(err *MgmtError) *MissingAttrRpcError {
	return &MissingAttrRpcError{
		MgmtError: err,
	}
}

// Rpc error when an expected attribute is missing
//
// badAttr is the name of the missing attribute
// badElem is the name of the element that is supposed to contain
// the missing attribute
func NewMissingAttrRpcError(badAttr, badElem string) *MissingAttrRpcError {
	return createMissingAttrRpcError(newMissingAttrError(rpc.String(), badAttr, badElem))
}

type MissingAttrProtocolError struct {
	*MgmtError
}

func (e *MissingAttrProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *MissingAttrProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createMissingAttrProtocolError(err *MgmtError) *MissingAttrProtocolError {
	return &MissingAttrProtocolError{
		MgmtError: err,
	}
}

// Protocol error when an expected attribute is missing
//
// badAttr is the name of the missing attribute
// badElem is the name of the element that is supposed to contain
// the missing attribute
func NewMissingAttrProtocolError(badAttr, badElem string) *MissingAttrProtocolError {
	return createMissingAttrProtocolError(newMissingAttrError(protocol.String(), badAttr, badElem))
}

type MissingAttrApplicationError struct {
	*MgmtError
}

func (e *MissingAttrApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *MissingAttrApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createMissingAttrApplicationError(err *MgmtError) *MissingAttrApplicationError {
	return &MissingAttrApplicationError{
		MgmtError: err,
	}
}

// Application error when an expected attribute is missing
//
// badAttr is the name of the missing attribute
// badElem is the name of the element that is supposed to contain
// the missing attribute
func NewMissingAttrApplicationError(badAttr, badElem string) *MissingAttrApplicationError {
	return createMissingAttrApplicationError(newMissingAttrError(application.String(), badAttr, badElem))
}

func newBadAttrError(typ, badAttr, badElem string) *MgmtError {
	return newAttrError(bad_attribute, typ, badAttr, badElem)
}

type BadAttrRpcError struct {
	*MgmtError
}

func (e *BadAttrRpcError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *BadAttrRpcError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createBadAttrRpcError(err *MgmtError) *BadAttrRpcError {
	return &BadAttrRpcError{
		MgmtError: err,
	}
}

// Rpc error when an attribute value is not correct
//
// badAttr is the name of the attribute with bad value
// badElem is the name of the element that contains the attribute with
// the bad value
func NewBadAttrRpcError(badAttr, badElem string) *BadAttrRpcError {
	return createBadAttrRpcError(newBadAttrError(rpc.String(), badAttr, badElem))
}

type BadAttrProtocolError struct {
	*MgmtError
}

func (e *BadAttrProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *BadAttrProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createBadAttrProtocolError(err *MgmtError) *BadAttrProtocolError {
	return &BadAttrProtocolError{
		MgmtError: err,
	}
}

// Protocol error when an attribute value is not correct
//
// badAttr is the name of the attribute with bad value
// badElem is the name of the element that contains the attribute with
// the bad value
func NewBadAttrProtocolError(badAttr, badElem string) *BadAttrProtocolError {
	return createBadAttrProtocolError(newBadAttrError(protocol.String(), badAttr, badElem))
}

type BadAttrApplicationError struct {
	*MgmtError
}

func (e *BadAttrApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *BadAttrApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createBadAttrApplicationError(err *MgmtError) *BadAttrApplicationError {
	return &BadAttrApplicationError{
		MgmtError: err,
	}
}

// Application error when an attribute value is not correct
//
// badAttr is the name of the attribute with bad value
// badElem is the name of the element that contains the attribute with
// the bad value
func NewBadAttrApplicationError(badAttr, badElem string) *BadAttrApplicationError {
	return createBadAttrApplicationError(newBadAttrError(application.String(), badAttr, badElem))
}

func newUnknownAttrError(typ, badAttr, badElem string) *MgmtError {
	return newAttrError(unknown_attribute, typ, badAttr, badElem)
}

type UnknownAttrRpcError struct {
	*MgmtError
}

func (e *UnknownAttrRpcError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *UnknownAttrRpcError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createUnknownAttrRpcError(err *MgmtError) *UnknownAttrRpcError {
	return &UnknownAttrRpcError{
		MgmtError: err,
	}
}

// Rpc error when an unexpected attribute is present
//
// badAttr is the name of the unexpected attribute
// badElem is the name of the element that contains the unexpected attribute
func NewUnknownAttrRpcError(badAttr, badElem string) *UnknownAttrRpcError {
	return createUnknownAttrRpcError(newUnknownAttrError(rpc.String(), badAttr, badElem))
}

type UnknownAttrProtocolError struct {
	*MgmtError
}

func (e *UnknownAttrProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *UnknownAttrProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createUnknownAttrProtocolError(err *MgmtError) *UnknownAttrProtocolError {
	return &UnknownAttrProtocolError{
		MgmtError: err,
	}
}

// Protocol error when an unexpected attribute is present
//
// badAttr is the name of the unexpected attribute
// badElem is the name of the element that contains the unexpected attribute
func NewUnknownAttrProtocolError(badAttr, badElem string) *UnknownAttrProtocolError {
	return createUnknownAttrProtocolError(newUnknownAttrError(protocol.String(), badAttr, badElem))
}

type UnknownAttrApplicationError struct {
	*MgmtError
}

func (e *UnknownAttrApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *UnknownAttrApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createUnknownAttrApplicationError(err *MgmtError) *UnknownAttrApplicationError {
	return &UnknownAttrApplicationError{
		MgmtError: err,
	}
}

// Application error when an unexpected attribute is present
//
// badAttr is the name of the unexpected attribute
// badElem is the name of the element that contains the unexpected attribute
func NewUnknownAttrApplicationError(badAttr, badElem string) *UnknownAttrApplicationError {
	return createUnknownAttrApplicationError(newUnknownAttrError(application.String(), badAttr, badElem))
}

func newMissingElemError(typ, badElem string) *MgmtError {
	return newElemError(missing_element, typ, badElem)
}

func missingElemErrorString(e *MgmtError) string {
	var b bytes.Buffer

	b.WriteString(strings.Title(e.Severity))
	b.WriteString(error_msg_separator)

	if e.Path != "" {
		b.WriteString(e.Path)
	}
	b.WriteByte('/')
	b.WriteString(e.Info[0].Value)

	if e.Message != "" {
		b.WriteString(error_msg_separator)
		b.WriteString(e.Message)
	}

	return b.String()
}

type MissingElementProtocolError struct {
	*MgmtError
}

func (e *MissingElementProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *MissingElementProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func (e *MissingElementProtocolError) Error() string {
	return missingElemErrorString(e.MgmtError)
}

func createMissingElementProtocolError(err *MgmtError) *MissingElementProtocolError {
	return &MissingElementProtocolError{
		MgmtError: err,
	}
}

// Protocol error when an expected element is missing
//
// badElem is the name of the missing element
func NewMissingElementProtocolError(badElem string) *MissingElementProtocolError {
	return createMissingElementProtocolError(newMissingElemError(protocol.String(), badElem))
}

type MissingElementApplicationError struct {
	*MgmtError
}

func (e *MissingElementApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *MissingElementApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func (e *MissingElementApplicationError) Error() string {
	return missingElemErrorString(e.MgmtError)
}

func createMissingElementApplicationError(err *MgmtError) *MissingElementApplicationError {
	return &MissingElementApplicationError{
		MgmtError: err,
	}
}

// Application error when an expected element is missing
//
// badElem is the name of the missing element
func NewMissingElementApplicationError(badElem string) *MissingElementApplicationError {
	return createMissingElementApplicationError(newMissingElemError(application.String(), badElem))
}

func newBadElemError(typ, badElem string) *MgmtError {
	return newElemError(bad_element, typ, badElem)
}

type BadElementProtocolError struct {
	*MgmtError
}

func (e *BadElementProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *BadElementProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createBadElementProtocolError(err *MgmtError) *BadElementProtocolError {
	return &BadElementProtocolError{
		MgmtError: err,
	}
}

// Protocol error when an element value is not correct
//
// badElem is the name of the element with bad value
func NewBadElementProtocolError(badElem string) *BadElementProtocolError {
	return createBadElementProtocolError(newBadElemError(protocol.String(), badElem))
}

type BadElementApplicationError struct {
	*MgmtError
}

func (e *BadElementApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *BadElementApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createBadElementApplicationError(err *MgmtError) *BadElementApplicationError {
	return &BadElementApplicationError{
		MgmtError: err,
	}
}

// Application error when an element value is not correct
//
// badElem is the name of the element with bad value
func NewBadElementApplicationError(badElem string) *BadElementApplicationError {
	return createBadElementApplicationError(newBadElemError(application.String(), badElem))
}

func newUnknownElemError(typ, badElem string) *MgmtError {
	return newElemError(unknown_element, typ, badElem)
}

func unknownElemErrorString(e *MgmtError) string {
	// TODO - all identical error functions should be using common code!!!
	var b bytes.Buffer

	b.WriteString(strings.Title(e.Severity))
	b.WriteString(error_msg_separator)

	if e.Path != "" {
		b.WriteString(e.Path)
	}
	b.WriteByte('/')
	b.WriteString(e.Info[0].Value)
	if e.Message != "" {
		b.WriteString(error_msg_separator)
		b.WriteString(e.Message)
	}

	return b.String()
}

type UnknownElementProtocolError struct {
	*MgmtError
}

// Too many copies of this, also getPathSlice is similar to makepath
func errpath(path []string) string {
	if len(path) < 2 {
		return fmt.Sprintf("%s", path)
	}
	path, val := path[:len(path)-1], path[len(path)-1]
	return fmt.Sprintf("%s [%s]", strings.Join(path, " "), val)
}

func (uepe *UnknownElementProtocolError) GetMessage() string {
	return fmt.Sprintf("%s is not valid",
		errpath(pathutil.Makepath(uepe.Path+"/"+uepe.Info[0].Value)))
}

func (e *UnknownElementProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *UnknownElementProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func (e *UnknownElementProtocolError) Error() string {
	return unknownElemErrorString(e.MgmtError)
}

func createUnknownElementProtocolError(err *MgmtError) *UnknownElementProtocolError {
	return &UnknownElementProtocolError{
		MgmtError: err,
	}
}

// Protocol error when an unexpected element is present
//
// badElem is the name of the unexpected element
func NewUnknownElementProtocolError(badElem string) *UnknownElementProtocolError {
	return createUnknownElementProtocolError(newUnknownElemError(protocol.String(), badElem))
}

type UnknownElementApplicationError struct {
	*MgmtError
}

func (ueae *UnknownElementApplicationError) GetMessage() string {
	return fmt.Sprintf("%s is not valid",
		errpath(pathutil.Makepath(ueae.Path+"/"+ueae.Info[0].Value)))
}

func (e *UnknownElementApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *UnknownElementApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func (e *UnknownElementApplicationError) Error() string {
	return unknownElemErrorString(e.MgmtError)
}

func createUnknownElementApplicationError(err *MgmtError) *UnknownElementApplicationError {
	return &UnknownElementApplicationError{
		MgmtError: err,
	}
}

// Application error when an unexpected element is present
//
// badElem is the name of the unexpected element
func NewUnknownElementApplicationError(badElem string) *UnknownElementApplicationError {
	return createUnknownElementApplicationError(newUnknownElemError(application.String(), badElem))
}

func newUnknownNamespaceError(typ, badElem, badNS string) *MgmtError {
	info := MgmtErrorInfo{
		MgmtErrorInfoTag{
			XMLName: xml.Name{
				Local: bad_element_info.String(),
			},
			Value: badElem,
		},
		MgmtErrorInfoTag{
			XMLName: xml.Name{
				Local: bad_namespace_info.String(),
			},
			Value: badNS,
		},
	}
	return newNcError(unknown_namespace, typ, "", "", &info)
}

type UnknownNamespaceProtocolError struct {
	*MgmtError
}

func (e *UnknownNamespaceProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *UnknownNamespaceProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createUnknownNamespaceProtocolError(err *MgmtError) *UnknownNamespaceProtocolError {
	return &UnknownNamespaceProtocolError{
		MgmtError: err,
	}
}

// Protocol error when an unexpected namespace is present
//
// badElem is the name of the element that contains the unexpected namespace
// badNS is the name of the unexpected namespace
func NewUnknownNamespaceProtocolError(badElem, badNS string) *UnknownNamespaceProtocolError {
	return createUnknownNamespaceProtocolError(newUnknownNamespaceError(protocol.String(), badElem, badNS))
}

type UnknownNamespaceApplicationError struct {
	*MgmtError
}

func (e *UnknownNamespaceApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *UnknownNamespaceApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createUnknownNamespaceApplicationError(err *MgmtError) *UnknownNamespaceApplicationError {
	return &UnknownNamespaceApplicationError{
		MgmtError: err,
	}
}

// Application error when an unexpected namespace is present
//
// badElem is the name of the element that contains the unexpected namespace
// badNS is the name of the unexpected namespace
func NewUnknownNamespaceApplicationError(badElem, badNS string) *UnknownNamespaceApplicationError {
	return createUnknownNamespaceApplicationError(newUnknownNamespaceError(application.String(), badElem, badNS))
}

func newAccessDeniedError(typ string) *MgmtError {
	return newNcError(access_denied, typ, "", "", nil)
}

type AccessDeniedProtocolError struct {
	*MgmtError
}

func (e *AccessDeniedProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *AccessDeniedProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createAccessDeniedProtocolError(err *MgmtError) *AccessDeniedProtocolError {
	return &AccessDeniedProtocolError{
		MgmtError: err,
	}
}

// Protocol error when access to the requested operation is denied
func NewAccessDeniedProtocolError() *AccessDeniedProtocolError {
	return createAccessDeniedProtocolError(newAccessDeniedError(protocol.String()))
}

type AccessDeniedApplicationError struct {
	*MgmtError
}

func (e *AccessDeniedApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *AccessDeniedApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createAccessDeniedApplicationError(err *MgmtError) *AccessDeniedApplicationError {
	return &AccessDeniedApplicationError{
		MgmtError: err,
	}
}

// Application error when access to the requested data model is denied
func NewAccessDeniedApplicationError() *AccessDeniedApplicationError {
	return createAccessDeniedApplicationError(newAccessDeniedError(application.String()))
}

type LockDeniedError struct {
	*MgmtError
}

func (e *LockDeniedError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *LockDeniedError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createLockDeniedError(err *MgmtError) *LockDeniedError {
	return &LockDeniedError{
		MgmtError: err,
	}
}

// Protocol error when access to the requested lock is denied
//
// sess is the session id that currently holds the lock or zero when a
// non-NETCONF entity holds the lock.
func NewLockDeniedError(sess string) *LockDeniedError {
	info := MgmtErrorInfo{
		MgmtErrorInfoTag{
			XMLName: xml.Name{
				Local: session_id_info.String(),
			},
			Value: sess,
		},
	}
	return createLockDeniedError(newNcError(lock_denied, protocol.String(), "", "", &info))
}

func newResourceDeniedError(typ string) *MgmtError {
	return newNcError(resource_denied, typ, "", "", nil)
}

type ResourceDeniedTransportError struct {
	*MgmtError
}

func (e *ResourceDeniedTransportError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *ResourceDeniedTransportError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createResourceDeniedTransportError(err *MgmtError) *ResourceDeniedTransportError {
	return &ResourceDeniedTransportError{
		MgmtError: err,
	}
}

// Transport error when request could not be completed because of
// insufficient resources.
func NewResourceDeniedTransportError() *ResourceDeniedTransportError {
	return createResourceDeniedTransportError(newResourceDeniedError(transport.String()))
}

type ResourceDeniedRpcError struct {
	*MgmtError
}

func (e *ResourceDeniedRpcError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *ResourceDeniedRpcError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createResourceDeniedRpcError(err *MgmtError) *ResourceDeniedRpcError {
	return &ResourceDeniedRpcError{
		MgmtError: err,
	}
}

// Rpc error when request could not be completed because of
// insufficient resources.
func NewResourceDeniedRpcError() *ResourceDeniedRpcError {
	return createResourceDeniedRpcError(newResourceDeniedError(rpc.String()))
}

type ResourceDeniedProtocolError struct {
	*MgmtError
}

func (e *ResourceDeniedProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *ResourceDeniedProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createResourceDeniedProtocolError(err *MgmtError) *ResourceDeniedProtocolError {
	return &ResourceDeniedProtocolError{
		MgmtError: err,
	}
}

// Protocol error when request could not be completed because of
// insufficient resources.
func NewResourceDeniedProtocolError() *ResourceDeniedProtocolError {
	return createResourceDeniedProtocolError(newResourceDeniedError(protocol.String()))
}

type ResourceDeniedApplicationError struct {
	*MgmtError
}

func (e *ResourceDeniedApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *ResourceDeniedApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createResourceDeniedApplicationError(err *MgmtError) *ResourceDeniedApplicationError {
	return &ResourceDeniedApplicationError{
		MgmtError: err,
	}
}

// Application error when request could not be completed because of
// insufficient resources.
func NewResourceDeniedApplicationError() *ResourceDeniedApplicationError {
	return createResourceDeniedApplicationError(newResourceDeniedError(application.String()))
}

func newRollbackFailedError(typ string) *MgmtError {
	return newNcError(rollback_failed, typ, "", "", nil)
}

type RollbackFailedProtocolError struct {
	*MgmtError
}

func (e *RollbackFailedProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *RollbackFailedProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createRollbackFailedProtocolError(err *MgmtError) *RollbackFailedProtocolError {
	return &RollbackFailedProtocolError{
		MgmtError: err,
	}
}

// Protocol error when request to roll back some configuration change
// (via rollback-on-error or <discard-changes> operations) was not
// completed.
func NewRollbackFailedProtocolError() *RollbackFailedProtocolError {
	return createRollbackFailedProtocolError(newRollbackFailedError(protocol.String()))
}

type RollbackFailedApplicationError struct {
	*MgmtError
}

func (e *RollbackFailedApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *RollbackFailedApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createRollbackFailedApplicationError(err *MgmtError) *RollbackFailedApplicationError {
	return &RollbackFailedApplicationError{
		MgmtError: err,
	}
}

// Application error when request to roll back some configuration
// change (via rollback-on-error or <discard-changes> operations) was
// not completed.
func NewRollbackFailedApplicationError() *RollbackFailedApplicationError {
	return createRollbackFailedApplicationError(newRollbackFailedError(application.String()))
}

type DataExistsError struct {
	*MgmtError
}

func (e *DataExistsError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *DataExistsError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createDataExistsError(err *MgmtError) *DataExistsError {
	return &DataExistsError{
		MgmtError: err,
	}
}

// Application error when the relevant data model content already
// exists.
func NewDataExistsError() *DataExistsError {
	return createDataExistsError(newNcError(data_exists, application.String(), "", "", nil))
}

type DataMissingError struct {
	*MgmtError
}

func (e *DataMissingError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *DataMissingError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createDataMissingError(err *MgmtError) *DataMissingError {
	return &DataMissingError{
		MgmtError: err,
	}
}

// Application error when the relevent data model content does not
// exist.
func NewDataMissingError() *DataMissingError {
	return createDataMissingError(newNcError(data_missing, application.String(), "", "", nil))
}

func newOperationNotSupportedError(typ string) *MgmtError {
	return newNcError(operation_not_supported, typ, "", "", nil)
}

type OperationNotSupportedProtocolError struct {
	*MgmtError
}

func (e *OperationNotSupportedProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *OperationNotSupportedProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createOperationNotSupportedProtocolError(err *MgmtError) *OperationNotSupportedProtocolError {
	return &OperationNotSupportedProtocolError{
		MgmtError: err,
	}
}

// Protocol error when the requested operation is not supported by
// this implementation.
func NewOperationNotSupportedProtocolError() *OperationNotSupportedProtocolError {
	return createOperationNotSupportedProtocolError(newOperationNotSupportedError(protocol.String()))
}

type OperationNotSupportedApplicationError struct {
	*MgmtError
}

func (e *OperationNotSupportedApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *OperationNotSupportedApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createOperationNotSupportedApplicationError(err *MgmtError) *OperationNotSupportedApplicationError {
	return &OperationNotSupportedApplicationError{
		MgmtError: err,
	}
}

// Application error when the requested operation is not supported by
// this implementation.
func NewOperationNotSupportedApplicationError() *OperationNotSupportedApplicationError {
	return createOperationNotSupportedApplicationError(newOperationNotSupportedError(application.String()))
}

func newOperationFailedError(typ string) *MgmtError {
	return newNcError(operation_failed, typ, "", "", nil)
}

type OperationFailedProtocolError struct {
	*MgmtError
}

func (e *OperationFailedProtocolError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *OperationFailedProtocolError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createOperationFailedProtocolError(err *MgmtError) *OperationFailedProtocolError {
	return &OperationFailedProtocolError{
		MgmtError: err,
	}
}

// Protocol error when the request could not be completed because the
// requested operation failed for some reason not covered by any other
// error condition.
func NewOperationFailedProtocolError() *OperationFailedProtocolError {
	return createOperationFailedProtocolError(newOperationFailedError(protocol.String()))
}

type OperationFailedApplicationError struct {
	*MgmtError
}

func (e *OperationFailedApplicationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *OperationFailedApplicationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createOperationFailedApplicationError(err *MgmtError) *OperationFailedApplicationError {
	return &OperationFailedApplicationError{
		MgmtError: err,
	}
}

// Application error when the request could not be completed because
// the requested operation failed for some reason not covered by any
// other error condition.
func NewOperationFailedApplicationError() *OperationFailedApplicationError {
	return createOperationFailedApplicationError(newOperationFailedError(application.String()))
}

type OperationFailedRpcError struct {
	*MgmtError
}

func (e *OperationFailedRpcError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *OperationFailedRpcError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createOperationFailedRpcError(err *MgmtError) *OperationFailedRpcError {
	return &OperationFailedRpcError{
		MgmtError: err,
	}
}

// Rpc error when the request could not be completed because the
// requested operation failed for some reason not covered by any other
// error condition.
func NewOperationFailedRpcError() *OperationFailedRpcError {
	return createOperationFailedRpcError(newOperationFailedError(rpc.String()))
}

type MalformedMessageError struct {
	*MgmtError
}

func (e *MalformedMessageError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *MalformedMessageError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createMalformedMessageError(err *MgmtError) *MalformedMessageError {
	return &MalformedMessageError{
		MgmtError: err,
	}
}

// Rpc error when a message could not be handled because it failed to
// be parsed correctly.
//
// This error is new in :base:1.1 and MUST NOT be sent to old clients.
func NewMalformedMessageError() *MalformedMessageError {
	return createMalformedMessageError(newNcError(malformed_message, "rpc", "", "", nil))
}

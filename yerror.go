// Copyright (c) 2017,2019 AT&T Intellectual Property
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
	"strings"
)

type yerrseverity uint

const (
	yang_severity_error yerrseverity = iota
	yang_severity_warning
)

const (
	needNodePath = ""
	noYangPath   = ""
	needYangPath = ""
)

var yerrseveritymap = map[yerrseverity]string{
	yang_severity_error:   "error",
	yang_severity_warning: "warning",
}

func (s yerrseverity) String() string {
	return yerrseveritymap[s]
}

type yerrtag uint

const (
	// RFC6020 Sect 13
	yang_operation_failed yerrtag = iota
	yang_data_missing
	yang_bad_attribute
)

var errtagmap = map[string]yerrtag{
	"operation-failed": yang_operation_failed,
	"data-missing":     yang_data_missing,
	"bad-attribute":    yang_bad_attribute,
}

func (t *yerrtag) set(tag string) error {
	if v, ok := errtagmap[tag]; ok {
		*t = v
		return nil
	}
	return errors.New("Invalid YANG error tag")
}

func (t yerrtag) String() string {
	for s, v := range errtagmap {
		if t == v {
			return s
		}
	}
	// Can not happen
	return ""
}

type yerrapptagid uint

const (
	// RFC6020 Section 13
	data_not_unique yerrapptagid = iota
	too_many_elements
	too_few_elements
	must_violation
	instance_required
	missing_choice
	missing_instance
)

var yerrapptagmap = map[string]yerrapptagid{
	"data-not-unique":   data_not_unique,
	"too-many-elements": too_many_elements,
	"too-few-elements":  too_few_elements,
	"must-violation":    must_violation,
	"instance-required": instance_required,
	"missing-choice":    missing_choice,
	"missing-instance":  missing_instance,
}

func (t yerrapptagid) String() string {
	for s, v := range yerrapptagmap {
		if t == v {
			return s
		}
	}
	// Can not happen
	return ""
}

type yErrAppTag string

func (t *yErrAppTag) set(tag string) {
	*t = yErrAppTag(tag)
}

type appTagMap map[yerrapptagid]interface{}

type yangErrTag struct {
	severity yerrseverity
	msg      string
	apptag   appTagMap
}

const (
	msg_yang_operation_failed = `The requested operation failed.`
	msg_yang_data_missing     = `Expected data is missing.`
	msg_yang_bad_attribute    = `An attribute value is not correct; e.g., wrong type, out of range, pattern mismatch.`
)

var yangErrTable map[yerrtag]yangErrTag

// TODO: add name for when converted to DBusError
// TODO: NewInstanceRequiredError and NewLeafrefMismatchError are indistinguishable
func init() {
	yangErrTable = map[yerrtag]yangErrTag{
		// RFC6241 Apdx A
		yang_operation_failed: {
			severity: yang_severity_error,
			msg:      msg_yang_operation_failed,
			apptag: appTagMap{
				data_not_unique:   createNonUniqueError,
				too_many_elements: createTooManyElementsError,
				too_few_elements:  createTooFewElementsError,
				must_violation:    createMustViolationError,
			},
		},
		yang_data_missing: {
			severity: yang_severity_error,
			msg:      msg_yang_data_missing,
			apptag: appTagMap{
				instance_required: createInstanceRequiredError,
				missing_choice:    createMissingChoiceError,
			},
		},
		yang_bad_attribute: {
			severity: yang_severity_error,
			msg:      msg_yang_bad_attribute,
			apptag: appTagMap{
				missing_instance: createInsertFailedError,
			},
		},
	}
}

func getYangError(err *MgmtError) error {
	tag, ok := errtagmap[err.Tag]
	if !ok {
		return nil
	}
	errtag, ok := yangErrTable[tag]
	if !ok {
		return nil
	}
	apptag, ok := yerrapptagmap[err.AppTag]
	if !ok {
		return nil
	}
	fn, ok := errtag.apptag[apptag]
	if !ok {
		return nil
	}
	return callCreate(fn, err)
}

type yangErrInfoId uint

// RFC6020 Sect 13
const (
	non_unique_info yangErrInfoId = iota
	missing_choice_info
)

var yangErrInfoIdMap = map[yangErrInfoId]string{
	non_unique_info:     "non-unique",
	missing_choice_info: "missing-choice",
}

func (i yangErrInfoId) String() string {
	if s, ok := yangErrInfoIdMap[i]; ok {
		return s
	}
	return ""
}

type YangError struct {
	*MgmtError
}

var invalid_error_app_tag = errors.New("invalid error app tag")

// setYangError - store details of YANG error
//
// 'path' is an overloaded term here.  ALL NETCONF errors must have a path,
// representing the 'node' in the data model on which the error has occurred.
// Separately, a YANG error *may* have an associated path, which represents
// the likes of a leafref reference that doesn't exist, ie it points to a
// different node.
func (e *MgmtError) setYangError(
	tag yerrtag,
	apptag, path, yangPath string,
	info *MgmtErrorInfo,
) error {

	yangErrTag, ok := yangErrTable[tag]
	if !ok {
		return invalid_error_tag
	}
	e.Tag = tag.String()
	e.Typ = "application"
	e.Severity = yangErrTag.severity.String()
	e.Message = yangErrTag.msg
	e.AppTag = apptag
	e.Path = path
	if info != nil {
		e.Info = *info
	}
	return nil
}

func newYangError(
	tag yerrtag,
	apptag, path, yangPath string,
	info *MgmtErrorInfo,
) *MgmtError {

	e := newMgmtError()
	if err := e.setYangError(tag, apptag, path, yangPath, info); err != nil {
		panic(err)
	}
	return e
}

// RFC6020 Sect 13.1
// Error Message for Data That Violates a unique Statement
type NonUniqueError struct {
	*MgmtError
}

func (e *NonUniqueError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *NonUniqueError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func (e *NonUniqueError) Error() string {
	if len(e.Info) < 2 {
		return e.MgmtError.Error()
	}
	var b bytes.Buffer
	b.WriteString(strings.Title(e.Severity))
	b.WriteString(error_msg_separator)
	b.WriteString(e.Path)
	b.WriteString(error_msg_separator)
	b.WriteString("Non-unique paths ")
	for i, p := range e.Info {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(strings.TrimPrefix(p.Value, e.Path+"/"))
	}
	return b.String()
}

func createNonUniqueError(err *MgmtError) *NonUniqueError {
	return &NonUniqueError{
		MgmtError: err,
	}
}

// When a NETCONF operation would result in configuration data where a
// unique constraint is invalidated.
//
// paths are a set of instance identifiers that points to a leaf that
// invalidates the unique constraint. This element is present once for
// each non-unique leaf.
func NewNonUniqueError(paths []string) *NonUniqueError {
	var info MgmtErrorInfo
	for _, p := range paths {
		i := MgmtErrorInfoTag{
			XMLName: xml.Name{
				Space: yang_namespace,
				Local: non_unique_info.String(),
			},
			Value: p,
		}
		info = append(info, i)
	}
	return createNonUniqueError(newYangError(yang_operation_failed,
		data_not_unique.String(), needNodePath, noYangPath, &info))
}

// RFC6020 Sect 13.2
// Error Message for Data That Violates a max-elements Statement
type TooManyElementsError struct {
	*MgmtError
}

func (e *TooManyElementsError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *TooManyElementsError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createTooManyElementsError(err *MgmtError) *TooManyElementsError {
	return &TooManyElementsError{
		MgmtError: err,
	}
}

// When a NETCONF operation would result in configuration data where a
// list or a leaf-list would have too many entries. This error should
// only be returned once even if there are more than one extra child
// present.
//
// path is the absolute XPath expression identifying the list node
func NewTooManyElementsError(path string) *TooManyElementsError {
	return createTooManyElementsError(newYangError(yang_operation_failed,
		too_many_elements.String(), path, noYangPath, nil))
}

// RFC6020 Sect 13.3
// Error Message for Data That Violates a min-elements Statement
type TooFewElementsError struct {
	*MgmtError
}

func (e *TooFewElementsError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *TooFewElementsError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createTooFewElementsError(err *MgmtError) *TooFewElementsError {
	return &TooFewElementsError{
		MgmtError: err,
	}
}

// When a NETCONF operation would result in configuration data where a
// list or a leaf-list would have too few entries. This error should
// only be returned once even if there are more than one child
// missing.
//
// path is the absolute XPath expression identifying the list node
func NewTooFewElementsError(path string) *TooFewElementsError {
	return createTooFewElementsError(newYangError(yang_operation_failed,
		too_few_elements.String(), path, noYangPath, nil))
}

// RFC6020 Sect 13.4
// Error Message for Data That Violates a must Statement
type MustViolationError struct {
	*MgmtError
}

func (e *MustViolationError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *MustViolationError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createMustViolationError(err *MgmtError) *MustViolationError {
	return &MustViolationError{
		MgmtError: err,
	}
}

// When a NETCONF operation would result in configuration data where
// the restrictions imposed by a "must" statement is violated
func NewMustViolationError() *MustViolationError {
	return createMustViolationError(newYangError(yang_operation_failed,
		must_violation.String(), needNodePath, noYangPath, nil))
}

// RFC6020 Sect 13.5
// Error Message for Data That Violates a require-instance Statement
type InstanceRequiredError struct {
	*MgmtError
}

func (e *InstanceRequiredError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *InstanceRequiredError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createInstanceRequiredError(err *MgmtError) *InstanceRequiredError {
	return &InstanceRequiredError{
		MgmtError: err,
	}
}

// When a NETCONF operation would result in configuration data where a
// leaf of type "instance-identifier" marked with require-instance
// "true" refers to a non-existing instance
//
// path is the absolute XPath expression identifying the
// instance-identifier leaf.
func NewInstanceRequiredError(path string) *InstanceRequiredError {
	return createInstanceRequiredError(newYangError(yang_data_missing,
		instance_required.String(), path, needYangPath, nil))
}

// RFC6020 Sect 13.6
// Error Message for Data That Does Not Match a leafref Type
type LeafrefMismatchError struct {
	*MgmtError
}

func (e *LeafrefMismatchError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *LeafrefMismatchError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createLeafrefMismatchError(err *MgmtError) *LeafrefMismatchError {
	return &LeafrefMismatchError{
		MgmtError: err,
	}
}

// When a NETCONF operation would result in configuration data where a
// leaf of type "leafref" refers to a non-existing instance
//
// path is the absolute XPath expression identifying the leafref leaf.
func NewLeafrefMismatchError(path, lrefPath string) *LeafrefMismatchError {
	return createLeafrefMismatchError(newYangError(yang_data_missing,
		instance_required.String(), path, lrefPath, nil))
}

// RFC6020 Sect 13.7
// Error Message for Data That Violates a mandatory choice Statement
type MissingChoiceError struct {
	*MgmtError
}

func (e *MissingChoiceError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *MissingChoiceError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createMissingChoiceError(err *MgmtError) *MissingChoiceError {
	return &MissingChoiceError{
		MgmtError: err,
	}
}

// When a NETCONF operation would result in configuration data where no
// nodes exists in a mandatory choice
//
// path is the absolute XPath expression identifying the element with
// the missing choice.
// name is the missing mandatory choice
func NewMissingChoiceError(path, name string) *MissingChoiceError {
	info := MgmtErrorInfo{
		MgmtErrorInfoTag{
			XMLName: xml.Name{
				Space: yang_namespace,
				Local: missing_choice_info.String(),
			},
			Value: name,
		},
	}
	return createMissingChoiceError(newYangError(yang_operation_failed,
		missing_choice.String(), path, needYangPath, &info))
}

// RFC6020 Sect 13.8
// Error Message for the "insert" Operation
type InsertFailedError struct {
	*MgmtError
}

func (e *InsertFailedError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *InsertFailedError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createInsertFailedError(err *MgmtError) *InsertFailedError {
	return &InsertFailedError{
		MgmtError: err,
	}
}

// When a NETCONF <edit-config> uses the "insert" and "key" or "value"
// attributes for a list or leaf-list node, and the "key" or "value"
// refers to a non-existing instance.
func NewInsertFailedError() *InsertFailedError {
	return createInsertFailedError(newYangError(yang_bad_attribute,
		missing_instance.String(), needNodePath, noYangPath, nil))
}

// Copyright (c) 2017-2019, AT&T Intellectual Property.  All rights reserved.
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
	"strings"

	"github.com/danos/utils/natsort"
	"github.com/danos/utils/pathutil"
)

type vyErrTag uint

const (
	vyatta_operation_failed vyErrTag = iota
)

var vyErrTagMap = map[string]vyErrTag{
	"operation-failed": vyatta_operation_failed,
}

func (t vyErrTag) String() string {
	for s, v := range vyErrTagMap {
		if t == v {
			return s
		}
	}
	// Can not happen
	return ""
}

type vyErrAppTagId uint

const (
	exec_failed vyErrAppTagId = iota
	path_ambig
)

var vyErrAppTagMap = map[string]vyErrAppTagId{
	"exec-failed":    exec_failed,
	"path-ambiguous": path_ambig,
}

func (t vyErrAppTagId) String() string {
	for s, v := range vyErrAppTagMap {
		if t == v {
			return s
		}
	}
	// Can not happen
	return ""
}

type vyAppTagMap map[vyErrAppTagId]interface{}

type vyError struct {
	severity yerrseverity
	msg      string
	apptag   vyAppTagMap
}

var vyErrTable map[vyErrTag]vyError

func init() {
	vyErrTable = map[vyErrTag]vyError{
		vyatta_operation_failed: {
			severity: yang_severity_error,
			msg:      msg_yang_operation_failed,
			apptag: vyAppTagMap{
				exec_failed: createExecError,
				path_ambig:  createPathAmbigError,
			},
		},
	}
}

func getVyattaError(err *MgmtError) error {
	tag, ok := vyErrTagMap[err.Tag]
	if !ok {
		return nil
	}
	errtag, ok := vyErrTable[tag]
	if !ok {
		return nil
	}
	apptag, ok := vyErrAppTagMap[err.AppTag]
	if !ok {
		return nil
	}
	fn, ok := errtag.apptag[apptag]
	if !ok {
		return nil
	}
	return callCreate(fn, err)
}

func (e *MgmtError) setVyattaError(tag vyErrTag, apptag, path string, info *MgmtErrorInfo) error {
	vyErr, ok := vyErrTable[tag]
	if !ok {
		return invalid_error_tag
	}
	e.Tag = tag.String()
	e.Typ = "application"
	e.Severity = vyErr.severity.String()
	e.Message = vyErr.msg
	e.AppTag = apptag
	e.Path = path
	if info != nil {
		e.Info = *info
	}
	return nil
}

func newVyattaError(tag vyErrTag, apptag, path string, info *MgmtErrorInfo) *MgmtError {
	e := newMgmtError()
	if err := e.setVyattaError(tag, apptag, path, info); err != nil {
		panic(err)
	}
	return e
}

type ExecError struct {
	*MgmtError
}

func (e *ExecError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *ExecError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func createExecError(err *MgmtError) *ExecError {
	return &ExecError{
		MgmtError: err,
	}
}

// A custom wrapper of a standard "operation failed" to represent an
// error when executing subtasks.
//
// path is the path of the subtask that was run
// out is the output of the subtask
func NewExecError(path []string, out string) *ExecError {
	err := newVyattaError(vyatta_operation_failed, exec_failed.String(),
		pathutil.Pathstr(path), nil)
	err.Message = out
	return createExecError(err)
}

type PathAmbiguousError struct {
	*MgmtError
}

func (e *PathAmbiguousError) UnmarshalJSON(value []byte) error {
	e.MgmtError = newMgmtError()
	return json.Unmarshal(value, e.MgmtError)
}

func (e *PathAmbiguousError) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.Encode(e.MgmtError)
}

func (e *PathAmbiguousError) GetMessage() string {
	var b bytes.Buffer
	b.WriteString(ErrPath(pathutil.Makepath(e.Path)))
	b.WriteString(" is ambiguous\n")

	b.WriteString("EZ9: Possible completions:\n")
	pathMap := make(map[string]string, len(e.Info))
	for _, elem := range e.Info {
		pathMap[elem.XMLName.Local] = elem.Value
	}
	sorted := make([]string, 0, len(pathMap))
	for n, _ := range pathMap {
		sorted = append(sorted, n)
	}

	natsort.Sort(sorted)
	for i, name := range sorted {
		b.WriteString(fmt.Sprintf("  %s\t%s", name, pathMap[name]))
		if i != len(sorted)-1 {
			b.WriteString("\n")
		}
	}
	return b.String()
}

func (e *PathAmbiguousError) Error() string {
	var mlist []string
	for _, t := range e.Info {
		mlist = append(mlist, t.XMLName.Local)
	}
	natsort.Sort(mlist)

	var b bytes.Buffer
	b.WriteString(strings.Title(e.Severity))
	b.WriteString(error_msg_separator)
	if len(e.Path) == 0 {
		b.WriteString("Ambiguous command")
	} else {
		b.WriteString(e.Path)
		b.WriteString(error_msg_separator)
		b.WriteString("Ambiguous path")
	}
	b.WriteString(", could be one of: ")
	for i, m := range mlist {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(m)
	}
	return b.String()
}

func createPathAmbigError(err *MgmtError) *PathAmbiguousError {
	return &PathAmbiguousError{
		MgmtError: err,
	}
}

// A custom wrapper of a standard "operation failed" to represent an
// error when a path is ambiguous
//
// path is the path that was ambiguous
// matches are the possible completions for path
func NewPathAmbiguousError(path []string, matches map[string]string) *PathAmbiguousError {
	var info MgmtErrorInfo
	for k, v := range matches {
		t := NewMgmtErrorInfoTag(VyattaNamespace, k, v)
		info = append(info, *t)
	}
	err := newVyattaError(vyatta_operation_failed, path_ambig.String(),
		pathutil.Pathstr(path), &info)
	return createPathAmbigError(err)
}

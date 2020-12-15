// Copyright (c) 2017-2020, AT&T Intellectual Property. All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0
//
// Useful test functions for validating mgmterrors.  Wraps the management
// errors to allow for different formatting for CLI, RPC over netconf etc.

package errtest

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/danos/mgmterror"
)

const (
	noPath = ""
)

var noInfo = []*mgmterror.MgmtErrorInfoTag{}

type ExpMgmtError struct {
	nameForDebug string
	// These fields will always be matched, even if empty
	expMsgContents []string
	expPath        string
	expInfo        []*mgmterror.MgmtErrorInfoTag
	// These fields will be ignored if empty.  Typically of less interest,
	// with default settings.
	expType     string
	expSeverity string
	expTag      string
	expAppTag   string
}

// Original constructor, sets Msg, Path and Info only.
func NewExpMgmtError(
	msgs []string,
	path string,
	info []*mgmterror.MgmtErrorInfoTag,
) *ExpMgmtError {

	return &ExpMgmtError{
		nameForDebug:   "Unspecified Error",
		expMsgContents: msgs, // Actual error should contain all these.
		expPath:        path, // Absolute match
		expInfo:        info, // May be empty
	}
}

func (eme *ExpMgmtError) SetName(name string) *ExpMgmtError {
	eme.nameForDebug = name
	return eme
}

func (eme *ExpMgmtError) SetType(typ string) *ExpMgmtError {
	eme.expType = typ
	return eme
}

func (eme *ExpMgmtError) SetTag(tag string) *ExpMgmtError {
	eme.expTag = tag
	return eme
}

func (eme *ExpMgmtError) SetAppTag(appTag string) *ExpMgmtError {
	eme.expAppTag = appTag
	return eme
}

func (eme *ExpMgmtError) SetSeverity(sev string) *ExpMgmtError {
	eme.expSeverity = sev
	return eme
}

// Constructors for some common errors.  Avoids repetition of common fields
// and allows for these to be modified in one place if needed.

// leafrefPath: space separated, no []
// nodeWithErrorPath: / separated, with leading /
func LeafrefMgmtErr(leafrefPath, nodeWithErrorPath string) *ExpMgmtError {
	return NewExpMgmtError(
		[]string{
			"The following path must exist:",
			fmt.Sprintf("[%s]", leafrefPath)},
		nodeWithErrorPath,
		noInfo).
		SetName("Non-existent Leafref").
		SetType("application").
		SetTag("operation-failed")
}

func MissingMandatoryNodeMgmtErr(missingNode, path string) *ExpMgmtError {
	return NewExpMgmtError(
		[]string{fmt.Sprintf("Missing mandatory node %s", missingNode)},
		path,
		noInfo).
		SetName("Missing Mandatory Node").
		SetType("application").
		SetTag("operation-failed")
}

func MustViolationMgmtErr(errmsg, path string) *ExpMgmtError {
	return NewExpMgmtError(
		[]string{errmsg},
		path,
		noInfo).
		SetName("Must Violation").
		SetType("application").
		SetTag("operation-failed")
}

func UniqueViolationMgmtErr(
	nonUniquePath, commonKeys, path string,
) *ExpMgmtError {
	return NewExpMgmtError(
		[]string{
			"The following path must be unique",
			fmt.Sprintf("[%s]", nonUniquePath),
			"but is defined in the following set of keys:",
			fmt.Sprintf("[%s]", commonKeys)},
		path,
		noInfo).
		SetName("Unique Violation").
		SetType("application").
		SetTag("operation-failed")
}

// Rough and ready check that all parts of all warnings appear at some point
// in the log
func CheckMgmtErrorsInLog(
	t *testing.T,
	log bytes.Buffer,
	expWarns []*ExpMgmtError,
) {
	logStr := log.String()
	for _, expWarn := range expWarns {
		if !strings.Contains(logStr, expWarn.expPath) {
			t.Fatalf("Syslog doesn't contain path: %s\n", expWarn.expPath)
		}
		for _, msg := range expWarn.expMsgContents {
			if !strings.Contains(logStr, msg) {
				t.Fatalf("Syslog doesn't contain msg: %s\n", msg)
			}
		}
		if len(expWarn.expInfo) > 0 {
			for _, info := range expWarn.expInfo {
				if !strings.Contains(logStr, info.XMLName.Space) {
					t.Fatalf("Syslog doesn't contain info space %s\n",
						info.XMLName.Space)
				}
				if !strings.Contains(logStr, info.XMLName.Local) {
					t.Fatalf("Syslog doesn't contain info local %s\n",
						info.XMLName.Local)
				}
				if !strings.Contains(logStr, info.Value) {
					t.Fatalf("Syslog doesn't contain info value %s\n",
						info.Value)
				}
			}
		}
	}
}

func setAndNoMatch(exp, act string) bool {
	return exp != "" && exp != act
}

func (eme *ExpMgmtError) Matches(actualErr mgmterror.Formattable) bool {
	if actualErr.GetPath() != eme.expPath {
		return false
	}
	if !checkInfoMatchesNonFatal(actualErr, eme.expInfo) {
		return false
	}
	for _, expMsg := range eme.expMsgContents {
		if !strings.Contains(actualErr.GetMessage(), expMsg) {
			return false
		}
	}
	if setAndNoMatch(eme.expType, actualErr.GetType()) {
		return false
	}
	if setAndNoMatch(eme.expAppTag, actualErr.GetAppTag()) {
		return false
	}
	if setAndNoMatch(eme.expTag, actualErr.GetTag()) {
		return false
	}
	if setAndNoMatch(eme.expSeverity, actualErr.GetSeverity()) {
		return false
	}
	return true
}

func CheckMgmtErrors(
	t *testing.T,
	expMgmtErrs []*ExpMgmtError,
	actualErrs []error,
) {
	// Check all actual errors were expected.  We assume all actual errors
	// are mgmterror.Formattable - if not then you're using the wrong test
	// function!
	for _, actErr := range actualErrs {
		me, _ := actErr.(mgmterror.Formattable)

		found := false
		for _, expErr := range expMgmtErrs {
			if !expErr.Matches(me) {
				continue
			}
			found = true
			break
		}
		if !found {
			expErr := expMgmtErrs[0]
			t.Logf("Expecting:\n"+
				"\tPath:\t%s\n\tMsg:\t%s\n\tTag:\t%s\n"+
				"\tType:\t%s\n\tSev:\t%s\n\tAppTag:\t%s\n",
				expErr.expPath, expErr.expMsgContents, expErr.expTag,
				expErr.expType, expErr.expSeverity, expErr.expAppTag)
			for _, info := range expErr.expInfo {
				t.Logf("\tInfo: NS %s:%s, Value %s\n",
					info.XMLName.Space, info.XMLName.Local, info.Value)
			}
			t.Fatalf(
				"Found unexpected error:\n"+
					"\tPath:\t%s\n\tMsg:\t%s\n\tTag:\t%s\n"+
					"\tType:\t%s\n\tSev:\t%s\n\tAppTag:\t%s\n"+
					"\tInfo:\t%s\n",
				me.GetPath(), me.GetMessage(), me.GetTag(),
				me.GetType(), me.GetSeverity(), me.GetAppTag(),
				me.GetInfo())
			return
		}
	}

	// Now check all expected errors were seen.
	for _, expErr := range expMgmtErrs {
		found := false
		for _, actErr := range actualErrs {
			me, _ := actErr.(mgmterror.Formattable)
			if !expErr.Matches(me) {
				continue
			}
			found = true
			break
		}
		if !found {
			t.Fatalf(
				"Error not found:\n\tPath:\t%s\n\tMsgs:\t%v\nInfo:\t%s\n",
				expErr.expPath, expErr.expMsgContents, expErr.expInfo)
			return
		}
	}
}

func CheckPath(t *testing.T, err error, expPath string) {
	me, ok := err.(mgmterror.Formattable)
	if !ok {
		t.Fatalf("Error does not meet Formattable interface!")
		return
	}

	if me.GetPath() != expPath {
		t.Fatalf("Path mismatch:\nExp:\t'%s'\nGot:\t'%s'\n",
			expPath, me.GetPath())
	}
}

func CheckMsg(t *testing.T, err error, expMsg string) {
	me, ok := err.(mgmterror.Formattable)
	if !ok {
		t.Fatalf("Error does not meet Formattable interface!")
		return
	}

	if me.GetMessage() != expMsg {
		t.Fatalf("Msg mismatch:\nExp:\t'%s'\nGot:\t'%s'\n",
			expMsg, me.GetMessage())
	}
}

func CheckInfo(t *testing.T, err error, expInfoVal string) {
	me, ok := err.(mgmterror.Formattable)
	if !ok {
		t.Fatalf("Error does not meet Formattable interface!")
		return
	}

	if expInfoVal == "" && len(me.GetInfo()) == 0 {
		// Nothing expected, nothing seen.  All clear.
		return
	}
	if expInfoVal == "" && len(me.GetInfo()) > 0 {
		t.Fatalf("Unexpected info value: '%s'\n", me.GetInfo()[0].Value)
		return
	}
	if expInfoVal != "" && len(me.GetInfo()) == 0 {
		t.Fatalf("No info value!\n")
		return
	}

	if me.GetInfo()[0].Value != expInfoVal {
		t.Fatalf("Info value mismatch:\nExp:\t'%s'\nGot:\t'%s'\n",
			expInfoVal, me.GetInfo()[0].Value)
	}
}

func checkInfoMatchesNonFatal(
	me mgmterror.Formattable,
	expInfo []*mgmterror.MgmtErrorInfoTag,
) bool {
	if len(expInfo) != len(me.GetInfo()) {
		return false
	}

	for _, expInfoTag := range expInfo {
		found := false
		for _, actInfoTag := range me.GetInfo() {
			if expInfoTag.XMLName.Space != actInfoTag.XMLName.Space {
				break
			}
			if expInfoTag.XMLName.Local != actInfoTag.XMLName.Local {
				break
			}
			if expInfoTag.Value != actInfoTag.Value {
				break
			}
			found = true
		}
		if !found {
			return false
		}
	}

	return true
}

type xpathType int

const (
	xpathMust xpathType = iota
	xpathWhen
)

type TestError struct {
	t         *testing.T
	path      string
	rawMsgs   []string
	cliMsgs   []string
	rpcMsgs   []string
	setMsg    string
	setSuffix string // used when set error doesn't end with 'is not valid'
}

func (te *TestError) CliErrorStrings() []string {

	pathSlice := getPathSlice(te.t, te.path, "generic error")

	retStr := []string{fmt.Sprintf("%s", mgmterror.ErrPath(pathSlice))}
	return append(retStr, te.cliMsgs...)
}

func (te *TestError) CommitCliErrorStrings() []string {
	return te.cliMsgs
}

func (te *TestError) RpcErrorStrings() []string {
	if len(te.rpcMsgs) == 0 {
		te.t.Fatalf("Test error message has no 'rpcMsgs'")
		return nil
	}

	pathSlice := getPathSlice(te.t, te.path, "rpc error")

	retStr := []string{fmt.Sprintf("%s", mgmterror.ErrPath(pathSlice))}
	return append(retStr, te.rpcMsgs...)
}

// Standard messages for set errors are:
//
// Configuration path: <path with last/only element in []> is not valid
//
// <setMsg>
//
// !!!DO NOT CHANGE THIS FORMAT WITHOUT CONSULTATION!!!
//
func (te *TestError) SetCliErrorStrings() []string {
	if te.setMsg == "" {
		te.t.Fatalf("Test error message has no 'setmsg'")
		return nil
	}

	pathSlice := getPathSlice(te.t, te.path, "generic error")
	if te.setMsg == noMsgPrinted {
		return []string{fmt.Sprintf("%s %s %s",
			configPathStr, mgmterror.ErrPath(pathSlice), isNotValidStr),
		}
	}
	if te.setSuffix == "" {
		return []string{fmt.Sprintf("%s %s %s",
			configPathStr, mgmterror.ErrPath(pathSlice), isNotValidStr),
			te.setMsg,
		}
	}

	return []string{fmt.Sprintf("%s %s %s",
		configPathStr, mgmterror.ErrPath(pathSlice), te.setSuffix),
		te.setMsg,
	}
}

func (te *TestError) RawErrorStrings() []string {

	retStr := []string{te.path}
	return append(retStr, te.rawMsgs...)
}

func (te *TestError) RawErrorStringsNoPath() []string {

	retStr := []string{}
	return append(retStr, te.rawMsgs...)
}

func NewAccessDeniedError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:    t,
		path: path,
		rawMsgs: []string{"Access to the requested protocol operation " +
			"or data model is denied because authorization failed."},
	}
}

func NewInterfaceMustExistError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{"Interface must exist"},
		cliMsgs: []string{"Interface must exist"},
	}
}

func NewInvalidNodeError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{"An unexpected element is present"},
		cliMsgs: []string{"Configuration path", "is not valid"},
		rpcMsgs: []string{"is not valid"},
		setMsg:  noMsgPrinted,
	}
}

// This is a bit of a special case, as this error is really meant for internal
// debug purposes if we get a prefix that doesn't match the command presented
// to cfgcli with it.  So, <setMsg> is a bit of a hack to let the test work.
func NewInvalidPrefixError(
	t *testing.T,
	path, prefix string,
) *TestError {
	return &TestError{
		t:         t,
		path:      path,
		rawMsgs:   []string{"An unexpected element is present"},
		cliMsgs:   []string{"Configuration path", "is not valid"},
		rpcMsgs:   []string{"is not valid"},
		setMsg:    "has invalid prefix",
		setSuffix: fmt.Sprintf("has invalid prefix '%s'", prefix),
	}
}

func NewInvalidNumElementsError(
	t *testing.T,
	path string,
	min, max int,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{fmt.Sprintf(wrongNumElementsFmtStr, min, max)},
		cliMsgs: []string{fmt.Sprintf(wrongNumElementsFmtStr, min, max)},
		setMsg:  noMsgPrinted,
	}
}

func NewInvalidRangeError(
	t *testing.T,
	path string,
	min, max int,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{fmt.Sprintf(wrongRangeFmtStr, min, max)},
		cliMsgs: []string{fmt.Sprintf(wrongRangeFmtStr, min, max)},
		setMsg:  fmt.Sprintf(wrongRangeFmtStr, min, max),
	}
}

func NewInvalidRangeCustomError(
	t *testing.T,
	path string,
	customErr string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{customErr},
		cliMsgs: []string{customErr},
		setMsg:  customErr,
	}
}

func NewInvalidPathError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{fmt.Sprintf("%s: %s", path, pathIsInvalidStr)},
		cliMsgs: []string{"TBD"},
		setMsg:  pathIsInvalidStr,
	}
}

func NewInvalidPatternError(
	t *testing.T,
	path string,
	pattern string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{fmt.Sprintf(mustMatchPatternFmtStr, pattern)},
		cliMsgs: []string{fmt.Sprintf(doesntMatchPatternFmtStr, pattern)},
		setMsg:  fmt.Sprintf(doesntMatchPatternFmtStr, pattern),
	}
}

func NewInvalidPatternCustomError(
	t *testing.T,
	path string,
	customErr string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{customErr},
		cliMsgs: []string{customErr},
		setMsg:  customErr,
	}
}

func NewInvalidTypeError(
	t *testing.T,
	path string,
	typ string,
) *TestError {
	pathSlice := getPathSlice(t, path, "invalid type")
	return &TestError{
		t:    t,
		path: path,
		rawMsgs: []string{fmt.Sprintf(
			wrongTypeFmtStr, pathSlice[len(pathSlice)-1], typ)},
		cliMsgs: []string{fmt.Sprintf(
			wrongTypeFmtStr, pathSlice[len(pathSlice)-1], typ)},
		setMsg: fmt.Sprintf(
			wrongTypeFmtStr, pathSlice[len(pathSlice)-1], typ),
	}
}

func NewInvalidLengthError(
	t *testing.T,
	path string,
	min, max int,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{fmt.Sprintf(wrongLengthFmtStr, min, max)},
		cliMsgs: []string{fmt.Sprintf(wrongLengthFmtStr, min, max)},
		setMsg:  fmt.Sprintf(wrongLengthFmtStr, min, max),
	}
}

func NewInvalidLengthCustomError(
	t *testing.T,
	path string,
	customErr string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{customErr},
		cliMsgs: []string{customErr},
		setMsg:  customErr,
	}
}

func NewLeafrefError(
	t *testing.T,
	path string,
	leafrefPath string,
) *TestError {
	return &TestError{
		t:    t,
		path: path,
		rawMsgs: []string{
			leafrefErrorStr, joinPathWithSpaces(
				getPathSlice(t, leafrefPath, "leafref"))},
		cliMsgs: []string{
			leafrefErrorStr, joinPathWithSpaces(
				getPathSlice(t, leafrefPath, "leafref"))},
	}
}

func NewMissingKeyError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{missingListKeyStr},
		cliMsgs: []string{missingListKeyStr},
		setMsg:  notYetTestedStr,
	}
}

func NewMissingMandatoryNodeError(
	t *testing.T,
	path string,
) *TestError {
	pathSlice := getPathSlice(t, path, "mandatory")
	if len(pathSlice) == 0 {
		t.Fatalf("Cannot have empty path for missing mandatory node error")
		return nil
	}
	return &TestError{
		t:    t,
		path: strings.Join(pathSlice[:len(pathSlice)-1], "/"),
		rawMsgs: []string{
			missingMandatoryStr + " " + pathSlice[len(pathSlice)-1]},
		cliMsgs: []string{
			missingMandatoryStr + " " + pathSlice[len(pathSlice)-1]},
	}
}

func NewNodeDoesntExistError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{nodeDoesntExistStr},
		cliMsgs: []string{nodeDoesntExistStr},
		setMsg:  nodeDoesntExistStr,
	}
}

func NewNodeExistsError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{nodeExistsStr},
		cliMsgs: []string{nodeExistsStr},
		setMsg:  nodeExistsStr,
	}
}

func NewNodeRequiresChildError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{notYetTestedStr},
		cliMsgs: []string{notYetTestedStr},
		setMsg:  nodeRequiresChildStr,
	}
}

func NewNodeRequiresValueError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{notYetTestedStr},
		cliMsgs: []string{notYetTestedStr},
		setMsg:  nodeRequiresValueStr,
	}
}

func NewNonUniquePathsError(
	t *testing.T,
	path string,
	keys []string,
	nonUniqueChildren []string,
) *TestError {
	return &TestError{
		t:    t,
		path: path,
		rawMsgs: []string{
			nonUniqueSetOfPathsStr,
			genChildPathsStr(nonUniqueChildren),
			nonUniqueSetOfKeysStr,
			genKeysStr(keys),
		},
		cliMsgs: []string{
			nonUniqueSetOfPathsStr,
			genChildPathsStr(nonUniqueChildren),
			nonUniqueSetOfKeysStr,
			genKeysStr(keys),
		},
	}
}

func NewPathAmbiguousError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:         t,
		path:      path,
		rawMsgs:   []string{"TBD"},
		cliMsgs:   []string{"TBD"},
		setSuffix: "is ambiguous",
		setMsg:    "Possible completions:",
	}
}

func NewSchemaMismatchError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{"Doesn't match schema"},
		cliMsgs: []string{"TBD"}, // TODO
	}
}

func NewSyntaxError(
	t *testing.T,
	path string,
	scriptErr string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{scriptErr},
		cliMsgs: []string{scriptErr},
	}
}

func NewUnknownElementError(
	t *testing.T,
	path string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{"Doesn't match schema"},
		cliMsgs: []string{"TBD"}, // TODO
		setMsg:  noMsgPrinted,
	}
}

func NewMustCustomError(
	t *testing.T,
	path,
	customError string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{path, customError},
		cliMsgs: []string{
			fmt.Sprintf("[%s]", joinPathWithSpaces(
				getPathSlice(t, path, "must custom"))),
			customError,
		},
	}
}

func NewWhenCustomError(
	t *testing.T,
	path,
	customError string,
) *TestError {
	return &TestError{
		t:       t,
		path:    path,
		rawMsgs: []string{path, customError},
		cliMsgs: []string{
			fmt.Sprintf("[%s]", joinPathWithSpaces(
				getPathSlice(t, path, "when custom"))),
			customError,
		},
	}
}

func NewMustDefaultError(
	t *testing.T,
	path,
	stmt string,
) *TestError {
	return &TestError{
		t:    t,
		path: path,
		rawMsgs: []string{
			path,
			fmt.Sprintf("'must' condition is false: '%s'", stmt),
		},
		cliMsgs: []string{
			fmt.Sprintf("[%s]", joinPathWithSpaces(
				getPathSlice(t, path, "must default"))),
			fmt.Sprintf("'must' condition is false: '%s'", stmt),
		},
	}
}

func NewWhenDefaultError(
	t *testing.T,
	path,
	stmt string,
) *TestError {
	return &TestError{
		t:    t,
		path: path,
		rawMsgs: []string{
			path,
			fmt.Sprintf("'when' condition is false: '%s'", stmt),
		},
		cliMsgs: []string{
			fmt.Sprintf("[%s]", joinPathWithSpaces(
				getPathSlice(t, path, "when default"))),
			fmt.Sprintf("'when' condition is false: '%s'", stmt),
		},
	}
}

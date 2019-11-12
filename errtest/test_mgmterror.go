// Copyright (c) 2017-2019, by AT&T Intellectual Property. All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

// Useful test functions for validating mgmterrors

package errtest

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

type ExpMgmtErrors struct {
	t          *testing.T
	nodes      []errNode
	endMessage string
}

type errNode struct {
	path string
	errs []*errDesc
}

type errDesc struct {
	typ  errType
	data []string
}

func NewErrDesc(typ errType, data ...string) *errDesc {
	return &errDesc{typ: typ, data: data}
}

type errType int

const (
	CustomMustError errType = iota
	DefaultMustError
	IntfMustExist
	InvalidPath
	LeafrefMissing
	MissingMandatory
	NodeDoesntExist
	NotUnique
)

// In order to insulate unit tests from changes to output format, we need
// to try to test content alone.  The aim of the functions here is to map
// between different formats, so that the content being verified is
// converted to the expected format before being validated.  Thus if we
// change format, we only need to change the functions here and ALL tests
// that are testing content not format will pass again.
const (
	commitNonFatalErrStr     = "Commit succeeded (non-fatal failures detected)."
	configPathStr            = "Configuration path:"
	doesntMatchPatternFmtStr = "Does not match pattern %s"
	isNotValidStr            = "is not valid"
	leafrefErrorStr          = "The following path must exist:"
	missingListKeyStr        = "List entry is missing key"
	missingMandatoryStr      = "Missing mandatory node"
	mustMatchPatternFmtStr   = "Must match %s"
	mustOrWhenDefaultFmtStr  = "'%s' condition is false: '%s'"
	// TODO remove this and when equivalent as unnecessary
	MustStmt               = "must"
	nodeExistsStr          = "Node exists"
	nodeDoesntExistStr     = "Node does not exist"
	nodeRequiresChildStr   = "Node requires a child"
	nodeRequiresValueStr   = "Node requires a value"
	noMsgPrinted           = "IGNORE"
	nonUniqueSetOfKeysStr  = "but is defined in the following set of keys:"
	nonUniqueSetOfPathsStr = "The following set of paths must be unique:"
	notYetTestedStr        = "This option hasn't been tested."
	pathIsInvalidStr       = "Path is invalid"
	TestCommitFailStr      = "\nCommit failed!\n"
	TestValidateFailStr    = "\nValidate failed!\n"
	WarningsGeneratedStr   = "Warnings were generated when applying " +
		"the configuration:"
	WhenStmt               = "when"
	wrongLengthFmtStr      = "Must have length between %d and %d"
	wrongNumElementsFmtStr = "Invalid number of nodes: " +
		"must be in the range %d to %d"
	wrongRangeFmtStr = "Must have value between %d and %d"
	wrongTypeFmtStr  = "'%s' is not %s"
)

func errpath(path []string) string {
	if len(path) < 2 {
		return fmt.Sprintf("%s", path)
	}
	path, val := path[:len(path)-1], path[len(path)-1]
	return fmt.Sprintf("%s [%s]", strings.Join(path, " "), val)
}

func getPathSlice(t *testing.T, path, desc string) []string {
	if len(path) == 0 || path == "/" {
		t.Fatalf("Cannot generate %s error with '' or '/' as path!",
			desc)
		return nil
	}
	if path[0] == '/' {
		path = path[1:]
	}
	return strings.Split(path, "/")
}

// joinPathWithSpaces - create space-separated path, deal with special chars
//
// Creates space-separate path from slice, and deals with the likes of '%2F'
// for a '/' within an IP-address-with-mask entry.
func joinPathWithSpaces(pathSlice []string) string {
	retStr := strings.Join(pathSlice, " ")
	return strings.Replace(retStr, "%2F", "/", -1)
}

// InvalidRangeErrorStrings - generate string(s) for invalid range error
//
// Expect to see:
//
// '[<path-to-node>'
// 'Must have value between <min> and <max>'
//
func InvalidRangeErrorStrings(
	t *testing.T,
	path string,
	min, max int,
) []string {

	pathSlice := getPathSlice(t, path, "invalid range")
	return []string{
		fmt.Sprintf("[%s]", strings.Join(pathSlice, " ")),
		fmt.Sprintf(wrongRangeFmtStr, min, max),
	}
}

// InvalidTypeErrorStrings - generate string(s) for invalid type error
//
// Expect to see:
//
// '[<path-to-node>'
// ''<last-elem-in-path>' is not <type>'
//
func InvalidTypeErrorStrings(
	t *testing.T,
	path string,
	typ string,
) []string {
	pathSlice := getPathSlice(t, path, "invalid type")
	return []string{
		fmt.Sprintf("[%s]", strings.Join(pathSlice, " ")),
		fmt.Sprintf(wrongTypeFmtStr, pathSlice[len(pathSlice)-1], typ),
	}
}

func NonFatalCommitErrorStrings(
	t *testing.T,
	path string,
) []string {

	pathSlice := getPathSlice(t, path, "non-fatal commit")
	return []string{
		fmt.Sprintf("[%s]", strings.Join(pathSlice, " ")),
		fmt.Sprintf(commitNonFatalErrStr),
	}
}

func genChildPathsStr(children []string) string {
	var retStr string
	first := true
	for _, child := range children {
		if !first {
			retStr += ", "
		}
		first = false
		retStr += fmt.Sprintf("[%s]", strings.Replace(child, "/", " ", -1))
	}
	return retStr
}

// Keys are in slashed-path format, and include list name.  Need to remove
// that in current error message output format.
func genKeysStr(keys []string) string {
	retStr := "["
	first := true
	for _, key := range keys {
		if !first {
			retStr += " "
		}
		first = false
		retStr += strings.Split(key, "/")[1]
	}
	return retStr
}

// Errors constructed using these functions are being tested for format as well
// as content,  ie we *do* care about format.  Previous function are all
// about content.
func NewExpectedFormattedErrors(t *testing.T) *ExpMgmtErrors {
	return &ExpMgmtErrors{t: t, nodes: make([]errNode, 1)}
}

func (eme *ExpMgmtErrors) AddNode(
	path string,
	errs ...*errDesc,
) *ExpMgmtErrors {
	eN := errNode{
		path: strings.Join(getPathSlice(eme.t, path, "add node"), " "),
		errs: errs}
	eme.nodes = append(eme.nodes, eN)
	return eme
}

func (eme *ExpMgmtErrors) AddEndMessage(msg string) *ExpMgmtErrors {
	eme.endMessage = msg
	return eme
}

func genCustomMustErrMsg(t *testing.T, data []string) string {
	if len(data) != 1 {
		t.Fatalf("Custom must error must have single error.\n")
		return ""
	}
	return fmt.Sprintf("%s\n\n", data[0])
}

func genDefaultMustErrMsg(t *testing.T, data []string) string {
	if len(data) != 1 {
		t.Fatalf("Default must error must have single error.\n")
		return ""
	}
	return fmt.Sprintf("'must' condition is false: '%s'\n\n", data[0])
}

func genIntfMustExistErrMsg(t *testing.T) string {
	return fmt.Sprintf("Interface must exist.\n\n")
}

func genInvalidPathErrMsg(t *testing.T, data []string) string {
	if len(data) != 1 {
		t.Fatalf("Invalid path error must have single data entry.\n")
		return ""
	}
	path := errpath(strings.Split(data[0], " "))
	return fmt.Sprintf("Configuration path: %s is not valid", path)
}

func genLeafrefErrMsg(t *testing.T, data []string) string {
	retStr := "The following path must exist:\n"
	if len(data) != 1 {
		t.Fatalf("Leafref error must have single path.\n")
		return ""
	}
	retStr += fmt.Sprintf("  [%s]\n\n", data[0])

	return retStr
}

func genMissingMandatoryErrMsg(t *testing.T, data []string) string {
	if len(data) != 1 {
		t.Fatalf("Missing mandatory error must have single node.\n")
		return ""
	}
	return fmt.Sprintf("Missing mandatory node %s\n\n", data[0])
}

func genNodeDoesntExistErrMsg(t *testing.T) string {
	return fmt.Sprintf("Node does not exist\n\n")
}

func genNotUniqueErrMsg(t *testing.T, data []string) string {

	if len(data) < 2 {
		t.Fatalf("Not unique error must have at least 2 strings.\n")
		return ""
	}
	retStr := "The following path must be unique:\n\n"
	retStr += fmt.Sprintf("  [%s]\n\n", data[0])
	retStr += "but is defined in the following set of keys:\n\n"
	for _, key := range data[1:] {
		retStr += fmt.Sprintf("  [%s]\n", key)
	}
	retStr += "\n"
	return retStr
}

func (eme *ExpMgmtErrors) String() string {
	var retStr string

	for _, node := range eme.nodes {
		for _, mgmtErr := range node.errs {
			// 'Unwrapped' errors
			switch mgmtErr.typ {
			case InvalidPath:
				retStr += genInvalidPathErrMsg(eme.t, mgmtErr.data)
				continue
			}
			retStr += fmt.Sprintf("[%s]\n\n", node.path)
			switch mgmtErr.typ {
			case CustomMustError:
				retStr += genCustomMustErrMsg(eme.t, mgmtErr.data)
			case DefaultMustError:
				retStr += genDefaultMustErrMsg(eme.t, mgmtErr.data)
			case IntfMustExist:
				retStr += genIntfMustExistErrMsg(eme.t)
			case LeafrefMissing:
				retStr += genLeafrefErrMsg(eme.t, mgmtErr.data)
			case MissingMandatory:
				retStr += genMissingMandatoryErrMsg(eme.t, mgmtErr.data)
			case NodeDoesntExist:
				retStr += genNodeDoesntExistErrMsg(eme.t)
			case NotUnique:
				retStr += genNotUniqueErrMsg(eme.t, mgmtErr.data)
			default:
				eme.t.Fatalf("Error type (%d) not supported\n", mgmtErr.typ)
				return ""
			}
			retStr += fmt.Sprintf("[[%s]] failed.\n", node.path)
		}
	}
	retStr += eme.endMessage
	return retStr
}

func (eme *ExpMgmtErrors) Matches(actual error) {
	if actual == nil {
		eme.t.Fatalf("Unexpected success")
	}

	CheckStringDivergence(eme.t, eme.String(), actual.Error())
}

// Very useful when debugging outputs that don't match up.
func CheckStringDivergence(t *testing.T, expOut, actOut string) {
	if expOut == actOut {
		return
	}

	var expOutCopy = expOut
	var act bytes.Buffer
	var charsToDump = 10
	var expCharsToDump = 10
	var actCharsLeft, expCharsLeft int
	for index, char := range actOut {
		if len(expOutCopy) > 0 {
			if char == rune(expOutCopy[0]) {
				act.WriteByte(byte(char))
			} else {
				act.WriteString("###") // Mark point of divergence.
				expCharsLeft = len(expOutCopy)
				actCharsLeft = len(actOut) - index
				if expCharsLeft < charsToDump {
					expCharsToDump = expCharsLeft
				}
				if actCharsLeft < charsToDump {
					charsToDump = actCharsLeft
				}
				act.WriteString(actOut[index : index+charsToDump])
				break
			}
		} else {
			t.Logf("Expected output terminates early.\n")
			t.Fatalf("Exp:\n%s\nGot extra:\n%s\n",
				expOut[:index], act.String()[index:])
		}
		expOutCopy = expOutCopy[1:]
	}

	// When expOut is longer than actOut, need to update the expCharsToDump
	if len(expOutCopy) < charsToDump {
		expCharsToDump = len(expOutCopy)
	}

	// Useful to print whole output first for reference (useful when debugging
	// when you don't want to have to construct the expected output up front).
	t.Logf("Actual output:\n%s\n--- ENDS ---\n", actOut)

	// After that we then print up to the point of divergence so it's easy to
	// work out what went wrong ...
	t.Fatalf("Unexpected output.\nGot:\n%s\nExp at ###:\n'%s ...'\n",
		act.String(), expOutCopy[:expCharsToDump])
}

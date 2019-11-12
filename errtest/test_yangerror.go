// Copyright (c) 2017,2019, by AT&T Intellectual Property. All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

// Useful test functions for validating mgmterrors

package errtest

import (
	"fmt"
	"testing"
)

type NodeLimits struct {
	Min string
	Max string
}

const (
	Bltin      = "builtin"
	BoolType   = "boolean"
	Dec64Type  = "decimal64"
	EmptyType  = "empty"
	Int8Type   = "int8"
	StringType = "string"
	Uint8Type  = "uint8"
)

func formatLimits(l NodeLimits, prefix, equality, suffix string) string {
	if l.Min == l.Max {
		return fmt.Sprintf("%s%s %s%s", prefix, equality, l.Min, suffix)
	}
	return fmt.Sprintf("%sbetween %s and %s%s", prefix, l.Min, l.Max, suffix)
}

func rangeValueString(t *testing.T, ranges []NodeLimits) string {
	if len(ranges) == 0 {
		t.Fatalf("Must have at least one range")
		return ""
	}
	if len(ranges) == 1 {
		return fmt.Sprintf("Must have value %s",
			formatLimits(ranges[0], "", "equal to", ""))
	}
	if len(ranges) > 1 {
		retStr := "Must have one of the following values: "
		first := true
		for _, r := range ranges {
			if !first {
				// TODO - why not picking up trailing ', '?
				retStr += ", "
			}
			first = false
			retStr += formatLimits(r, "", "equal to", "")
		}
		return retStr
	}
	return ""
}

// YangInvalidDefaultValueErrorStrings - string(s) for invalid default value
//
// Expect to see:
//
// TBD
//
func YangInvalidDefaultValueErrorStrings(
	t *testing.T,
	yangNS, yangType,
	invalidDefault string,
	ranges []NodeLimits,
) []string {
	return []string{
		fmt.Sprintf("type %s: Invalid default '%s' for {%s %s}",
			yangType, invalidDefault, yangNS, yangType),
		rangeValueString(t, ranges),
	}
}

func lengthValueString(t *testing.T, ranges []NodeLimits) string {
	if len(ranges) == 0 {
		t.Fatalf("Must have at least one range")
		return ""
	}
	if len(ranges) == 1 {
		return fmt.Sprintf("Must have length %s",
			formatLimits(ranges[0], "", "of", " characters"))
	}
	if len(ranges) > 1 {
		retStr := "Must be one of the following: "
		first := true
		for _, r := range ranges {
			if !first {
				retStr += ", "
			}
			first = false
			retStr += formatLimits(r, "have length ", "of", " characters")
		}
		return retStr
	}
	return ""
}

// YangInvalidDefaultLengthErrorStrings - string(s) for invalid default length
//
// Expect to see:
//
// TBD
//
func YangInvalidDefaultLengthErrorStrings(
	t *testing.T,
	yangNS, yangType,
	invalidDefault string,
	lengths []NodeLimits,
) []string {
	return []string{
		fmt.Sprintf("type %s: Invalid default '%s' for {%s %s}",
			yangType, invalidDefault, yangNS, yangType),
		lengthValueString(t, lengths),
	}
}

// YangInvalidDefaultTypeErrorStrings - string(s) for invalid default type
//
// Expect to see:
//
// TBD
//
func YangInvalidDefaultTypeErrorStrings(
	t *testing.T,
	yangNS, yangType,
	invalidDefault string,
) []string {
	return []string{
		fmt.Sprintf("type %s: Invalid default '%s' for {%s %s}",
			yangType, invalidDefault, yangNS, yangType),
		fmt.Sprintf("'%s' is not an %s", invalidDefault, yangType),
	}
}

func enumString(validValues []string) string {
	if len(validValues) == 1 {
		return fmt.Sprintf("Must have value %s", validValues[0])
	}
	var retStr = "Must have one of the following values: "
	first := true
	for _, value := range validValues {
		if !first {
			retStr += ", "
		}
		first = false
		retStr += value
	}
	return retStr
}

// YangInvalidDefaultEnumOrBoolErrorStrings - string(s) for invalid bool/enum
//
// Expect to see:
//
// TBD
//
func YangInvalidDefaultEnumOrBoolErrorStrings(
	t *testing.T,
	yangNS, yangType,
	invalidDefault string,
	validValues []string,
) []string {
	return []string{
		fmt.Sprintf("type %s: Invalid default '%s' for {%s %s}",
			yangType, invalidDefault, yangNS, yangType),
		enumString(validValues),
	}
}

// YangInvalidDefaultEmptyErrorStrings - string(s) for invalid empty leaf
//
// Expect to see:
//
// TBD
//
func YangInvalidDefaultEmptyErrorStrings(
	t *testing.T,
	yangNS, yangType,
	invalidDefault string,
) []string {
	return []string{
		fmt.Sprintf("type %s: Invalid default '%s' for {%s %s}",
			yangType, invalidDefault, yangNS, yangType),
	}
}

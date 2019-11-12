// Copyright (c) 2017,2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package mgmterror

import (
	"fmt"
)

func ExampleExecError() {
	path := []string{"usr", "bin", "app"}
	err := NewExecError(path, "core dumped")
	fmt.Println(err.Error())

	//Output:
	// Error: /usr/bin/app: core dumped
}

func ExamplePathAmbiguousError() {
	path := []string{"s"}
	matches := map[string]string{
		"system":   "System parameters",
		"service":  "Services",
		"security": "Security",
	}
	err := NewPathAmbiguousError(path, matches)
	fmt.Println(err.Error())

	//Output:
	// Error: /s: Ambiguous path, could be one of: security, service, system
}

func ExamplePathAmbiguousError_command() {
	matches := map[string]string{
		"show": "Show the configuration (default values may be suppressed)",
		"set":  "Set the value of a parameter or create a new element",
		"save": "Save configuration to a file",
	}
	err := NewPathAmbiguousError([]string{}, matches)
	fmt.Println(err.Error())

	//Output:
	// Error: Ambiguous command, could be one of: save, set, show
}

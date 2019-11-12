// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2017 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

/*

	Package mgmterror provides the standard NETCONF and YANG errors
	as defined in RFC6241 and RFC6020 respectively.

	While RFC6241 section 4.3 specifies 7 elements of an
	rpc-error, not all of them are mandatory for a specific
	error. For example, Path (error-path), will not be present if
	no appropriate payload element or datastore node can be
	associated with a particular error condition.

	Since the marshal and unmarshal interfaces require the
	MgmtError structure's members be public, they can be set as
	needed after an error's constructor is called.

	The specifications allow the Info (error-info) field to
	include additional elements to provide extended and/or
	implementation-specific debugging information. For example, an
	error from a low level library can add additional information
	about the error so that a higher layer can use that
	information.

*/
package mgmterror

// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated (@generated) by entc, DO NOT EDIT.

package user

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldURL holds the string denoting the url vertex property in the database.
	FieldURL = "url"
	// FieldRaw holds the string denoting the raw vertex property in the database.
	FieldRaw = "raw"
	// FieldDirs holds the string denoting the dirs vertex property in the database.
	FieldDirs = "dirs"
	// FieldInts holds the string denoting the ints vertex property in the database.
	FieldInts = "ints"
	// FieldFloats holds the string denoting the floats vertex property in the database.
	FieldFloats = "floats"
	// FieldStrings holds the string denoting the strings vertex property in the database.
	FieldStrings = "strings"

	// Table holds the table name of the user in the database.
	Table = "users"
)

// Columns holds all SQL columns are user fields.
var Columns = []string{
	FieldID,
	FieldURL,
	FieldRaw,
	FieldDirs,
	FieldInts,
	FieldFloats,
	FieldStrings,
}

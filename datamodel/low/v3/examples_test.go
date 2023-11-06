// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package v3

import (
	"fmt"
	"os"

	"github.com/speakeasy-api/libopenapi/datamodel"
)

// How to create a low-level OpenAPI 3+ Document from an OpenAPI specification
func Example_createLowLevelOpenAPIDocument() {
	// How to create a low-level OpenAPI 3 Document

	// load petstore into bytes
	petstoreBytes, _ := os.ReadFile("../../../test_specs/petstorev3.json")

	// read in specification
	info, _ := datamodel.ExtractSpecInfo(petstoreBytes)

	// build low-level document model
	document, errors := CreateDocument(info)

	// if something went wrong, a slice of errors is returned
	if len(errors) > 0 {
		for i := range errors {
			fmt.Printf("error: %s\n", errors[i].Error())
		}
		panic("cannot build document")
	}

	// print out email address from the info > contact object.
	fmt.Print(document.Info.Value.Contact.Value.Email.Value)
	// Output: apiteam@swagger.io
}

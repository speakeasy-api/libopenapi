// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"github.com/speakeasy-api/libopenapi/datamodel/low"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestContact_Hash(t *testing.T) {

	left := `url: https://pb33f.io
description: the ranch
email: buckaroo@pb33f.io`

	right := `url: https://pb33f.io
description: the ranch
email: buckaroo@pb33f.io`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc Contact
	var rDoc Contact
	_ = low.BuildModel(lNode.Content[0], &lDoc)
	_ = low.BuildModel(rNode.Content[0], &rDoc)

	assert.Equal(t, lDoc.Hash(), rDoc.Hash())
}

// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package base

import (
	"context"
	"fmt"
	"testing"

	lowmodel "github.com/pb33f/libopenapi/datamodel/low"
	lowbase "github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestNewInfo(t *testing.T) {
	var cNode yaml.Node

	yml := `title: chicken
summary: a chicken nugget
description: nugget
termsOfService: chicken soup
contact:
  name: buckaroo
license:
  name: pb33f
  url: https://pb33f.io
version: 99.99
x-cli-name: chicken cli`

	_ = yaml.Unmarshal([]byte(yml), &cNode)

	var lowInfo lowbase.Info
	_ = lowmodel.BuildModel(cNode.Content[0], &lowInfo)
	_, _ = lowInfo.Build(context.Background(), nil, cNode.Content[0], nil)

	highInfo := NewInfo(&lowInfo)

	var xCliName string
	_ = highInfo.Extensions.GetOrZero("x-cli-name").Decode(&xCliName)

	assert.Equal(t, "chicken", highInfo.Title)
	assert.Equal(t, "a chicken nugget", highInfo.Summary)
	assert.Equal(t, "nugget", highInfo.Description)
	assert.Equal(t, "chicken soup", highInfo.TermsOfService)
	assert.Equal(t, "buckaroo", highInfo.Contact.Name)
	assert.Equal(t, "pb33f", highInfo.License.Name)
	assert.Equal(t, "https://pb33f.io", highInfo.License.URL)
	assert.Equal(t, "99.99", highInfo.Version)
	assert.Equal(t, "chicken cli", xCliName)

	wentLow := highInfo.GoLow()
	assert.Equal(t, 10, wentLow.Version.ValueNode.Line)

	wentLower := highInfo.License.GoLow()
	assert.Equal(t, 9, wentLower.URL.ValueNode.Line)
}

func ExampleNewInfo() {
	// create an example info object (including contact and license)
	// this can be either JSON or YAML.
	yml := `title: some spec by some company
summary: this is a summary
description: this is a specification, for an API, by a company.
termsOfService: https://pb33f.io/tos
contact:
  name: buckaroo
license:
  name: MIT
  url: https://opensource.org/licenses/MIT
version: 1.2.3`

	// unmarshal raw bytes
	var node yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &node)

	// build out the low-level model
	var lowInfo lowbase.Info
	_ = lowmodel.BuildModel(&node, &lowInfo)
	_, _ = lowInfo.Build(context.Background(), nil, node.Content[0], nil)

	// build the high level model
	highInfo := NewInfo(&lowInfo)

	// print out the contact name.
	fmt.Print(highInfo.Contact.Name)
	// Output: buckaroo
}

func ExampleNewLicense() {
	// create an example license object
	// this can be either JSON or YAML.
	yml := `name: MIT
url: https://opensource.org/licenses/MIT`

	// unmarshal raw bytes
	var node yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &node)

	// build out the low-level model
	var lowLicense lowbase.License
	_ = lowmodel.BuildModel(node.Content[0], &lowLicense)
	_, _ = lowLicense.Build(context.Background(), nil, node.Content[0], nil)

	// build the high level model
	highLicense := NewLicense(&lowLicense)

	// print out the contact name.
	fmt.Print(highLicense.Name)
	// Output: MIT
}

func TestInfo_Render(t *testing.T) {
	ext := orderedmap.New[string, *yaml.Node]()
	ext.Set("x-pizza", utils.CreateStringNode("pepperoni"))
	ext.Set("x-cake", utils.CreateYamlNode(&License{
		Name: "someone",
		URL:  "nowhere",
	}))
	highI := &Info{
		Title:          "hey",
		Description:    "there you",
		TermsOfService: "have you got any money",
		Contact: &Contact{
			Name:  "buckaroo",
			Email: "buckaroo@pb33f.io",
		},
		License: &License{
			Name: "MIT",
			URL:  "https://opensource.org/licenses/MIT",
		},
		Version:    "1.2.3",
		Extensions: ext,
	}

	dat, _ := highI.Render()

	// unmarshal yaml into a *yaml.Node instance
	var cNode yaml.Node
	_ = yaml.Unmarshal(dat, &cNode)

	// build low
	var lowInfo lowbase.Info
	_ = lowmodel.BuildModel(cNode.Content[0], &lowInfo)
	_, _ = lowInfo.Build(context.Background(), nil, cNode.Content[0], nil)

	// build high
	highInfo := NewInfo(&lowInfo)

	var xPizza string
	_ = highInfo.Extensions.GetOrZero("x-pizza").Decode(&xPizza)

	assert.Equal(t, "hey", highInfo.Title)
	assert.Equal(t, "there you", highInfo.Description)
	assert.Equal(t, "have you got any money", highInfo.TermsOfService)
	assert.Equal(t, "buckaroo", highInfo.Contact.Name)
	assert.Equal(t, "buckaroo@pb33f.io", highInfo.Contact.Email)
	assert.Equal(t, "MIT", highInfo.License.Name)
	assert.Equal(t, "https://opensource.org/licenses/MIT", highInfo.License.URL)
	assert.Equal(t, "1.2.3", highInfo.Version)
	assert.Equal(t, "pepperoni", xPizza)
	assert.NotNil(t, highInfo.GoLowUntyped())
}

func TestInfo_RenderOrder(t *testing.T) {
	yml := `title: hey
description: there you
termsOfService: have you got any money
contact:
    name: buckaroo
    email: buckaroo@pb33f.io
license:
    name: MIT
    url: https://opensource.org/licenses/MIT
version: 1.2.3
x-pizza: pepperoni
x-cake:
    name: someone
    url: nowhere`

	// unmarshal yaml into a *yaml.Node instance
	var cNode yaml.Node
	_ = yaml.Unmarshal([]byte(yml), &cNode)

	// build low
	var lowInfo lowbase.Info
	_ = lowmodel.BuildModel(cNode.Content[0], &lowInfo)
	_, _ = lowInfo.Build(context.Background(), nil, cNode.Content[0], nil)

	// build high
	highInfo := NewInfo(&lowInfo)

	var xPizza string
	_ = highInfo.Extensions.GetOrZero("x-pizza").Decode(&xPizza)

	assert.Equal(t, "hey", highInfo.Title)
	assert.Equal(t, "there you", highInfo.Description)
	assert.Equal(t, "have you got any money", highInfo.TermsOfService)
	assert.Equal(t, "buckaroo", highInfo.Contact.Name)
	assert.Equal(t, "buckaroo@pb33f.io", highInfo.Contact.Email)
	assert.Equal(t, "MIT", highInfo.License.Name)
	assert.Equal(t, "https://opensource.org/licenses/MIT", highInfo.License.URL)
	assert.Equal(t, "1.2.3", highInfo.Version)
	assert.Equal(t, "pepperoni", xPizza)

	// marshal high back to yaml, should be the same as the original, in same order.
	bytes, _ := highInfo.Render()
	assert.Len(t, bytes, 275)
}

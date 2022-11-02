// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"github.com/pb33f/libopenapi/datamodel/low"
	v2 "github.com/pb33f/libopenapi/datamodel/low/v2"
	v3 "github.com/pb33f/libopenapi/datamodel/low/v3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestCompareOperations_V2(t *testing.T) {

	left := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	right := left

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Nil(t, extChanges)
}

func TestCompareOperations_V2_ModifyParam(t *testing.T) {

	left := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	right := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: fridge`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
}

func TestCompareOperations_V2_AddParamProperty(t *testing.T) {

	left := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam`

	right := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
	assert.Equal(t, PropertyAdded, extChanges.Changes[0].ChangeType)
	assert.Equal(t, "parameters", extChanges.Changes[0].Property)

}

func TestCompareOperations_V2_RemoveParamProperty(t *testing.T) {

	left := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam`

	right := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&rDoc, &lDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
	assert.Equal(t, PropertyRemoved, extChanges.Changes[0].ChangeType)
	assert.Equal(t, "parameters", extChanges.Changes[0].Property)

}

func TestCompareOperations_V2_AddParam(t *testing.T) {

	left := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	right := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer
  - name: jummy
    in: oven`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
	assert.Equal(t, ObjectAdded, extChanges.Changes[0].ChangeType)
}

func TestCompareOperations_V2_RemoveParam(t *testing.T) {

	left := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	right := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer
  - name: jummy
    in: oven`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&rDoc, &lDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
	assert.Equal(t, ObjectRemoved, extChanges.Changes[0].ChangeType)
}

func TestCompareOperations_V2_ModifyTag(t *testing.T) {

	left := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	right := `tags:
  - one
  - twenty
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 2, extChanges.TotalChanges())
	assert.Equal(t, 0, extChanges.TotalBreakingChanges())
	assert.Equal(t, PropertyRemoved, extChanges.Changes[0].ChangeType)
	assert.Equal(t, PropertyAdded, extChanges.Changes[1].ChangeType)
}

func TestCompareOperations_V2_AddTag(t *testing.T) {

	left := `tags:
  - one
  - two
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	right := `tags:
  - one
  - two
  - three
summary: hello
description: hello there my pal
operationId: mintyFresh
consumes:
  - pizza
  - cake
produces: 
  - toast
  - jam
parameters:
  - name: jimmy
  - name: jammy
    in: freezer`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 0, extChanges.TotalBreakingChanges())
	assert.Equal(t, PropertyAdded, extChanges.Changes[0].ChangeType)
	assert.Equal(t, v3.TagsLabel, extChanges.Changes[0].Property)
}

func TestCompareOperations_V2_Modify_ProducesConsumesSchemes(t *testing.T) {

	left := `produces:
  - electricity
consumes:
  - oil
schemes:
  - burning`

	right := `produces:
  - heat
consumes:
  - wind
schemes:
  - blowing`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 6, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
}

func TestCompareOperations_V2_Modify_ExtDocs_Responses(t *testing.T) {

	left := `externalDocs:
  url: https://pb33f.io/old
responses:
  200:
    description: OK matey!
`

	right := `externalDocs:
  url: https://pb33f.io/new
responses:
  200:
    description: OK me matey!`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 2, extChanges.TotalChanges())
	assert.Equal(t, 0, extChanges.TotalBreakingChanges())
	assert.Equal(t, Modified, extChanges.ResponsesChanges.ResponseChanges["200"].Changes[0].ChangeType)
	assert.Equal(t, Modified, extChanges.ExternalDocChanges.Changes[0].ChangeType)

}

func TestCompareOperations_V2_Add_ExtDocs_Responses(t *testing.T) {

	left := `operationId: nuggets`

	right := `operationId: nuggets
externalDocs:
  url: https://pb33f.io/new
responses:
  200:
    description: OK me matey!`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 2, extChanges.TotalChanges())
	assert.Equal(t, 0, extChanges.TotalBreakingChanges())
	assert.Equal(t, PropertyAdded, extChanges.Changes[0].ChangeType)
	assert.Equal(t, PropertyAdded, extChanges.Changes[1].ChangeType)
}

func TestCompareOperations_V2_Remove_ExtDocs_Responses(t *testing.T) {

	left := `operationId: nuggets`

	right := `operationId: nuggets
externalDocs:
  url: https://pb33f.io/new
responses:
  200:
    description: OK me matey!`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&rDoc, &lDoc)
	assert.Equal(t, 2, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
	assert.Equal(t, PropertyRemoved, extChanges.Changes[0].ChangeType)
	assert.Equal(t, PropertyRemoved, extChanges.Changes[1].ChangeType)
}

func TestCompareOperations_V2_AddSecurityReq_Role(t *testing.T) {

	left := `operationId: nuggets
security:
  - things:
    - stuff
    - junk`

	right := `operationId: nuggets
security:
  - things:
    - stuff
    - junk
    - crap`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 0, extChanges.TotalBreakingChanges())
	assert.Equal(t, ObjectAdded, extChanges.SecurityRequirementChanges[0].Changes[0].ChangeType)
	assert.Equal(t, "crap", extChanges.SecurityRequirementChanges[0].Changes[0].New)

}

func TestCompareOperations_V2_RemoveSecurityReq_Role(t *testing.T) {

	left := `operationId: nuggets
security:
  - things:
    - stuff
    - junk`

	right := `operationId: nuggets
security:
  - things:
    - stuff
    - junk
    - crap`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&rDoc, &lDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
	assert.Equal(t, ObjectRemoved, extChanges.SecurityRequirementChanges[0].Changes[0].ChangeType)
	assert.Equal(t, "crap", extChanges.SecurityRequirementChanges[0].Changes[0].Original)

}

func TestCompareOperations_V2_AddSecurityRequirement(t *testing.T) {

	left := `operationId: nuggets
security:
  - things:
    - stuff
    - junk`

	right := `operationId: nuggets
security:
  - things:
    - stuff
    - junk
  - thongs:
    - small
    - smelly`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&lDoc, &rDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 0, extChanges.TotalBreakingChanges())
	assert.Equal(t, ObjectAdded, extChanges.Changes[0].ChangeType)
	assert.Equal(t, "thongs", extChanges.Changes[0].New)

}

func TestCompareOperations_V2_RemoveSecurityRequirement(t *testing.T) {

	left := `operationId: nuggets
security:
  - things:
    - stuff
    - junk`

	right := `operationId: nuggets
security:
  - things:
    - stuff
    - junk
  - thongs:
    - small
    - smelly`

	var lNode, rNode yaml.Node
	_ = yaml.Unmarshal([]byte(left), &lNode)
	_ = yaml.Unmarshal([]byte(right), &rNode)

	// create low level objects
	var lDoc v2.Operation
	var rDoc v2.Operation
	_ = low.BuildModel(&lNode, &lDoc)
	_ = low.BuildModel(&rNode, &rDoc)
	_ = lDoc.Build(lNode.Content[0], nil)
	_ = rDoc.Build(rNode.Content[0], nil)

	// compare.
	extChanges := CompareOperations(&rDoc, &lDoc)
	assert.Equal(t, 1, extChanges.TotalChanges())
	assert.Equal(t, 1, extChanges.TotalBreakingChanges())
	assert.Equal(t, ObjectRemoved, extChanges.Changes[0].ChangeType)
	assert.Equal(t, "thongs", extChanges.Changes[0].Original)

}

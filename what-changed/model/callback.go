// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/v3"
)

type CallbackChanges struct {
	PropertyChanges
	ExpressionChanges map[string]*PathItemChanges `json:"expressions,omitempty" yaml:"expressions,omitempty"`
	ExtensionChanges  *ExtensionChanges           `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

func (c *CallbackChanges) TotalChanges() int {
	d := c.PropertyChanges.TotalChanges()
	for k := range c.ExpressionChanges {
		d += c.ExpressionChanges[k].TotalChanges()
	}
	if c.ExtensionChanges != nil {
		d += c.ExtensionChanges.TotalChanges()
	}
	return d
}

func (c *CallbackChanges) TotalBreakingChanges() int {
	d := c.PropertyChanges.TotalBreakingChanges()
	for k := range c.ExpressionChanges {
		d += c.ExpressionChanges[k].TotalBreakingChanges()
	}
	if c.ExtensionChanges != nil {
		d += c.ExtensionChanges.TotalBreakingChanges()
	}
	return d
}

func CompareCallback(l, r *v3.Callback) *CallbackChanges {

	cc := new(CallbackChanges)
	var changes []*Change

	lHashes := make(map[string]string)
	rHashes := make(map[string]string)

	lValues := make(map[string]low.ValueReference[*v3.PathItem])
	rValues := make(map[string]low.ValueReference[*v3.PathItem])

	for k := range l.Expression.Value {
		lHashes[k.Value] = low.GenerateHashString(l.Expression.Value[k].Value)
		lValues[k.Value] = l.Expression.Value[k]
	}

	for k := range r.Expression.Value {
		rHashes[k.Value] = low.GenerateHashString(r.Expression.Value[k].Value)
		rValues[k.Value] = r.Expression.Value[k]
	}

	expChanges := make(map[string]*PathItemChanges)

	// check left path item hashes
	for k := range lHashes {
		rhash := rHashes[k]
		if rhash == "" {
			CreateChange(&changes, ObjectRemoved, k,
				lValues[k].GetValueNode(), nil, true,
				lValues[k].GetValue(), nil)
			continue
		}
		if lHashes[k] == rHashes[k] {
			continue
		}
		// run comparison.
		expChanges[k] = ComparePathItems(lValues[k].Value, rValues[k].Value)
	}

	//check right path item hashes
	for k := range rHashes {
		lhash := lHashes[k]
		if lhash == "" {
			CreateChange(&changes, ObjectAdded, k,
				nil, rValues[k].GetValueNode(), false,
				nil, rValues[k].GetValue())
			continue
		}
	}
	cc.ExpressionChanges = expChanges
	cc.ExtensionChanges = CompareExtensions(l.Extensions, r.Extensions)
	cc.Changes = changes
	if cc.TotalChanges() <= 0 {
		return nil
	}
	return cc
}
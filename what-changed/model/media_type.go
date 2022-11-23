// Copyright 2022 Princess B33f Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
    "fmt"
    "github.com/pb33f/libopenapi/datamodel/low"
    "github.com/pb33f/libopenapi/datamodel/low/v3"
    "github.com/pb33f/libopenapi/utils"
    "gopkg.in/yaml.v3"
)

// MediaTypeChanges represent changes made between two OpenAPI MediaType instances.
type MediaTypeChanges struct {
    *PropertyChanges
    SchemaChanges    *SchemaChanges              `json:"schemas,omitempty" yaml:"schemas,omitempty"`
    ExtensionChanges *ExtensionChanges           `json:"extensions,omitempty" yaml:"extensions,omitempty"`
    ExampleChanges   map[string]*ExampleChanges  `json:"examples,omitempty" yaml:"examples,omitempty"`
    EncodingChanges  map[string]*EncodingChanges `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// TotalChanges returns the total number of changes between two MediaType instances.
func (m *MediaTypeChanges) TotalChanges() int {
    c := m.PropertyChanges.TotalChanges()
    for k := range m.ExampleChanges {
        c += m.ExampleChanges[k].TotalChanges()
    }
    if m.SchemaChanges != nil {
        c += m.SchemaChanges.TotalChanges()
    }
    if len(m.EncodingChanges) > 0 {
        for i := range m.EncodingChanges {
            c += m.EncodingChanges[i].TotalChanges()
        }
    }
    if m.ExtensionChanges != nil {
        c += m.ExtensionChanges.TotalChanges()
    }
    return c
}

// TotalBreakingChanges returns the total number of breaking changes made between two MediaType instances.
func (m *MediaTypeChanges) TotalBreakingChanges() int {
    c := m.PropertyChanges.TotalBreakingChanges()
    for k := range m.ExampleChanges {
        c += m.ExampleChanges[k].TotalBreakingChanges()
    }
    if m.SchemaChanges != nil {
        c += m.SchemaChanges.TotalBreakingChanges()
    }
    if len(m.EncodingChanges) > 0 {
        for i := range m.EncodingChanges {
            c += m.EncodingChanges[i].TotalBreakingChanges()
        }
    }
    return c
}

// CompareMediaTypes compares a left and a right MediaType object for any changes. If found, a pointer to a
// MediaTypeChanges instance is returned, otherwise nothing is returned.
func CompareMediaTypes(l, r *v3.MediaType) *MediaTypeChanges {

    var props []*PropertyCheck
    var changes []*Change

    mc := new(MediaTypeChanges)

    if low.AreEqual(l, r) {
        return nil
    }

    // Example
    if !l.Example.IsEmpty() && !r.Example.IsEmpty() {
        if (utils.IsNodeMap(l.Example.ValueNode) && utils.IsNodeMap(r.Example.ValueNode)) ||
            (utils.IsNodeArray(l.Example.ValueNode) && utils.IsNodeArray(r.Example.ValueNode)) {
            render, err := yaml.Marshal(l.Example.ValueNode)
            if err != nil {
                l.Example.ValueNode.Value = fmt.Sprintf("%s", render)
            }
            render, err = yaml.Marshal(r.Example.ValueNode)
            if err == nil {
                r.Example.ValueNode.Value = fmt.Sprintf("%s", render)
            }
        }
        addPropertyCheck(&props, l.Example.ValueNode, r.Example.ValueNode,
            l.Example.Value, r.Example.Value, &changes, v3.ExampleLabel, false)

    } else {

        if utils.IsNodeMap(l.Example.ValueNode) || utils.IsNodeArray(l.Example.ValueNode) {
            render, err := yaml.Marshal(l.Example.ValueNode)
            if err != nil {
                l.Example.ValueNode.Value = fmt.Sprint(render)
            }
        }

        if utils.IsNodeMap(r.Example.ValueNode) || utils.IsNodeArray(r.Example.ValueNode) {
            render, err := yaml.Marshal(r.Example.ValueNode)
            if err == nil {
                r.Example.ValueNode.Value = fmt.Sprintf("%s", render)
            }
        }

        addPropertyCheck(&props, l.Example.ValueNode, r.Example.ValueNode,
            l.Example.Value, r.Example.Value, &changes, v3.ExampleLabel, false)
    }

    CheckProperties(props)

    // schema
    if !l.Schema.IsEmpty() && !r.Schema.IsEmpty() {
        mc.SchemaChanges = CompareSchemas(l.Schema.Value, r.Schema.Value)
    }
    if !l.Schema.IsEmpty() && r.Schema.IsEmpty() {
        CreateChange(&changes, ObjectRemoved, v3.SchemaLabel, l.Schema.ValueNode,
            nil, true, l.Schema.Value, nil)
    }
    if l.Schema.IsEmpty() && !r.Schema.IsEmpty() {
        CreateChange(&changes, ObjectAdded, v3.SchemaLabel, nil,
            r.Schema.ValueNode, true, nil, r.Schema.Value)
    }

    // examples
    mc.ExampleChanges = CheckMapForChanges(l.Examples.Value, r.Examples.Value,
        &changes, v3.ExamplesLabel, CompareExamples)

    // encoding
    mc.EncodingChanges = CheckMapForChanges(l.Encoding.Value, r.Encoding.Value,
        &changes, v3.EncodingLabel, CompareEncoding)

    mc.ExtensionChanges = CompareExtensions(l.Extensions, r.Extensions)
    mc.PropertyChanges = NewPropertyChanges(changes)

    if mc.TotalChanges() <= 0 {
        return nil
    }
    return mc
}

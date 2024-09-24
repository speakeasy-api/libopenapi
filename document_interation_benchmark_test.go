package libopenapi_test

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/stretchr/testify/require"
)

type loopFrameBenchmark struct {
	Type       string
	Restricted bool
}

type contextBenchmark struct {
	visited []string
	stack   []loopFrameBenchmark
}

func Benchmark_Docusign_Document_Iteration(b *testing.B) {
	// Setup code: read the spec file
	spec, err := os.ReadFile("test_specs/docusignv3.1.json")
	if err != nil {
		b.Fatalf("Failed to read file: %v", err)
	}

	config := &datamodel.DocumentConfiguration{
		BasePath:                            "./test_specs",
		IgnorePolymorphicCircularReferences: true,
		IgnoreArrayCircularReferences:       true,
		AllowFileReferences:                 true,
	}

	b.ResetTimer() // Reset the timer after setup

	for n := 0; n < b.N; n++ {
		// Code to benchmark

		doc, err := libopenapi.NewDocumentWithConfiguration(spec, config)
		if err != nil {
			b.Fatalf("Failed to create new document: %v", err)
		}

		m, errs := doc.BuildV3Model()
		if len(errs) > 0 {
			b.Fatalf("Failed to build V3 model with errors: %v", errs)
		}

		for path, pathItem := range m.Model.Paths.PathItems.FromOldest() {
			// Optional logging
			if b.N == 1 || testing.Verbose() {
				b.Logf("Path: %s", path)
			}

			iterateOperationsBenchmark(b, pathItem.GetOperations())
		}

		for path, pathItem := range m.Model.Webhooks.FromOldest() {
			if b.N == 1 || testing.Verbose() {
				b.Logf("Webhook Path: %s", path)
			}

			iterateOperationsBenchmark(b, pathItem.GetOperations())
		}

		for name, schemaProxy := range m.Model.Components.Schemas.FromOldest() {
			if b.N == 1 || testing.Verbose() {
				b.Logf("Schema Name: %s", name)
			}

			handleSchemaBenchmark(b, schemaProxy, contextBenchmark{})
		}
	}
}

func iterateOperationsBenchmark(b *testing.B, ops *orderedmap.Map[string, *v3.Operation]) {
	for method, op := range ops.FromOldest() {
		b.Log(method)

		for i, param := range op.Parameters {
			b.Log("param", i, param.Name)

			if param.Schema != nil {
				handleSchemaBenchmark(b, param.Schema, contextBenchmark{})
			}
		}

		if op.RequestBody != nil {
			b.Log("request body")

			for contentType, mediaType := range op.RequestBody.Content.FromOldest() {
				b.Log(contentType)

				if mediaType.Schema != nil {
					handleSchemaBenchmark(b, mediaType.Schema, contextBenchmark{})
				}
			}
		}

		if orderedmap.Len(op.Responses.Codes) > 0 {
			b.Log("responses")
		}

		for code, response := range op.Responses.Codes.FromOldest() {
			b.Log(code)

			for contentType, mediaType := range response.Content.FromOldest() {
				b.Log(contentType)

				if mediaType.Schema != nil {
					handleSchemaBenchmark(b, mediaType.Schema, contextBenchmark{})
				}
			}
		}

		if orderedmap.Len(op.Responses.Codes) > 0 {
			b.Log("callbacks")
		}

		for callbackName, callback := range op.Callbacks.FromOldest() {
			b.Log(callbackName)

			for name, pathItem := range callback.Expression.FromOldest() {
				b.Log(name)

				iterateOperationsBenchmark(b, pathItem.GetOperations())
			}
		}
	}
}

func handleSchemaBenchmark(b *testing.B, schProxy *base.SchemaProxy, ctx contextBenchmark) {
	if checkCircularReferenceBenchmark(b, &ctx, schProxy) {
		return
	}

	sch, err := schProxy.BuildSchema()
	require.NoError(b, err)

	typ, subTypes := getResolvedType(sch)

	b.Log("schema", typ, subTypes)

	if len(sch.Enum) > 0 {
		switch typ {
		case "string":
			return
		case "integer":
			return
		default:
			// handle as base type
		}
	}

	switch typ {
	case "allOf":
		fallthrough
	case "anyOf":
		fallthrough
	case "oneOf":
		if len(subTypes) > 0 {
			return
		}

		handleAllOfAnyOfOneOfBenchmark(b, sch, ctx)
	case "array":
		handleArrayBenchmark(b, sch, ctx)
	case "object":
		handleObject(b, sch, ctx)
	default:
		return
	}
}

func getResolvedType(sch *base.Schema) (string, []string) {
	subTypes := []string{}

	for _, t := range sch.Type {
		if t == "" { // treat empty type as any
			subTypes = append(subTypes, "any")
		} else if t != "null" {
			subTypes = append(subTypes, t)
		}
	}

	if len(sch.AllOf) > 0 {
		return "allOf", nil
	}

	if len(sch.AnyOf) > 0 {
		return "anyOf", nil
	}

	if len(sch.OneOf) > 0 {
		return "oneOf", nil
	}

	if len(subTypes) == 0 {
		if len(sch.Enum) > 0 {
			return "string", nil
		}

		if orderedmap.Len(sch.Properties) > 0 {
			return "object", nil
		}

		if sch.AdditionalProperties != nil {
			return "object", nil
		}

		if sch.Items != nil {
			return "array", nil
		}

		return "any", nil
	}

	if len(subTypes) == 1 {
		return subTypes[0], nil
	}

	return "oneOf", subTypes
}

func handleAllOfAnyOfOneOfBenchmark(b *testing.B, sch *base.Schema, ctx contextBenchmark) {
	var schemas []*base.SchemaProxy

	switch {
	case len(sch.AllOf) > 0:
		schemas = sch.AllOf
	case len(sch.AnyOf) > 0:
		schemas = sch.AnyOf
		ctx.stack = append(ctx.stack, loopFrameBenchmark{Type: "anyOf", Restricted: len(sch.AnyOf) == 1})
	case len(sch.OneOf) > 0:
		schemas = sch.OneOf
		ctx.stack = append(ctx.stack, loopFrameBenchmark{Type: "oneOf", Restricted: len(sch.OneOf) == 1})
	}

	for _, s := range schemas {
		handleSchemaBenchmark(b, s, ctx)
	}
}

func handleArrayBenchmark(b *testing.B, sch *base.Schema, ctx contextBenchmark) {
	ctx.stack = append(ctx.stack, loopFrameBenchmark{Type: "array", Restricted: sch.MinItems != nil && *sch.MinItems > 0})

	if sch.Items != nil && sch.Items.IsA() {
		handleSchemaBenchmark(b, sch.Items.A, ctx)
	}

	if sch.Contains != nil {
		handleSchemaBenchmark(b, sch.Contains, ctx)
	}

	if sch.PrefixItems != nil {
		for _, s := range sch.PrefixItems {
			handleSchemaBenchmark(b, s, ctx)
		}
	}
}

func handleObject(b *testing.B, sch *base.Schema, ctx contextBenchmark) {
	for name, schemaProxy := range sch.Properties.FromOldest() {
		ctx.stack = append(ctx.stack, loopFrameBenchmark{Type: "object", Restricted: slices.Contains(sch.Required, name)})
		handleSchemaBenchmark(b, schemaProxy, ctx)
	}

	if sch.AdditionalProperties != nil && sch.AdditionalProperties.IsA() {
		handleSchemaBenchmark(b, sch.AdditionalProperties.A, ctx)
	}
}

func checkCircularReferenceBenchmark(b *testing.B, ctx *contextBenchmark, schProxy *base.SchemaProxy) bool {
	loopRef := getSimplifiedRef(schProxy.GetReference())

	if loopRef != "" {
		if slices.Contains(ctx.visited, loopRef) {
			isRestricted := true
			containsObject := false

			for _, v := range ctx.stack {
				if v.Type == "object" {
					containsObject = true
				}

				if v.Type == "array" && !v.Restricted {
					isRestricted = false
				} else if !v.Restricted {
					isRestricted = false
				}
			}

			if !containsObject {
				isRestricted = true
			}

			require.False(b, isRestricted, "circular reference: %s", append(ctx.visited, loopRef))
			return true
		}

		ctx.visited = append(ctx.visited, loopRef)
	}

	return false
}

// getSimplifiedRef will return the reference without the preceding file path
// caveat is that if a spec has the same ref in two different files they include this may identify them incorrectly
// but currently a problem anyway as libopenapi when returning references from an external file won't include the file path
// for a local reference with that file and so we might fail to distinguish between them that way.
// The fix needed is for libopenapi to also track which file the reference is in so we can always prefix them with the file path
func getSimplifiedRef(ref string) string {
	if ref == "" {
		return ""
	}

	refParts := strings.Split(ref, "#/")
	return "#/" + refParts[len(refParts)-1]
}

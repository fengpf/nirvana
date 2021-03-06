/*
Copyright 2017 Caicloud Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generators

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/namer"
	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

func construct(t *testing.T, files map[string]string, testNamer namer.Namer) (*parser.Builder, types.Universe, []*types.Type) {
	b := parser.New()
	for name, src := range files {
		if err := b.AddFileForTest(filepath.Dir(name), name, []byte(src)); err != nil {
			t.Fatal(err)
		}
	}
	u, err := b.FindTypes()
	if err != nil {
		t.Fatal(err)
	}
	orderer := namer.Orderer{Namer: testNamer}
	o := orderer.OrderUniverse(u)
	return b, u, o
}

func constructType(t *testing.T, code string) (*generator.Context, *types.Type, namer.ImportTracker) {
	var testFiles = map[string]string{
		"base/foo/bar.go": code,
	}
	it := generator.NewImportTracker()
	rawNamer := namer.NewRawNamer("o", nil)
	namers := namer.NameSystems{
		"raw": namer.NewRawNamer("", it),
	}
	builder, universe, _ := construct(t, testFiles, rawNamer)
	context, err := generator.NewContext(builder, namers, "raw")
	if err != nil {
		t.Fatal(err)
	}
	blahT := universe.Type(types.Name{Package: "base/foo", Name: "Blah"})
	return context, blahT, it
}

func newWriter(ctx *generator.Context) (*openAPITypeWriter, *bytes.Buffer) {
	buffer := &bytes.Buffer{}
	sw := generator.NewSnippetWriter(buffer, ctx, "$", "$")
	return newOpenAPITypeWriter(sw), buffer
}

// TODO(liubog2008): add more unit test

func TestSimple(t *testing.T) {
	ctx, typ, _ := constructType(t, `
package foo

// Blah is a test.
// +caicloud:openapi-gen=true
type Blah struct {
	// A simple string
	String string
	// A simple int
	Int int `+"`"+`json:",omitempty"`+"`"+`
	// An int considered string simple int
	IntString int `+"`"+`json:",string"`+"`"+`
	// A simple int64
	Int64 int64
	// A simple int32
	Int32 int32
	// A simple int16
	Int16 int16
	// A simple int8
	Int8 int8
	// A simple int
	Uint uint
	// A simple int64
	Uint64 uint64
	// A simple int32
	Uint32 uint32
	// A simple int16
	Uint16 uint16
	// A simple int8
	Uint8 uint8
	// A simple byte
	Byte byte
	// A simple boolean
	Bool bool
	// A simple float64
	Float64 float64
	// A simple float32
	Float32 float32
	// a base64 encoded characters
	ByteArray []byte
}
		`)
	sw, buf := newWriter(ctx)
	if err := sw.generate(typ); err != nil {
		t.Fatal(err)
	}
	res := trimPrefixSpace(buf.Bytes())
	assert.Equal(t, `"base/foo.Blah": {
Schema: spec.Schema{
SchemaProps: spec.SchemaProps{
Description: "Blah is a test.",
Properties: map[string]spec.Schema{
"String": {
SchemaProps: spec.SchemaProps{
Description: "A simple string",
Type: []string{"string"},
Format: "",
},
},
"Int64": {
SchemaProps: spec.SchemaProps{
Description: "A simple int64",
Type: []string{"integer"},
Format: "int64",
},
},
"Int32": {
SchemaProps: spec.SchemaProps{
Description: "A simple int32",
Type: []string{"integer"},
Format: "int32",
},
},
"Int16": {
SchemaProps: spec.SchemaProps{
Description: "A simple int16",
Type: []string{"integer"},
Format: "int16",
},
},
"Int8": {
SchemaProps: spec.SchemaProps{
Description: "A simple int8",
Type: []string{"integer"},
Format: "uint8",
},
},
"Uint": {
SchemaProps: spec.SchemaProps{
Description: "A simple int",
Type: []string{"integer"},
Format: "uint",
},
},
"Uint64": {
SchemaProps: spec.SchemaProps{
Description: "A simple int64",
Type: []string{"integer"},
Format: "uint64",
},
},
"Uint32": {
SchemaProps: spec.SchemaProps{
Description: "A simple int32",
Type: []string{"integer"},
Format: "uint32",
},
},
"Uint16": {
SchemaProps: spec.SchemaProps{
Description: "A simple int16",
Type: []string{"integer"},
Format: "uint16",
},
},
"Uint8": {
SchemaProps: spec.SchemaProps{
Description: "A simple int8",
Type: []string{"integer"},
Format: "uint8",
},
},
"Byte": {
SchemaProps: spec.SchemaProps{
Description: "A simple byte",
Type: []string{"integer"},
Format: "uint8",
},
},
"Bool": {
SchemaProps: spec.SchemaProps{
Description: "A simple boolean",
Type: []string{"boolean"},
Format: "",
},
},
"Float64": {
SchemaProps: spec.SchemaProps{
Description: "A simple float64",
Type: []string{"number"},
Format: "double",
},
},
"Float32": {
SchemaProps: spec.SchemaProps{
Description: "A simple float32",
Type: []string{"number"},
Format: "float",
},
},
"ByteArray": {
SchemaProps: spec.SchemaProps{
Description: "a base64 encoded characters",
Type: []string{"string"},
Format: "byte",
},
},
},
Required: []string{
"String",
"Int64",
"Int32",
"Int16",
"Int8",
"Uint",
"Uint64",
"Uint32",
"Uint16",
"Uint8",
"Byte",
"Bool",
"Float64",
"Float32",
"ByteArray",
},
},
},
},
`, string(res))
}

func TestPointer(t *testing.T) {
	ctx, typ, it := constructType(t, `
package foo

// PointerSample demonstrate pointer's properties
type Blah struct {
	// A string pointer
	StringPointer *string
	// A struct pointer
	StructPointer *Blah
	// A slice pointer
	SlicePointer *[]string
	// A map pointer
	MapPointer *map[string]string
}
	`)
	sw, buf := newWriter(ctx)

	err := sw.generate(typ)
	if err != nil {
		t.Fatal(err)
	}

	res := trimPrefixSpace(buf.Bytes())

	assert.Equal(t, `"base/foo.Blah": {
Schema: spec.Schema{
SchemaProps: spec.SchemaProps{
Description: "PointerSample demonstrate pointer's properties",
Properties: map[string]spec.Schema{
"StringPointer": {
SchemaProps: spec.SchemaProps{
Description: "A string pointer",
Type: []string{"string"},
Format: "",
},
},
"StructPointer": {
SchemaProps: spec.SchemaProps{
Description: "A struct pointer",
Ref: ref("base/foo.Blah"),
},
},
"SlicePointer": {
SchemaProps: spec.SchemaProps{
Description: "A slice pointer",
Type: []string{"array"},
Items: &spec.SchemaOrArray{
Schema: &spec.Schema{
SchemaProps: spec.SchemaProps{
Type: []string{"string"},
Format: "",
},
},
},
},
},
"MapPointer": {
SchemaProps: spec.SchemaProps{
Description: "A map pointer",
Type: []string{"object"},
AdditionalProperties: &spec.SchemaOrBool{
Schema: &spec.Schema{
SchemaProps: spec.SchemaProps{
Type: []string{"string"},
Format: "",
},
},
},
},
},
},
Required: []string{
"StringPointer",
"StructPointer",
"SlicePointer",
"MapPointer",
},
},
},
Dependencies: []string{
"base/foo.Blah",
},
},
`, string(res))

	imports := it.ImportLines()
	assert.Equal(t, []string{
		`foo "base/foo"`,
		`spec "github.com/go-openapi/spec"`,
	}, imports, "imports should be equal")
}

func trimPrefixSpace(in []byte) []byte {
	res := []byte{}
	first := true
	for _, b := range in {
		if !first {
			res = append(res, b)
			first = b == '\n'
		} else if b != '\t' && b != ' ' {
			first = false
			res = append(res, b)
		}
	}
	return res
}

func TestContext(t *testing.T) {
	ctx, typ, it := constructType(t, `
package foo
	`)
	sw, _ := newWriter(ctx)
	err := sw.generate(typ)
	if err != nil {
		t.Fatal(err)
	}

	imports := it.ImportLines()
	assert.Equal(t, []string{}, imports, "imports should be equal")
}

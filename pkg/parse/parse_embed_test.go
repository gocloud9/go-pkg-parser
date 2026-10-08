package parse_test

import (
	"embed"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/gocloud9/go-pkg-parser/pkg/parse"
)

//go:embed _testdata/simple
var simpleFS embed.FS

//go:embed _testdata/embedded
var embeddedFS embed.FS

func TestParser_ParseEmbed(t *testing.T) {
	tests := []struct {
		name    string
		fsys    embed.FS
		path    string
		want    *parse.Results
		wantErr bool
	}{
		{
			name: "simple struct",
			fsys: simpleFS,
			path: "_testdata/simple",
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"simple": {
						Name:         "simple",
						Constants:    map[string]*parse.ConstantInfo{},
						Functions:    map[string]*parse.FuncInfo{},
						Interfaces:   map[string]*parse.InterfaceInfo{},
						Vars:         map[string]*parse.VarInfo{},
						DefinedTypes: map[string]*parse.DefinedTypeInfo{},
						Aliases:      map[string]*parse.AliasTypeInfo{},
						Structs: map[string]*parse.StructInfo{
							"User": {
								Name: "User",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								Fields: map[string]*parse.FieldInfo{
									"ID": {
										Name: "ID",
										Tags: map[string][]string{
											"json": {"id"},
										},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "string",
											ExternalTypeName: "string",
										},
										Markers: map[string]string{
											"+something:id": "true",
										},
									},
									"DisplayName": {
										Name: "DisplayName",
										Tags: map[string][]string{
											"json": {"display_name"},
										},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "*string",
											ExternalTypeName: "*string",
											IsPointer:        true,
											Pointer: &parse.TypeInfo{
												TypeName:         "string",
												ExternalTypeName: "string",
											},
										},
										Markers: map[string]string{},
									},
									"Email": {
										Name: "Email",
										Tags: map[string][]string{
											"json": {"email"},
										},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "string",
											ExternalTypeName: "string",
										},
										Markers: map[string]string{},
									},
									"Age": {
										Name: "Age",
										Tags: map[string][]string{
											"json": {"age"},
										},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "int",
											ExternalTypeName: "int",
										},
										Markers: map[string]string{},
									},
								},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods:        map[string]*parse.FuncInfo{},
							},
						},
					},
				},
			},
		},
		{
			name: "embedded structs and interfaces",
			fsys: embeddedFS,
			path: "_testdata/embedded",
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"embedded": {
						Name:      "embedded",
						Constants: map[string]*parse.ConstantInfo{},
						Vars:      map[string]*parse.VarInfo{},
						Functions: map[string]*parse.FuncInfo{},
						DefinedTypes: map[string]*parse.DefinedTypeInfo{
							"ParentStruct": {
								Name:    "ParentStruct",
								Markers: map[string]string{},
								TypeInfo: &parse.TypeInfo{
									TypeName:         "func()",
									ExternalTypeName: "func()",
									IsFunc:           true,
									Func: &parse.FuncDefInfo{
										Params:  []*parse.ParamInfo{},
										Results: []*parse.ResultInfo{},
									},
								},
							},
						},
						Aliases: map[string]*parse.AliasTypeInfo{},
						Interfaces: map[string]*parse.InterfaceInfo{
							"ParentInterface": {
								Name:          "ParentInterface",
								Markers:       map[string]string{},
								Methods:       map[string]*parse.FuncInfo{},
								EmbeddedTypes: map[string]*parse.EmbeddedTypeInfo{},
							},
							"ChildInterface": {
								Name: "ChildInterface",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								Methods: map[string]*parse.FuncInfo{},
								EmbeddedTypes: map[string]*parse.EmbeddedTypeInfo{
									"ParentInterface": {
										Name:     "ParentInterface",
										TypeName: "ParentInterface",
										Markers:  map[string]string{"+Bar": "123", "+Foo": "true"},
									},
									"Parent": {
										Name:     "Parent",
										TypeName: "Parent",
										Markers:  map[string]string{"+Bar": "123", "+Foo": "true"},
									},
								},
							},
						},
						Structs: map[string]*parse.StructInfo{
							"Child": {
								Name: "Child",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								Fields: map[string]*parse.FieldInfo{},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{
									"ParentInterface": {
										Name:     "ParentInterface",
										TypeName: "ParentInterface",
										Markers:  map[string]string{"+Bar": "123", "+Foo": "true"},
										Tags:     map[string][]string{"yaml": {"", "inline"}},
									},
									"Parent": {
										Name:     "Parent",
										TypeName: "Parent",
										Markers:  map[string]string{"+Bar": "123", "+Foo": "true"},
										Tags:     map[string][]string{"yaml": {"", "inline"}},
									},
								},
								Methods: map[string]*parse.FuncInfo{},
							},
							"Parent": {
								Name:           "Parent",
								Markers:        map[string]string{},
								Fields:         map[string]*parse.FieldInfo{},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods:        map[string]*parse.FuncInfo{},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parse.Parser{}

			got, err := p.ParseEmbed(tt.fsys, parse.Options{Path: tt.path})
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseEmbed() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.want != nil {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

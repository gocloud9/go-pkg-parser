package parse_test

import (
	"github.com/gocloud9/go-pkg-parser/pkg/parse"
	"github.com/google/go-cmp/cmp"
	"os"
	"path/filepath"
	"regexp"

	"testing"
)

func TestParser_ParseDirectory(t *testing.T) {
	type args struct {
		path      string
		skipFiles []*regexp.Regexp
	}
	tests := []struct {
		name    string
		args    args
		want    *parse.Results
		wantErr bool
	}{
		{
			name: "simple struct",
			args: args{
				path: "./_testdata/simple",
			},
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
			wantErr: false,
		},
		{
			name: "permutations with multiple packages",
			args: args{
				path: "./_testdata/permutation",
				skipFiles: []*regexp.Regexp{
					regexp.MustCompile(`Skipped file in permutation test`),
				},
			},
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"package1": {
						Name:         "package1",
						DefinedTypes: map[string]*parse.DefinedTypeInfo{},
						Aliases:      map[string]*parse.AliasTypeInfo{},
						Structs: map[string]*parse.StructInfo{
							"AnotherUser": {
								Name: "AnotherUser",
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
										Markers: map[string]string{"+ID": "true"},
									},
									"DisplayName": {
										Name: "DisplayName",
										Tags: map[string][]string{
											"json": {"display_name"},
										},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "string",
											ExternalTypeName: "string",
										},
										Markers: map[string]string{},
									},
									"Duration": {
										Name:    "Duration",
										Markers: map[string]string{},
										Tags:    map[string][]string{"json": {"duration"}},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "time.Duration",
											ExternalTypeName: "time.Duration",
											IsImported:       true,
											ImportedType: &parse.ImportedTypeInfo{
												TypeName:            "Duration",
												ImportRaw:           "\"time\"",
												PackagePath:         "time",
												PackageDefaultAlias: "time",
											},
										},
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
									"Time": {
										Name:    "Time",
										Markers: map[string]string{},
										Tags:    map[string][]string{"json": {"time"}},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "time.Time",
											ExternalTypeName: "time.Time",
											IsImported:       true,
											ImportedType: &parse.ImportedTypeInfo{
												TypeName:            "Time",
												ImportRaw:           "\"time\"",
												PackagePath:         "time",
												PackageDefaultAlias: "time",
											},
										},
									},
									"Timestamp": {
										Name:    "Timestamp",
										Markers: map[string]string{},
										Tags:    map[string][]string{"json": {"timestamp"}},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "timestamppb.Timestamp",
											ExternalTypeName: "timestamppb.Timestamp",
											IsImported:       true,
											ImportedType: &parse.ImportedTypeInfo{
												TypeName:            "Timestamp",
												ImportRaw:           "\"google.golang.org/protobuf/types/known/timestamppb\"",
												PackagePath:         "google.golang.org/protobuf/types/known/timestamppb",
												PackageDefaultAlias: "timestamppb",
											},
										},
									},
								},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods:        map[string]*parse.FuncInfo{},
							},
						},
						Constants:  map[string]*parse.ConstantInfo{},
						Functions:  map[string]*parse.FuncInfo{},
						Interfaces: map[string]*parse.InterfaceInfo{},
						Vars:       map[string]*parse.VarInfo{},
					},
					"package2": {
						Name:         "package2",
						DefinedTypes: map[string]*parse.DefinedTypeInfo{},
						Aliases:      map[string]*parse.AliasTypeInfo{},
						Structs: map[string]*parse.StructInfo{
							"SomeStruct": {
								Name:    "SomeStruct",
								Markers: map[string]string{},
								Fields: map[string]*parse.FieldInfo{
									"StringField": {
										Name: "StringField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "string",
											ExternalTypeName: "string",
										},
										Tags:    map[string][]string{"json": {"string_field"}, "yaml": {"stringField"}},
										Markers: map[string]string{},
									},
									"IntField": {
										Name: "IntField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "int",
											ExternalTypeName: "int",
										},
										Tags:    map[string][]string{"json": {"int_field"}, "yaml": {"intField"}},
										Markers: map[string]string{},
									},
									"BoolField": {
										Name: "BoolField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "bool",
											ExternalTypeName: "bool",
										},
										Tags:    map[string][]string{"json": {"bool_field"}, "yaml": {"boolField"}},
										Markers: map[string]string{},
									},
									"ChanField": {
										Name: "ChanField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "chan int",
											ExternalTypeName: "chan int",
											IsChan:           true,
											Chan: &parse.TypeInfo{
												TypeName:         "int",
												ExternalTypeName: "int",
											},
										},
										Tags:    map[string][]string{"json": {"chan_field"}, "yaml": {"chanField"}},
										Markers: map[string]string{},
									},
									"MapField": {
										Name: "MapField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "map[string]int",
											ExternalTypeName: "map[string]int",
											IsMap:            true,
											MapKey: &parse.TypeInfo{
												TypeName:         "string",
												ExternalTypeName: "string",
											},
											MapValue: &parse.TypeInfo{
												TypeName:         "int",
												ExternalTypeName: "int",
											},
										},
										Tags:    map[string][]string{"json": {"map_field"}, "yaml": {"mapField"}},
										Markers: map[string]string{},
									},
									"SliceField": {
										Name: "SliceField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "[]int",
											ExternalTypeName: "[]int",
											IsSlice:          true,
											Slice: &parse.TypeInfo{
												TypeName:         "int",
												ExternalTypeName: "int",
											},
										},
										Tags:    map[string][]string{"json": {"slice_field"}, "yaml": {"sliceField"}},
										Markers: map[string]string{},
									},
									"SubStructField": {
										Name: "SubStructField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "SubStruct",
											ExternalTypeName: "package2.SubStruct",
											IsStruct:         true,
											IsType:           true,
											TypeOf: &parse.TypeInfo{
												IsStruct: true,
											},
										},
										Tags:    map[string][]string{"json": {"sub_struct_field"}, "yaml": {"subStructField"}},
										Markers: map[string]string{},
									},
									"SubStructMapField": {
										Name: "SubStructMapField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "map[string]SubStruct",
											ExternalTypeName: "map[string]package2.SubStruct",
											IsMap:            true,
											MapKey: &parse.TypeInfo{
												TypeName:         "string",
												ExternalTypeName: "string",
											},
											MapValue: &parse.TypeInfo{
												TypeName:         "SubStruct",
												ExternalTypeName: "package2.SubStruct",
												IsStruct:         true,
												IsType:           true,
												TypeOf: &parse.TypeInfo{
													IsStruct: true,
												},
											},
										},
										Tags:    map[string][]string{"json": {"sub_struct_map_field"}, "yaml": {"subStructMapField"}},
										Markers: map[string]string{},
									},
									"SubStructSliceField": {
										Name: "SubStructSliceField",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "[]SubStruct",
											ExternalTypeName: "[]package2.SubStruct",
											IsSlice:          true,
											Slice: &parse.TypeInfo{
												TypeName:         "SubStruct",
												ExternalTypeName: "package2.SubStruct",
												IsStruct:         true,
												IsType:           true,
												TypeOf: &parse.TypeInfo{
													IsStruct: true,
												},
											},
										},
										Tags:    map[string][]string{"json": {"sub_struct_slice_field"}, "yaml": {"subStructSliceField"}},
										Markers: map[string]string{},
									},
								},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods:        map[string]*parse.FuncInfo{},
							},
							"SubStruct": {
								Name: "SubStruct", Markers: map[string]string{}, Fields: map[string]*parse.FieldInfo{
									"Something": {
										Name: "Something",
										TypeInfo: &parse.TypeInfo{
											TypeName:         "string",
											ExternalTypeName: "string",
										},
										Tags:    map[string][]string{},
										Markers: map[string]string{},
									},
								},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods:        map[string]*parse.FuncInfo{},
							},
						},
						Constants:  map[string]*parse.ConstantInfo{},
						Functions:  map[string]*parse.FuncInfo{},
						Interfaces: map[string]*parse.InterfaceInfo{},
						Vars:       map[string]*parse.VarInfo{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "functions",
			args: args{
				path: "./_testdata/functions",
			},
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"functions": {
						Name:         "functions",
						Constants:    map[string]*parse.ConstantInfo{},
						DefinedTypes: map[string]*parse.DefinedTypeInfo{},
						Aliases:      map[string]*parse.AliasTypeInfo{},
						Structs: map[string]*parse.StructInfo{
							"Field": {
								Name: "Field",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								Fields:         map[string]*parse.FieldInfo{},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods: map[string]*parse.FuncInfo{
									"Test5": {
										Name: "Test5",
										Markers: map[string]string{
											"+Foo": "true",
											"+Bar": "123",
										},
										FuncDefInfo: &parse.FuncDefInfo{
											Params: []*parse.ParamInfo{
												{
													Name: "arg",
													TypeInfo: &parse.TypeInfo{
														TypeName:         "Field",
														ExternalTypeName: "functions.Field",
														IsStruct:         true,
														IsType:           true,
														TypeOf:           &parse.TypeInfo{IsStruct: true},
													},
												},
											},
											Results: []*parse.ResultInfo{},
										},
										HasReciver:  true,
										ReciverName: "Field",
									},
								},
							},
							"Reference": {
								Name: "Reference",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								Methods: map[string]*parse.FuncInfo{
									"Test6": {
										Name: "Test6",
										Markers: map[string]string{
											"+Foo": "true",
											"+Bar": "123",
										},
										FuncDefInfo: &parse.FuncDefInfo{
											Params: []*parse.ParamInfo{
												{
													Name: "arg",
													TypeInfo: &parse.TypeInfo{
														TypeName:         "Field",
														ExternalTypeName: "functions.Field",
														IsStruct:         true,
														IsType:           true,
														TypeOf:           &parse.TypeInfo{IsStruct: true},
													},
												},
											},
											Results: []*parse.ResultInfo{},
										},
										HasReciver:  true,
										ReciverName: "Reference",
									},
								},
								Fields:         map[string]*parse.FieldInfo{},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
							},
						},
						Interfaces: map[string]*parse.InterfaceInfo{},
						Vars: map[string]*parse.VarInfo{
							"myFunc": {
								Name: "myFunc",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
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
							"myGroupedFunc": {
								Name: "myGroupedFunc",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
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
						Functions: map[string]*parse.FuncInfo{
							"Test1": {
								Name:    "Test1",
								Markers: map[string]string{"+Bar": "123", "+Foo": "true"},
								FuncDefInfo: &parse.FuncDefInfo{
									IsVariadic: false,
									Params:     []*parse.ParamInfo{},
									Results:    []*parse.ResultInfo{},
								},
							},
							"Test2": {
								Name:    "Test2",
								Markers: map[string]string{"+Bar": "123", "+Foo": "true"},
								FuncDefInfo: &parse.FuncDefInfo{
									IsVariadic: false,
									Params:     []*parse.ParamInfo{},
									Results: []*parse.ResultInfo{
										{
											TypeInfo: &parse.TypeInfo{
												TypeName:         "error",
												ExternalTypeName: "error",
											},
										},
									},
								},
							},
							"Test3": {
								Name: "Test3",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								FuncDefInfo: &parse.FuncDefInfo{
									Params: []*parse.ParamInfo{
										{
											Name: "arg",
											TypeInfo: &parse.TypeInfo{
												TypeName:         "string",
												ExternalTypeName: "string",
											},
										},
									},
									Results: []*parse.ResultInfo{},
								},
							},
							"Test4": {
								Name: "Test4",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								FuncDefInfo: &parse.FuncDefInfo{
									Params: []*parse.ParamInfo{
										{
											Name: "arg",
											TypeInfo: &parse.TypeInfo{
												TypeName:         "Field",
												ExternalTypeName: "functions.Field",
												IsStruct:         true,
												IsType:           true,
												TypeOf:           &parse.TypeInfo{IsStruct: true},
											},
										},
									},
									Results: []*parse.ResultInfo{},
								},
							},
							"Field.Test5": {
								Name: "Test5",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								FuncDefInfo: &parse.FuncDefInfo{
									Params: []*parse.ParamInfo{
										{
											Name: "arg",
											TypeInfo: &parse.TypeInfo{
												TypeName:         "Field",
												ExternalTypeName: "functions.Field",
												IsStruct:         true,
												IsType:           true,
												TypeOf:           &parse.TypeInfo{IsStruct: true},
											},
										},
									},
									Results: []*parse.ResultInfo{},
								},
								HasReciver:  true,
								ReciverName: "Field",
							},
							"Reference.Test6": {
								Name: "Test6",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								FuncDefInfo: &parse.FuncDefInfo{
									Params: []*parse.ParamInfo{
										{
											Name: "arg",
											TypeInfo: &parse.TypeInfo{
												TypeName:         "Field",
												ExternalTypeName: "functions.Field",
												IsStruct:         true,
												IsType:           true,
												TypeOf:           &parse.TypeInfo{IsStruct: true},
											},
										},
									},
									Results: []*parse.ResultInfo{},
								},
								HasReciver:  true,
								ReciverName: "Reference",
							},
							"Variadic": {
								Name:    "Variadic",
								Markers: map[string]string{},
								FuncDefInfo: &parse.FuncDefInfo{
									IsVariadic: true,
									Params: []*parse.ParamInfo{
										{
											Name: "vArg",
											TypeInfo: &parse.TypeInfo{
												TypeName:         "...Field",
												ExternalTypeName: "...functions.Field",
												IsEllipsis:       true,
												Ellipsis: &parse.TypeInfo{
													TypeName:         "Field",
													ExternalTypeName: "functions.Field",
													IsStruct:         true,
													IsType:           true,
													TypeOf:           &parse.TypeInfo{IsStruct: true},
												},
											},
										},
									},
									Results: []*parse.ResultInfo{},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "interfaces",
			args: args{
				path: "./_testdata/interfaces",
			},
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"interfaces": {
						Name:         "interfaces",
						Constants:    map[string]*parse.ConstantInfo{},
						DefinedTypes: map[string]*parse.DefinedTypeInfo{},
						Aliases:      map[string]*parse.AliasTypeInfo{},
						Structs: map[string]*parse.StructInfo{
							"TestStruct": {
								Name:           "TestStruct",
								Markers:        map[string]string{},
								Fields:         map[string]*parse.FieldInfo{},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods:        map[string]*parse.FuncInfo{},
							},
						},
						Vars:      map[string]*parse.VarInfo{},
						Functions: map[string]*parse.FuncInfo{},
						Interfaces: map[string]*parse.InterfaceInfo{
							"MyInterface": {
								Name:    "MyInterface",
								Markers: map[string]string{"+Bar": "123", "+Foo": "true"},
								Methods: map[string]*parse.FuncInfo{
									"DoSomething": {
										Name:    "DoSomething",
										Markers: map[string]string{"+Bar": "123", "+Foo": "true"},
										FuncDefInfo: &parse.FuncDefInfo{
											IsVariadic: false,
											Params: []*parse.ParamInfo{
												{
													Name: "input",
													TypeInfo: &parse.TypeInfo{
														TypeName:         "*string",
														ExternalTypeName: "*string",
														IsPointer:        true,
														Pointer: &parse.TypeInfo{
															TypeName:         "string",
															ExternalTypeName: "string",
														},
													},
												},
												{
													Name: "f",
													TypeInfo: &parse.TypeInfo{
														TypeName:         "func([]TestInterface) map[string]TestStruct",
														ExternalTypeName: "func([]interfaces.TestInterface) map[string]interfaces.TestStruct",
														IsFunc:           true,
														Func: &parse.FuncDefInfo{
															Params: []*parse.ParamInfo{
																{
																	Name: "",
																	TypeInfo: &parse.TypeInfo{
																		TypeName:         "[]TestInterface",
																		ExternalTypeName: "[]interfaces.TestInterface",
																		IsSlice:          true,
																		Slice: &parse.TypeInfo{
																			TypeName:         "TestInterface",
																			ExternalTypeName: "interfaces.TestInterface",
																			IsInterface:      true,
																			IsType:           true,
																			TypeOf:           &parse.TypeInfo{IsInterface: true},
																		},
																	},
																},
															},
															Results: []*parse.ResultInfo{
																{
																	TypeInfo: &parse.TypeInfo{
																		TypeName:         "map[string]TestStruct",
																		ExternalTypeName: "map[string]interfaces.TestStruct",
																		IsMap:            true,
																		MapKey: &parse.TypeInfo{
																			TypeName:         "string",
																			ExternalTypeName: "string",
																		},
																		MapValue: &parse.TypeInfo{
																			TypeName:         "TestStruct",
																			ExternalTypeName: "interfaces.TestStruct",
																			IsStruct:         true,
																			IsType:           true,
																			TypeOf:           &parse.TypeInfo{IsStruct: true},
																		},
																	},
																},
															},
														},
													},
												},
											},
											Results: []*parse.ResultInfo{
												{
													Name: "output",
													TypeInfo: &parse.TypeInfo{
														TypeName:         "string",
														ExternalTypeName: "string",
													},
												},
												{
													Name: "err",
													TypeInfo: &parse.TypeInfo{
														TypeName:         "error",
														ExternalTypeName: "error",
													},
												},
											},
										},
									},
								},
								EmbeddedTypes: map[string]*parse.EmbeddedTypeInfo{},
							},
							"TestInterface": {
								Name:          "TestInterface",
								Markers:       map[string]string{},
								Methods:       map[string]*parse.FuncInfo{},
								EmbeddedTypes: map[string]*parse.EmbeddedTypeInfo{},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "globals",
			args: args{
				path: "./_testdata/globals",
			},
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"globals": {
						Name:    "globals",
						Structs: map[string]*parse.StructInfo{},
						DefinedTypes: map[string]*parse.DefinedTypeInfo{
							"MyString": {
								Name:    "MyString",
								Markers: map[string]string{},
								TypeInfo: &parse.TypeInfo{
									TypeName:         "string",
									ExternalTypeName: "string",
								},
							},
						},
						Aliases: map[string]*parse.AliasTypeInfo{},
						Constants: map[string]*parse.ConstantInfo{
							"myConstant": {
								Name: "myConstant",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								TypeName: "string",
								Value:    `"test"`,
							},
							"myStringType": {
								Name:     "myStringType",
								Markers:  map[string]string{},
								TypeName: "MyString",
								Value:    `"test"`,
							},
						},
						Vars: map[string]*parse.VarInfo{
							"myFunc": {
								Name: "myFunc",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								TypeInfo: &parse.TypeInfo{
									TypeName:         "func(string) error",
									ExternalTypeName: "func(string) error",
									IsFunc:           true,
									Func: &parse.FuncDefInfo{
										Params: []*parse.ParamInfo{
											{
												Name: "arg",
												TypeInfo: &parse.TypeInfo{
													TypeName:         "string",
													ExternalTypeName: "string",
												},
											},
										},
										Results: []*parse.ResultInfo{
											{
												TypeInfo: &parse.TypeInfo{
													TypeName:         "error",
													ExternalTypeName: "error",
												},
											},
										},
									},
								},
							},
						},
						Functions:  map[string]*parse.FuncInfo{},
						Interfaces: map[string]*parse.InterfaceInfo{},
					},
				},
			},
		},
		{
			name: "embedded",
			args: args{
				path: "./_testdata/embedded",
			},
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
		{
			name: "typing",
			args: args{
				path: "./_testdata/typing",
			},
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"typing": {
						Name: "typing",
						Structs: map[string]*parse.StructInfo{
							"AStruct": {
								Name: "AStruct",
								Markers: map[string]string{
									"+Foo": "true",
									"+Bar": "123",
								},
								Fields: map[string]*parse.FieldInfo{
									"Field1": {
										Name:    "Field1",
										Markers: map[string]string{},
										Tags:    map[string][]string{},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "StringType",
											ExternalTypeName: "typing.StringType",
											IsType:           true,
											TypeOf: &parse.TypeInfo{
												TypeName:         "string",
												ExternalTypeName: "string",
											},
										},
									},
									"Field2": {
										Name:    "Field2",
										Markers: map[string]string{},
										Tags:    map[string][]string{},
										TypeInfo: &parse.TypeInfo{
											TypeName:         "IntType",
											ExternalTypeName: "typing.IntType",
											IsType:           true,
											TypeOf: &parse.TypeInfo{
												TypeName:         "int",
												ExternalTypeName: "int",
											},
										},
									},
								},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods:        map[string]*parse.FuncInfo{},
							},
						},
						Constants:  map[string]*parse.ConstantInfo{},
						Functions:  map[string]*parse.FuncInfo{},
						Interfaces: map[string]*parse.InterfaceInfo{},
						Vars:       map[string]*parse.VarInfo{},
						DefinedTypes: map[string]*parse.DefinedTypeInfo{
							"StringType": {
								Name:    "StringType",
								Markers: map[string]string{"+Bar": "123", "+Foo": "true"},
								TypeInfo: &parse.TypeInfo{
									TypeName:         "string",
									ExternalTypeName: "string",
								},
							},
							"AStructType": {
								Name:     "AStructType",
								Markers:  map[string]string{"+Bar": "123", "+Foo": "true"},
								TypeInfo: &parse.TypeInfo{TypeName: "AStruct", ExternalTypeName: "typing.AStruct", IsStruct: true, IsType: true, TypeOf: &parse.TypeInfo{IsStruct: true}},
							},
							"IntType": {
								Name:     "IntType",
								Markers:  map[string]string{"+Bar": "123", "+Foo": "true"},
								TypeInfo: &parse.TypeInfo{TypeName: "int", ExternalTypeName: "int"},
							},
							"OfAType": {
								Name:    "OfAType",
								Markers: map[string]string{"+Bar": "123", "+Foo": "true"},
								TypeInfo: &parse.TypeInfo{
									TypeName:         "IntType",
									ExternalTypeName: "typing.IntType",
									IsType:           true,
									TypeOf:           &parse.TypeInfo{TypeName: "int", ExternalTypeName: "int"},
								},
							},
							"SliceType": {
								Name:    "SliceType",
								Markers: map[string]string{"+Bar": "123", "+Foo": "true"},
								TypeInfo: &parse.TypeInfo{
									TypeName:         "[]AStruct",
									ExternalTypeName: "[]typing.AStruct",
									IsSlice:          true,
									Slice:            &parse.TypeInfo{TypeName: "AStruct", ExternalTypeName: "typing.AStruct", IsStruct: true, IsType: true, TypeOf: &parse.TypeInfo{IsStruct: true}},
								},
							},
						},
						Aliases: map[string]*parse.AliasTypeInfo{
							"AliasStringType": {
								Name:     "AliasStringType",
								Markers:  map[string]string{"+Bar": "123", "+Foo": "true"},
								TypeInfo: &parse.TypeInfo{TypeName: "string", ExternalTypeName: "string"},
							},
						},
					},
				},
			},
		},
		{
			name: "heredoc markers",
			args: args{
				path: "./_testdata/heredoc",
			},
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"heredoc": {
						Name:         "heredoc",
						Constants:    map[string]*parse.ConstantInfo{},
						DefinedTypes: map[string]*parse.DefinedTypeInfo{},
						Aliases:      map[string]*parse.AliasTypeInfo{},
						Structs:      map[string]*parse.StructInfo{},
						Interfaces:   map[string]*parse.InterfaceInfo{},
						Vars:         map[string]*parse.VarInfo{},
						Functions: map[string]*parse.FuncInfo{
							"Example": {
								Name: "Example",
								Markers: map[string]string{
									"+test:word2": "value",
									"+test:word":  "This is the use for blah blah blah blah",
								},
								FuncDefInfo: &parse.FuncDefInfo{
									Params:  []*parse.ParamInfo{},
									Results: []*parse.ResultInfo{},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "method name collisions across structs",
			args: args{
				path: "./_testdata/methodcollision",
			},
			want: &parse.Results{
				Packages: map[string]*parse.PackageInfo{
					"methodcollision": {
						Name:         "methodcollision",
						Constants:    map[string]*parse.ConstantInfo{},
						DefinedTypes: map[string]*parse.DefinedTypeInfo{},
						Aliases:      map[string]*parse.AliasTypeInfo{},
						Interfaces:   map[string]*parse.InterfaceInfo{},
						Vars:         map[string]*parse.VarInfo{},
						Structs: map[string]*parse.StructInfo{
							"A": {
								Name:           "A",
								Markers:        map[string]string{},
								Fields:         map[string]*parse.FieldInfo{},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods: map[string]*parse.FuncInfo{
									"Validate": {
										Name:    "Validate",
										Markers: map[string]string{},
										FuncDefInfo: &parse.FuncDefInfo{
											Params:  []*parse.ParamInfo{},
											Results: []*parse.ResultInfo{},
										},
										HasReciver:  true,
										ReciverName: "A",
									},
								},
							},
							"B": {
								Name:           "B",
								Markers:        map[string]string{},
								Fields:         map[string]*parse.FieldInfo{},
								EmbeddedFields: map[string]parse.EmbeddedFieldInfo{},
								Methods: map[string]*parse.FuncInfo{
									"Validate": {
										Name:    "Validate",
										Markers: map[string]string{},
										FuncDefInfo: &parse.FuncDefInfo{
											Params:  []*parse.ParamInfo{},
											Results: []*parse.ResultInfo{},
										},
										HasReciver:  true,
										ReciverName: "B",
									},
								},
							},
						},
						Functions: map[string]*parse.FuncInfo{
							"A.Validate": {
								Name:    "Validate",
								Markers: map[string]string{},
								FuncDefInfo: &parse.FuncDefInfo{
									Params:  []*parse.ParamInfo{},
									Results: []*parse.ResultInfo{},
								},
								HasReciver:  true,
								ReciverName: "A",
							},
							"B.Validate": {
								Name:    "Validate",
								Markers: map[string]string{},
								FuncDefInfo: &parse.FuncDefInfo{
									Params:  []*parse.ParamInfo{},
									Results: []*parse.ResultInfo{},
								},
								HasReciver:  true,
								ReciverName: "B",
							},
							"Validate": {
								Name:    "Validate",
								Markers: map[string]string{},
								FuncDefInfo: &parse.FuncDefInfo{
									Params:  []*parse.ParamInfo{},
									Results: []*parse.ResultInfo{},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parse.Parser{}
			wd, err := os.Getwd()
			if err != nil {
				t.Fatalf("failed to get working directory: %v", err)
			}

			path := filepath.Join(wd, tt.args.path)

			got, err := p.ParseDirectory(parse.Options{Path: path, SkipFilesWithContentsRegex: tt.args.skipFiles})
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

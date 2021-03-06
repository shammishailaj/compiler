package compiler

import (
	"fmt"
	"github.com/gqlc/graphql/ast"
	"github.com/gqlc/graphql/parser"
	"github.com/gqlc/graphql/token"
	"strings"
	"testing"
)

func TestValue(t *testing.T) {

	testCases := []struct {
		Name         string
		CName        string
		C            interface{}
		Val, ValType interface{}
		Items        map[string]*ast.TypeDecl
		Errs         []string
	}{
		{
			Name:    "Basic:Int",
			CName:   "intField",
			Val:     &ast.BasicLit{Kind: int64(token.INT), Value: "2"},
			ValType: &ast.Ident{Name: "Int"},
		},
		{
			Name:    "Basic:Float:AsInt",
			CName:   "intAsFloatField",
			Val:     &ast.BasicLit{Kind: int64(token.INT), Value: "2"},
			ValType: &ast.Ident{Name: "Float"},
		},
		{
			Name:    "Basic:Float:AsFloat",
			CName:   "floatAsFloatField",
			Val:     &ast.BasicLit{Kind: int64(token.FLOAT), Value: "2.0"},
			ValType: &ast.Ident{Name: "Float"},
		},
		{
			Name:    "Basic:String",
			CName:   "stringField",
			ValType: &ast.BasicLit{Kind: int64(token.STRING), Value: `"hello"`},
		},
		{
			Name:    "Basic:Boolean",
			CName:   "boolField",
			Val:     &ast.BasicLit{Kind: int64(token.BOOL), Value: "true"},
			ValType: &ast.Ident{Name: "Boolean"},
		},
		{
			Name:    "Basic:ID:String",
			CName:   "stringAsIDField",
			Val:     &ast.BasicLit{Kind: int64(token.STRING), Value: `"erbgoayueboguyvabef"`},
			ValType: &ast.Ident{Name: "ID"},
		},
		{
			Name:    "Basic:ID:Int",
			CName:   "intAsIDField",
			Val:     &ast.BasicLit{Kind: int64(token.INT), Value: "2"},
			ValType: &ast.Ident{Name: "ID"},
		},
		{
			Name:    "Basic:Ident:EnumValue",
			CName:   "enumValue",
			Val:     &ast.BasicLit{Kind: int64(token.IDENT), Value: "ONE"},
			ValType: &ast.Ident{Name: "Test"},
			Items: map[string]*ast.TypeDecl{
				"Test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.TypeSpec_Enum{Enum: &ast.EnumType{
							Values: &ast.FieldList{
								List: []*ast.Field{
									{Name: &ast.Ident{Name: "ONE"}},
								},
							},
						}},
					}},
				},
			},
		},
		{
			Name:    "Basic:Ident:UnknownEnumValue",
			CName:   "unknownEnumValue",
			Val:     &ast.BasicLit{Kind: int64(token.IDENT), Value: "TWO"},
			ValType: &ast.Ident{Name: "Test"},
			Items: map[string]*ast.TypeDecl{
				"Test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.TypeSpec_Enum{Enum: &ast.EnumType{
							Values: &ast.FieldList{
								List: []*ast.Field{
									{Name: &ast.Ident{Name: "ONE"}},
								},
							},
						}},
					}},
				},
			},
			Errs: []string{
				fmt.Sprintf("%s:%s: enum: %s has no value named: %s", "Basic:Ident:UnknownEnumValue", "unknownEnumValue", "Test", "TWO"),
			},
		},
		{
			Name:    "Basic:InvalidInputValue",
			CName:   "invalidInputValueForField",
			Val:     &ast.BasicLit{Kind: int64(token.INT), Value: "2"},
			ValType: &ast.Ident{Name: "String"},
			Errs: []string{
				fmt.Sprintf("%s:%s: %s is not coercible to: %s", "Basic:InvalidInputValue", "invalidInputValueForField", token.INT, "String"),
			},
		},
		{
			Name:  "Composite",
			CName: "inputObject",
			Val: &ast.CompositeLit{Value: &ast.CompositeLit_ObjLit{ObjLit: &ast.ObjLit{
				Fields: []*ast.ObjLit_Pair{
					{Key: &ast.Ident{Name: "a"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}}},
					{Key: &ast.Ident{Name: "b"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.FLOAT), Value: "2.0"}}}},
					{Key: &ast.Ident{Name: "c"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.STRING), Value: `"2"`}}}},
				},
			}}},
			ValType: &ast.Ident{Name: "Test"},
			Items: map[string]*ast.TypeDecl{
				"Test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.TypeSpec_Input{Input: &ast.InputType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Name: &ast.Ident{Name: "a"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Int"}}},
									{Name: &ast.Ident{Name: "b"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Float"}}},
									{Name: &ast.Ident{Name: "c"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "String"}}},
								},
							},
						}},
					}},
				},
			},
		},
		{
			Name:    "Composite:NotAnObjectLit",
			CName:   "notAnObjectLit",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_ListLit{}},
			ValType: &ast.Ident{},
			Errs: []string{
				fmt.Sprintf("%s:%s: input object must be provided", "Composite:NotAnObjLit", "notAnObjectLit"),
			},
		},
		{
			Name:    "Composite:UndefinedObject",
			CName:   "undefinedObject",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_ObjLit{}},
			ValType: &ast.Ident{Name: "Test"},
			Items:   map[string]*ast.TypeDecl{},
			Errs: []string{
				fmt.Sprintf("%s:%s: undefined input object: %s", "Composite:UndefinedObject", "undefinedObject", "Test"),
			},
		},
		{
			Name:    "Composite:OnlyExtensionProvided",
			CName:   "onlyExtensionProvided",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_ObjLit{}},
			ValType: &ast.Ident{Name: "Test"},
			Items:   map[string]*ast.TypeDecl{"Test": {Spec: &ast.TypeDecl_TypeExtSpec{}}},
			Errs: []string{
				fmt.Sprintf("%s:%s: could not find type spec for input object: %s", "Composite:OnlyExtensionProvided", "onlyExtensionProvided", "Test"),
			},
		},
		{
			Name:    "Composite:ExpectedValueNotAnInputObject",
			CName:   "expectedValueNotAnInputObject",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_ObjLit{}},
			ValType: &ast.Ident{Name: "Test"},
			Items: map[string]*ast.TypeDecl{
				"Test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.TypeSpec_Scalar{},
					}},
				},
			},
			Errs: []string{
				fmt.Sprintf("%s:%s: %s is not an input object", "Composite:ExpectedValueNotAnInputObject", "expectedValueNotAnInputObject", "Test"),
			},
		},
		{
			Name:  "Composite:UndefinedField",
			CName: "undefinedField",
			Val: &ast.CompositeLit{Value: &ast.CompositeLit_ObjLit{ObjLit: &ast.ObjLit{
				Fields: []*ast.ObjLit_Pair{
					{Key: &ast.Ident{Name: "a"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}}},
					{Key: &ast.Ident{Name: "b"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.FLOAT), Value: "2.0"}}}},
					{Key: &ast.Ident{Name: "c"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.STRING), Value: `"2"`}}}},
					{Key: &ast.Ident{Name: "d"}},
				},
			}}},
			ValType: &ast.Ident{Name: "Test"},
			Items: map[string]*ast.TypeDecl{
				"Test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.TypeSpec_Input{Input: &ast.InputType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Name: &ast.Ident{Name: "a"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Int"}}},
									{Name: &ast.Ident{Name: "b"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Float"}}},
									{Name: &ast.Ident{Name: "c"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "String"}}},
								},
							},
						}},
					}},
				},
			},
			Errs: []string{
				fmt.Sprintf("%s:%s: undefined field: %s", "Composite:UndefinedField", "undefinedField", "d"),
			},
		},
		{
			Name:  "Composite:NonUniqueField",
			CName: "nonUniqueField",
			Val: &ast.CompositeLit{Value: &ast.CompositeLit_ObjLit{ObjLit: &ast.ObjLit{
				Fields: []*ast.ObjLit_Pair{
					{Key: &ast.Ident{Name: "a"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}}},
					{Key: &ast.Ident{Name: "b"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.FLOAT), Value: "2.0"}}}},
					{Key: &ast.Ident{Name: "c"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.STRING), Value: `"2"`}}}},
					{Key: &ast.Ident{Name: "a"}},
				},
			}}},
			ValType: &ast.Ident{Name: "Test"},
			Items: map[string]*ast.TypeDecl{
				"Test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.TypeSpec_Input{Input: &ast.InputType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Name: &ast.Ident{Name: "a"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Int"}}},
									{Name: &ast.Ident{Name: "b"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Float"}}},
									{Name: &ast.Ident{Name: "c"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "String"}}},
								},
							},
						}},
					}},
				},
			},
			Errs: []string{
				fmt.Sprintf("%s:%s: field must be unique: %s", "Composite:NonUniqueField", "nonUniqueField", "a"),
			},
		},
		{
			Name:  "Composite:MissingRequiredField",
			CName: "missingRequiredField",
			Val: &ast.CompositeLit{Value: &ast.CompositeLit_ObjLit{ObjLit: &ast.ObjLit{
				Fields: []*ast.ObjLit_Pair{
					{Key: &ast.Ident{Name: "a"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}}},
					{Key: &ast.Ident{Name: "b"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.FLOAT), Value: "2.0"}}}},
					{Key: &ast.Ident{Name: "c"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.STRING), Value: `"2"`}}}},
				},
			}}},
			ValType: &ast.Ident{Name: "Test"},
			Items: map[string]*ast.TypeDecl{
				"Test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.TypeSpec_Input{Input: &ast.InputType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Name: &ast.Ident{Name: "a"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Int"}}},
									{Name: &ast.Ident{Name: "b"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Float"}}},
									{Name: &ast.Ident{Name: "c"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "String"}}},
									{Name: &ast.Ident{Name: "d"}, Type: &ast.Field_NonNull{NonNull: &ast.NonNull{}}},
								},
							},
						}},
					}},
				},
			},
			Errs: []string{
				fmt.Sprintf("%s: non-null field must be present in: %s", "d", "Composite:MissingRequiredField"),
			},
		},
		{
			Name:    "Composite:Basic:Int",
			CName:   "intField",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}},
			ValType: &ast.Ident{Name: "Int"},
		},
		{
			Name:    "Composite:Basic:Float:AsInt",
			CName:   "intAsFloatField",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}},
			ValType: &ast.Ident{Name: "Float"},
		},
		{
			Name:    "Composite:Basic:Float:AsFloat",
			CName:   "floatAsFloatField",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.FLOAT), Value: "2.0"}}},
			ValType: &ast.Ident{Name: "Float"},
		},
		{
			Name:    "Composite:Basic:String",
			CName:   "stringField",
			ValType: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.STRING), Value: `"hello"`}}},
		},
		{
			Name:    "Composite:Basic:Boolean",
			CName:   "boolField",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.BOOL), Value: "true"}}},
			ValType: &ast.Ident{Name: "Boolean"},
		},
		{
			Name:    "Composite:Basic:ID:String",
			CName:   "stringAsIDField",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.STRING), Value: `"erbgoayueboguyvabef"`}}},
			ValType: &ast.Ident{Name: "ID"},
		},
		{
			Name:    "Composite:Basic:ID:Int",
			CName:   "intAsIDField",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}},
			ValType: &ast.Ident{Name: "ID"},
		},
		{
			Name:    "Composite:Basic:InvalidInputValue",
			CName:   "invalidInputValueForField",
			Val:     &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}},
			ValType: &ast.Ident{Name: "String"},
			Errs: []string{
				fmt.Sprintf("%s:%s: %s is not coercible to: %s", "Basic:InvalidInputValue", "invalidInputValueForField", token.INT, "String"),
			},
		},
		{
			Name:    "List:BasicLitAsList",
			CName:   "basicLitAsList",
			C:       &ast.Arg{Value: &ast.Arg_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}},
			Val:     &ast.BasicLit{Kind: int64(token.INT), Value: "2"},
			ValType: &ast.List{Type: &ast.List_Ident{Ident: &ast.Ident{Name: "Int"}}},
		},
		{
			Name:    "List:BasicLitAsListAsList",
			CName:   "basicLitAsListAsList",
			C:       &ast.Arg{Value: &ast.Arg_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}},
			Val:     &ast.BasicLit{Kind: int64(token.INT), Value: "2"},
			ValType: &ast.List{Type: &ast.List_List{List: &ast.List{Type: &ast.List_Ident{Ident: &ast.Ident{Name: "Int"}}}}},
		},
		{
			Name:  "List:ObjLitAsList",
			CName: "objLitAsList",
			C: &ast.Arg{Value: &ast.Arg_CompositeLit{CompositeLit: &ast.CompositeLit{
				Value: &ast.CompositeLit_ObjLit{
					ObjLit: &ast.ObjLit{
						Fields: []*ast.ObjLit_Pair{
							{Key: &ast.Ident{Name: "a"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}}},
							{Key: &ast.Ident{Name: "b"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.FLOAT), Value: "2.0"}}}},
							{Key: &ast.Ident{Name: "c"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.STRING), Value: `"2"`}}}},
						},
					},
				},
			}}},
			Val: &ast.CompositeLit{Value: &ast.CompositeLit_ObjLit{ObjLit: &ast.ObjLit{
				Fields: []*ast.ObjLit_Pair{
					{Key: &ast.Ident{Name: "a"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.INT), Value: "2"}}}},
					{Key: &ast.Ident{Name: "b"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.FLOAT), Value: "2.0"}}}},
					{Key: &ast.Ident{Name: "c"}, Val: &ast.CompositeLit{Value: &ast.CompositeLit_BasicLit{BasicLit: &ast.BasicLit{Kind: int64(token.STRING), Value: `"2"`}}}},
				},
			}}},
			ValType: &ast.List{Type: &ast.List_Ident{Ident: &ast.Ident{Name: "Test"}}},
			Items: map[string]*ast.TypeDecl{
				"Test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.TypeSpec_Input{Input: &ast.InputType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Name: &ast.Ident{Name: "a"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Int"}}},
									{Name: &ast.Ident{Name: "b"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "Float"}}},
									{Name: &ast.Ident{Name: "c"}, Type: &ast.Field_Ident{Ident: &ast.Ident{Name: "String"}}},
								},
							},
						}},
					}},
				},
			},
		},
		{
			Name:    "NonNull:InvalidValue",
			CName:   "invalidValue",
			Val:     &ast.BasicLit{Kind: int64(token.NULL), Value: "null"},
			ValType: &ast.NonNull{Type: &ast.NonNull_Ident{Ident: &ast.Ident{Name: "String"}}},
			Errs: []string{
				fmt.Sprintf("%s:%s: non-null arg cannot be the null value", "NonNull:InvalidValue", "invalidValue"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(subT *testing.T) {
			var errs []*TypeError
			validateValue(testCase.Name, testCase.CName, testCase.C, testCase.Val, testCase.ValType, testCase.Items, &errs)

			var count int
			for _, terr := range errs {
				for _, serr := range testCase.Errs {
					if terr.Msg == serr {
						count++
					}
				}
			}

			if count != len(errs) && len(errs) != len(testCase.Errs) {
				subT.Fail()
				return
			}
		})
	}
}

func TestDirectives(t *testing.T) {

	testCases := []struct {
		Name  string
		Dirs  []*ast.DirectiveLit
		Loc   ast.DirectiveLocation_Loc
		Items map[string]*ast.TypeDecl
		Errs  []string
	}{
		{
			Name: "Undefined",
			Dirs: []*ast.DirectiveLit{{Name: "asfadfbdfba"}},
			Errs: []string{
				fmt.Sprintf("%s: undefined directive", "asfadfbdfba"),
			},
		},
		{
			Name: "InvalidLocation",
			Dirs: []*ast.DirectiveLit{{Name: "test"}},
			Loc:  ast.DirectiveLocation_FIELD,
			Items: map[string]*ast.TypeDecl{
				"test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "test"},
						Type: &ast.TypeSpec_Directive{Directive: &ast.DirectiveType{
							Locs: []*ast.DirectiveLocation{{Loc: ast.DirectiveLocation_NoPos}},
						}},
					}},
				},
			},
			Errs: []string{
				fmt.Sprintf("%s: invalid location for directive: %s", "test", ast.DirectiveLocation_FIELD),
			},
		},
		{
			Name: "MustBeUnique",
			Dirs: []*ast.DirectiveLit{{Name: "test"}, {Name: "test"}},
			Loc:  ast.DirectiveLocation_FIELD,
			Items: map[string]*ast.TypeDecl{
				"test": {
					Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{
						Name: &ast.Ident{Name: "test"},
						Type: &ast.TypeSpec_Directive{Directive: &ast.DirectiveType{
							Locs: []*ast.DirectiveLocation{{Loc: ast.DirectiveLocation_FIELD}},
						}},
					}},
				},
			},
			Errs: []string{
				fmt.Sprintf("%s: directive cannot be applied more than once per location: %s", "test", ast.DirectiveLocation_FIELD),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(subT *testing.T) {
			var errs []*TypeError
			validateDirectives(testCase.Dirs, testCase.Loc, testCase.Items, &errs)

			var count int
			for _, terr := range errs {
				for _, serr := range testCase.Errs {
					if terr.Msg == serr {
						count++
					}
				}
			}

			if count != len(errs) && len(errs) != len(testCase.Errs) {
				subT.Fail()
				return
			}
		})
	}
}

func TestCompareTypes(t *testing.T) {
	items := map[string]*ast.TypeDecl{
		"TestInterface": {
			Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{Name: &ast.Ident{Name: "TestInterface"}, Type: &ast.TypeSpec_Interface{}}},
		},
		"TestUnion": {
			Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{Name: &ast.Ident{Name: "TestUnion"}, Type: &ast.TypeSpec_Union{
				Union: &ast.UnionType{
					Members: []*ast.Ident{{Name: "TestObjA"}, {Name: "TestObjB"}},
				},
			}}},
		},
		"TestObjA": {
			Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{Name: &ast.Ident{Name: "TestObjA"}, Type: &ast.TypeSpec_Object{
				Object: &ast.ObjectType{
					Interfaces: []*ast.Ident{{Name: "TestInterface"}},
				},
			}}},
		},
		"TestObjB": {
			Spec: &ast.TypeDecl_TypeSpec{TypeSpec: &ast.TypeSpec{Name: &ast.Ident{Name: "TestObjB"}, Type: &ast.TypeSpec_Object{}}},
		},
	}

	testCases := []struct {
		Name     string
		A, B     interface{}
		Expected bool
	}{
		{
			Name:     "Eq",
			A:        &ast.Ident{Name: "Test"},
			B:        &ast.Ident{Name: "Test"},
			Expected: true,
		},
		{
			Name:     "Interface",
			A:        &ast.Ident{Name: "TestObjA"},
			B:        &ast.Ident{Name: "TestInterface"},
			Expected: true,
		},
		{
			Name:     "Union",
			A:        &ast.Ident{Name: "TestObjA"},
			B:        &ast.Ident{Name: "TestUnion"},
			Expected: true,
		},
		{
			Name:     "List",
			A:        &ast.List{Type: &ast.List_Ident{Ident: &ast.Ident{Name: "Test"}}},
			B:        &ast.List{Type: &ast.List_Ident{Ident: &ast.Ident{Name: "Test"}}},
			Expected: true,
		},
		{
			Name:     "NonNull",
			A:        &ast.NonNull{Type: &ast.NonNull_Ident{Ident: &ast.Ident{Name: "Test"}}},
			B:        &ast.NonNull{Type: &ast.NonNull_Ident{Ident: &ast.Ident{Name: "Test"}}},
			Expected: true,
		},
		{
			Name:     "Fail",
			A:        &ast.Ident{Name: "TestObjA"},
			B:        &ast.Ident{Name: "TestObjB"},
			Expected: false,
		},
		{
			Name:     "Fail:MisMatchTypes",
			A:        &ast.Ident{},
			B:        &ast.NonNull{},
			Expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(subT *testing.T) {
			ok := compareTypes(testCase.A, testCase.B, items)
			if ok != testCase.Expected {
				subT.Fail()
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		Name string
		Src  string
		Errs []string
	}{
		{
			Name: "Scalar",
			Src:  `scalar Int`,
		},
		{
			Name: "InvalidName",
			Src:  `scalar __Int`,
			Errs: []string{fmt.Sprintf("%s is an invalid name for type: %s", "__Int", token.SCALAR)},
		},
		{
			Name: "InvalidDirectives",
			Src: `scalar Int @asdads @a @b @b

directive @a on FIELD

directive @b on SCALAR`,
			Errs: []string{
				fmt.Sprintf("%s: undefined directive", "asdads"),
				fmt.Sprintf("%s: invalid location for directive: %s", "a", ast.DirectiveLocation_SCALAR),
				fmt.Sprintf("%s: directive cannot be applied more than once per location: %s", "b", ast.DirectiveLocation_SCALAR),
			},
		},
		{
			Name: "Enum",
			Src: `enum A {}

enum B {
	One
	Two
	Two
}`,
			Errs: []string{
				fmt.Sprintf("%s: enum type must define one or more unique enum values", "A"),
				fmt.Sprintf("%s:%s: enum value must be unique", "B", "Two"),
			},
		},
		{
			Name: "Union",
			Src: `union B = Undefined | String | String

scalar String`,
			Errs: []string{
				fmt.Sprintf("%s:%s: undefined type", "B", "Undefined"),
				fmt.Sprintf("%s:%s: member type must be an object type", "B", "String"),
				fmt.Sprintf("%s:%s: member type must be unique", "B", "String"),
			},
		},
		{
			Name: "Interface",
			Src: `interface A {}

interface B {
	__one: String
	__one: String
	two(__a: String, __a: String): Input
}

scalar String

directive @Input on FIELD`,
			Errs: []string{
				fmt.Sprintf("%s: interface type must one or more fields", "A"),
				fmt.Sprintf("%s:%s: field must be unique", "B", "__one"),
				fmt.Sprintf("%s:%s: field name cannot start with \"__\" (double underscore)", "B", "__one"),
				fmt.Sprintf("%s:%s: argument must be unique", "B:two", "__a"),
				fmt.Sprintf("%s:%s: argument name cannot start with \"__\" (double underscore)", "B:two", "__a"),
				fmt.Sprintf("%s:%s: field type must be a valid output type, not: %s", "B", "two", "Input"),
			},
		},
		{
			Name: "Input",
			Src: `input A {}

input B {
	one(a: String): String
}

scalar String`,
			Errs: []string{
				fmt.Sprintf("%s: input object type must define one or more input fields", "A"),
				fmt.Sprintf("%s:%s: input object fields cannot have arguments", "B", "one"),
			},
		},
		{
			Name: "Schema",
			Src: `schema {}

schema {
	mutation: Mutation
}

schema {
	query: String
}

scalar String`,
			Errs: []string{
				fmt.Sprintf("schema: at minimum query object must be provided"),
				fmt.Sprintf("schema: query object must be provided"),
				fmt.Sprintf("schema:%s: unknown type: %s", "mutation", "Mutation"),
				fmt.Sprintf("schema:%s: root operation return type must be an object type", "query"),
			},
		},
		{
			Name: "Object",
			Src: `type A {}

scalar Int

scalar String

type B implements One & String {
	a: String
}

union Thr = C | D

interface Two {
	id: String
	edges: [Two]
	u: Thr
}

interface Four {
	a(i: Int, s: String): Thr
	b(i: Int): String
}

type C implements Two & Four {
	id: String
	edges: [D]
	u: B
	a(s: Int, ni: Int!): D
	b: String
}

type D implements Two {
	id: String
	edges: [Two]
}`,
			Errs: []string{
				fmt.Sprintf("%s: an object type must define one or more fields", "A"),
				fmt.Sprintf("%s: undefined interface: %s", "B", "One"),
				fmt.Sprintf("%s:%s: non-interface type can not be used as interface", "B", "String"),
				fmt.Sprintf("%s:%s: object type must include field: %s", "D", "Two", "u"),
				fmt.Sprintf("%s:%s: object field type must be a sub-type of interface field type", "C", "u"),
				fmt.Sprintf("%s:%s: object field is missing interface field argument: %s", "C", "a", "i"),
				fmt.Sprintf("%s:%s:%s: object argument and interface argument must be the same type", "C", "a", "s"),
				fmt.Sprintf("%s:%s:%s: additional arguments to interface field implementation must be non-null", "C", "a", "ni"),
				fmt.Sprintf("%s:%s: object field must include the same argument definitions that the interface field has", "C", "b"),
			},
		},
		{
			Name: "Directive",
			Src: `directive @test(__a: Test @test) on ARGUMENT_DEFINITION

scalar String

interface Test {
	a: String
}`,
			Errs: []string{
				fmt.Sprintf("%s:%s: argument name cannot start with \"__\" (double underscore)", "test", "__a"),
				fmt.Sprintf("%s:%s: directive argument must be a valid input type, not: %s", "test", "__a", "Test"),
				fmt.Sprintf("%s:%s: directive argument cannont reference its own directive definition", "test", "__a"),
			},
		},
		{
			Name: "Extend:NoDefinitionFound",
			Src:  `extend scalar String`,
			Errs: []string{
				fmt.Sprintf("extend:%s: no definition found for this type", "String"),
			},
		},
		{
			Name: "Extend:Scalar",
			Src: `enum Test {
	A
}

extend scalar Test`,
			Errs: []string{
				fmt.Sprintf("extend:scalar:%s: original type defintion must be a scalar", "Test"),
			},
		},
		{
			Name: "Extend:Object",
			Src: `scalar String

extend type String

type Test {
	a: String
	b: String
}

interface A {
	a: String
}

extend type Test implements A & B & String {
	b: String
}`,
			Errs: []string{
				fmt.Sprintf("extend:object:%s: original type defintion must be a object", "String"),
				fmt.Sprintf("%s:%s: field defintion already exists in original object definition", "extend:object:Test", "b"),
				fmt.Sprintf("%s: undefined interface: %s", "extend:object:Test", "B"),
				fmt.Sprintf("%s:%s: non-interface type can not be used as interface", "extend:object:Test", "String"),
			},
		},
		{
			Name: "Extend:Interface",
			Src: `scalar String

extend interface String

interface Test {
	a: String
}

extend interface Test {
	a: String
}`,
			Errs: []string{
				fmt.Sprintf("extend:interface:%s: original type defintion must be a interface", "String"),
				fmt.Sprintf("%s:%s: field already exists in original interface definition", "extend:interface:Test", "a"),
			},
		},
		{
			Name: "Extend:Union",
			Src: `scalar String

type A {
	a: String
}

type B {
	a: String
}

union Test = A | B

extend union String

extend union Test = A`,
			Errs: []string{
				fmt.Sprintf("extend:union:%s: original type defintion must be a union", "String"),
				fmt.Sprintf("%s:%s: union member already exists in original union definition", "extend:union:Test", "A"),
			},
		},
		{
			Name: "Extend:Enum",
			Src: `scalar String

extend enum String

enum Test {
	A
}

extend enum Test {
	A
}`,
			Errs: []string{
				fmt.Sprintf("extend:enum:%s: original type defintion must be a enum", "String"),
				fmt.Sprintf("%s:%s: enum value already exists in original enum definition", "extend:enum:Test", "A"),
			},
		},
		{
			Name: "Extend:Input",
			Src: `scalar String

extend input String

input Test {
	a: String
}

extend input Test {
	a: String
}`,
			Errs: []string{
				fmt.Sprintf("extend:input:%s: original type defintion must be a input", "String"),
				fmt.Sprintf("%s:%s: field defintion already exists in original input definition", "extend:input:Test", "a"),
			},
		},
		{
			Name: "Extend:RepeatedDirectives",
			Src: `scalar String @a

directive @a on SCALAR

directive @b on SCALAR

extend scalar String @a @b`,
			Errs: []string{
				fmt.Sprintf("%s:%s: directive is already applied to original type definition", "String", "a"),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(subT *testing.T) {
			doc, err := parser.ParseDoc(token.NewDocSet(), testCase.Name, strings.NewReader(testCase.Src), 0)
			if err != nil {
				subT.Error(err)
				return
			}

			errs := Validate(doc)

			var count int
			for _, terr := range errs {
				for _, serr := range testCase.Errs {
					if terr.Msg == serr {
						count++
					}
				}
			}

			if count != len(testCase.Errs) || count != len(errs) || len(errs) != len(testCase.Errs) {
				fmt.Println("----------------------")
				for _, terr := range errs {
					fmt.Println("terr:", terr.Msg)
				}
				fmt.Println("----------------------")
				for _, serr := range testCase.Errs {
					fmt.Println("serr:", serr)
				}
				fmt.Println("----------------------")

				subT.Fail()
				return
			}
		})
	}
}

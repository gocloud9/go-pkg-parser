package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"path"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

type PackageInfo struct {
	Name         string
	Structs      map[string]*StructInfo
	Constants    map[string]*ConstantInfo
	Functions    map[string]*FuncInfo
	Interfaces   map[string]*InterfaceInfo
	Vars         map[string]*VarInfo
	DefinedTypes map[string]*DefinedTypeInfo
	Aliases      map[string]*AliasTypeInfo
}

type DefinedTypeInfo struct {
	Name    string
	Markers map[string]string
	*TypeInfo
}

type AliasTypeInfo struct {
	Name    string
	Markers map[string]string
	*TypeInfo
}

type ConstantInfo struct {
	Name     string
	TypeName string
	Markers  map[string]string
	Value    string
}

type VarInfo struct {
	Name    string
	Markers map[string]string
	*TypeInfo
}

type ImportedTypeInfo struct {
	TypeName            string
	ImportRaw           string
	PackagePath         string
	PackageDefaultAlias string
}

type TypeInfo struct {
	TypeName         string
	ExternalTypeName string
	PackageName      string
	IsPointer        bool
	IsMap            bool
	IsSlice          bool
	IsStruct         bool
	IsChan           bool
	IsFunc           bool
	IsInterface      bool
	IsEllipsis       bool
	IsImported       bool
	IsType           bool
	MapKey           *TypeInfo
	MapValue         *TypeInfo
	Slice            *TypeInfo
	Chan             *TypeInfo
	Func             *FuncDefInfo
	Interface        *InterfaceInfo
	Pointer          *TypeInfo
	Ellipsis         *TypeInfo
	ImportedType     *ImportedTypeInfo
	TypeOf           *TypeInfo
}
type FieldInfo struct {
	Name    string
	Markers map[string]string
	Tags    map[string][]string
	*TypeInfo
}

type InterfaceInfo struct {
	Name          string
	Markers       map[string]string
	Methods       map[string]*FuncInfo
	EmbeddedTypes map[string]*EmbeddedTypeInfo
}

type FuncDefInfo struct {
	IsVariadic bool
	IsReceiver bool
	Receiver   string
	Params     []*ParamInfo
	Results    []*ResultInfo
}

type FuncInfo struct {
	Name        string
	Markers     map[string]string
	HasReciver  bool
	ReciverName string
	*FuncDefInfo
}

type ParamInfo struct {
	*TypeInfo
	Name       string
	IsVariadic bool
}

type ResultInfo struct {
	*TypeInfo
	Name string
}

type EmbeddedFieldInfo struct {
	Name     string
	TypeName string
	Markers  map[string]string
	Tags     map[string][]string
}

type EmbeddedTypeInfo struct {
	Name     string
	TypeName string
	Markers  map[string]string
}

type StructInfo struct {
	Name           string
	Markers        map[string]string
	Fields         map[string]*FieldInfo
	Methods        map[string]*FuncInfo
	EmbeddedFields map[string]EmbeddedFieldInfo
}

type Results struct {
	Packages map[string]*PackageInfo
}

type Parser struct {
}

func exprToTypeInfo(e ast.Expr, fileCache fileCachedData) *TypeInfo {
	varInfo := &TypeInfo{}

	switch expr := e.(type) {
	case *ast.MapType:
		varInfo.IsMap = true
		varInfo.MapKey = exprToTypeInfo(expr.Key, fileCache)
		varInfo.MapValue = exprToTypeInfo(expr.Value, fileCache)
		varInfo.TypeName = fmt.Sprintf("map[%s]%s", varInfo.MapKey.TypeName, varInfo.MapValue.TypeName)
		varInfo.ExternalTypeName = fmt.Sprintf("map[%s]%s", varInfo.MapKey.ExternalTypeName, varInfo.MapValue.ExternalTypeName)
	case *ast.SliceExpr:
		varInfo.IsSlice = true
		varInfo.Slice = exprToTypeInfo(expr.X, fileCache)
		varInfo.TypeName = "[]" + varInfo.Slice.TypeName
		varInfo.ExternalTypeName = "[]" + varInfo.Slice.ExternalTypeName
	case *ast.StarExpr:
		varInfo.IsPointer = true
		varInfo.Pointer = exprToTypeInfo(expr.X, fileCache)
		varInfo.TypeName = "*" + varInfo.Pointer.TypeName
		varInfo.ExternalTypeName = "*" + varInfo.Pointer.ExternalTypeName
	case *ast.ArrayType:
		varInfo.IsSlice = true
		varInfo.Slice = exprToTypeInfo(expr.Elt, fileCache)
		varInfo.TypeName = "[]" + varInfo.Slice.TypeName
		varInfo.ExternalTypeName = "[]" + varInfo.Slice.ExternalTypeName
	case *ast.ChanType:
		varInfo.IsChan = true
		varInfo.Chan = exprToTypeInfo(expr.Value, fileCache)
		varInfo.TypeName = "chan " + varInfo.Chan.TypeName
		varInfo.ExternalTypeName = "chan " + varInfo.Chan.ExternalTypeName
	case *ast.StructType:
		varInfo.IsStruct = true

	case *ast.InterfaceType:
		varInfo.IsInterface = true
	case *ast.Ellipsis:
		varInfo.IsEllipsis = true
		varInfo.Ellipsis = exprToTypeInfo(expr.Elt, fileCache)
		varInfo.TypeName = "..." + varInfo.Ellipsis.TypeName
		varInfo.ExternalTypeName = "..." + varInfo.Ellipsis.ExternalTypeName
	case *ast.FuncType:
		varInfo.IsFunc = true
		params := fieldListToParamInfoList(expr.Params, fileCache)
		results := fieldListToResultInfoList(expr.Results, fileCache)

		varInfo.TypeName = funcTypeNameFromParamsAndResults(params, results, false)
		varInfo.ExternalTypeName = funcTypeNameFromParamsAndResults(params, results, true)

		varInfo.Func = &FuncDefInfo{
			IsVariadic: isVariadicFunc(params),
			Params:     params,
			Results:    results,
		}

	case *ast.SelectorExpr:
		varInfo = exprToTypeInfo(expr.Sel, fileCache)
		imported, ok := fileCache.imports[expr.X.(*ast.Ident).Name]
		if ok {
			packagePath := strings.ReplaceAll(imported.Path.Value, "\"", "")
			parts := strings.Split(packagePath, "/")
			varInfo.IsImported = true
			varInfo.ImportedType = &ImportedTypeInfo{
				TypeName:            varInfo.TypeName,
				ImportRaw:           imported.Path.Value,
				PackagePath:         packagePath,
				PackageDefaultAlias: parts[len(parts)-1],
			}
			varInfo.TypeName = varInfo.ImportedType.PackageDefaultAlias + "." + varInfo.TypeName
			varInfo.ExternalTypeName = varInfo.TypeName
		}

	case *ast.Ident:
		if expr.Obj != nil && expr.Obj.Decl != nil {
			switch s := expr.Obj.Decl.(type) {
			case *ast.TypeSpec:
				varInfo = exprToTypeInfo(s.Type, fileCache)
				varInfo.TypeName = s.Name.Name
				varInfo.ExternalTypeName = fileCache.packageName + "." + varInfo.TypeName

			case *ast.Field:
				varInfo = exprToTypeInfo(s.Type, fileCache)
				if varInfo.IsStruct {
					varInfo.ExternalTypeName = fileCache.packageName + "." + varInfo.TypeName
				}
			}
		}
		if expr.Obj == nil {
			varInfo.TypeName = expr.Name
			varInfo.ExternalTypeName = varInfo.TypeName
		} else {
			if expr.Obj.Decl != nil {
				t, ok := expr.Obj.Decl.(*ast.TypeSpec)
				if ok {
					varInfo.TypeOf = exprToTypeInfo(t.Type, fileCache)
					varInfo.IsType = true
				}
			}
		}
	}

	return varInfo
}

func handleValueSpec(node *ast.ValueSpec, pi *PackageInfo, fileCache fileCachedData) {
	if len(node.Names) != 1 {
		return
	}
	name := node.Names[0]

	for _, val := range node.Values {
		switch v := val.(type) {
		case *ast.BasicLit:
			vs := name.Obj.Decl.(*ast.ValueSpec)
			if name.Obj.Kind == ast.Con {
				typeName := strings.ToLower(v.Kind.String())
				vsType, ok := vs.Type.(*ast.Ident)
				if ok && vsType != nil && vsType.Name != "" {
					typeName = vsType.Name
				}
				pi.Constants[name.Obj.Name] = &ConstantInfo{
					Name:     name.Obj.Name,
					TypeName: typeName,
					Markers:  markerValues(vs.Doc),
					Value:    v.Value,
				}
			} else {
				pi.Vars[name.Obj.Name] = &VarInfo{
					Name:     name.Obj.Name,
					Markers:  markerValues(vs.Doc),
					TypeInfo: exprToTypeInfo(v, fileCache),
				}
			}

		case *ast.FuncLit:
			doc, foundDoc := fileCache.commentGroups[node.Pos()-5] // Magic number to get comment group before type name
			if foundDoc && node.Doc == nil {
				node.Doc = doc
			}
			params := fieldListToParamInfoList(v.Type.Params, fileCache)
			results := fieldListToResultInfoList(v.Type.Results, fileCache)

			pi.Vars[name.Obj.Name] = &VarInfo{
				Name:    name.Obj.Name,
				Markers: markerValues(node.Doc),
				TypeInfo: &TypeInfo{
					TypeName:         funcTypeNameFromParamsAndResults(params, results, false),
					ExternalTypeName: funcTypeNameFromParamsAndResults(params, results, true),
					IsFunc:           true,
					Func: &FuncDefInfo{
						IsVariadic: isVariadicFunc(params),
						Params:     params,
						Results:    results,
					},
				},
			}
		}
	}
}

func funcTypeNameFromParamsAndResults(params []*ParamInfo, results []*ResultInfo, external bool) string {
	paramTypes := []string{}
	for _, p := range params {
		if external {
			paramTypes = append(paramTypes, p.ExternalTypeName)
		} else {
			paramTypes = append(paramTypes, p.TypeName)
		}
	}
	resultTypes := []string{}
	for _, r := range results {
		if external {
			resultTypes = append(resultTypes, r.ExternalTypeName)
		} else {
			resultTypes = append(resultTypes, r.TypeName)
		}
	}

	resultsPrefix := ""
	if len(resultTypes) == 1 {
		resultsPrefix = " "
	} else if len(resultTypes) > 1 {
		resultsPrefix = " ("
	}
	resultsSuffix := ""
	if len(resultTypes) > 1 {
		resultsSuffix = ")"
	}

	return fmt.Sprintf("func(%s)%s%s%s", strings.Join(paramTypes, ", "), resultsPrefix, strings.Join(resultTypes, ", "), resultsSuffix)
}

func handleFuncDecl(node *ast.FuncDecl, pi *PackageInfo, fileCache fileCachedData) {
	params := fieldListToParamInfoList(node.Type.Params, fileCache)
	results := fieldListToResultInfoList(node.Type.Results, fileCache)

	receiverTypeName := ""
	if node.Recv != nil && len(node.Recv.List) > 0 {
		receiverType, ok := node.Recv.List[0].Type.(*ast.Ident)
		if ok {
			receiverTypeName = receiverType.Name
		}
	}

	pi.Functions[node.Name.Name] = &FuncInfo{
		Name:        node.Name.Name,
		HasReciver:  receiverTypeName != "",
		ReciverName: receiverTypeName,

		Markers: markerValues(node.Doc),
		FuncDefInfo: &FuncDefInfo{
			IsVariadic: isVariadicFunc(params),
			Params:     params,
			Results:    results,
		},
	}

	for i := range pi.Structs {
		if receiverTypeName == pi.Structs[i].Name {
			pi.Structs[i].Methods[node.Name.Name] = pi.Functions[node.Name.Name]
		}
	}
}

type fileCachedData struct {
	commentGroups map[token.Pos]*ast.CommentGroup
	imports       map[string]*ast.ImportSpec
	packageName   string
}

type Options struct {
	Path                       string
	SkipFilesWithContentsRegex []*regexp.Regexp
	IncludeEmptyPackages       bool
}

func (p *Parser) ParseDirectory(opts Options) (*Results, error) {
	results := &Results{
		Packages: map[string]*PackageInfo{},
	}

	cfg := &packages.Config{
		Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedDeps | packages.NeedFiles | packages.NeedName |
			packages.NeedTypesInfo | packages.NeedModule | packages.NeedEmbedFiles | packages.NeedEmbedPatterns |
			packages.NeedTarget | packages.NeedCompiledGoFiles | packages.NeedExportFile | packages.NeedImports |
			packages.NeedForTest,
		Dir: opts.Path,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Fatal(err)
	}

	for _, pkg := range pkgs {
		pi := newPackageInfo(pkg.Name)

		for _, file := range pkg.Syntax {
			inspectFile(file, pi, opts.SkipFilesWithContentsRegex)
		}

		results.Packages[pi.Name] = pi
	}

	return results, err
}

// ParseEmbed parses Go source files from an fs.FS (e.g. an embed.FS populated
// via //go:embed) and returns the same *Results structure as ParseDirectory.
// opts.Path selects the root directory within the FS; an empty string defaults
// to ".". opts.SkipFilesWithContentsRegex is applied identically to ParseDirectory.
func (p *Parser) ParseEmbed(fsys fs.FS, opts Options) (*Results, error) {
	root := opts.Path
	if root == "" {
		root = "."
	}

	results := &Results{
		Packages: map[string]*PackageInfo{},
	}

	// pkgFiles groups parsed AST files by package name, mirroring how
	// packages.Load groups by pkg.Name in ParseDirectory.
	pkgFiles := map[string][]*ast.File{}
	fset := token.NewFileSet()

	err := fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !isGoSourceFile(p) {
			return nil
		}

		data, readErr := fs.ReadFile(fsys, p)
		if readErr != nil {
			return fmt.Errorf("failed to read %s: %w", p, readErr)
		}

		file, parseErr := parser.ParseFile(fset, p, data, parser.ParseComments)
		if parseErr != nil {
			return fmt.Errorf("failed to parse %s: %w", p, parseErr)
		}

		pkgName := file.Name.Name
		pkgFiles[pkgName] = append(pkgFiles[pkgName], file)

		return nil
	})
	if err != nil {
		return nil, err
	}

	for pkgName, files := range pkgFiles {
		pi := newPackageInfo(pkgName)

		for _, file := range files {
			inspectFile(file, pi, opts.SkipFilesWithContentsRegex)
		}

		results.Packages[pi.Name] = pi
	}

	return results, nil
}

// isGoSourceFile reports whether a path names a Go source file that should be
// parsed. Test files (ending in _test.go) are excluded to match the default
// non-test behaviour of ParseDirectory.
func isGoSourceFile(p string) bool {
	base := path.Base(p)

	return strings.HasSuffix(base, ".go") && !strings.HasSuffix(base, "_test.go")
}

// newPackageInfo allocates a PackageInfo with all maps initialised.
func newPackageInfo(name string) *PackageInfo {
	return &PackageInfo{
		Name:         name,
		Structs:      map[string]*StructInfo{},
		Functions:    map[string]*FuncInfo{},
		Interfaces:   map[string]*InterfaceInfo{},
		Vars:         map[string]*VarInfo{},
		Constants:    map[string]*ConstantInfo{},
		DefinedTypes: map[string]*DefinedTypeInfo{},
		Aliases:      map[string]*AliasTypeInfo{},
	}
}

// inspectFile runs the AST dispatch over a single parsed file, populating pi.
// It returns without modifying pi when any SkipFilesWithContentsRegex matches a
// comment in the file.
func inspectFile(file *ast.File, pi *PackageInfo, skipRegexps []*regexp.Regexp) {
	fileCache := fileCachedData{
		commentGroups: map[token.Pos]*ast.CommentGroup{},
		imports:       map[string]*ast.ImportSpec{},
		packageName:   pi.Name,
	}

	for i := range file.Comments {
		for j := range skipRegexps {
			if skipRegexps[j].MatchString(file.Comments[i].Text()) {
				return
			}
		}
	}

	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.ValueSpec:
			handleValueSpec(node, pi, fileCache)
		case *ast.FuncDecl:
			handleFuncDecl(node, pi, fileCache)
		case *ast.CommentGroup:
			fileCache.commentGroups[node.End()] = node
		case *ast.ImportSpec:
			if node.Name == nil {
				fileCache.imports[strings.ReplaceAll(node.Path.Value, "\"", "")] = node
			} else {
				fileCache.imports[node.Name.Name] = node
			}
		case *ast.TypeSpec:
			handleTypeSpec(node, pi, fileCache)
		default:
		}

		return true
	})
}

func handleInterfaceType(ts *ast.TypeSpec, node *ast.InterfaceType, pi *PackageInfo, fileCache fileCachedData) {
	doc, foundDoc := fileCache.commentGroups[ts.Name.Pos()-6] // Magic number to get comment group before type name
	if foundDoc && ts.Doc == nil {
		ts.Doc = doc
	}

	ii := &InterfaceInfo{
		Name:          ts.Name.Name,
		Markers:       markerValues(ts.Doc),
		Methods:       map[string]*FuncInfo{},
		EmbeddedTypes: map[string]*EmbeddedTypeInfo{},
	}
	for i, m := range node.Methods.List {
		if len(node.Methods.List[i].Names) == 0 {
			eti := &EmbeddedTypeInfo{
				Name:     m.Type.(*ast.Ident).Name,
				TypeName: m.Type.(*ast.Ident).Name,
				Markers:  markerValues(m.Doc),
			}

			ii.EmbeddedTypes[eti.TypeName] = eti
		} else {
			params := fieldListToParamInfoList(m.Type.(*ast.FuncType).Params, fileCache)
			results := fieldListToResultInfoList(m.Type.(*ast.FuncType).Results, fileCache)

			funcName := node.Methods.List[i].Names[0].Name
			ii.Methods[funcName] = &FuncInfo{
				Name:    funcName,
				Markers: markerValues(node.Methods.List[i].Doc),

				FuncDefInfo: &FuncDefInfo{
					IsVariadic: isVariadicFunc(params),
					Params:     params,
					Results:    results,
				},
			}
		}
	}

	pi.Interfaces[ii.Name] = ii
}

func isVariadicFunc(params []*ParamInfo) bool {
	if len(params) == 0 {
		return false
	}

	lastParam := params[len(params)-1]
	return lastParam.IsEllipsis
}

func fieldListToParamInfoList(params *ast.FieldList, fileCache fileCachedData) []*ParamInfo {
	ps := []*ParamInfo{}

	if params == nil {
		return ps
	}

	for _, param := range params.List {
		for _, paramNameProperties := range param.Names {
			ps = append(ps, &ParamInfo{
				Name:       paramNameProperties.Name,
				TypeInfo:   exprToTypeInfo(paramNameProperties, fileCache),
				IsVariadic: false,
			})
		}
		if param.Type != nil && len(param.Names) == 0 {
			ps = append(ps, &ParamInfo{
				TypeInfo: exprToTypeInfo(param.Type, fileCache),
			})
		}
	}

	return ps
}

func fieldListToResultInfoList(results *ast.FieldList, fileCache fileCachedData) []*ResultInfo {
	rs := []*ResultInfo{}

	if results == nil {
		return rs
	}

	for _, result := range results.List {
		for _, paramNameProperties := range result.Names {
			rs = append(rs, &ResultInfo{
				Name:     paramNameProperties.Name,
				TypeInfo: exprToTypeInfo(paramNameProperties, fileCache),
			})
		}
		if result.Type != nil && len(result.Names) == 0 {
			rs = append(rs, &ResultInfo{
				TypeInfo: exprToTypeInfo(result.Type, fileCache),
			})
		}
	}

	return rs
}

func handleTypeSpec(ts *ast.TypeSpec, pi *PackageInfo, fileCache fileCachedData) {
	switch t := ts.Type.(type) {
	case *ast.StructType:
		handleStructType(ts, t, pi, fileCache)
	case *ast.InterfaceType:
		handleInterfaceType(ts, t, pi, fileCache)
	default:
		if ts.Doc == nil {
			ts.Doc, _ = fileCache.commentGroups[ts.Pos()-6] // Magic number to get comment group before type name
		}

		if ts.Assign == token.NoPos {
			dti := &DefinedTypeInfo{
				Name:     ts.Name.Name,
				Markers:  markerValues(ts.Doc),
				TypeInfo: exprToTypeInfo(t, fileCache),
			}

			pi.DefinedTypes[dti.Name] = dti
		} else {
			ati := &AliasTypeInfo{
				Name:     ts.Name.Name,
				Markers:  markerValues(ts.Doc),
				TypeInfo: exprToTypeInfo(t, fileCache),
			}

			pi.Aliases[ati.Name] = ati
		}

	}
}

func handleStructType(t *ast.TypeSpec, s *ast.StructType, pi *PackageInfo, fileCache fileCachedData) {
	si := &StructInfo{
		Name: t.Name.Name,
	}

	doc, foundDoc := fileCache.commentGroups[t.Name.Pos()-6] // Magic number to get comment group before type name
	if foundDoc && t.Doc == nil {
		t.Doc = doc
	}

	si.Markers = markerValues(t.Doc)
	si.Fields = map[string]*FieldInfo{}
	si.Methods = map[string]*FuncInfo{}
	si.EmbeddedFields = map[string]EmbeddedFieldInfo{}

	for _, f := range s.Fields.List {
		if f.Doc == nil {
			f.Doc, _ = fileCache.commentGroups[f.Pos()]
		}

		if len(f.Names) == 0 {
			ti := exprToTypeInfo(f.Type, fileCache)
			name := embeddedFieldName(f.Type)
			efi := EmbeddedFieldInfo{
				Name:     name,
				TypeName: ti.TypeName,
				Markers:  markerValues(f.Doc),
				Tags:     parseTags(f.Tag),
			}

			si.EmbeddedFields[efi.TypeName] = efi
		} else {
			fi := &FieldInfo{
				Name:     f.Names[0].Name,
				TypeInfo: exprToTypeInfo(f.Type, fileCache),
				Tags:     parseTags(f.Tag),
				Markers:  markerValues(f.Doc),
			}

			si.Fields[fi.Name] = fi
		}
	}

	for i := range pi.Functions {
		if pi.Functions[i].HasReciver && pi.Functions[i].ReciverName == si.Name {

			si.Methods[pi.Functions[i].Name] = pi.Functions[i]
		}
	}

	pi.Structs[si.Name] = si
}

// embeddedFieldName returns the field name introduced by an embedded field.
// For an embedded type the field name is the unqualified type name, e.g.
// `network.SubVpcReservation` is accessed as `SubVpcReservation`, and a pointer
// embed `*Foo` as `Foo`. It handles plain identifiers, selector expressions
// (qualified types from another package), and pointers to either.
func embeddedFieldName(e ast.Expr) string {
	switch expr := e.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.SelectorExpr:
		return expr.Sel.Name
	case *ast.StarExpr:
		return embeddedFieldName(expr.X)
	default:
		return ""
	}
}

func parseTags(tagLit *ast.BasicLit) map[string][]string {
	if tagLit == nil || tagLit.Kind != token.STRING {
		return map[string][]string{}
	}

	tags := map[string][]string{}

	raw := strings.ReplaceAll(tagLit.Value, "`", "")
	tagGroups := strings.Split(raw, " ")
	for _, group := range tagGroups {
		parts := strings.SplitN(group, ":", 2)
		if len(parts) != 2 {
			continue
		}
		tags[parts[0]] = strings.Split(strings.ReplaceAll(parts[1], "\"", ""), ",")
	}

	_ = reflect.StructTag(raw)

	return tags
}

// heredocDelimiter opens and closes a multi-line marker value. A marker whose
// value is exactly this delimiter begins a heredoc that spans the following
// comment lines until a line consisting solely of the delimiter closes it.
const heredocDelimiter = "`"

func markerValues(cg *ast.CommentGroup) map[string]string {
	if cg == nil {
		return map[string]string{}
	}

	values := map[string]string{}

	for i := 0; i < len(cg.List); i++ {
		txt := cleanCommentText(cg.List[i].Text)
		parts := strings.SplitN(txt, "=", 2)
		key := strings.TrimSpace(parts[0])

		if len(parts) != 2 {
			values[key] = ""

			continue
		}

		value := strings.TrimSpace(parts[1])
		if value != heredocDelimiter {
			values[key] = value

			continue
		}

		body, consumed := collectHeredoc(cg.List[i+1:])
		values[key] = body
		i += consumed
	}

	return values
}

// collectHeredoc concatenates comment lines up to the closing delimiter and
// returns the joined body along with the number of comment lines consumed
// (including the closing delimiter line). An unterminated heredoc consumes the
// remaining lines.
func collectHeredoc(lines []*ast.Comment) (string, int) {
	body := make([]string, 0, len(lines))

	for idx, c := range lines {
		txt := cleanCommentText(c.Text)
		if txt == heredocDelimiter {
			return strings.Join(body, "\n"), idx + 1
		}
		body = append(body, txt)
	}

	return strings.Join(body, "\n"), len(lines)
}

// cleanCommentText strips the surrounding comment markers from a single raw
// comment line and trims surrounding whitespace.
func cleanCommentText(raw string) string {
	txt := strings.TrimSpace(strings.TrimPrefix(raw, "//"))
	txt = strings.TrimSpace(strings.TrimPrefix(txt, "/*"))
	txt = strings.TrimSpace(strings.TrimSuffix(txt, "*/"))

	return txt
}

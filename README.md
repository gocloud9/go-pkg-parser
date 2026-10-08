# go-pkg-parser

A powerful parser that makes Go's built-in parser easier to use. It wraps
`golang.org/x/tools/go/packages` to extract structured type, function,
interface, constant, variable, and alias information from a Go module directory.

## Install

```sh
go get github.com/gocloud9/go-pkg-parser
```

## Import

```go
import "github.com/gocloud9/go-pkg-parser/pkg/parse"
```

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/gocloud9/go-pkg-parser/pkg/parse"
)

func main() {
    p := parse.Parser{}

    results, err := p.ParseDirectory(parse.Options{
        Path: "./",
    })
    if err != nil {
        log.Fatal(err)
    }

    for pkgPath, pkg := range results.Packages {
        fmt.Printf("package %s (%s)\n", pkg.Name, pkgPath)
        for _, s := range pkg.Structs {
            fmt.Printf("  struct %s — %d fields\n", s.Name, len(s.Fields))
        }
    }
}
```

### Options

| Field | Type | Description |
|---|---|---|
| `Path` | `string` | Directory to parse (passed to `packages.Load("./...")`) |
| `SkipFilesWithContentsRegex` | `[]*regexp.Regexp` | Skip any file whose contents match one of these patterns |
| `IncludeEmptyPackages` | `bool` | Include packages with no exported declarations |

### Results

`Results.Packages` is a `map[string]*PackageInfo` keyed by import path. Each
`PackageInfo` exposes:

- `Structs` — `[]*StructInfo` (fields, tags, embedded types)
- `Functions` — `[]*FuncInfo` (receiver, params, results, variadic flag)
- `Interfaces` — `[]*InterfaceInfo` (methods)
- `Constants` — `[]*ConstantInfo`
- `Vars` — `[]*VarInfo`
- `DefinedTypes` — `[]*DefinedTypeInfo`
- `AliasTypes` — `[]*AliasTypeInfo`

## License

MIT — see [LICENSE](LICENSE).

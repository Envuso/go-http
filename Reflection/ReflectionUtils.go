package Reflection

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// FuncPathAndName Get the name and path of a func
func FuncPathAndName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// FuncName Get the name of a func (with package path)
func FuncName(f interface{}) string {
	splitFuncName := strings.Split(FuncPathAndName(f), ".")
	return splitFuncName[len(splitFuncName)-1]
}

// FuncDescription Get description of a func
func FuncDescription(f interface{}) string {
	fileName, _ := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).FileLine(0)
	funcName := FuncName(f)
	fset := token.NewFileSet()

	// Parse src
	parsedAst, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	pkg := &ast.Package{
		Name:  "Any",
		Files: make(map[string]*ast.File),
	}
	pkg.Files[fileName] = parsedAst

	importPath, _ := filepath.Abs("/")
	myDoc := doc.New(pkg, importPath, doc.AllDecls)
	for _, theFunc := range myDoc.Funcs {
		if theFunc.Name == funcName {
			return theFunc.Doc
		}
	}
	return ""
}

func GetStructName(structType reflect.Type) string {
	if structType.Kind() == reflect.Ptr {
		return structType.Elem().Name()
	}
	return structType.Name()
}

func StructHasMethod(structType reflect.Type, method string) (methodType reflect.Method, ok bool) {
	indirect := IndirectType(structType)

	if indirect.Kind() != reflect.Struct {
		log.Printf("structType passed to StructHasMethod is not a struct.")
		return
	}

	// My fucking brain hurts

	// If we have a pointer version, check via that, since we could have pointer receiver...
	if structType.Kind() == reflect.Ptr {
		methodType, ok = structType.MethodByName(method)
		if ok {
			return
		}
	}

	// Check if we have it on the raw type of the struct...
	methodType, ok = indirect.MethodByName(method)
	if ok {
		return
	}

	// Last ditch hope, maybe we passed a non pointer version, but our methods require the pointer receiver?
	structType = reflect.PointerTo(structType)
	methodType, ok = structType.MethodByName(method)

	return
}

func IndirectType(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return typ.Elem()
	}
	return typ
}

func PkgPathOfType(elemTyp reflect.Type) string {
	typName := elemTyp.Name()
	pkgPath := elemTyp.PkgPath()

	return pkgPath[strings.LastIndexByte(pkgPath, '/')+1:] + "." + typName
}

func NameOf(v interface{}) string {
	elemTyp := IndirectType(reflect.ValueOf(v).Type())
	return PkgPathOfType(elemTyp)
}

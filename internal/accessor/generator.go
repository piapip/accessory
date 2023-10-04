package accessor

import (
	"bytes"
	"fmt"
	"go/types"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/go/packages"
)

type generator struct {
	writer   *writer
	typ      string
	output   string
	receiver string
	lock     string
}

type methodGenParameters struct {
	Receiver     string
	Struct       string
	Field        string
	GetterMethod string
	SetterMethod string
	Type         string
	ZeroValue    string // used only when generating getter
	EmptyValue   string // used only when generating stuff for tester
	Lock         string
}

type testGenParameters struct {
	Struct        string
	Package       string
	WantStruct    string
	AssertTest    string
	NilTestData   string
	EmptyTestData string
}

func newGenerator(fs afero.Fs, pkg *Package, options ...Option) *generator {
	g := new(generator)
	for _, opt := range options {
		opt(g)
	}

	path := g.outputFilePath(pkg.Dir)
	g.writer = newWriter(fs, path)

	return g
}

// Generate generates a file and accessor methods.
func Generate(fs afero.Fs, pkg *Package, options ...Option) error {
	g := newGenerator(fs, pkg, options...)

	fmt.Printf("starting pkg\n")
	for _, field := range pkg.Structs[0].Fields {
		fmt.Printf("field: %+v\n", field)
		if field.Tag != nil {
			fmt.Printf("getter: %s\n", *(field.Tag.Getter))
		}
	}

	accessors := make([]string, 0)
	usedPkgs := make([]string, 0, len(pkg.Imports))

	testParameters := &testGenParameters{
		Struct: pkg.Structs[0].Name,
		// this package is hard coded.
		Package: "models",
	}

	// for iStruct, st := range pkg.Structs {
	for _, st := range pkg.Structs {
		if st.Name != g.typ {
			continue
		}

		// for iField, field := range st.Fields {
		for _, field := range st.Fields {
			if field.Tag == nil {
				continue
			}

			params := g.setupParameters(pkg, st, field)

			// // I only generate getters for one struct at the time
			// // so only check this once
			// if iStruct == 0 && iField == 0 {
			// 	testParameters.Struct = params.Struct
			// }

			if field.Tag.Getter != nil {
				getter, err := g.generateGetter(params)
				if err != nil {
					return err
				}
				accessors = append(accessors, getter)
			}
			if field.Tag.Setter != nil {
				setter, err := g.generateSetter(params)
				if err != nil {
					return err
				}
				accessors = append(accessors, setter)
			}

			err := g.updateTestComponent(params, testParameters)
			if err != nil {
				return err
			}

			replacer := strings.NewReplacer(
				"[]", "", // trim []
				"*", "", // trim *
			)
			replaced := replacer.Replace(params.Type)
			if typePaths := strings.Split(replaced, "."); len(typePaths) > 1 {
				usedPkgs = append(usedPkgs, typePaths[0])
			}
		}
	}

	generatedTest, err := g.assembleTest(testParameters)
	if err != nil {
		return err
	}

	accessors = append(accessors, generatedTest)

	imports := g.generateImportStrings(pkg.Imports, usedPkgs)
	return g.writer.write(pkg.Name, imports, accessors)
}

func (g *generator) outputFilePath(dir string) string {
	output := g.output
	if output == "" {
		// Use snake_case name of type as output file if output file is not specified.
		// type TestStruct will be test_struct_accessor.go
		var firstCapMatcher = regexp.MustCompile("(.)([A-Z][a-z]+)")
		var articleCapMatcher = regexp.MustCompile("([a-z0-9])([A-Z])")

		name := firstCapMatcher.ReplaceAllString(g.typ, "${1}_${2}")
		name = articleCapMatcher.ReplaceAllString(name, "${1}_${2}")
		output = strings.ToLower(fmt.Sprintf("%s_accessor.go", name))
	}

	return filepath.Join(dir, output)
}

func (g *generator) generateSetter(
	params *methodGenParameters,
) (string, error) {
	var lockingCode string
	if params.Lock != "" {
		lockingCode = ` {{.Receiver}}.{{.Lock}}.Lock()
		defer {{.Receiver}}.{{.Lock}}.Unlock()
		`
	}

	var tpl = `
	func ({{.Receiver}} *{{.Struct}}) {{.SetterMethod}}(val {{.Type}}) {
		if {{.Receiver}} == nil {
			return
		}
	` +
		lockingCode + // inject locing code
		`{{.Receiver}}.{{.Field}} = val
	}`

	t := template.Must(template.New("setter").Parse(tpl))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) generateGetter(
	params *methodGenParameters,
) (string, error) {
	var lockingCode string
	if params.Lock != "" {
		lockingCode = `{{.Receiver}}.{{.Lock}}.Lock()
		defer {{.Receiver}}.{{.Lock}}.Unlock()
		`
	}

	var getterTemplate = `
	// {{.GetterMethod}} returns the {{.Struct}}'s {{.Field}}.
	func ({{.Receiver}} *{{.Struct}}) {{.GetterMethod}}() {{.Type}} {
		if {{.Receiver}} == nil {
			return {{.ZeroValue}}
		}

		` +
		lockingCode + // inject locing code
		`return {{.Receiver}}.{{.Field}}
	}`

	t := template.Must(template.New("getter").Parse(getterTemplate))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// updateTestComponent fills in the content using for testing.
func (g *generator) updateTestComponent(
	params *methodGenParameters,
	testParameters *testGenParameters,
) error {
	var (
		wantStructTemplate = `want{{.Field}} {{.Type}}`
		assertTemplate     = `got{{.Field}} := ctx.testData.args.Get{{.Field}}()
			assert.Equal(t, ctx.testData.want{{.Field}}, got{{.Field}})
		`
		nilTestDataTemplate   = `want{{.Field}}: {{.ZeroValue}},`
		emptyTestDataTemplate = `want{{.Field}}: {{.EmptyValue}},`
	)

	wantStructTemplateExecutor := template.Must(template.New("wantStruct").Parse(wantStructTemplate))
	bufWantStruct := new(bytes.Buffer)

	if err := wantStructTemplateExecutor.Execute(bufWantStruct, params); err != nil {
		return err
	}

	assertTemplateExecutor := template.Must(template.New("assert").Parse(assertTemplate))
	bufAssert := new(bytes.Buffer)

	if err := assertTemplateExecutor.Execute(bufAssert, params); err != nil {
		return err
	}

	nilTestDataTemplateExecutor := template.Must(template.New("nilTest").Parse(nilTestDataTemplate))
	bufNilTestData := new(bytes.Buffer)

	if err := nilTestDataTemplateExecutor.Execute(bufNilTestData, params); err != nil {
		return err
	}

	emptyTestDataExecutor := template.Must(template.New("emptyTest").Parse(emptyTestDataTemplate))
	bufEmptyTestData := new(bytes.Buffer)

	if err := emptyTestDataExecutor.Execute(bufEmptyTestData, params); err != nil {
		return err
	}

	if testParameters.WantStruct == "" {
		testParameters.WantStruct = testParameters.WantStruct + bufWantStruct.String()
		testParameters.AssertTest = testParameters.AssertTest + bufAssert.String()
		testParameters.NilTestData = testParameters.NilTestData + bufNilTestData.String()
		testParameters.EmptyTestData = testParameters.EmptyTestData + bufEmptyTestData.String()
	} else {
		testParameters.WantStruct = testParameters.WantStruct + "\n" + bufWantStruct.String()
		testParameters.AssertTest = testParameters.AssertTest + "\n" + bufAssert.String()
		testParameters.NilTestData = testParameters.NilTestData + "\n" + bufNilTestData.String()
		testParameters.EmptyTestData = testParameters.EmptyTestData + "\n" + bufEmptyTestData.String()
	}

	return nil
}

func (g *generator) assembleTest(
	params *testGenParameters,
) (string, error) {
	var getTestTemplate = `
	func Test{{.Struct}}_GetFunctions(t *testing.T) {
		type want struct {
			args *{{.Package}}.{{.Struct}}
			{{.WantStruct}}
			wantProto *replaceMe.{{.Struct}}
		}
	
		type Context struct {
			testData *want
		}
	
		contextInitiateFunction := func(t *testing.T) *Context {
			return &Context{}
		}

		gt.Begin(t,
			contextInitiateFunction,
			gt.Run("Get functions return proper value", func(t *testing.T, ctx *Context) {
				// GET functions
				{{.AssertTest}}

				// Convert from models to Proto.
				gotProto := ctx.testData.args.ToProto()
				assert.Equal(t, ctx.testData.wantProto, gotProto)

				// Then convert from Proto back to model
				gotModel := models.ProtoToRetentionAvoid(gotProto)
				assert.Equal(t, ctx.testData.args, gotModel)
			}).
				Using("given nil value", func(t *testing.T, ctx *Context) {
					ctx.testData = &want{
						args: nil,
						{{.NilTestData}}
						wantProto: nil,
					}
				}).
				Using("given empty value", func(t *testing.T, ctx *Context) {
					ctx.testData = &want{
						args: &{{.Package}}.{{.Struct}}{},
						{{.EmptyTestData}}
						wantProto: &replaceMe.{{.Struct}}{},
					}
				}).
				Using("given NON nil value", func(t *testing.T, ctx *Context) {
					ctx.testData = &want{

					}
				}),
		)
	}`

	t := template.Must(template.New("tester").Parse(getTestTemplate))
	buf := new(bytes.Buffer)

	if err := t.Execute(buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (g *generator) setupParameters(
	pkg *Package,
	st *Struct,
	field *Field,
) *methodGenParameters {
	typeName := g.typeName(pkg.Types, field.Type)
	fmt.Println("pkg.Types: ", pkg.Types)
	fmt.Println("field.Type: ", field.Type)
	fmt.Println("typeName: ", typeName)
	getter, setter := g.methodNames(field)
	return &methodGenParameters{
		Receiver:     g.receiverName(st.Name),
		Struct:       st.Name,
		Field:        field.Name,
		GetterMethod: getter,
		SetterMethod: setter,
		Type:         typeName,
		ZeroValue:    g.zeroValue(field.Type, typeName),
		EmptyValue:   g.emptyValue(field.Type, typeName),
		Lock:         g.lock,
	}
}

func (g *generator) receiverName(structName string) string {
	if g.receiver != "" {
		// Do nothing if receiver name specified in args.
		return g.receiver
	}

	// Use the first letter of struct as receiver if receiver name is not specified.
	return strings.ToLower(string(structName[0]))
}

func (g *generator) methodNames(field *Field) (getter, setter string) {
	if getterName := field.Tag.Getter; getterName != nil && *getterName != "" {
		getter = *getterName
	} else {
		getter = "Get" + cases.Title(language.Und, cases.NoLower).String(field.Name)
	}

	if setterName := field.Tag.Setter; setterName != nil && *setterName != "" {
		setter = *setterName
	} else {
		setter = "Set" + cases.Title(language.Und, cases.NoLower).String(field.Name)
	}

	return getter, setter
}

func (g *generator) typeName(pkg *types.Package, t types.Type) string {
	return types.TypeString(t, func(p *types.Package) string {
		fmt.Println("pkg: ", pkg)
		fmt.Println("p: ", p)
		// type is defined in the same package
		if pkg == p {
			return ""
		}
		// path string(like example.com/user/project/package) into slice
		return p.Name()
	})
}

func (g *generator) zeroValue(t types.Type, typeString string) string {
	switch t := t.(type) {
	case *types.Pointer:
		return "nil"
	case *types.Array:
		return "nil"
	case *types.Slice:
		return "nil"
	case *types.Chan:
		return "nil"
	case *types.Interface:
		return "nil"
	case *types.Map:
		return "nil"
	case *types.Signature:
		return "nil"
	case *types.Struct:
		return "nil"
	case *types.Basic:
		info := types.Typ[t.Kind()].Info()
		switch {
		case types.IsNumeric&info != 0:
			return "0"
		case types.IsBoolean&info != 0:
			return "false"
		case types.IsString&info != 0:
			return `""`
		}
	case *types.Named:
		if types.Identical(t, types.Universe.Lookup("error").Type()) {
			return "nil"
		}

		return g.zeroValue(t.Underlying(), typeString)
	}

	return "nil"
}

func (g *generator) emptyValue(t types.Type, typeString string) string {
	switch t := t.(type) {
	case *types.Pointer:
		return "nil"
	case *types.Array:
		return "nil"
	case *types.Slice:
		return "nil"
	case *types.Chan:
		return "nil"
	case *types.Interface:
		return "nil"
	case *types.Map:
		return "nil"
	case *types.Signature:
		return "nil"
	case *types.Struct:
		return "&" + typeString + "{}"
	case *types.Basic:
		info := types.Typ[t.Kind()].Info()
		switch {
		case types.IsNumeric&info != 0:
			return "0"
		case types.IsBoolean&info != 0:
			return "false"
		case types.IsString&info != 0:
			return `""`
		}
	case *types.Named:
		if types.Identical(t, types.Universe.Lookup("error").Type()) {
			return "nil"
		}

		return g.zeroValue(t.Underlying(), typeString)
	}

	return "nil"
}

func (g *generator) generateImportStrings(
	pkgs map[string]*packages.Package,
	usedPkgs []string,
) []string {
	usedMap := make(map[string]struct{}, 0)
	for i := range usedPkgs {
		usedMap[usedPkgs[i]] = struct{}{}
	}

	imports := make([]string, 0, len(usedMap))
	for _, pkg := range pkgs {
		if _, ok := usedMap[pkg.Name]; ok {
			imports = append(imports, pkg.PkgPath)
		}
	}
	sort.Strings(imports)

	return imports
}

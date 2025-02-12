package genaep

import (
	"fmt"
	"strings"

	"github.com/blaberg/aep-go/resourcepath"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

const (
	stringsPackage     = protogen.GoImportPath("strings")
	fmtPackage         = protogen.GoImportPath("fmt")
	resourcepathImport = protogen.GoImportPath("github.com/blaberg/aep-go/resourcepath")
)

func generateResourcePath(_ *protogen.Plugin, g *protogen.GeneratedFile, file *protogen.File) error {
	for _, m := range file.Messages {
		resource := proto.GetExtension(
			m.Desc.Options(), annotations.E_Resource,
		).(*annotations.ResourceDescriptor)
		if resource == nil || len(resource.GetPattern()) == 0 {
			continue
		}
		hasMultipattern := len(resource.Pattern) > 1
		generators := make([]PathGenerator, 0, len(resource.Pattern))
		for _, pattern := range resource.Pattern {
			pg := PathGenerator{
				Name:    m.GoIdent.GoName,
				Pattern: pattern,
			}
			if hasMultipattern {
				pg = newPatternGenerator(pattern)
			}
			generators = append(generators, pg)
			g.Unskip()
			g.P("")
			g.P("type ", pg.Name, "ResourcePath struct{")
			g.P("  path *", resourcepathImport.Ident("ResourcePath"))
			g.P("}")
			g.P("")
			if err := pg.generateParseStringMethod(g); err != nil {
				return err
			}
			if err := pg.generateNewResourcePathMethod(g); err != nil {
				return err
			}
			if err := pg.generateStringMethod(g); err != nil {
				return err
			}
			if err := pg.generateGetterMethods(g); err != nil {
				return err
			}
			g.P("")
			g.P("func (*", pg.Name, "ResourcePath) isMultipattern() {}")
			g.P("")
		}
		if hasMultipattern {
			g.P("")
			g.P("type isMultipattern interface {")
			g.P("    isMultipattern()")
			g.P("}")
			g.P("")
			g.P("func ParseMultipattern(p string) (isMultipattern, error) {")
			g.P("    switch {")
			for _, gen := range generators {
				g.P("        case ", resourcepathImport.Ident("Matches"), "(\"", gen.Pattern, "\", p):")
				g.P("                 return Parse", gen.Name, "ResourcePath(p)")
			}
			g.P("    }")
			g.P("  return nil, ", fmtPackage.Ident("Errorf"), "(\"failed to match pattern\")")
			g.P("}")
		}
	}
	return nil
}

type PathGenerator struct {
	Pattern string
	Name    string
}

func newPatternGenerator(pattern string) PathGenerator {
	var scanner resourcepath.Scanner
	scanner.Init(pattern)
	var name string
	for scanner.Scan() {
		if !scanner.Segment().IsVariable() {
			continue
		}
		literal := scanner.Segment().Literal().ResourceID()
		c := cases.Title(language.English)
		name = fmt.Sprintf("%s%s", name, c.String(literal))
	}
	return PathGenerator{
		Pattern: pattern,
		Name:    name,
	}
}

func (p PathGenerator) generateParseStringMethod(g *protogen.GeneratedFile) error {
	g.P("func Parse", p.Name, "ResourcePath(p string) (*", p.Name, "ResourcePath, error) {")
	g.P("  path, err := ", resourcepathImport.Ident("ParseString"), "(p, \"", p.Pattern, "\")")
	g.P("  if err != nil {")
	g.P("    return nil, err")
	g.P("  }")
	g.P("  return &", p.Name, "ResourcePath{")
	g.P("    path: path,")
	g.P("  }, nil")
	g.P("}")
	g.P("")
	return nil
}

func (p PathGenerator) generateStringMethod(g *protogen.GeneratedFile) error {
	g.P("func(p *", p.Name, "ResourcePath) String() string {")
	g.P("  return ", stringsPackage.Ident("Join"), "(")
	g.P("    []string{")
	var scanner resourcepath.Scanner
	scanner.Init(p.Pattern)
	for scanner.Scan() {
		if !scanner.Segment().IsVariable() {
			g.P("      \"", scanner.Segment().Literal().ResourceID(), "\",")
			continue
		}
		g.P("      p.path.Get(\"", scanner.Segment().Literal(), "\"),")
	}
	g.P("    },")
	g.P("    \"/\",")
	g.P("  )")
	g.P("}")
	g.P("")
	return nil
}

func (p PathGenerator) generateNewResourcePathMethod(g *protogen.GeneratedFile) error {
	var scanner resourcepath.Scanner
	scanner.Init(p.Pattern)
	g.P("func New", p.Name, "Path(")
	for scanner.Scan() {
		if !scanner.Segment().IsVariable() {
			continue
		}
		g.P("  ", scanner.Segment().Literal().ResourceID(), " string,")
	}
	g.P("  ) *", p.Name, "ResourcePath {")
	g.P("    segments := map[string]string{")
	scanner.Init(p.Pattern)
	for scanner.Scan() {
		if !scanner.Segment().IsVariable() {
			continue
		}
		id := scanner.Segment().Literal().ResourceID()
		g.P("    \"", id, "\": ", id, ",")
	}
	g.P("    }")
	g.P("  return &", p.Name, "ResourcePath{")
	g.P("    path: ", resourcepathImport.Ident("NewResourcePath"), "(segments),")
	g.P("  }")
	g.P("}")
	g.P("")
	return nil
}

func (p PathGenerator) generateGetterMethods(g *protogen.GeneratedFile) error {
	var scanner resourcepath.Scanner
	scanner.Init(p.Pattern)
	for scanner.Scan() {
		if !scanner.Segment().IsVariable() {
			continue
		}
		literal := scanner.Segment().Literal().ResourceID()
		g.P("func(p *", p.Name, "ResourcePath) Get", strings.ToUpper(literal[:1])+literal[1:], "() string {")
		g.P("  return p.path.Get(\"", literal, "\")")
		g.P("}")
		g.P("")
	}
	return nil
}

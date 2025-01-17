package genaep

import (
	"fmt"
	"strings"

	"github.com/blaberg/aep-go/resourcepath"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func generateResourcePath(_ *protogen.Plugin, g *protogen.GeneratedFile, file *protogen.File) error {
	for _, m := range file.Messages {
		resource := proto.GetExtension(
			m.Desc.Options(), annotations.E_Resource,
		).(*annotations.ResourceDescriptor)
		if resource == nil || len(resource.GetPattern()) == 0 {
			continue
		}
		if len(resource.Pattern) > 1 {
			return fmt.Errorf("generator does not support multipatterns yet")
		}
		g.Unskip()
		resourcepathImport := protogen.GoImportPath("github.com/blaberg/aep-go/resourcepath")
		g.Import(resourcepathImport)
		g.P("")
		g.P("type ", m.GoIdent.GoName, "ResourcePath struct{")
		g.P("  path *", resourcepathImport.Ident("ResourcePath"))
		g.P("}")
		g.P("")
		g.P("func Parse", m.GoIdent.GoName, "ResourcePath(p string) (*", m.GoIdent.GoName, "ResourcePath, error) {")
		g.P("  path, err := ", resourcepathImport.Ident("ParseString"), "(p, \"", resource.GetPattern()[0], "\")")
		g.P("  if err != nil {")
		g.P("    return nil, err")
		g.P("  }")
		g.P("  return &", m.GoIdent.GoName, "ResourcePath{")
		g.P("    path: path,")
		g.P("  }, nil")
		g.P("}")
		var scanner resourcepath.Scanner
		scanner.Init(resource.GetPattern()[0])
		for scanner.Scan() {
			if !scanner.Segment().IsVariable() {
				continue
			}
			literal := scanner.Segment().Literal().ResourceID()
			g.P("func(p *", m.GoIdent.GoName, "ResourcePath) Get", strings.ToUpper(literal[:1])+literal[1:], "() string {")
			g.P("  return p.path.Get(\"", literal, "\")")
			g.P("}")
		}

	}
	return nil
}

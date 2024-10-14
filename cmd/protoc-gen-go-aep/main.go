package main

import (
	"github.com/blaberg/aep-go/cmd/protoc-gen-go-aep/internal/genaep"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		return genaep.Run(gen)
	})
}

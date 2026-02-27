package docs

import (
	_ "embed"
	"sync"

	"github.com/swaggo/swag"
)

//go:embed swagger.json
var openAPIDoc string

var registerOnce sync.Once

type embeddedSwagger struct {
	doc string
}

func (s *embeddedSwagger) ReadDoc() string {
	return s.doc
}

func init() {
	registerOnce.Do(func() {
		swag.Register(swag.Name, &embeddedSwagger{doc: openAPIDoc})
	})
}

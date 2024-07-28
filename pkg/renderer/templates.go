package renderer

import (
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Renderer interface {
	Header(string, int) string
	PropertyHeader(string, int) string
	TableHeader() string
	TableFooter() string
	PropertyRow(string, jsonschema.Schema) string
}

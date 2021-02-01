package parser

import "github.com/elliotchance/pie/pie"

type Parser interface {
	ParserText(body string) pie.Strings
}

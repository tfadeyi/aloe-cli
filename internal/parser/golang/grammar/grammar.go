package grammar

import (
	"strings"

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/juju/errors"
	"github.com/tfadeyi/go-aloe/pkg/api"
)

type (
	ErrorGrammar struct {
		Stmts []*attribute `@@*`
	}
	InfoGrammar struct {
		Stmts []*attribute `@@*`
	}
	attribute struct {
		Key   string `Aloe Whitespace* @("name"|"title"|"description"|"url"|"version"|"summary"|"details"|"code")`
		Value string `Whitespace* @(String (Whitespace|EOL)*)+`
	}
)

const (
	versionAttr     = "version"
	nameAttr        = "name"
	descriptionAttr = "description"
	baseURLAttr     = "url"
	summaryAttr     = "summary"
	detailsAttr     = "details"
	codeAttr        = "code"
	titleAttr       = "title"
)

var (
	ErrMissingRequiredField = errors.New("missing required application field(s)")
	ErrParseSource          = errors.New("error parsing source material")
)

func (g ErrorGrammar) getAttribute(attribute string) (string, bool) {
	for _, attr := range g.Stmts {
		if strings.ToLower(attr.Key) == attribute {
			return attr.Value, true
		}
	}
	return "", false
}

func (g InfoGrammar) getAttribute(attribute string) (string, bool) {
	for _, attr := range g.Stmts {
		if strings.ToLower(attr.Key) == attribute {
			return attr.Value, true
		}
	}
	return "", false
}

var lexerDefinition = lexer.MustSimple([]lexer.SimpleRule{
	{"EOL", `[\n\r]+`},
	{"Aloe", `@aloe`},
	{"String", `([a-zA-Z_0-9/.//:,-/'])\w*`},
	{"Whitespace", `[ \t]+`},
})

func evalError(filename, source string, options ...participle.ParseOption) (*ErrorGrammar, error) {
	ast, err := participle.Build[ErrorGrammar](
		participle.Lexer(lexerDefinition),
	)
	if err != nil {
		return nil, err
	}

	return ast.ParseString(filename, source, options...)
}

func evalInfo(filename, source string, options ...participle.ParseOption) (*InfoGrammar, error) {
	ast, err := participle.Build[InfoGrammar](
		participle.Lexer(lexerDefinition),
	)
	if err != nil {
		return nil, err
	}

	return ast.ParseString(filename, source, options...)
}

func EvalError(filename, source string, options ...participle.ParseOption) (map[string]api.Error, error) {
	grammar, err := evalError(filename, source, options...)
	if err != nil {
		return nil, ErrParseSource
	}

	foundErrs := make(map[string]api.Error)
	newErr := api.Error{}
	if title, ok := grammar.getAttribute(titleAttr); ok {
		newErr.Title = strings.TrimSpace(title)
	}
	if code, ok := grammar.getAttribute(codeAttr); ok {
		newErr.Code = strings.TrimSpace(code)
	}
	if summary, ok := grammar.getAttribute(summaryAttr); ok {
		newErr.Summary = strings.TrimSpace(summary)
	}
	if details, ok := grammar.getAttribute(detailsAttr); ok {
		var val = strings.TrimSpace(details)
		newErr.Details = &val
	}

	// TODO checks on the required fields
	if newErr.Code == "" {
		return nil, errors.New("error code is missing")
	}

	foundErrs[newErr.Code] = newErr
	return foundErrs, nil
}

func EvalInfo(source string, options ...participle.ParseOption) (*api.Application, error) {
	info, err := evalInfo("", source, options...)
	if err != nil {
		return nil, ErrParseSource
	}

	app := api.Application{
		ErrorsDefinitions: map[string]api.Error{},
	}

	if title, ok := info.getAttribute(titleAttr); ok {
		var val = strings.TrimSpace(title)
		app.Title = &val
	}

	name, ok := info.getAttribute(nameAttr)
	if !ok {
		return nil, ErrMissingRequiredField
	}
	app.Name = strings.TrimSpace(name)

	version, ok := info.getAttribute(versionAttr)
	if !ok {
		return nil, ErrMissingRequiredField
	}
	app.Version = strings.TrimSpace(version)

	if description, ok := info.getAttribute(descriptionAttr); ok {
		var val = strings.TrimSpace(description)
		app.Description = &val
	}

	baseURL, ok := info.getAttribute(baseURLAttr)
	if !ok {
		return nil, ErrMissingRequiredField
	}
	app.BaseUrl = strings.TrimSpace(baseURL)

	return &app, nil
}

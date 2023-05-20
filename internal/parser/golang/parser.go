package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	"github.com/tfadeyi/aloe-cli/internal/logging"
	"github.com/tfadeyi/aloe-cli/internal/parser/golang/grammar"
	errhandler "github.com/tfadeyi/go-aloe"
	"github.com/tfadeyi/go-aloe/pkg/api"
)

type Parser struct {
	Spec                *api.Application
	GeneralInfoSource   string
	IncludedDirs        []string
	ApplicationPackages map[string]*ast.Package
	Logger              logging.Logger
}

// New client parser performs all checks at initialization time
func New(logger logging.Logger, source string, dirs ...string) *Parser {
	pkgs := map[string]*ast.Package{}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			// skip if dir doesn't exists
			continue
		}

		foundPkgs, err := getPackages(dir)
		if err != nil {
			logger.Info(err.Error())
			continue
		}

		for pkgName, pkg := range foundPkgs {
			if _, ok := pkgs[pkgName]; !ok {
				pkgs[pkgName] = pkg
			}
		}
	}

	return &Parser{
		Spec: &api.Application{
			ErrorsDefinitions: map[string]api.Error{},
		},
		GeneralInfoSource:   source,
		IncludedDirs:        dirs,
		ApplicationPackages: pkgs,
		Logger:              logger,
	}
}

func getPackages(dir string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return map[string]*ast.Package{}, err
	}

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			foundPkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			for pkgName, pkg := range foundPkgs {
				if _, ok := pkgs[pkgName]; !ok {
					pkgs[pkgName] = pkg
				}
			}
		}
		return err
	})
	return pkgs, err
}

func getFile(file string) (*ast.File, error) {
	fset := token.NewFileSet()
	return parser.ParseFile(fset, file, nil, parser.ParseComments)
}

func (p Parser) Parse() (*api.Application, error) {
	// collect application for the given source file
	f, err := getFile(p.GeneralInfoSource)
	if err != nil {
		// @aloe code parsing_golang_error
		// @aloe title parsing_golang_error
		// @aloe summary there was an error parsing the application source file
		// @aloe details The parser was processing the source file when it encountered an error
		// This could be because the file doesn't exist or the file not being a golang source file
		// Make sure the file is a golang source file

		return nil, errhandler.DefaultOrDie().Error(err, "parsing_golang_error")
	}
	if err := p.parseApplicationComments(f.Comments...); err != nil {
		p.Logger.Info(err.Error())
	}

	// collect all aloe error comments from packages and add them to the spec struct
	for _, pkg := range p.ApplicationPackages {
		for _, file := range pkg.Files {
			if err := p.parseErrorComments(file.Comments...); err != nil {
				p.Logger.Info(err.Error())
				continue
			}
		}
	}

	return p.Spec, nil
}

func (p Parser) parseApplicationComments(comments ...*ast.CommentGroup) error {
	for _, comment := range comments {
		app, err := grammar.EvalInfo(strings.TrimSpace(comment.Text()))
		switch {
		case errors.Is(err, grammar.ErrParseSource):
			continue
		case err != nil:
			p.Logger.Error(err, "")
			continue
		}

		p.Spec.Name = app.Name
		p.Spec.Version = app.Version
		p.Spec.BaseUrl = app.BaseUrl
		p.Spec.Description = app.Description
	}
	return nil
}

func (p Parser) parseErrorComments(comments ...*ast.CommentGroup) error {
	for _, comment := range comments {
		newErrs, err := grammar.EvalError("", strings.TrimSpace(comment.Text()))
		switch {
		case errors.Is(err, grammar.ErrParseSource):
			continue
		case err != nil:
			p.Logger.Error(err, "")
			continue
		}

		for code, err := range newErrs {
			if _, ok := p.Spec.ErrorsDefinitions[code]; !ok {
				p.Spec.ErrorsDefinitions[code] = err
			}
		}
	}
	return nil
}

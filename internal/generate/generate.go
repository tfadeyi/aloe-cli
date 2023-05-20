package generate

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	"github.com/tfadeyi/go-aloe/pkg/api"
	yaml "gopkg.in/yaml.v3"
)

//go:embed templates/markdown/error.md.tmpl
var errorDefinitionMarkdownTmpl string

//go:embed templates/markdown/info.md.tmpl
var applicationInfoMarkdownTmpl string

// defaultOutputFile is the default filename for the output file
const defaultOutputFile = "default.aloe"

var ErrUnsupportedFormat = errors.New("the specification is in an invalid format")

func IsValidOutputFormat(format string) bool {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "json", "yaml", "markdown":
		return true
	}
	return false
}

// WriteSpecification write the application bytes to a given writer
func WriteSpecification(spec *api.Application, stdout bool, formats ...string) error {
	// remove all previous output files
	err := cleanAll(spec.Name, formats...)
	if err != nil {
		return err
	}

	for _, format := range formats {
		var files = make(map[string][]byte)
		var body []byte
		var err error

		format = strings.ToLower(strings.TrimSpace(format))
		switch format {
		case "json":
			body, err = json.Marshal(spec)
			file := fmt.Sprintf("%s.%s", defaultOutputFile, format)
			files[file] = body
		case "yaml":
			body, err = yaml.Marshal(spec)
			file := fmt.Sprintf("%s.%s", defaultOutputFile, format)
			files[file] = body
		case "markdown":
			files, err = getMarkdownFromSpec(spec)
		default:
			err = ErrUnsupportedFormat
		}
		if err != nil {
			return err
		}

		if err := writeAll(stdout, files); err != nil {
			return err
		}
	}

	return nil
}

func cleanAll(applicationName string, formats ...string) error {
	for _, format := range formats {
		if !IsValidOutputFormat(format) {
			continue
		}

		if format == "markdown" {
			if _, err := os.Stat(applicationName); !errors.Is(err, os.ErrNotExist) {
				// delete spec file
				err = os.RemoveAll(applicationName)
				if err != nil {
					return errors.Annotate(err, "could not delete existing markdown directory")
				}
			}
			continue
		}

		file := fmt.Sprintf("%s.%s", defaultOutputFile, format)
		if _, err := os.Stat(file); !errors.Is(err, os.ErrNotExist) {
			// delete spec file
			err = os.RemoveAll(file)
			if err != nil {
				return errors.Annotate(err, "could not delete existing file")
			}
		}
	}
	return nil
}

func writeAll(stdout bool, files map[string][]byte) error {
	for path, body := range files {
		var err error
		var w = io.WriteCloser(os.Stdout)

		dirpath := filepath.Dir(path)
		if err := os.MkdirAll(dirpath, 0755); err != nil {
			return err
		}

		// decide which writer to use to print the application spec
		if !stdout {
			w, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			}
		}
		// write to writer
		_, err = w.Write(body)
		if err != nil {
			return err
		}
		w.Close()
	}
	return nil
}

func getMarkdownFromSpec(spec *api.Application) (map[string][]byte, error) {
	files := make(map[string][]byte)
	root := fmt.Sprintf("./%s/info.md", spec.Name)
	// parse application general information
	tmpl, err := template.New(spec.Name).Parse(applicationInfoMarkdownTmpl)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, spec)
	if err != nil {
		return nil, err
	}
	if _, ok := files[root]; !ok {
		files[root] = buf.Bytes()
	}

	for code, def := range spec.ErrorsDefinitions {
		tmpl, err := template.New(code).Parse(errorDefinitionMarkdownTmpl)
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer([]byte{})
		err = tmpl.Execute(buf, def)
		if err != nil {
			return nil, err
		}
		path := fmt.Sprintf("./%s/errors_definitions/%s.md", spec.Name, code)
		if _, ok := files[path]; !ok {
			files[path] = buf.Bytes()
		}
	}
	return files, nil
}

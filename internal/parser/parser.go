package parser

import (
	"github.com/tfadeyi/go-aloe/pkg/api"
)

type Parser interface {
	Parse() (*api.Application, error)
}

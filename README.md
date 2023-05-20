<div align="center">

# Aloe CLI

[![Continuous Integration](https://img.shields.io/github/actions/workflow/status/tfadeyi/aloe-cli/ci.yml?branch=main&style=flat-square)](https://github.com/tfadeyi/aloe-cli/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/License-MIT-yellowgreen.svg?style=flat-square)](https://github.com/tfadeyi/aloe-cli/blob/main/LICENSE)
[![Language](https://img.shields.io/github/go-mod/go-version/tfadeyi/aloe-cli?style=flat-square)](https://github.com/tfadeyi/aloe-cli)
[![GitHub release](https://img.shields.io/github/v/release/tfadeyi/aloe-cli?color=green&style=flat-square)](https://github.com/tfadeyi/aloe-cli/releases)
[![Code size](https://img.shields.io/github/languages/code-size/tfadeyi/aloe-cli?color=orange&style=flat-square)](https://github.com/tfadeyi/aloe-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/tfadeyi/aloe-cli?style=flat-square)](https://goreportcard.com/report/github.com/tfadeyi/aloe-cli)
</div>

---

Companion CLI tool for the [Aloe specification and library](https://github.com/tfadeyi/go-aloe).
Generate a specification from comments in the source code.

```shell
Generates the aloe specification from a given source code

Usage:
  aloe-cli spec generate [flags]

Flags:
      --dirs strings     Comma separated list of directories to be parses by the tool (default [./internal/parser/golang])
      --format strings   Output format (yaml,json,markdown) (default [yaml])
  -h, --help             help for generate
      --stdout           Print output to standard output.
```

Download the binary and run:

```shell
aloe-cli spec generate main.go
```

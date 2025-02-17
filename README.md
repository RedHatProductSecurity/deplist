# deplist

[![Tests](https://github.com/RedHatProductSecurity/deplist/actions/workflows/go.yml/badge.svg)](https://github.com/RedHatProductSecurity/deplist/actions/workflows/go.yml)

Scan and list the dependencies in a source code repository.

Supports:

- Go
- Java
- NodeJS
- Python
- Ruby
- Rust

Dependencies are printed in PackageURL format.

## Requirements

On Fedora:

```bash
dnf install golang-bin yarnpkg rubygem-bundler ruby-devel npm
```

## Command Line

### Build from source

```bash
$ make
go build cmd/deplist/deplist.go
```

### Run

```bash
$ ./deplist test/testRepo
pkg:npm/d3-scale-chromatic@2.0.0
pkg:npm/d3-time@2.0.0
pkg:npm/prop-types@15.7.2
pkg:npm/react@16.13.1
...
```

Verbose/debug output:

```bash
 deplist -debug ./test/testRepo/
DEBU[0000] Checking ./test/testRepo/
DEBU[0000] GetRubyDeps test/testRepo/Gemfile
DEBU[0000] Running env [--chdir=test/testRepo ruby /tmp/gemfile-parser.rb927489446 .]
DEBU[0000] GetGoPkgDeps test/testRepo/Gopkg.lock
DEBU[0000] GetGlideDeps test/testRepo/glide.lock
DEBU[0000] GetGolangDeps test/testRepo/go.mod
...
```

## API

The api functions as follows:

```go
func GetDeps(fullPath string, ignoreDirs ...string) ([]Dependency, Bitmask, error) {
```

### Parameters

- **fullPath:**

  Path to directory with source code.

### Returns

- **Dependency:**

  Array of Dependency structs from [dependencies.go](dependencies.go)

- **Bitmask:**

  A bitmask of found languages:

```go
const (
  LangGolang  = 1 << iota
  LangJava
  LangNodeJS
  LangPython
  LangRuby
  LangRust
)
```

- **error:**

  Standard Go error handling

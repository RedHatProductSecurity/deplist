[![Tests](https://github.com/RedHatProductSecurity/deplist/actions/workflows/go.yml/badge.svg)](https://github.com/RedHatProductSecurity/deplist/actions/workflows/go.yml)

# deplist

Scan and list the dependencies in a source code repository.

* Supports:
 - Go
 - NodeJS
 - Python
 - Ruby
 - Java
* Dependencies are printed in PackageURL format

## Requirements

On Fedora:

```bash
$ dnf install -y golang-bin yarnpkg maven rubygem-bundler npm
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

## API

The api functions as follows:

```
func GetDeps(fullPath string) ([]Dependency, Bitmask, error) {
```

### Parameters

* **fullPath:**

  Path to directory with source code.

### Returns

* **Depenency:**

  Array of Dependency structs from [dependencies.go](dependencies.go)


* **Bitmask:**

  A bitmask of found languages:

```
const (
	LangGolang = 1 << iota
	LangNodeJS
	LangPython
	LangRuby
)
```

* **error:**

  Standard Go error handling

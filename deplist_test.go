package deplist

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func BuildWant() []Dependency {
	var deps []Dependency

	golangPaths := []string{
		"errors",
		"fmt",
		"github.com/RedHatProductSecurity/deplist",
		"github.com/openshift/api/config/v1",
		"golang.org/x/text/unicode",
		"internal/abi",
		"internal/bytealg",
		"internal/cpu",
		"internal/fmtsort",
		"internal/goexperiment",
		"internal/itoa",
		"internal/oserror",
		"internal/poll",
		"internal/race",
		"internal/reflectlite",
		"internal/syscall/execenv",
		"internal/syscall/unix",
		"internal/testlog",
		"internal/unsafeheader",
		"io",
		"io/fs",
		"math",
		"math/bits",
		"os",
		"path",
		"reflect",
		"runtime",
		"runtime/internal/atomic",
		"runtime/internal/math",
		"runtime/internal/sys",
		"sort",
		"strconv",
		"sync",
		"sync/atomic",
		"syscall",
		"time",
		"unicode",
		"unicode/utf8",
		"unsafe",
	}

	glidePaths := []string{
		"github.com/beorn7/perks",
		"github.com/beorn7/perks/quantile",
		"github.com/bgentry/speakeasy",
		"github.com/boltdb/bolt",
		"github.com/cockroachdb/cmux",
		"github.com/coreos/go-semver",
		"github.com/coreos/go-semver/semver",
		"github.com/coreos/go-systemd",
		"github.com/coreos/go-systemd/daemon",
		"github.com/coreos/go-systemd/journal",
		"github.com/coreos/go-systemd/util",
	}

	gopkgPaths := []string{
		"github.com/BurntSushi/toml",
		"github.com/aws/aws-sdk-go",
		"github.com/aws/aws-sdk-go/aws",
		"github.com/aws/aws-sdk-go/aws/awserr",
		"github.com/aws/aws-sdk-go/aws/awsutil",
		"github.com/aws/aws-sdk-go/aws/client",
		"github.com/aws/aws-sdk-go/aws/client/metadata",
		"github.com/aws/aws-sdk-go/aws/corehandlers",
	}

	npmSet1 := []string{
		"loose-envify",
		"iconv-lite",
		"d3-brush",
		"d3-zoom",
		"rw",
		"d3-ease",
		"object-assign",
		"commander",
		"d3-dsv",
		"d3-scale",
		"is-plain-object",
		"d3-quadtree",
		"tiny-warning",
		"d3-hierarchy",
		"d3-scale-chromatic",
		"d3-axis",
		"d3-color",
		"prismjs",
		"iconv-lite",
		"angular",
		"d3-delaunay",
		"rxjs",
		"d3-path",
		"d3-array",
		"js-tokens",
		"d3-contour",
		"safer-buffer",
		"react-is",
		"d3-dispatch",
		"d3-force",
		"prop-types",
		"tiny-emitter",
		"d3-polygon",
		"d3-chord",
		"d3-fetch",
		"tslib",
		"good-listener",
		"d3",
		"delegate",
		"d3-drag",
		"delaunator",
		"d3-timer",
		"d3-geo",
		"slate",
		"select",
		"esrever",
		"d3-transition",
		"clipboard",
		"d3-format",
		"d3-random",
		"d3-shape",
		"d3-time",
		"immer",
		"@types/esrever",
		"d3-time-format",
		"d3-selection",
		"react",
		"tether",
		"d3-interpolate",
	}

	rubySet := []string{
		"concurrent-ruby",
		"lru_redux",
		"zeitwerk",
		"async",
		"fluent-plugin-systemd",
		"http-parser",
		"ltsv",
		"public_suffix",
		"faraday-multipart",
		"fluent-config-regexp-type",
		"recursive-open-struct",
		"unf_ext",
		"aws-eventstream",
		"webrick",
		"faraday-em_http",
		"fluentd",
		"yajl-ruby",
		"fluent-plugin-elasticsearch",
		"faraday-patron",
		"mini_mime",
		"tzinfo",
		"connection_pool",
		"fluent-plugin-kubernetes_metadata_filter",
		"fluent-plugin-prometheus",
		"nio4r",
		"oj",
		"openid_connect",
		"rack",
		"sigdump",
		"digest-crc",
		"ethon",
		"multipart-post",
		"addressable",
		"faraday-net_http_persistent",
		"rack-oauth2",
		"excon",
		"fluent-plugin-label-router",
		"bindata",
		"fluent-plugin-record-modifier",
		"http",
		"systemd-journal",
		"faraday-retry",
		"ruby2_keywords",
		"mime-types",
		"timers",
		"unf",
		"fluent-plugin-detect-exceptions",
		"jsonpath",
		"rake",
		"validate_email",
		"aws-sdk-cloudwatchlogs",
		"jmespath",
		"prometheus-client",
		"protocol-http1",
		"ffi",
		"fluent-plugin-grafana-loki",
		"bigdecimal",
		"protocol-http",
		"aws-partitions",
		"faraday-httpclient",
		"fluent-plugin-multi-format-parser",
		"http_parser.rb",
		"protocol-http2",
		"rest-client",
		"activesupport",
		"ffi-compiler",
		"fluent-plugin-splunk-hec",
		"json-jwt",
		"msgpack",
		"protocol-hpack",
		"strptime",
		"validate_url",
		"faraday",
		"async-pool",
		"faraday-net_http",
		"fluent-plugin-concat",
		"fluent-plugin-kafka",
		"multi_json",
		"net-http-persistent",
		"uuidtools",
		"activemodel",
		"elasticsearch-transport",
		"mail",
		"ruby-kafka",
		"serverengine",
		"tzinfo-data",
		"webfinger",
		"aws-sigv4",
		"elasticsearch-api",
		"fiber-local",
		"fluent-plugin-remote-syslog",
		"attr_required",
		"http-form_data",
		"syslog_protocol",
		"faraday-em_synchrony",
		"httpclient",
		"fluent-mixin-config-placeholders",
		"fluent-plugin-cloudwatch-logs",
		"i18n",
		"async-io",
		"elasticsearch",
		"http-cookie",
		"kubeclient",
		"minitest",
		"aes_key_wrap",
		"mime-types-data",
		"netrc",
		"console",
		"cool.io",
		"domain_name",
		"async-http",
		"http-accept",
		"traces",
		"typhoeus",
		"faraday-rack",
		"faraday-excon",
		"fluent-plugin-rewrite-tag-filter",
		"swd",
		"aws-sdk-core",
	}

	pythonSet := []Dependency{
		Dependency{DepType: LangPython, Path: "cotyledon"},
		Dependency{DepType: LangPython, Path: "Flask"},
		Dependency{DepType: LangPython, Path: "kuryr-lib"},
		Dependency{DepType: LangPython, Path: "docutils"},
		Dependency{DepType: LangPython, Path: "python-dateutil"},
		Dependency{DepType: LangPython, Path: "unittest2", Version: "0.5.1"},
		Dependency{DepType: LangPython, Path: "cryptography", Version: "2.3.0"},
		Dependency{DepType: LangPython, Path: "suds-py3"},
		Dependency{DepType: LangPython, Path: "suds"},
		Dependency{DepType: LangPython, Path: "git+https://github.com/candlepin/subscription-manager#egg=subscription_manager"},
		Dependency{DepType: LangPython, Path: "git+https://github.com/candlepin/python-iniparse#egg=iniparse"},
		Dependency{DepType: LangPython, Path: "iniparse"},
		Dependency{DepType: LangPython, Path: "requests"},
		Dependency{DepType: LangPython, Path: "m2crypto"},
	}

	for _, n := range golangPaths {
		d := Dependency{
			DepType: 1,
			Path:    n,
		}

		deps = append(deps, d)
	}

	deps[4].Version = "v0.3.3" // test golang.org/x/text/unicode version

	for _, n := range glidePaths {
		deps = append(deps, Dependency{
			DepType: 1,
			Path:    n,
		})
	}

	for i, n := range gopkgPaths {
		ver := ""
		if i > 0 {
			ver = "v1.13.49"
		}
		deps = append(deps, Dependency{
			DepType: 1,
			Path:    n,
			Version: ver,
		})
	}

	for _, n := range npmSet1 {
		d := Dependency{
			DepType: LangNodeJS,
			Path:    n,
		}
		deps = append(deps, d)
	}
	deps = append(deps, Dependency{DepType: 2, Path: "com.amazonaws:aws-lambda-java-core:jar", Version: "1.0.0"}) // java

	for _, n := range rubySet {
		d := Dependency{
			DepType: LangRuby,
			Path:    n,
		}
		deps = append(deps, d)
	}

	deps = append(deps, pythonSet...)

	return deps
}

func depToKey(pkg Dependency) string {
	key := fmt.Sprintf("%s:%s", GetLanguageStr(pkg.DepType), pkg.Path)
	// fmt.Println(key)
	// return key
	return key
}

func TestGetDeps(t *testing.T) {
	want := BuildWant()

	got, gotBitmask, err := GetDeps("test/testRepo")
	if err != nil {
		t.Errorf("GetDeps failed: %s", err)
		return
	}

	if gotBitmask != 31 {
		t.Errorf("GotBitmask() != 31; got: %d", gotBitmask)
	}

	gotMap := make(map[string]Dependency)
	wantMap := make(map[string]Dependency)

	for _, pkg := range got {
		key := depToKey(pkg)
		if _, ok := gotMap[key]; !ok {
			gotMap[key] = pkg
		}
	}

	for _, pkg := range want {
		key := depToKey(pkg)
		if _, ok := wantMap[key]; !ok {
			wantMap[key] = pkg
		}
	}

	for _, w := range want {
		key := depToKey(w)
		if g, ok := gotMap[key]; !ok {
			t.Errorf("GetDeps() wanted: %s - not found", key)
		} else {
			if w.Version != "" && w.Version != g.Version {
				t.Errorf("%s version mismatch: wanted %s but got %s", key, w.Version, g.Version)
			}
		}
	}

	if len(want) != len(got) {
		if len(got) > len(want) {
			for _, pkg := range got {
				if _, ok := wantMap[depToKey(pkg)]; !ok {
					t.Errorf("GetDeps() got unexpected: %s", pkg.Path)
				}
			}
		}
		t.Errorf("GetDeps() = %d; want %d", len(got), len(want))
	}
}

func TestFindBaseDir(t *testing.T) {
	type TestCase struct {
		Input    string
		Expected string
		Err      bool
	}

	tests := make([]TestCase, 5)

	top := t.TempDir()
	tests[0] = TestCase{
		Input:    "non-existent directory",
		Expected: "",
		Err:      true,
	}

	dirpath := filepath.Join(top, "baz")
	os.MkdirAll(dirpath, 0755)
	tests[1] = TestCase{
		Input:    dirpath,
		Expected: dirpath,
		Err:      false,
	}

	tempFile, err := ioutil.TempFile(top, "bar")
	if err != nil {
		t.Error(err)
	}
	tests[2] = TestCase{
		Input:    tempFile.Name(),
		Expected: "",
		Err:      true,
	}

	dirpath = filepath.Join(top, "foo/bar/foo/bar/foo/bar")
	err = os.MkdirAll(dirpath, 0755)
	if err != nil {
		t.Error(err)
	}
	tests[3] = TestCase{
		Input:    filepath.Join(top, "foo"),
		Expected: dirpath,
		Err:      false,
	}

	top = t.TempDir()
	dirpath = filepath.Join(top, "foo/bar/foo/bar/foo/bar")
	err = os.MkdirAll(dirpath, 0755)
	if err != nil {
		t.Error(err)
	}
	tempFile, err = ioutil.TempFile(filepath.Join(top, "foo/bar/foo"), "baz")
	if err != nil {
		t.Error(err)
	}
	tests[4] = TestCase{
		Input:    filepath.Join(top, "foo"),
		Expected: filepath.Join(top, "foo/bar/foo"),
		Err:      false,
	}

	for i, test := range tests {
		dir, err := findBaseDir(test.Input)

		if test.Err {
			if err == nil {
				t.Errorf("%d: Expected error reading directory: %s but didn't get one", i, dir)
			}
		}
		if test.Expected != dir {
			t.Errorf("%d: Expected %s, got %s", i, test.Expected, dir)
		}
	}
}

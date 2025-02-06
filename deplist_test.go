package deplist

import (
	"fmt"
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
		"internal/coverage/rtcov",
		"internal/cpu",
		"internal/fmtsort",
		"internal/goarch",
		"internal/godebugs",
		"internal/goexperiment",
		"internal/goos",
		"internal/itoa",
		"internal/oserror",
		"internal/poll",
		"internal/race",
		"internal/reflectlite",
		"internal/safefilepath",
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
		"runtime/internal/syscall",
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

	rubySet := []Dependency{
		{DepType: LangRuby, Path: "concurrent-ruby"},
		{DepType: LangRuby, Path: "lru_redux"},
		{DepType: LangRuby, Path: "zeitwerk"},
		{DepType: LangRuby, Path: "async"},
		{DepType: LangRuby, Path: "fluent-plugin-systemd"},
		{DepType: LangRuby, Path: "http-parser"},
		{DepType: LangRuby, Path: "ltsv"},
		{DepType: LangRuby, Path: "public_suffix"},
		{DepType: LangRuby, Path: "faraday-multipart"},
		{DepType: LangRuby, Path: "fluent-config-regexp-type"},
		{DepType: LangRuby, Path: "recursive-open-struct"},
		{DepType: LangRuby, Path: "unf_ext"},
		{DepType: LangRuby, Path: "aws-eventstream"},
		{DepType: LangRuby, Path: "webrick"},
		{DepType: LangRuby, Path: "faraday-em_http"},
		{DepType: LangRuby, Path: "fluentd"},
		{DepType: LangRuby, Path: "yajl-ruby"},
		{DepType: LangRuby, Path: "fluent-plugin-elasticsearch"},
		{DepType: LangRuby, Path: "faraday-patron"},
		{DepType: LangRuby, Path: "mini_mime"},
		{DepType: LangRuby, Path: "tzinfo"},
		{DepType: LangRuby, Path: "connection_pool"},
		{DepType: LangRuby, Path: "fluent-plugin-kubernetes_metadata_filter"},
		{DepType: LangRuby, Path: "fluent-plugin-prometheus"},
		{DepType: LangRuby, Path: "nio4r"},
		{DepType: LangRuby, Path: "oj"},
		{DepType: LangRuby, Path: "openid_connect"},
		{DepType: LangRuby, Path: "rack"},
		{DepType: LangRuby, Path: "sigdump"},
		{DepType: LangRuby, Path: "digest-crc"},
		{DepType: LangRuby, Path: "ethon"},
		{DepType: LangRuby, Path: "multipart-post"},
		{DepType: LangRuby, Path: "addressable"},
		{DepType: LangRuby, Path: "faraday-net_http_persistent"},
		{DepType: LangRuby, Path: "rack-oauth2"},
		{DepType: LangRuby, Path: "excon"},
		{DepType: LangRuby, Path: "fluent-plugin-label-router"},
		{DepType: LangRuby, Path: "bindata"},
		{DepType: LangRuby, Path: "fluent-plugin-record-modifier"},
		{DepType: LangRuby, Path: "http"},
		{DepType: LangRuby, Path: "systemd-journal"},
		{DepType: LangRuby, Path: "faraday-retry"},
		{DepType: LangRuby, Path: "ruby2_keywords"},
		{DepType: LangRuby, Path: "mime-types"},
		{DepType: LangRuby, Path: "timers"},
		{DepType: LangRuby, Path: "unf"},
		{DepType: LangRuby, Path: "fluent-plugin-detect-exceptions"},
		{DepType: LangRuby, Path: "jsonpath"},
		{DepType: LangRuby, Path: "rake"},
		{DepType: LangRuby, Path: "validate_email"},
		{DepType: LangRuby, Path: "aws-sdk-cloudwatchlogs"},
		{DepType: LangRuby, Path: "jmespath"},
		{DepType: LangRuby, Path: "prometheus-client"},
		{DepType: LangRuby, Path: "protocol-http1"},
		{DepType: LangRuby, Path: "ffi"},
		{DepType: LangRuby, Path: "fluent-plugin-grafana-loki"},
		{DepType: LangRuby, Path: "bigdecimal"},
		{DepType: LangRuby, Path: "protocol-http"},
		{DepType: LangRuby, Path: "aws-partitions"},
		{DepType: LangRuby, Path: "faraday-httpclient"},
		{DepType: LangRuby, Path: "fluent-plugin-multi-format-parser"},
		{DepType: LangRuby, Path: "http_parser.rb"},
		{DepType: LangRuby, Path: "protocol-http2"},
		{DepType: LangRuby, Path: "rest-client"},
		{DepType: LangRuby, Path: "activesupport"},
		{DepType: LangRuby, Path: "ffi-compiler"},
		{DepType: LangRuby, Path: "fluent-plugin-splunk-hec"},
		{DepType: LangRuby, Path: "json-jwt"},
		{DepType: LangRuby, Path: "msgpack"},
		{DepType: LangRuby, Path: "protocol-hpack"},
		{DepType: LangRuby, Path: "strptime"},
		{DepType: LangRuby, Path: "validate_url"},
		{DepType: LangRuby, Path: "faraday"},
		{DepType: LangRuby, Path: "async-pool"},
		{DepType: LangRuby, Path: "faraday-net_http"},
		{DepType: LangRuby, Path: "fluent-plugin-concat"},
		{DepType: LangRuby, Path: "fluent-plugin-kafka"},
		{DepType: LangRuby, Path: "multi_json"},
		{DepType: LangRuby, Path: "net-http-persistent"},
		{DepType: LangRuby, Path: "uuidtools"},
		{DepType: LangRuby, Path: "activemodel"},
		{DepType: LangRuby, Path: "elasticsearch-transport"},
		{DepType: LangRuby, Path: "mail"},
		{DepType: LangRuby, Path: "ruby-kafka"},
		{DepType: LangRuby, Path: "serverengine"},
		{DepType: LangRuby, Path: "tzinfo-data"},
		{DepType: LangRuby, Path: "webfinger"},
		{DepType: LangRuby, Path: "aws-sigv4"},
		{DepType: LangRuby, Path: "elasticsearch-api"},
		{DepType: LangRuby, Path: "fiber-local"},
		{DepType: LangRuby, Path: "fluent-plugin-remote-syslog"},
		{DepType: LangRuby, Path: "attr_required"},
		{DepType: LangRuby, Path: "http-form_data"},
		{DepType: LangRuby, Path: "syslog_protocol"},
		{DepType: LangRuby, Path: "faraday-em_synchrony"},
		{DepType: LangRuby, Path: "httpclient"},
		{DepType: LangRuby, Path: "fluent-mixin-config-placeholders"},
		{DepType: LangRuby, Path: "fluent-plugin-cloudwatch-logs"},
		{DepType: LangRuby, Path: "i18n"},
		{DepType: LangRuby, Path: "async-io"},
		{DepType: LangRuby, Path: "elasticsearch"},
		{DepType: LangRuby, Path: "http-cookie"},
		{DepType: LangRuby, Path: "kubeclient"},
		{DepType: LangRuby, Path: "minitest"},
		{DepType: LangRuby, Path: "aes_key_wrap"},
		{DepType: LangRuby, Path: "mime-types-data"},
		{DepType: LangRuby, Path: "netrc"},
		{DepType: LangRuby, Path: "console"},
		{DepType: LangRuby, Path: "cool.io"},
		{DepType: LangRuby, Path: "domain_name"},
		{DepType: LangRuby, Path: "async-http"},
		{DepType: LangRuby, Path: "http-accept"},
		{DepType: LangRuby, Path: "traces"},
		{DepType: LangRuby, Path: "typhoeus"},
		{DepType: LangRuby, Path: "faraday-rack"},
		{DepType: LangRuby, Path: "faraday-excon"},
		{DepType: LangRuby, Path: "fluent-plugin-rewrite-tag-filter"},
		{DepType: LangRuby, Path: "swd"},
		{DepType: LangRuby, Path: "aws-sdk-core"},
	}

	pythonSet := []Dependency{
		{DepType: LangPython, Path: "cotyledon"},
		{DepType: LangPython, Path: "Flask"},
		{DepType: LangPython, Path: "kuryr-lib"},
		{DepType: LangPython, Path: "docutils"},
		{DepType: LangPython, Path: "python-dateutil"},
		{DepType: LangPython, Path: "unittest2", Version: "0.5.1"},
		{DepType: LangPython, Path: "cryptography", Version: "2.3.0"},
		{DepType: LangPython, Path: "suds-py3"},
		{DepType: LangPython, Path: "suds"},
		{DepType: LangPython, Path: "git+https://github.com/candlepin/subscription-manager#egg=subscription_manager"},
		{DepType: LangPython, Path: "git+https://github.com/candlepin/python-iniparse#egg=iniparse"},
		{DepType: LangPython, Path: "iniparse"},
		{DepType: LangPython, Path: "requests"},
		{DepType: LangPython, Path: "m2crypto"},
	}

	rustSet := []Dependency{
		{DepType: LangRust, Path: "cc", Version: "1.0.79"},
		{DepType: LangRust, Path: "hermit-abi", Version: "0.3.1"},
		{DepType: LangRust, Path: "num-traits", Version: "0.2.15"},
		{DepType: LangRust, Path: "windows-sys", Version: "0.48.0"},
		{DepType: LangRust, Path: "bitflags", Version: "1.3.2"},
		{DepType: LangRust, Path: "io-uring", Version: "0.6.0"},
		{DepType: LangRust, Path: "memmap2", Version: "0.5.10"},
		{DepType: LangRust, Path: "rustix", Version: "0.37.15"},
		{DepType: LangRust, Path: "autocfg", Version: "1.1.0"},
		{DepType: LangRust, Path: "paste", Version: "1.0.12"},
		{DepType: LangRust, Path: "virtio-bindings", Version: "0.2.0"},
		{DepType: LangRust, Path: "windows_aarch64_gnullvm", Version: "0.48.0"},
		{DepType: LangRust, Path: "libc", Version: "0.2.142"},
		{DepType: LangRust, Path: "lazy_static", Version: "1.4.0"},
		{DepType: LangRust, Path: "virtio-driver", Version: "0.5.0"},
		{DepType: LangRust, Path: "windows-targets", Version: "0.48.0"},
		{DepType: LangRust, Path: "errno", Version: "0.3.1"},
		{DepType: LangRust, Path: "io-lifetimes", Version: "1.0.10"},
		{DepType: LangRust, Path: "windows_aarch64_msvc", Version: "0.48.0"},
		{DepType: LangRust, Path: "windows_i686_msvc", Version: "0.48.0"},
		{DepType: LangRust, Path: "blkio", Version: "0.4.0"},
		{DepType: LangRust, Path: "linux-raw-sys", Version: "0.3.4"},
		{DepType: LangRust, Path: "windows_i686_gnu", Version: "0.48.0"},
		{DepType: LangRust, Path: "errno-dragonfly", Version: "0.1.2"},
		{DepType: LangRust, Path: "windows_x86_64_gnu", Version: "0.48.0"},
		{DepType: LangRust, Path: "libblkio", Version: "1.3.0"},
		{DepType: LangRust, Path: "windows_x86_64_gnullvm", Version: "0.48.0"},
		{DepType: LangRust, Path: "windows_x86_64_msvc", Version: "0.48.0"},
		{DepType: LangRust, Path: "pci-driver", Version: "0.1.3"},
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
	deps = append(deps, rubySet...)
	deps = append(deps, pythonSet...)
	deps = append(deps, rustSet...)

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

	expectedBitmask := 63
	if gotBitmask != Bitmask(expectedBitmask) {
		t.Errorf("GotBitmask() != %d; got: %d", expectedBitmask, gotBitmask)
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
	os.MkdirAll(dirpath, 0o755)
	tests[1] = TestCase{
		Input:    dirpath,
		Expected: dirpath,
		Err:      false,
	}

	tempFile, err := os.CreateTemp(top, "bar")
	if err != nil {
		t.Error(err)
	}
	tests[2] = TestCase{
		Input:    tempFile.Name(),
		Expected: "",
		Err:      true,
	}

	dirpath = filepath.Join(top, "foo/bar/foo/bar/foo/bar")
	err = os.MkdirAll(dirpath, 0o755)
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
	err = os.MkdirAll(dirpath, 0o755)
	if err != nil {
		t.Error(err)
	}
	_, err = os.CreateTemp(filepath.Join(top, "foo/bar/foo"), "baz")
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

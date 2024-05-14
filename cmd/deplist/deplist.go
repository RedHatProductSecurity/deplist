package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/RedHatProductSecurity/deplist"
	purl "github.com/mcoops/packageurl-go"
	log "github.com/sirupsen/logrus"
)

type searchModeFlag []string

func (i *searchModeFlag) String() string {
	return strings.Join(*i, ", ")
}

func (i *searchModeFlag) Set(value string) error {
	values := strings.Split(value, ",")
	for _, value := range values {
		switch value {
		case "deps":
			*i = append(*i, "deps")
		case "bundled":
			*i = append(*i, "bundled")
		}
	}
	return nil
}

func (i *searchModeFlag) ShouldHandleDep(dep deplist.Dependency) bool {
	res := false
	for _, searchMode := range *i {
		switch searchMode {
		case "deps":
			res = res || !dep.IsBundled
		case "bundled":
			res = res || dep.IsBundled
		}
	}
	return res
}

func main() {
	deptypePtr := flag.Int("deptype", -1, "golang, nodejs, python etc")
	debugPtr := flag.Bool("debug", false, "debug logging (default false)")
	var searchModes searchModeFlag
	flag.Var(&searchModes, "modes", "search mode (bundled, deps)")

	flag.Parse()

	if len(searchModes) == 0 {
		searchModes = []string{"deps", "bundled"}
	}
	if *debugPtr == true {
		log.SetLevel(log.DebugLevel)
	}

	if flag.Args() == nil || len(flag.Args()) == 0 {
		fmt.Println("Not path to scan was specified, i.e. deplist /tmp/files/")
		return
	}

	path := flag.Args()[0]

	deps, _, err := deplist.GetDeps(path)
	if err != nil {
		fmt.Println(err.Error())
	}

	if *deptypePtr == -1 {
		for _, dep := range deps {
			if !searchModes.ShouldHandleDep(dep) {
				continue
			}
			version := dep.Version

			inst, _ := purl.FromString(fmt.Sprintf("pkg:%s/%s@%s", deplist.GetLanguageStr(dep.DepType), dep.Path, version))
			fmt.Print(inst)
			fmt.Println()
		}
	} else {
		deptype := deplist.Bitmask(*deptypePtr)
		for _, dep := range deps {
			if !searchModes.ShouldHandleDep(dep) {
				continue
			}

			if (dep.DepType & deptype) == deptype {
				fmt.Printf("%s@%s", dep.Path, dep.Version)
				fmt.Println()
			}
		}
	}
}

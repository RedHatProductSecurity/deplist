package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/RedHatProductSecurity/deplist"
	purl "github.com/package-url/packageurl-go"
	log "github.com/sirupsen/logrus"
)

func main() {
	deptypePtr := flag.Int("deptype", -1, "golang, nodejs, python etc")
	debugPtr := flag.Bool("debug", false, "debug logging (default false)")
	ignorePtr := flag.String("ignore", "", "comma separated list of directory names to ignore (default '')")

	flag.Parse()

	if *debugPtr == true {
		log.SetLevel(log.DebugLevel)
	}

	var ignoreDirs []string
	if ignorePtr != nil {
		ignoreDirs = strings.Split(*ignorePtr, ",")
	}

	if flag.Args() == nil || len(flag.Args()) == 0 {
		fmt.Println("No path to scan was specified, i.e. deplist /tmp/files/")
		return
	}

	path := flag.Args()[0]

	deps, _, err := deplist.GetDeps(path, ignoreDirs...)
	if err != nil {
		fmt.Println(err.Error())
	}

	if *deptypePtr == -1 {
		for _, dep := range deps {
			inst, _ := purl.FromString(dep.ToString())
			fmt.Println(inst)
		}
	} else {
		deptype := deplist.Bitmask(*deptypePtr)
		for _, dep := range deps {
			if (dep.DepType & deptype) == deptype {
				fmt.Printf("%s@%s\n", dep.Path, dep.Version)
			}
		}
	}
}

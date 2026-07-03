package main

import (
	"flag"

	"github.com/curioswitch/go-build"
	"github.com/goyek/x/boot"
	"github.com/wasilibs/tools/tasks"
)

func main() {
	// TODO: Enable after migrating go-build to ryl
	_ = flag.Lookup("skip").Value.Set("lint-yaml,format-yaml")
	tasks.Define(tasks.Params{
		LibraryName: "ryl",
		LibraryRepo: "owenlamont/ryl",
		BuildOpts:   []build.Option{build.DisableCoverage()},
	})
	boot.Main()
}

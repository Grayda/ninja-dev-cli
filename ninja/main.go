package main

import (
	"github.com/docopt/docopt-go"
)

const usage = `Ninja Sphere Developer CLI.

Usage:
  ninja [options] build <path>
  ninja -h | --help

Options:
  build                                 Builds a package for a driver or app, ready for distribution to the Ninja Sphere.
	--debian-install-path=<buildpath>     The deb version of this driver or app will be installed to this path. Include trailing slash [default: /opt/ninjablocks/]
	--build-type=<buildtype>              What type of package should be built? Valid options are "deb", "snappy" or "all" [default: all]
	--package-type=[auto|apps|drivers]    If the binary name starts with driver-, it'll deploy to <buildpath>/drivers, otherwise <buildpath>/apps. Use this to overwrite [default: auto]
  --snappy-namespace=<ns>               Set the namespace for the Snappy package.
  -v --verbose                          Verbose
  -h --help                             Show this screen.
`

func main() {

	arguments, _ = docopt.Parse(usage, nil, true, "Ninja Sphere Developer CLI", false)
	verifyDockerImages()
	if arguments["build"].(bool) {
		doBuild(arguments)
	}
}

package main

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"log"
	"os"
	"path/filepath"
)

type buildContext struct {
	stagingHost   string
	stagingDocker string

	outputDocker string

	archBinaries map[string]string

	dockerVolumeArgs []string
}

func doBuild(arguments map[string]interface{}) {
	path := arguments["<path>"].(string)

	pkg, err := NewNinjaPackage(path)
	if err != nil {
		log.Fatalf("Could not parse package metadata: %v\n", err)
	}

	color.Printf("@bPackage name: @|%v\n", pkg.Name())
	color.Printf("@bShort name:   @|%v\n", pkg.ShortName())
	color.Printf("@bVersion:      @|%v\n", pkg.Version())

	dockerArgs := []string{
		"-v", pkg.BasePath + ":/home/builder/go/src/ninja-dev-cli/" + pkg.ShortName(),
	}

	fmt.Printf("\n")

	/** Compile **/

	fullPath := "/home/builder/go/src/ninja-dev-cli/" + pkg.ShortName()

	// Perform Intel magic
	cmdIntel := []string{
		"/bin/bash", "-c",
		"GOPATH=" + fullPath + "/.gopath-cli:$GOPATH; (" +
			"go get -d ninja-dev-cli/" + pkg.ShortName() + ";" +
			"go build -o " + fullPath + "/.gopath-cli/autoclibinary-amd64 ninja-dev-cli/" + pkg.ShortName() + ";" +
			")",
	}

	color.Printf("@cCompiling for architecture: amd64\n")
	runNativeDockerCrossCommand(dockerArgs, cmdIntel...)

	// Perform ARM magic
	cmdARM := []string{
		"/bin/bash", "-c",
		"(" +
			"sed -i '/source-root-users/d' /etc/schroot/chroot.d/click-ubuntu-sdk-14.04-armhf; " +
			"click chroot -aarmhf -fubuntu-sdk-14.04 -s trusty run " +
			"CGO_ENABLED=1 GOARCH=arm GOARM=7 " +
			"PKG_CONFIG_LIBDIR=/usr/lib/arm-linux-gnueabihf/pkgconfig:/usr/lib/pkgconfig:/usr/share/pkgconfig " +
			"CC=arm-linux-gnueabihf-gcc " +
			"GOPATH=" + fullPath + "/.gopath-cli:$GOPATH " +
			"go build -ldflags '-extld=arm-linux-gnueabihf-gcc' " +
			"-o " + fullPath + "/.gopath-cli/autoclibinary-armhf ninja-dev-cli/" + pkg.ShortName() + ";" +
			")",
	}

	color.Printf("@cCompiling for architecture: armhf\n")
	runNativeDockerCrossCommand(dockerArgs, cmdARM...)

	fmt.Printf("\n")

	/** Package **/

	context := &buildContext{}

	context.archBinaries = map[string]string{}
	context.archBinaries["amd64"] = filepath.Join(pkg.BasePath, ".gopath-cli/autoclibinary-amd64")
	context.archBinaries["armhf"] = filepath.Join(pkg.BasePath, ".gopath-cli/autoclibinary-armhf")

	context.stagingHost = filepath.Join(pkg.BasePath, ".staging-cli")
	context.stagingDocker = filepath.Join(fullPath, ".staging-cli")

	context.outputDocker = fullPath

	context.dockerVolumeArgs = dockerArgs

	// log.Printf("Out: %v\n", context.stagingHost)
	// log.Printf("In : %v\n", context.stagingDocker)

	// remove if we can
	os.RemoveAll(context.stagingHost)
	// and recreate ready to fill
	os.MkdirAll(context.stagingHost, 0750)

	// build for all archs for now
	for _, arch := range []string{"amd64", "armhf"} {
		color.Printf("@mCreating Snappy package for architecture: %v\n", arch)
		buildSnappyPackage(pkg, context, arch)

		color.Printf("@mCreating Debian package for architecture: %v\n", arch)
		buildDebianPackage(pkg, context, arch)
	}
}

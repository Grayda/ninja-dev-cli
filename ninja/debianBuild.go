package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/termie/go-shutil"
)

type upstartTemplateVars struct {
	PackageName string
	TargetDir   string
	BaseDir     string // Where should dpkg install to?
}

func buildDebianPackage(pkg *NinjaPackage, ctx *buildContext, arch string) {
	dockerCurr := filepath.Join(ctx.stagingDocker, "debian-"+arch)
	stagingCurr := filepath.Join(ctx.stagingHost, "debian-"+arch)
	os.MkdirAll(stagingCurr, 0755)

	var targetDir string

	switch {
	case (arguments["--package-type"] == "auto" && strings.HasPrefix(pkg.ShortName(), "driver-") == true) || arguments["--package-type"].(string) == "driver":
		targetDir = "drivers"
	case (arguments["--package-type"] == "auto" && strings.HasPrefix(pkg.ShortName(), "driver-") == false) || arguments["--package-type"].(string) == "app":
		targetDir = "apps"
	}

	targetPath := filepath.Join(stagingCurr, deployPath, targetDir, pkg.ShortName())
	os.MkdirAll(targetPath, 0755)

	// binary itself
	binCurr := filepath.Join(targetPath, pkg.ShortName())
	fmt.Println(ctx.archBinaries[arch], binCurr)
	shutil.Copy(ctx.archBinaries[arch], binCurr, false)

	// package.json
	srcFile := filepath.Join(pkg.BasePath, "package.json")
	pkgCurr := filepath.Join(targetPath, "package.json")
	shutil.Copy(srcFile, pkgCurr, true) // don't copy symlink itself, just the real file

	// Copy all of the files in package.json into the same dir as the binary. TO-DO: Fix this so you can copy to the binary folder, or other places like the icons folder in /data
	for _, fn := range pkg.PathsToCopy() {
		srcPath := filepath.Join(pkg.BasePath, fn)
		dstPath := filepath.Join(targetPath, fn)
		err := copyAnything(srcPath, dstPath)

		if err != nil {
			// Can't copy the files? Panic!
			panic(err)
		}

	}

	tplData := &upstartTemplateVars{
		PackageName: pkg.ShortName(),
		TargetDir:   targetDir,
		BaseDir:     deployPath,
	}

	// prepare /etc/init
	initPath := filepath.Join(stagingCurr, "etc", "init")
	os.MkdirAll(initPath, 0755)

	// upstart file
	initConfPath := filepath.Join(initPath, pkg.ShortName()+".conf")
	file, err := os.Create(initConfPath)
	if err != nil {
		panic(err)
	}
	template.Must(template.New("ninjaDebUpstart").Parse(ninjaDebUpstart)).Execute(file, tplData)
	file.Close()

	// post install
	postInstallPath := filepath.Join(ctx.stagingHost, "postinstall-deb-"+arch)
	file, err = os.Create(postInstallPath)
	if err != nil {
		panic(err)
	}
	template.Must(template.New("ninjaDebPostInstall").Parse(ninjaDebPostInstall)).Execute(file, tplData)
	file.Close()

	// prepare the build command
	fpm := []string{
		"fpm", "-s", "dir", "-t", "deb",
		"--deb-compression", "xz",
		"--deb-user", "root",
		"--deb-group", "root",
		"--category", "web",
		"-m", "Ninja CLI Builder <builder@ninjablocks.com>",
		"--url", "http://ninjablocks.com/",
		"-n", pkg.ShortName(),
		"-v", pkg.Version(),
		"--description", pkg.Name(),
		"-C", dockerCurr,
		"-a", arch,
		"-d", "ninjasphere-minimal",
		"-p", filepath.Join(ctx.outputDocker, fmt.Sprintf("%s_%s_%s.deb", pkg.ShortName(), pkg.Version(), arch)),
		"--after-install", filepath.Join(ctx.stagingDocker, "postinstall-deb-"+arch),
		".",
	}

	runNativeDockerCommand(ctx.dockerVolumeArgs, fpm...)
}

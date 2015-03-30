package main

import (
	"fmt"
	"github.com/termie/go-shutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type upstartTemplateVars struct {
	PackageName string
	TargetDir   string
}

func buildDebianPackage(pkg *NinjaPackage, ctx *buildContext, arch string) {
	dockerCurr := filepath.Join(ctx.stagingDocker, "debian-"+arch)
	stagingCurr := filepath.Join(ctx.stagingHost, "debian-"+arch)
	os.MkdirAll(stagingCurr, 0755)

	targetDir := "apps"
	if strings.HasPrefix(pkg.ShortName(), "driver-") {
		targetDir = "drivers"
	}

	targetPath := filepath.Join(stagingCurr, "opt", "ninjablocks", targetDir, pkg.ShortName())
	os.MkdirAll(targetPath, 0755)

	// binary itself
	binCurr := filepath.Join(targetPath, pkg.ShortName())
	shutil.Copy(ctx.archBinaries[arch], binCurr, false)

	// package.json
	srcFile := filepath.Join(pkg.BasePath, "package.json")
	pkgCurr := filepath.Join(stagingCurr, "package.json")
	shutil.Copy(srcFile, pkgCurr, true) // don't copy symlink itself, just the real file

	tplData := &upstartTemplateVars{
		PackageName: pkg.ShortName(),
		TargetDir:   targetDir,
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

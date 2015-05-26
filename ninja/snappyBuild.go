package main

import (
	"github.com/termie/go-shutil"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type snappyPackageMetaService struct {
	Name        string
	Description string
	Start       string
}

type snappyPackageMetaIntegration struct {
	AppArmorProfile string `yaml:"apparmor-profile"`
}

type snappyPackageMeta struct {
	Name         string
	Vendor       string
	Architecture string
	Version      string
	Icon         string
	Frameworks   string
	Services     []snappyPackageMetaService
	Integration  map[string]snappyPackageMetaIntegration `yaml:"integration,omitempty"`
}

func buildSnappyPackage(pkg *NinjaPackage, ctx *buildContext, arch string, arguments map[string]interface{}) {
	dockerCurr := filepath.Join(ctx.stagingDocker, "snappy-"+arch)
	stagingCurr := filepath.Join(ctx.stagingHost, "snappy-"+arch)
	os.MkdirAll(stagingCurr, 0750)

	metaCurr := filepath.Join(stagingCurr, "meta")
	os.MkdirAll(metaCurr, 0750)

	pkgSuffix := ""
	if ns := arguments["--snappy-namespace"]; ns != nil {
		pkgSuffix = "." + ns.(string)
	}

	if arch != "multi" {
		// binary itself
		binCurr := filepath.Join(stagingCurr, pkg.ShortName())
		shutil.Copy(ctx.archBinaries[arch], binCurr, false)

		// package.json
		srcFile := filepath.Join(pkg.BasePath, "package.json")
		pkgCurr := filepath.Join(stagingCurr, "package.json")
		shutil.Copy(srcFile, pkgCurr, true) // don't copy symlink itself, just the real file

		meta := &snappyPackageMeta{
			Name:         pkg.ShortName() + pkgSuffix,
			Vendor:       pkg.Author(),
			Architecture: arch,
			Version:      pkg.Version(),
			Icon:         "meta/null.svg",
			Frameworks:   "ninjasphere",
			Services: []snappyPackageMetaService{
				{
					Name:        pkg.ShortName(),
					Description: pkg.ShortName() + " service",
					Start:       "ninja-shim ./" + pkg.ShortName(),
				},
			},
			Integration: map[string]snappyPackageMetaIntegration{
				pkg.ShortName(): {
					AppArmorProfile: "meta/" + pkg.ShortName() + ".profile",
				},
			},
		}

		metaBytes, err := yaml.Marshal(&meta)
		if err != nil {
			panic(err) // FIXME: meh
		}

		// meta/package.yaml
		metaPackageFile := filepath.Join(metaCurr, "package.yaml")
		ioutil.WriteFile(metaPackageFile, metaBytes, 0644)

		// meta/readme.md
		metaReadmeFile := filepath.Join(metaCurr, "readme.md")
		ioutil.WriteFile(metaReadmeFile, []byte(pkg.Description()), 0644)

		// meta/<pkg-name>.profile
		metaProfileFile := filepath.Join(metaCurr, pkg.ShortName()+".profile")
		ioutil.WriteFile(metaProfileFile, []byte(ninjaAppProfileRediculouslyPermissive), 0644)

		// ninja-shim
		stagingShimFile := filepath.Join(stagingCurr, "ninja-shim")
		ioutil.WriteFile(stagingShimFile, []byte(ninjaLaunchShim), 0755)

		// and all the files specified
		for _, fn := range pkg.PathsToCopy() {
			srcPath := filepath.Join(pkg.BasePath, fn)
			dstPath := filepath.Join(stagingCurr, fn)
			copyAnything(srcPath, dstPath)
		}

		runNativeDockerCommand(ctx.dockerVolumeArgs, "sh", "-c", "cd "+ctx.outputDocker+"; snappy build "+dockerCurr+"")
	} else {
		panic("Multi-arch snaps not supported yet, FIXME add shim wrapper here")
	}
}

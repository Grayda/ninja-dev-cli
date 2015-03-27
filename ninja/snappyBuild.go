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
	Services     []snappyPackageMetaService
	Integration  map[string]snappyPackageMetaIntegration `yaml:"integration,omitempty"`
}

func buildSnappyPackage(pkg *NinjaPackage, ctx *buildContext, arch string) {
	dockerCurr := filepath.Join(ctx.stagingDocker, "snappy-"+arch)
	stagingCurr := filepath.Join(ctx.stagingHost, "snappy-"+arch)
	os.MkdirAll(stagingCurr, 0750)

	metaCurr := filepath.Join(stagingCurr, "meta")
	os.MkdirAll(metaCurr, 0750)

	if arch != "multi" {
		binCurr := filepath.Join(stagingCurr, pkg.ShortName())

		shutil.CopyFile(ctx.archBinaries[arch], binCurr, false)

		meta := &snappyPackageMeta{
			Name:         pkg.ShortName(),
			Vendor:       pkg.Author(),
			Architecture: arch,
			Version:      pkg.Version(),
			Icon:         "meta/null.svg",
			Services: []snappyPackageMetaService{
				{
					Name:        pkg.ShortName(),
					Description: pkg.ShortName() + " service",
					Start:       pkg.ShortName(),
				},
			},
		}

		metaBytes, err := yaml.Marshal(&meta)
		if err != nil {
			panic(err) // FIXME: meh
		}

		metaPackageFile := filepath.Join(metaCurr, "package.yaml")
		ioutil.WriteFile(metaPackageFile, metaBytes, 0644)

		metaReadmeFile := filepath.Join(metaCurr, "readme.md")
		ioutil.WriteFile(metaReadmeFile, []byte(pkg.Description()), 0644)

		runNativeDockerCommand(ctx.dockerVolumeArgs, "sh", "-c", "cd "+ctx.outputDocker+"; snappy build "+dockerCurr+"")
	} else {
		panic("Multi-arch snaps not supported yet, FIXME add shim wrapper here")
	}
}

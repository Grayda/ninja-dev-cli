package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type packageJson struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Files       []string `json:"files"`
}

type NinjaPackage struct {
	BasePath string
	info     packageJson
}

func NewNinjaPackage(basePath string) (*NinjaPackage, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	pkg := &NinjaPackage{absPath, packageJson{}}

	packageJsonFilename := filepath.Join(absPath, "package.json")
	if _, err := os.Stat(packageJsonFilename); os.IsNotExist(err) {
		return nil, errors.New("package.json not found")
	}

	data, err := ioutil.ReadFile(packageJsonFilename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &pkg.info); err != nil {
		return nil, err
	}

	return pkg, nil
}

func (pkg *NinjaPackage) ShortName() string {
	return filepath.Base(pkg.BasePath)
}

func (pkg *NinjaPackage) Name() string {
	if pkg.info.Name == "" {
		return pkg.ShortName()
	}
	return pkg.info.Name
}

func (pkg *NinjaPackage) Version() string {
	return pkg.info.Version
}

func (pkg *NinjaPackage) Author() string {
	return pkg.info.Author
}

func (pkg *NinjaPackage) Description() string {
	return pkg.info.Description
}

func (pkg *NinjaPackage) PathsToCopy() []string {
	return pkg.info.Files
}

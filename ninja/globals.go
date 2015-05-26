package main

import (
	"github.com/termie/go-shutil"
	"os"
)

var verbose = false
var arguments map[string]interface{}
var deployPath string

func copyAnything(from string, to string) error {
	fromStat, err := os.Stat(from)
	if err != nil {
		return err
	}

	if fromStat.IsDir() {
		return shutil.CopyTree(from, to, nil)
	} else {
		_, err = shutil.Copy(from, to, true)
		return err
	}
}

// +build debug

package main

import "io/ioutil"

func init() {
	LoadAsset = ioutil.ReadFile
}

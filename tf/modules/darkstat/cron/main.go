package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/darklab8/fl-data-discovery/autopatcher"
)

var LatestPatch autopatcher.Patch

func GetLatestPatch() (autopatcher.Patch, error) {
	discovery_url := "https://patch.discoverygc.com/"
	discovery_path_url := discovery_url + "patchlist.xml"
	resp := autopatcher.Request(discovery_path_url)
	patches := autopatcher.ParseForPatches(discovery_url, resp.Body)

	if len(patches) > 0 {
		return patches[len(patches)-1], nil
	}
	return autopatcher.Patch{}, errors.New("not found")
}

func main() {

	if value, err := GetLatestPatch(); err == nil {
		LatestPatch = value
		fmt.Println("found patch=", value.Name, value.Hash)
	} else {
		panic("can't grab patch on start")
	}

	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		panic("ENVIRONMENT is not defined")
	}

	for {
		var latest_patch autopatcher.Patch

		if value, err := GetLatestPatch(); err == nil {
			latest_patch = value
		} else {
			fmt.Println("can't grab patch, skipping")
			continue
		}

		if latest_patch != LatestPatch {
			fmt.Println("patch changed. new patch ", latest_patch.Name, latest_patch.Hash, " updating")
			args := strings.Split(fmt.Sprintf("service update --force %s-darkstat-app", environment), " ")
			cmd := exec.Command("docker", args...)
			stdout, err := cmd.Output()
			if err != nil {
				fmt.Println(err.Error())
				panic("can't launch service restart")
			}
			fmt.Println(string(stdout))
			LatestPatch = latest_patch
		} else {
			fmt.Println("patch remained same ", latest_patch.Name, latest_patch.Hash, " skipping update")
		}
		time.Sleep(time.Minute * 30)
	}
}

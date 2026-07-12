package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/darklab8/fl-data-discovery/autopatcher"
	"github.com/darklab8/go-utils/typelog"
)

var LatestPatch autopatcher.Patch

var Log *typelog.Logger = typelog.NewLogger("discoverycron", typelog.WithLogLevel(typelog.LEVEL_INFO))

func GetLatestPatchWeb() (autopatcher.Patch, error) {
	discovery_url := "https://patch.discoverygc.com/"
	discovery_path_url := discovery_url + "patchlist.xml"
	resp, err := autopatcher.Request(discovery_path_url)
	if Log.CheckError(err, "failed to get patchlist.xml") {
		return autopatcher.Patch{}, err
	}
	patches := autopatcher.ParseForPatches(discovery_url, resp.Body)

	if len(patches) == 0 {
		return autopatcher.Patch{}, errors.New("not found")
	}

	return patches[len(patches)-1], nil
}

func GetLatestPatchLocal() (autopatcher.Patch, error) {
	os.Chdir(GameFolderPath)
	patchhistory, _ := autopatcher.ReadLauncherConfig()

	if len(patchhistory.Patches) == 0 {
		return autopatcher.Patch{}, errors.New("not found")
	}

	latst_patch_hash := patchhistory.Patches[len(patchhistory.Patches)-1]

	return autopatcher.Patch{
		Hash: latst_patch_hash,
	}, nil
}

func Autopatcher(workdir string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintln("paniced autopatcher", r))
		}
	}()

	err = autopatcher.RunAutopatcher("/data/freelancer_folder")
	return err
}

func runCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("command %q failed: %w", command, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {

	if value, err := GetLatestPatchLocal(); err == nil {
		LatestPatch = value
		Log.Info(fmt.Sprintln("found local patch=", value.Name, value.Hash))
	} else {
		Log.Panic("can't grab patch on start")
	}

	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		Log.Panic("ENVIRONMENT is not defined")
	}

	for {
		var latest_patch autopatcher.Patch

		if value, err := GetLatestPatchWeb(); err == nil {
			latest_patch = value
		} else {
			Log.Error("can't grab patch, skipping")
			continue
		}

		if latest_patch.Hash != LatestPatch.Hash {
			Log.Info(fmt.Sprintln("patch changed. new patch ", latest_patch.Name, latest_patch.Hash, " updating. But first sleep for 10 minutes"))

			err := Autopatcher(GameFolderPath)

			if Log.CheckError(err, "failed to run autopatcher, sleeping 5 minutes") {
				time.Sleep(time.Minute * 5)
				continue
			}

			cmds := []string{
				fmt.Sprintf("chown -R 1001:1001 %s", GameFolderPath),
			}
			for _, c := range cmds {
				if out, err := runCommand(c); err != nil {
					Log.CheckErrorln(err, "error running ", c, out)
				}
			}

			args := strings.Split(fmt.Sprintf("service update --force %s-darkstat-app", environment), " ")
			cmd := exec.Command("docker", args...)
			stdout, err := cmd.Output()
			if Log.CheckError(err, "can't launch service restart") {
				fmt.Println(string(stdout))
				panic("can't launch server restart")
			}
			LatestPatch = latest_patch
			Log.Info("new patch was applied", typelog.NestedStruct("patch", latest_patch))
		} else {
			Log.Info(fmt.Sprintln("patch remained same ", latest_patch.Name, latest_patch.Hash, " skipping update"))
		}
		time.Sleep(time.Minute * 30)
	}
}

const GameFolderPath = "/data/freelancer_folder"

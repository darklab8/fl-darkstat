package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	darkmap_token := os.Getenv("DARKMAP_REFRESH_GH_TOKEN")
	if darkmap_token == "" {
		Log.Warn("darkmap refresh token is empty")
	}

	args := os.Args[1:]
	command := args[0]

	switch command {
	case "force_map_patch":
		PatchMap(darkmap_token)
	case "main":
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
			var latest_map_patch autopatcher.PatchHash

			if value, err := GetLatestPatchWeb(); err == nil {
				latest_patch = value
				latest_map_patch = value.Hash
			} else {
				Log.Error("can't grab patch, skipping")
				continue
			}

			if latest_patch.Hash != latest_map_patch {
				Log.Info(fmt.Sprintln("map patch changed. new patch ", latest_patch.Name, latest_patch.Hash, " updating."))

				err := PatchMap(darkmap_token)
				if Log.CheckError(err, "failed to trigger darkmap refresh") {

				} else {
					Log.Info("succesfully triggered map patch refresh")
					latest_map_patch = latest_patch.Hash
				}
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

}

func PatchMap(darkmap_token string) error {
	err := triggerWorkflow(
		darkmap_token,
		"darklab8/fl-data-discovery",
		"publish.yaml",
		"fl-darkstat/cron_job",
		"master",
	)
	return err
}

const GameFolderPath = "/data/freelancer_folder"

type dispatchRequest struct {
	Ref    string            `json:"ref"`
	Inputs map[string]string `json:"inputs"`
}

func triggerWorkflow(token, repository, workflowFile, callerRepository, ref string) error {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/actions/workflows/%s/dispatches",
		repository, workflowFile,
	)

	body := dispatchRequest{
		Ref: ref,
		Inputs: map[string]string{
			"repository": callerRepository,
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// GitHub returns 204 No Content on success.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

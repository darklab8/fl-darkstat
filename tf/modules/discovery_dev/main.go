package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"net/http"

	_ "time/tzdata" // load fallback timezone data for use in docker images.

	_ "golang.org/x/crypto/x509roots/fallback" // CA bundle for FROM Scratch
)

func runCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("command %q failed: %w", command, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func getDockerEnv(serviceName string) (map[string]string, error) {
	out, err := runCommand("docker service inspect " + serviceName)
	if err != nil {
		return nil, err
	}

	var services []struct {
		Spec struct {
			TaskTemplate struct {
				ContainerSpec struct {
					Env []string `json:"Env"`
				} `json:"ContainerSpec"`
			} `json:"TaskTemplate"`
		} `json:"Spec"`
	}
	if err := json.Unmarshal([]byte(out), &services); err != nil {
		return nil, fmt.Errorf("failed to parse docker inspect output: %w", err)
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("no service found for %s", serviceName)
	}

	environ := make(map[string]string)
	for _, env := range services[0].Spec.TaskTemplate.ContainerSpec.Env {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			environ[parts[0]] = parts[1]
		}
	}
	return environ, nil
}

func getDockerUpdateState(serviceName string) (string, error) {
	out, err := runCommand("docker service inspect " + serviceName)
	if err != nil {
		return out, err
	}

	var services []struct {
		UpdateStatus struct {
			State string `json:"State"`
		} `json:"UpdateStatus"`
	}
	if err := json.Unmarshal([]byte(out), &services); err != nil {
		return out, fmt.Errorf("failed to parse docker inspect output: %w", err)
	}
	if len(services) == 0 {
		return out, fmt.Errorf("no service found")
	}
	return services[0].UpdateStatus.State, nil
}

func sendDiscordWebhook(webhookURL, username, content string) error {
	payload := map[string]string{
		"username": username,
		"content":  content,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}
	return nil
}

func main() {
	const serviceName = "dev-darkstat-app"
	const repoDir = "freelancer_folder"
	const branch = "5.3"

	oldPassword := "smth"

	for {
		log.Println("Finding out old password...")
		environ, err := getDockerEnv(serviceName)
		if err != nil {
			log.Printf("ERROR getting docker env: %v", err)
		} else {
			oldPassword = environ["DARKCORE_PASSWORD"]
		}

		log.Println("Attempting to clone/update freelancer folder...")
		if _, statErr := os.Stat(repoDir); os.IsNotExist(statErr) {
			if out, err := runCommand("git clone git@git.discoverygc.com:dgcrepository/game-repository.git " + repoDir); err != nil {
				log.Printf("ERROR cloning repo: %v, out=%v", err, out)
				time.Sleep(30 * time.Second)
				continue
			}
		}
		log.Println("Proceeding to git checkout pull reset...")
		cmds := []string{
			fmt.Sprintf("cd %s && git checkout %s", repoDir, branch),
			fmt.Sprintf("cd %s && git pull", repoDir),
			fmt.Sprintf("cd %s && git reset --hard origin/%s", repoDir, branch),
		}
		for _, c := range cmds {
			if out, err := runCommand(c); err != nil {
				log.Printf("ERROR running %q: %v, out=%v", c, err, out)
			}
		}

		log.Println("Getting last commit hash...")

		newPassword, err := runCommand(fmt.Sprintf("cd %s && git rev-parse HEAD", repoDir))
		if err != nil {
			log.Printf("ERROR getting git HEAD: %v", err)
		}

		if oldPassword != newPassword {
			log.Println("Detected new password: ", newPassword, " running docker service update. it takes a minute to that, have patience...")
			updateCmd := fmt.Sprintf(
				`docker service update --env-add DARKCORE_PASSWORD=%s --env-add "FLDARKSTAT_HEADING=commit:%s" %s`,
				newPassword, newPassword, serviceName,
			)
			if out, err := runCommand(updateCmd); err != nil {
				log.Printf("ERROR updating docker service: %v, out=%s", err, out)
			}
			log.Println("finished running docker service update")
			oldPassword = newPassword

			for {
				state, err := getDockerUpdateState(serviceName)
				if err != nil {
					log.Printf("ERROR polling update state: %v", err)
					break
				}
				if !strings.Contains(state, "updating") {
					contentMsg := fmt.Sprintf(
						"https://darkstat-dev.dd84ai.com/?password=%s state=%s",
						newPassword, state,
					)
					if !strings.Contains(contentMsg, "state=completed") {
						contentMsg += " <@370435997974134785>"
					}

					webhookURL := os.Getenv("DISCO_DEV_WEBHOOK")
					if webhookURL == "" {
						log.Println("WARNING: DISCO_DEV_WEBHOOK env var not set")
					} else if err := sendDiscordWebhook(webhookURL, "Darkstat", contentMsg); err != nil {
						log.Printf("ERROR sending Discord webhook: %v", err)
					}
					break
				}
				time.Sleep(10 * time.Second)
			}
		}

		log.Println("completed run. Sleeping")
		time.Sleep(30 * time.Second)
	}
}

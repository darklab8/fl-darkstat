package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"net/http"

	_ "time/tzdata" // load fallback timezone data for use in docker images.

	"github.com/darklab8/go-utils/typelog"
	_ "golang.org/x/crypto/x509roots/fallback" // CA bundle for FROM Scratch

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Log *typelog.Logger = typelog.NewLogger(
	"disco_patch",
	typelog.WithLogLevel(typelog.LEVEL_INFO),
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

var IsAlerting = false

func main() {
	args := os.Args[1:]
	command := args[0]

	switch command {
	case "health":
		tr := &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    10 * time.Second,
			DisableCompression: true,
		}

		client := &http.Client{Transport: tr}
		resp, err := client.Get("http://localhost:8000/metrics")
		Log.CheckPanic(err, "failed to health check")
		if resp.StatusCode != 200 {
			Log.Panic("status code is not 200", typelog.Any("code", resp.StatusCode))
		}
		Log.Debug("service is healthy")
	case "main":
		mux := http.NewServeMux()
		tcp_server := http.Server{Handler: mux}
		var err error
		tcp_listener, err := net.Listen("tcp", "0.0.0.0:8000")
		Log.CheckPanic(err, "http error")
		go func() {
			err := tcp_server.Serve(tcp_listener) // if serving over Http
			Log.CheckFatal(err, "http error")
		}()
		metronom := NewMetronom(mux)
		go metronom.Run()

		const serviceName = "dev-darkstat-app"
		const repoDir = "freelancer_folder"
		const branch = "5.3"

		oldPassword := "smth"

		for {
			Log.Info("Finding out old password...")
			environ, err := getDockerEnv(serviceName)
			if !Log.CheckError(err, "error getting docker env:") {
				oldPassword = environ["DARKCORE_PASSWORD"]
			}

			Log.Info("Attempting to clone/update freelancer folder...")
			if _, statErr := os.Stat(repoDir); os.IsNotExist(statErr) {
				if out, err := runCommand("git clone git@git.discoverygc.com:dgcrepository/game-repository.git " + repoDir); err != nil {
					Log.CheckErrorln(err, "error cloning repo: out=", out)
					time.Sleep(30 * time.Second)
					continue
				}
			}
			Log.Info("Proceeding to git checkout pull reset...")
			cmds := []string{
				fmt.Sprintf("cd %s && git config --global --add safe.directory %s", repoDir, repoDir),
				fmt.Sprintf("cd %s && git checkout %s", repoDir, branch),
				fmt.Sprintf("cd %s && git pull", repoDir),
				fmt.Sprintf("cd %s && git reset --hard origin/%s", repoDir, branch),
				fmt.Sprintf("chown -R 1001:1001 %s", repoDir),
			}
			for _, c := range cmds {
				if out, err := runCommand(c); err != nil {
					Log.CheckErrorln(err, "error running ", c, out)
				}
			}

			Log.Info("Getting last commit hash...")

			newPassword, err := runCommand(fmt.Sprintf("cd %s && git rev-parse HEAD", repoDir))
			Log.CheckError(err, "error getting git head")

			if oldPassword != newPassword {
				Log.Infoln("Detected new password: ", newPassword, " running docker service update. it takes a minute to that, have patience...")
				updateCmd := fmt.Sprintf(
					`docker service update --env-add DARKCORE_PASSWORD=%s --env-add "FLDARKSTAT_HEADING=commit:%s" %s`,
					newPassword, newPassword, serviceName,
				)
				if out, err := runCommand(updateCmd); err != nil {
					Log.Errorln("error updating docker service: ", err, out)
				}
				Log.Info("finished running docker service update")
				oldPassword = newPassword

				for {
					state, err := getDockerUpdateState(serviceName)
					if Log.CheckError(err, "error polling update state") {
						break
					}
					if !strings.Contains(state, "updating") {
						contentMsg := fmt.Sprintf(
							"https://darkstat-dev.dd84ai.com/?password=%s state=%s",
							newPassword, state,
						)
						if !strings.Contains(contentMsg, "state=completed") {

							if !IsAlerting {
								contentMsg += " <@370435997974134785>"
							}
							IsAlerting = true

						} else {
							IsAlerting = false
						}

						webhookURLs := []string{
							os.Getenv("DISCO_DEV_WEBHOOK"),
							os.Getenv("MY_SERVER_WEBHOOK"),
						}

						for index, webhookURL := range webhookURLs {
							if webhookURL == "" {
								Log.Warnln("webhook env var not set for hook #", index)
							} else if err := sendDiscordWebhook(webhookURL, "Darkstat", contentMsg); err != nil {
								Log.CheckError(err, "sending Discord webhook")
							}
						}
						break
					}
					time.Sleep(10 * time.Second)
				}
			}

			Log.Info("completed run. Sleeping")
			time.Sleep(30 * time.Second)
		}
	}
}

//////////////// METRICS STUFF

type Metronom struct {
	Reg *prometheus.Registry
}

func NewMetronom(mux *http.ServeMux) *Metronom {

	newreg := prometheus.NewRegistry()
	reg := prometheus.WrapRegistererWith(
		prometheus.Labels{
			"environment": os.Getenv("ENVIRONMENT"),
			"project":     "darkstat",
			"component":   "dev.data",
		}, newreg)

	metrics := []prometheus.Collector{
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		version.NewCollector("darkstat"),
		typelog.MetricMsgsTotal,
		upTime,
	}

	reg.MustRegister(
		metrics...,
	)

	mux.Handle(
		"/metrics", promhttp.HandlerFor(
			newreg,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			}),
	)

	return &Metronom{
		Reg: newreg,
	}
}

var (
	upTime prometheus.Counter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "darkstat_uptime_seconds",
		Help: "Up time in seconds",
	})
)

func (m *Metronom) Run() {
	for {
		seconds := 60
		upTime.Add(float64(seconds))
		time.Sleep(time.Second * time.Duration(seconds))
	}
}

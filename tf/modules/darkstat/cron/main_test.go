package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestTrigger(t *testing.T) {
	Main()
}

func getFlagOrEnv(flagVal, envKey string) string {
	if flagVal != "" {
		return flagVal
	}
	return os.Getenv(envKey)
}

func Main() {
	var (
		token            = flag.String("token", "", "Fine-grained PAT with actions:write permission (or GH_TOKEN env)")
		repository       = flag.String("repository", "darklab8/fl-data-discovery", "Target repository, e.g. darklab8/fl-data-flsr (or TARGET_REPOSITORY env)")
		workflowFile     = flag.String("workflow-file", "publish.yaml", "Workflow file name, e.g. publish.yaml (or WORKFLOW_FILE env)")
		callerRepository = flag.String("caller-repository", "fl-darkstat/cron_job", "Repository that is triggering this call (or GITHUB_REPOSITORY env)")
		ref              = flag.String("ref", "master", "Git ref to run the workflow on (or REF env)")
	)
	flag.Parse()

	*token = getFlagOrEnv(*token, "DARKMAP_REFRESH_GH_TOKEN")
	*repository = getFlagOrEnv(*repository, "TARGET_REPOSITORY")
	*workflowFile = getFlagOrEnv(*workflowFile, "WORKFLOW_FILE")
	*callerRepository = getFlagOrEnv(*callerRepository, "GITHUB_REPOSITORY")
	if envRef := os.Getenv("REF"); envRef != "" && *ref == "master" {
		*ref = envRef
	}

	var missing []string
	if *token == "" {
		missing = append(missing, "token")
	}
	if *repository == "" {
		missing = append(missing, "repository")
	}
	if *workflowFile == "" {
		missing = append(missing, "workflow-file")
	}
	if *callerRepository == "" {
		missing = append(missing, "caller-repository")
	}
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "missing required inputs: %v\n", missing)
		os.Exit(1)
	}

	fmt.Println("Called by next repository")
	fmt.Println(*callerRepository)

	if err := triggerWorkflow(*token, *repository, *workflowFile, *callerRepository, *ref); err != nil {
		fmt.Fprintf(os.Stderr, "error triggering workflow: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Workflow dispatch triggered successfully.")
}

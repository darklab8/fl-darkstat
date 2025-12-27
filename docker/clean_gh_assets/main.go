package main

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/v80/github"
	"golang.org/x/oauth2"
)

const (
	owner      = "darklab8"
	repo       = "fl-darkstat"
	keepLatest = 8
	dryRun     = false
)

func main() {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("CLEAN_ASSETS_GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.ListOptions{PerPage: 100}
	var allReleases []*github.RepositoryRelease

	for {
		releases, resp, err := client.Repositories.ListReleases(
			ctx,
			owner,
			repo,
			opt,
		)
		if err != nil {
			log.Fatalf("list releases failed: %v", err)
		}

		allReleases = append(allReleases, releases...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	if len(allReleases) <= keepLatest {
		log.Println("No releases to clean up")
		return
	}

	releasesToClean := allReleases[keepLatest:]

	for _, release := range releasesToClean {
		assets := release.Assets
		if len(assets) == 0 {
			continue
		}

		log.Printf(
			"Cleaning release %s (%d assets)",
			release.GetTagName(),
			len(assets),
		)

		for _, asset := range assets {

			if dryRun {
				log.Printf("Would delete asset %s", asset.GetName())
				continue
			}
			_, err := client.Repositories.DeleteReleaseAsset(
				ctx,
				owner,
				repo,
				asset.GetID(),
			)
			if err != nil {
				log.Fatalf(
					"failed to delete asset %s: %v",
					asset.GetName(),
					err,
				)
			}

			log.Printf("  Deleted asset: %s", asset.GetName())
		}
	}

	log.Println("Cleanup completed successfully")
}

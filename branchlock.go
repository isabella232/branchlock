package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	"fmt"

	"github.com/xoom/stash"
)

var (
	stashBaseURL = flag.String("stash-rest-base-url", "https://git.corp.xoom.com:8080", "Stash REST Base URL")
	project      = flag.String("project-key", "", "Stash project key")
	slug         = flag.String("project-slug", "", "Stash project slug")
	userName     = flag.String("username", "", "Stash username")
	password     = flag.String("password", "", "Stash password")
	versionFlag  = flag.Bool("version", false, "Print build info")

	Log *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	buildInfo string
)

func init() {
	flag.Parse()
	if *versionFlag {
		log.Printf("%s\n", buildInfo)
		os.Exit(0)
	}
}

func main() {

	u, err := url.Parse(*stashBaseURL)
	if err != nil {
		Log.Fatalf("%v\n", err)
	}
	client := stash.NewClient(*userName, *password, u)

	branches, err := client.GetBranches(*project, *slug)
	if err != nil {
		Log.Fatalf("%v\n", err)
	}

	fmt.Fprintf(os.Stderr, "Fetching pull requests...open...")
	open, err := client.GetPullRequests(*project, *slug, "OPEN")
	if err != nil {
		Log.Fatalf("%v\n", err)
	}

	fmt.Fprintf(os.Stderr, "merged...")
	merged, err := client.GetPullRequests(*project, *slug, "MERGED")
	if err != nil {
		Log.Fatalf("%v\n", err)
	}

	fmt.Fprintf(os.Stderr, "declined...")
	declined, err := client.GetPullRequests(*project, *slug, "DECLINED")
	if err != nil {
		Log.Fatalf("%v\n", err)
	}
	fmt.Fprintf(os.Stderr, "\n")

	prs := make([]stash.PullRequest, 0)
	prs = append(prs, merged...)
	prs = append(prs, declined...)

	for _, pr := range prs {
		if branch, present := branches[pr.FromRef.DisplayID]; present {
			var branchIsOpenElsewhere bool
			for _, v := range open {
				if v.FromRef.DisplayID == branch.DisplayID {
					branchIsOpenElsewhere = true
					break
				}
			}
			if !branchIsOpenElsewhere && branch.DisplayID != "develop" && branch.DisplayID != "master" {
				fmt.Printf("%s %s\n", branch.DisplayID, pr.State)
			}
		}
	}
}

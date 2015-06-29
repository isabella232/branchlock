package main

import (
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/xoom/stash"
)

var (
	stashBaseURL = flag.String("stash-rest-base-url", "https://git.corp.xoom.com", "Stash REST Base URL")
	project = flag.String("project-key", "", "Stash project key")
	slug = flag.String("project-slug", "", "Stash project slug")
	userName = flag.String("username", "", "Stash username")
	password = flag.String("password", "", "Stash password")
	versionFlag = flag.Bool("version", false, "Print build info")

	Log *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	buildInfo string
)

func init() {
	flag.Parse()
	if *versionFlag {
		log.Printf("%s\n", buildInfo)
		os.Exit(0)
	} else if *stashBaseURL == "" || *project == "" || *slug == "" || *userName == "" || *password == "" {
		Log.Fatal("Usage: branchlock -stash-rest-base-url <Stash URL> -project-key <project> -project-slug <repo> -username <user name> -password <password> [-version true|false]")
	}
}

func main() {

	u, err := url.Parse(*stashBaseURL)
	if err != nil {
		Log.Fatalf("%v\n", err)
	}
	client := stash.NewClient(*userName, *password, u)

	branchPermissions, err := client.GetBranchPermissions(*project, *slug)
	if err != nil {
		Log.Fatalf("%v\n", err)
	}

	log.Printf("%v", branchPermissions)
}

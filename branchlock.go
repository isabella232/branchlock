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
	project = flag.String("project-key", "", "Stash project key (the short one in the URL")
	slug = flag.String("project-slug", "", "Stash project slug (repo)")
	userName = flag.String("username", "", "Stash username")
	password = flag.String("password", "", "Stash password")
	versionFlag = flag.Bool("version", false, "Print build info")
	lock = flag.Bool("lock", false, "Set to true to lock the specified branch, false to unlock it.")
	branch = flag.String("branch", "", "The branch to lock")
	permittedUser = flag.String("permitted-user", "", "The user with permission to write to the branch.")

	Log *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

	buildInfo string
)

func init() {
	flag.Parse()
	if *versionFlag {
		log.Printf("%s\n", buildInfo)
		os.Exit(0)
	} else if *stashBaseURL == "" || *project == "" || *slug == "" || *userName == "" || *password == "" || *branch == "" {
		log.Printf("stash-rest-base-url, project-key, project-slug, username, password, and branch are always required.")
		flag.PrintDefaults()
		os.Exit(1)
	} else if (*lock && *permittedUser == "") {
		log.Printf("If you are locking the branch, you must specify a user with write access to the branch.  This should be your build user.  Someone has to do the release build, right?")
		log.Printf("permitted-user: %+v", *permittedUser)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {

	u, err := url.Parse(*stashBaseURL)
	if err != nil {
		Log.Fatalf("%v\n", err)
	}
	client := stash.NewClient(*userName, *password, u)

	if *lock {
				client.CreateBranchRestriction(*project, *slug, *branch, *permittedUser)
	} else {
		branchRestrictions, err := client.GetBranchRestrictions(*project, *slug)
		if err != nil {
			Log.Fatalf("%v\n", err)
		}

		for _, branchRestriction := range branchRestrictions.BranchRestriction {
			if branchRestriction.Branch.ID == *branch {
				client.DeleteBranchRestriction(*project, *slug, branchRestriction.Id)
			}
		}

	}
}

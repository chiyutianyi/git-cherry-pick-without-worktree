package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatalf("usage: %s <upstream> <commit>\n", os.Args[0])
	}

	upstream := args[1]
	commit := args[2]

	parent := upstream

	// get upstream commit
	upstreamC := getCommit(upstream)
	commitC := getCommit(commit)

	// create temp commit use given commit's parent and upstream's tree
	// so that we will merge given commit's parent..commit and given commit's parent..upstream
	c := exec.Command("git", "commit-tree", upstreamC.treeID, "-p", commitC.parentID, "-m", "temp")
	c.Env = os.Environ()
	out, err := c.CombinedOutput()
	if err != nil {
		log.Fatalf("git commit-tree failed: %v %v", string(out), err)
	}

	tempCommitID := strings.Trim(string(out), "\n")

	// range each commitID and merge
	c = exec.Command("git", "merge-tree", "--write-tree", tempCommitID, commit)
	c.Env = os.Environ()

	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("git merge-tree failed: %v %v", string(out), err)
	}

	// get tree id from the result of merge-tree
	treeID := strings.Trim(string(out), "\n")

	c = exec.Command("git", "commit-tree", treeID, "-p", parent, "-m", commitC.body)
	// use original author and committer
	c.Env = append(
		os.Environ(),
		fmt.Sprintf("GIT_AUTHOR_NAME=%s", commitC.author),
		fmt.Sprintf("GIT_AUTHOR_EMAIL=%s", commitC.authorEmail),
		fmt.Sprintf("GIT_AUTHOR_DATE=%s", commitC.authorDate),
		fmt.Sprintf("GIT_COMMITTER_NAME=%s", commitC.committer),
		fmt.Sprintf("GIT_COMMITTER_EMAIL=%s", commitC.committerEmail),
	)

	out, err = c.CombinedOutput()
	if err != nil {
		log.Fatalf("git commit-tree failed: %v %v", string(out), err)
	}
	// use current merge result for next parent
	parent = strings.Trim(string(out), "\n")

	fmt.Fprint(os.Stdout, parent)
}

type GitCommit struct {
	treeID         string
	parentID       string
	author         string
	authorEmail    string
	authorDate     string
	committer      string
	committerEmail string
	body           string
}

func getCommit(commitID string) *GitCommit {
	c := exec.Command("git", "show", "--pretty=format:%T%n%P%n%an%n%ae%n%ai%n%cn%n%ce%n%B", "-s", commitID)
	c.Env = os.Environ()

	out, err := c.CombinedOutput()
	if err != nil {
		log.Fatalf("git show failed: %v %v", string(out), err)
	}
	outs := strings.SplitN(string(out), "\n", 8)
	return &GitCommit{
		treeID:         outs[0],
		parentID:       outs[1],
		author:         outs[2],
		authorEmail:    outs[3],
		authorDate:     outs[4],
		committer:      outs[5],
		committerEmail: outs[6],
		body:           outs[7],
	}
}

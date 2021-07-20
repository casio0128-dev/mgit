package main

import (
	"fmt"
	"strings"
	"log"
	"flag"
	"os/exec"
)

func main() {
	flag.Parse()

	current, err := getGitCurrentBranchName()
	if err != nil {
		log.Fatal(err)
	}

	commitCommentPrefix := ""
	if strings.HasPrefix(current, "r") {
		tickedId := current[1:]
		commitCommentPrefix = fmt.Sprintf("#%s ", tickedId)
	}

	var commitComment = fmt.Sprintf(`"%s"`, commitCommentPrefix + getCommitMessage(flag.Args()))

	var gitCommand = flag.Args()
	if commitMessageIndex, isCommit := getCommitMessageIndex(flag.Args()); isCommit && commitMessageIndex != -1 {
		gitCommand[commitMessageIndex] = commitComment
	}

	output, err := exec.Command("git", gitCommand...).CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		log.Fatal(err)
	}
}

func getCommitMessageIndex(gitArgs []string) (int, bool) {
	for index, arg := range gitArgs {
		if strings.EqualFold(arg, "commit") {
			if len(gitArgs) >= index + 1 && strings.EqualFold(gitArgs[index + 1], "-m") {
				return index + 2, true
			}
		}
	}
	return -1, false
}

func getCommitMessage(gitArgs []string) string {
	if index, isCommit := getCommitMessageIndex(gitArgs); isCommit {
		return gitArgs[index]
	}
	return ""
}

func getGitCurrentBranchName() (string, error) {
	// git branch | grep -E '^\*' | sed "s/\* //1"
	result, err := exec.Command("git", "branch").Output()

	if err != nil {
		return "", err
	}

	orgCurrentBranchName := string(result)	// exam)  master\n* hogehoge\n  fugafuga
	currentBranchName := strings.Split(orgCurrentBranchName, "* ")[1]	// exam)hogehoge\n  fugafuga
	parsedBranchName := strings.Split(currentBranchName, "\n")[0]	// exam)hogehoge

	return parsedBranchName, nil
}
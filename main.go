package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	flag.Parse()

	// カレントのブランチ名を取得
	current, err := getGitCurrentBranchName()
	if err != nil {
		log.Fatal(err)
	}

	// コミットコメントの接頭辞を取得
	commitCommentPrefix := ""
	if strings.HasPrefix(current, "r") {
		tickedId := current[1:]
		commitCommentPrefix = fmt.Sprintf("#%s ", tickedId)
	}

	// コミットコメントを作成（接頭辞付与済）
	commitComment := fmt.Sprintf(`"%s"`, commitCommentPrefix+getCommitMessage(flag.Args()))
	// Gitコマンドを取得
	gitCommand := flag.Args()
	// コミットメッセージがあった場合に、既存のコミットメッセージと置換（接頭辞付与済）
	if commitMessageIndex, isCommit := getCommitMessageIndex(flag.Args()); isCommit && commitMessageIndex != -1 {
		gitCommand[commitMessageIndex] = commitComment
	}

	// Gitコマンドを実行
	output, err := exec.Command("git", gitCommand...).CombinedOutput()
	fmt.Println(string(output))

	if err != nil {
		log.Fatal(err)
	}
}

func getCommitMessageIndex(gitArgs []string) (int, bool) {
	for index, arg := range gitArgs {
		if strings.EqualFold(arg, "commit") {
			if len(gitArgs) > index+1 && strings.EqualFold(gitArgs[index+1], "-m") {
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

	orgCurrentBranchName := string(result)                            // exam)  master\n* hogehoge\n  fugafuga
	currentBranchName := strings.Split(orgCurrentBranchName, "* ")[1] // exam)hogehoge\n  fugafuga
	parsedBranchName := strings.Split(currentBranchName, "\n")[0]     // exam)hogehoge

	return parsedBranchName, nil
}

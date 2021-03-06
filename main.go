package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"unicode"
)

func main() {
	flag.Parse()

	// Gitコマンドを取得
	gitCommand := flag.Args()

	if _, isCommit := isCommit(flag.Args()); isCommit {
		// カレントのブランチ名を取得
		currentBranchName := getGitCurrentBranchName()

		// コミットコメントの接頭辞を取得
		commitCommentPrefix := ""
		if isTicketIdBranch(currentBranchName) {
			tickedId := currentBranchName[1:]
			commitCommentPrefix = fmt.Sprintf("#%s ", tickedId)
		}

		// コミットコメントを作成（接頭辞付与済）
		commitComment := fmt.Sprintf(`%s`, commitCommentPrefix+getCommitMessage(gitCommand))

		// コミットメッセージがあった場合に、既存のコミットメッセージと置換（接頭辞付与済）
		if commitMessageIndex, isCommit := getCommitMessageIndex(gitCommand); isCommit && commitMessageIndex != -1 {
			gitCommand[commitMessageIndex] = commitComment
		}
	}

	// Gitコマンドを実行
	output, _ := exec.Command("git", gitCommand...).CombinedOutput()
	fmt.Println(string(output))
}

func isCommit(gitArgs []string) (int, bool) {
	for index, arg := range gitArgs {
		if strings.EqualFold(arg, "commit") {
			return index, true
		}
	}
	return -1, false
}

func getCommitMessageIndex(gitArgs []string) (int, bool) {
	if index, isCommit := isCommit(gitArgs); isCommit {
		if len(gitArgs) > index+1 && strings.EqualFold(gitArgs[index+1], "-m") {
			return index + 2, true
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

func getGitCurrentBranchName() (string) {
	// git branch | grep -E '^\*' | sed "s/\* //1"
	result, err := exec.Command("git", "branch").CombinedOutput()

	if err != nil {
		return ""
	}

	orgCurrentBranchName := string(result)                            // exam)  master\n* hogehoge\n  fugafuga
	if strings.Contains(orgCurrentBranchName, "* ") {
		orgCurrentBranchName = strings.Split(orgCurrentBranchName, "* ")[1] // exam)hogehoge\n  fugafuga
	}
	parsedBranchName := strings.Split(orgCurrentBranchName, "\n")[0]     // exam)hogehoge

	return parsedBranchName
}

func isTicketIdBranch(branch string) bool {
	if strings.HasPrefix(branch, "r") {
		for _, b := range branch[1:] {
			if !unicode.IsDigit(b) {
				return false
			}
		}
		return true
	}
	return false
}
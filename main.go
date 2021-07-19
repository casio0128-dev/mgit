package main

import (
	"fmt"
	"os/exec"
)

func main() {
	const getCurrentBranchName = "git symbolic-ref --short HEAD"
	exec.Command(getCurrentBranchName, "")
}
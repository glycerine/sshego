package sshego

import "fmt"

const ProgramName = "sshego-library"

var LAST_GIT_COMMIT_HASH string
var NEAREST_GIT_TAG string
var GIT_BRANCH string
var GO_VERSION string

// SourceVersion returns the git source code version this code was built from.
func SourceVersion() string {
	return fmt.Sprintf("%s commit: %s / nearest-git-"+
		"tag: %s / branch: %s / %s\n",
		ProgramName, LAST_GIT_COMMIT_HASH,
		NEAREST_GIT_TAG, GIT_BRANCH, GO_VERSION)
}

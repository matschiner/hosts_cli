package main

import (
	"os/exec"
	"regexp"
	"strings"
)

func containsPart(haystack, needle string) bool {
	return strings.Contains(haystack, "\t"+needle) || strings.Contains(haystack, needle+"\t")
}

func condenseNewlines(s string) string {
	var re = regexp.MustCompile(`\n+`)
	return re.ReplaceAllString(s, "\n\n")
}

func reverse(a []int) []int {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

func amIRoot() bool {
	cmd := exec.Command("whoami")
	user, err := cmd.Output()
	if err != nil {
		// couldn't determine root due to error, run anyway - the user won't be able
		// to mod anything without root rights anyway
		return true
	}

	return strings.TrimSpace(string(user)) == "root"
}

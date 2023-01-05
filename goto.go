package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var chromeBinaries = [...]string{
	"/usr/bin/google-chrome",
	"/usr/bin/google-chrome-stable",
}

func getChromeBinary() (string, error) {
	for _, bin := range chromeBinaries {
		_, err := os.Stat(bin)
		if err == nil {
			return bin, nil
		}
	}
	return "", errors.New("no workable chrome binary found")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "error: %v\n", fmt.Errorf("specify exactly one go link"))
		os.Exit(-1)
	}

	bin, err := getChromeBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}

	link := os.Args[1]
	if !strings.HasPrefix(link, "go/") {
		link = fmt.Sprintf("go/%s", link)
	}

	cmd := exec.Command(bin, fmt.Sprintf("--app=http://%s", link))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	_ = cmd.Run()
}

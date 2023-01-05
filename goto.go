package main

import (
	"errors"
	"flag"
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
	urlToggle := flag.Bool("u", false, "don't interpret link as go/link")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "error: %v\n", fmt.Errorf("specify exactly one go/link or use '-u URL'"))
		os.Exit(-1)
	}

	bin, err := getChromeBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}

	link := flag.Arg(0)

	if *urlToggle {
		if !strings.HasPrefix(link, "https://") {
			link = fmt.Sprintf("https:///%s", link)
		}
	} else {
		if !strings.HasPrefix(link, "go/") {
			link = fmt.Sprintf("go/%s", link)
		}
		link = fmt.Sprintf("http:///%s", link)
	}

	cmd := exec.Command(bin, fmt.Sprintf("--app=%s", link))
	cmd.Stdin = os.Stdin

	_ = cmd.Run()
}

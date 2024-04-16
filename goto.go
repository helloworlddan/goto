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

const corpAppSuffix = "corp.google.com" // example

func main() {
	urlToggle := flag.Bool("u", false, "don't interpret link as go/link")
	corpToggle := flag.Bool("c", false, "interpret as corp app")
	googleToggle := flag.Bool("g", false, "interpret as google app")
	profileOverride := flag.Int("p", 1, "override chrome profile index")

	flag.Parse()

	bin, err := getChromeBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(-1)
	}

	link := prepLink(flag.Arg(0), *urlToggle, *corpToggle, *googleToggle)

	prompt := fmt.Sprintf("%s --profile-directory='Profile %d'", bin, *profileOverride)
	if len(link) > 0 {
		prompt = fmt.Sprintf("%s --app='%s'", prompt, link)
	}

	// Apparently chrome inspects it's process parent and simply spawning a process directly fails.
	// Spawning from with a shell as parent is somehow OK.
	cmd := exec.Command("/bin/sh", "-c", prompt)
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}

func prepLink(link string, url, corp, google bool) string {
	if strings.HasPrefix(link, "localhost:") {
		return fmt.Sprintf("http://%s", link)
	}
	if url {
		if !strings.HasPrefix(link, "https://") {
			link = fmt.Sprintf("https://%s", link)
		}
		return link
	}
	if corp {
		return fmt.Sprintf("https://%s.%s", link, corpAppSuffix)
	}
	if google {
		return fmt.Sprintf("https://%s.%s", link, "google.com")
	}
	if len(link) > 0 {
		if !strings.HasPrefix(link, "go/") {
			link = fmt.Sprintf("go/%s", link)
		}
		return fmt.Sprintf("http://%s", link)
	}
	return link
}

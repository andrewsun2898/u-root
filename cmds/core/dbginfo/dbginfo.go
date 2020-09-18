package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func ls(re *regexp.Regexp, root, lsargs string) {
	matched := ""
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if re.MatchString(path) {
			matched = fmt.Sprintf("%s %s", matched, path)
		}
		return nil
	})
	if matched != "" {
		log.Printf("DEBUG: ls %s\n %s", lsargs, matched)
	} else {
		log.Printf("DEBUG: ls %s\n No files match!", lsargs)
	}
}

func grep(path, target string) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("DEBUG: grep %s %s\n Unable to open file %s", target, path, path)
		return
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	found := ""
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, target) {
			found = fmt.Sprintf("%s\n%s", found, line)
		}
	}
	if found == "" {
		log.Printf("DEBUG: grep %s %s\n %s not found in file %s", target, path, target, path)
	} else {
		log.Printf("DEBUG: grep: %s %s\n%s", target, path, found)
	}
}

func main() {
	r, _ := regexp.Compile("(\\/dev\\/nvme.*)|(\\/dev\\/sd.*)")
	ls(r, "/dev", "/dev/nvme* /dev/sd*")

	grep("/proc/mounts", "/disk")

	r, _ = regexp.Compile("\\/disk\\/boot\\/grub\\/entry\\-.*\\.cfg")
	ls(r, "/disk/boot/grub", "/disk/boot/grub/entry-*.cfg")

	grep("/disk/boot/grub/entry-1.cfg", "image=")

	ip := exec.Command("ip", "a")
	b, e := ip.Output()
	log.Printf("DEBUG: ip a\n %s, error: %v", string(b), e)

	ip = exec.Command("ip", "route")
	b, e = ip.Output()
	log.Printf("DEBUG: ip route\n%s, error: %v", string(b), e)

	ip = exec.Command("ip", "-6", "route")
	b, e = ip.Output()
	log.Printf("DEBUG: ip -6 route\n%s, error: %v", string(b), e)
}

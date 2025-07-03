package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func cmd(args ...string) string {
	out, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		die("Error: Failed to execute: " + strings.Join(args, " "))
	}
	return strings.TrimSpace(string(out))
}

func checkGit() {
	if cmd("git", "status", "--porcelain") != "" {
		die("Error: There are uncommitted changes. Please commit or stash them before proceeding.")
	}
	if cmd("git", "tag", "--points-at", "HEAD") != "" {
		die("Error: HEAD is already tagged.")
	}
	if cmd("git", "branch", "--show-current") != "main" {
		die("Error: Not on main branch.")
	}
}

func getCurrentVersion() (int, int, int) {
	content, err := os.ReadFile("lara/httpclient.go")
	if err != nil {
		die("Error: Cannot read httpclient.go")
	}

	re := regexp.MustCompile(`sdkVersion:\s*"(\d+)\.(\d+)\.(\d+)"`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) != 4 {
		die("Error: Cannot parse version from httpclient.go")
	}

	major, _ := strconv.Atoi(matches[1])
	minor, _ := strconv.Atoi(matches[2])
	patch, _ := strconv.Atoi(matches[3])

	return major, minor, patch
}

func updateVersionInFile(filename, oldVersion, newVersion string) {
	file, err := os.Open(filename)
	if err != nil {
		die("Error: Cannot open " + filename)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, oldVersion) {
			line = strings.ReplaceAll(line, oldVersion, newVersion)
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		die("Error: Cannot read " + filename)
	}

	output := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(filename, []byte(output), 0644); err != nil {
		die("Error: Cannot write " + filename)
	}
}

func bumpVersion(component string) {
	checkGit()

	major, minor, patch := getCurrentVersion()
	oldVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	switch component {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	default:
		die("Error: Invalid component. Use 'major', 'minor', or 'patch'.")
	}

	newVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)

	// Update version in httpclient.go
	updateVersionInFile("lara/httpclient.go", oldVersion, newVersion)

	// Git operations
	cmd("git", "add", ".")
	cmd("git", "commit", "-m", "v"+newVersion)
	cmd("git", "tag", "-a", "v"+newVersion, "-m", "v"+newVersion)

	fmt.Printf("Tag v%s created.\n", newVersion)
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("Usage: go run version.go [major|minor|patch]")
		os.Exit(1)
	}

	component := args[0]
	if component != "major" && component != "minor" && component != "patch" {
		fmt.Println("Usage: go run version.go [major|minor|patch]")
		os.Exit(1)
	}

	bumpVersion(component)
}

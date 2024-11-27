package main

import (
	"flag"
	"log"
	"os"
	"strings"
)

var logger *log.Logger
var addChangedLine string
var changelogPath string

type Changelog struct {
	Title    string    `json:"title"`
	Notes    []string  `json:"notes"`
	Versions []Version `json:"versions"`
	Refs     []string  `json:"refs"`
}

type Version struct {
	Name            string   `json:"name"`
	Changed         []string `json:"changed"`
	Added           []string `json:"added"`
	Fixed           []string `json:"fixed"`
	Removed         []string `json:"removed"`
	Notes           []string `json:"notes"`
	BreakingChanges []string `json:"breakingchanges"`
}

func init() {
	logger = log.New(os.Stdout, "", 644)
}

func main() {
	// Flags
	flag.StringVar(&addChangedLine, "add-changed", "", "Adds a new Changed entry to the Unreleased section of the changelog")
	flag.StringVar(&changelogPath, "changelog-path", "CHANGELOG.md", "Path to the changelog file")
	flag.Parse()

	data, err := os.ReadFile(changelogPath)
	checkError(err)

	lines := strings.Split(string(data), "\n")
	lines = lines[:len(lines)-1]

	// Convert the markdown file into a changelog struct
	changelog := parseMarkdown(lines)

	if addChangedLine != "" {
		// Add a new Changed entry to the Unreleased section
		logger.Println("Adding new Changed entry to the Unreleased section")
		addChangedEntry(&changelog.Versions[0], addChangedLine)
	}

	// Write the new changelog file
	err = writeChangelogFile(changelog, changelogPath)
	checkError(err)

	logger.Println("Changelog updated successfully")
}

func addChangedEntry(version *Version, entry string) {
	version.Changed = append(version.Changed, entry)
}

func writeChangelogFile(newChangelog Changelog, path string) error {
	lines := []string{}

	lines = append(lines, newChangelog.Title+"\n")
	if len(newChangelog.Notes) > 0 {
		lines = append(lines, newChangelog.Notes...)
	}
	for _, version := range newChangelog.Versions {
		lines = append(lines, version.Name+"\n")

		if len(version.Added) != 0 {
			header := "### Added\n"
			lines = append(lines, header)
			lines = append(lines, version.Added...)
			lines = append(lines, "")
		}
		if len(version.Changed) != 0 {
			header := "### Changed\n"
			lines = append(lines, header)
			lines = append(lines, version.Changed...)
			lines = append(lines, "")
		}
		if len(version.Fixed) != 0 {
			header := "### Fixed\n"
			lines = append(lines, header)
			lines = append(lines, version.Fixed...)
			lines = append(lines, "")
		}
		if len(version.Removed) != 0 {
			header := "### Removed\n"
			lines = append(lines, header)
			lines = append(lines, version.Removed...)
			lines = append(lines, "")
		}
		if len(version.Notes) != 0 {
			header := "### Notes\n"
			lines = append(lines, header)
			lines = append(lines, version.Notes...)
			lines = append(lines, "")
		}
		if len(version.BreakingChanges) != 0 {
			header := "### Breaking Changes\n"
			lines = append(lines, header)
			lines = append(lines, version.BreakingChanges...)
			lines = append(lines, "")
		}
	}

	lines = append(lines, newChangelog.Refs...)
	// Append new line at the end of the file
	lines = append(lines, "")

	newFile := strings.Join(lines, "\n")

	err := os.WriteFile(path, []byte(newFile), 0666)

	return err
}

func parseMarkdown(markdown []string) Changelog {
	newChangelog := Changelog{}
	newChangelog.Versions = []Version{}
	populateVersion := false
	currentlyPopulating := ""
	currentVersion := &Version{}
	for _, line := range markdown {
		if len(line) > 0 {
			header := strings.Split(line, " ")
			// Identify headers
			if header[0] == "#" {
				populateVersion = true
				newChangelog.Title = line
			} else if header[0] == "##" {
				populateVersion = false
				newVersion := Version{}
				newVersion.Name = line
				newChangelog.Versions = append(newChangelog.Versions, newVersion)
			} else if header[0] == "###" {
				populateVersion = true
				currentlyPopulating = header[1]
			} else if strings.Contains(header[1], "https://") {
				populateVersion = false
				newChangelog.Refs = append(newChangelog.Refs, line)
			} else {
				if populateVersion {
					if len(newChangelog.Versions) != 0 {
						currentVersion = &newChangelog.Versions[len(newChangelog.Versions)-1]
					}
					switch currentlyPopulating {
					case "":
						newChangelog.Notes = append(newChangelog.Notes, line)
						// Append a new line if the last char is a dot
						if line[len(line)-1] == '.' {
							newChangelog.Notes = append(newChangelog.Notes, "")
						}
					case "Added":
						currentVersion.Added = append(currentVersion.Added, line)
					case "Changed":
						currentVersion.Changed = append(currentVersion.Changed, line)
					case "Fixed":
						currentVersion.Fixed = append(currentVersion.Fixed, line)
					case "Removed":
						currentVersion.Removed = append(currentVersion.Removed, line)
					case "Notes":
						currentVersion.Notes = append(currentVersion.Notes, line)
					case "Breaking":
						currentVersion.BreakingChanges = append(currentVersion.BreakingChanges, line)
					default:
						continue
					}
				}
			}
		}
	}

	return newChangelog
}

func checkError(e error) {
	if e != nil {
		logger.Panic(e)
	}
}

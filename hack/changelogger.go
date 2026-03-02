package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var logger *log.Logger

type multiFlag []string

func (m *multiFlag) String() string {
	return strings.Join(*m, ", ")
}

func (m *multiFlag) Set(value string) error {
	*m = append(*m, value)
	return nil
}

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
	logger = log.New(os.Stdout, "", 0)
}

func main() {
	var (
		addAdded     multiFlag
		addChanged   multiFlag
		addFixed     multiFlag
		addRemoved   multiFlag
		addNotes     multiFlag
		addBreaking  multiFlag
		changelogPath string
	)

	flag.Var(&addAdded, "add-added", "Adds a new Added entry to the Unreleased section of the changelog")
	flag.Var(&addChanged, "add-changed", "Adds a new Changed entry to the Unreleased section of the changelog")
	flag.Var(&addFixed, "add-fixed", "Adds a new Fixed entry to the Unreleased section of the changelog")
	flag.Var(&addRemoved, "add-removed", "Adds a new Removed entry to the Unreleased section of the changelog")
	flag.Var(&addNotes, "add-notes", "Adds a new Notes entry to the Unreleased section of the changelog")
	flag.Var(&addBreaking, "add-breaking", "Adds a new Breaking Changes entry to the Unreleased section of the changelog")
	flag.StringVar(&changelogPath, "changelog-path", "CHANGELOG.md", "Path to the changelog file")
	flag.Parse()

	data, err := os.ReadFile(changelogPath)
	checkError(err)

	lines := strings.Split(string(data), "\n")
	lines = lines[:len(lines)-1]

	// Convert the markdown file into a changelog struct
	changelog := parseMarkdown(lines)

	if len(changelog.Versions) == 0 {
		logger.Fatal("No versions found in changelog; cannot add entries")
	}

	v := &changelog.Versions[0]
	addEntries(&v.Added, addAdded, "Added")
	addEntries(&v.Changed, addChanged, "Changed")
	addEntries(&v.Fixed, addFixed, "Fixed")
	addEntries(&v.Removed, addRemoved, "Removed")
	addEntries(&v.Notes, addNotes, "Notes")
	addEntries(&v.BreakingChanges, addBreaking, "Breaking Changes")

	// Write the new changelog file
	err = writeChangelogFile(changelog, changelogPath)
	checkError(err)

	logger.Println("Changelog updated successfully")
}

func addEntries(dest *[]string, entries []string, sectionName string) {
	for _, entry := range entries {
		if slices.Contains(*dest, entry) {
			logger.Printf("Skipping duplicate entry in %s: %s", sectionName, entry)
			continue
		}
		logger.Printf("Adding entry to %s: %s", sectionName, entry)
		*dest = append(*dest, entry)
	}
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
			lines = append(lines, "### Added\n")
			lines = append(lines, version.Added...)
			lines = append(lines, "")
		}
		if len(version.Changed) != 0 {
			lines = append(lines, "### Changed\n")
			lines = append(lines, version.Changed...)
			lines = append(lines, "")
		}
		if len(version.Fixed) != 0 {
			lines = append(lines, "### Fixed\n")
			lines = append(lines, version.Fixed...)
			lines = append(lines, "")
		}
		if len(version.Removed) != 0 {
			lines = append(lines, "### Removed\n")
			lines = append(lines, version.Removed...)
			lines = append(lines, "")
		}
		if len(version.Notes) != 0 {
			lines = append(lines, "### Notes\n")
			lines = append(lines, version.Notes...)
			lines = append(lines, "")
		}
		if len(version.BreakingChanges) != 0 {
			lines = append(lines, "### Breaking Changes\n")
			lines = append(lines, version.BreakingChanges...)
			lines = append(lines, "")
		}
	}

	lines = append(lines, newChangelog.Refs...)
	// Append new line at the end of the file
	lines = append(lines, "")

	newFile := strings.Join(lines, "\n")

	// Atomic write: write to a temp file in the same directory, then rename
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".changelog-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName) // no-op if rename succeeded

	if _, err = tmp.WriteString(newFile); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, path)
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
				if len(header) > 1 {
					currentlyPopulating = strings.Join(header[1:], " ")
				}
			} else if len(header) > 1 && strings.Contains(header[1], "https://") {
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
					case "Breaking Changes":
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
		logger.Fatal(e)
	}
}

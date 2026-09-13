// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxpassfile

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// inspired by https://github.com/jackc/pgpassfile

// Entry represents a line in a MX passfile.
type Entry struct {
	Matrixhost string
	Localpart  string
	Domain     string
	Token      string
}

// Passfile is the in memory data structure representing a MX passfile.
type Passfile struct {
	Entries []*Entry
}

// ReadPassfile reads the file at path and parses it into a Passfile.
func readPassfile(path string) (*Passfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParsePassfile(f)
}

// ParsePassfile reads r and parses it into a Passfile.
func ParsePassfile(r io.Reader) (*Passfile, error) {
	passfile := &Passfile{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		entry := parseLine(scanner.Text())
		if entry != nil {
			passfile.Entries = append(passfile.Entries, entry)
		}
	}

	return passfile, scanner.Err()
}

// parseLine parses a line into an *Entry. It returns nil on comment lines or any other unparsable
// line.
func parseLine(line string) *Entry {
	const (
		tmpBackslash = "\r"
		tmpPipe      = "\n"
	)

	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "#") {
		return nil
	}

	line = strings.ReplaceAll(line, `\\`, tmpBackslash)
	line = strings.ReplaceAll(line, `\|`, tmpPipe)

	parts := strings.Split(line, "|")
	if len(parts) != 4 {
		return nil
	}

	// Unescape escaped colons and backslashes
	for i := range parts {
		parts[i] = strings.ReplaceAll(parts[i], tmpBackslash, `\`)
		parts[i] = strings.ReplaceAll(parts[i], tmpPipe, `|`)
		parts[i] = strings.TrimSpace(parts[i])
	}

	return &Entry{
		Matrixhost: parts[0],
		Localpart:  parts[1],
		Domain:     parts[2],
		Token:      parts[3],
	}
}

func superCMP(fileitem, filteritem string) bool {
	if filteritem == "*" {
		return true
	}
	return fileitem == filteritem
}

// FindPassword finds the password for the provided synapsehost, localpart, and domain. An empty
// string will be returned if no match is found.
func (pf *Passfile) FindPassword(matrixhost, localpart, domain string) string {
	for _, e := range pf.Entries {
		if superCMP(e.Matrixhost, matrixhost) &&
			superCMP(e.Localpart, localpart) &&
			superCMP(e.Domain, domain) {
			return e.Token
		}
	}
	return ""
}

func (pf *Passfile) FindPasswordFill(matrixhost, localpart, domain string) *Entry {
	for _, e := range pf.Entries {
		if superCMP(e.Matrixhost, matrixhost) &&
			superCMP(e.Localpart, localpart) &&
			superCMP(e.Domain, domain) {
			return e
		}
	}
	return nil
}

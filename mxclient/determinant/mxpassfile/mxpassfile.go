// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxpassfile

import (
	"bufio"
	"errors"
	"fmt"
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
func ReadPassfile(path string) (*Passfile, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	permissions := fileInfo.Mode().Perm()

	if permissions != 0o400 && permissions != 0o600 {
		return nil, errors.New("Too wide permissions, ignore file")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParsePassfile(f)
}

// ParsePassfile reads r and parses it into a Passfile.
// ignores invalid lines
func ParsePassfile(r io.Reader) (*Passfile, error) {
	passfile := &Passfile{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		entry, err := parseLine(scanner.Text())
		if err == nil && entry != nil {
			passfile.Entries = append(passfile.Entries, entry)
		}
	}

	return passfile, scanner.Err()
}

// ParsePassfile reads r and parses it into a Passfile.
// returns an empty Passfile on errors
func ParsePassfileValidate(r io.Reader) (*Passfile, error) {
	passfile := &Passfile{}

	line := 0

	var parseErrs error

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line++
		entry, err := parseLine(scanner.Text())
		if err != nil {
			if parseErrs == nil {
				parseErrs = fmt.Errorf("invalid line %d", line)
			} else {
				parseErrs = errors.Join(parseErrs, fmt.Errorf("invalid line %d", line))
			}
		} else if entry != nil {
			passfile.Entries = append(passfile.Entries, entry)
		}
	}
	parseErrs = errors.Join(parseErrs, scanner.Err())
	if parseErrs != nil {
		return &Passfile{}, parseErrs
	}
	return passfile, nil
}

var ErrorInvalidMXPassLine = errors.New("invalid MXPassfile line")

// parseLine parses a line into an *Entry. It returns nil on empty or comment lines or an error on unparsable lines.
func parseLine(line string) (*Entry, error) {
	if line == "" {
		return nil, nil
	}

	const (
		tmpBackslash = "\r"
		tmpPipe      = "\n"
	)

	line = strings.TrimSpace(line)

	if strings.HasPrefix(line, "#") {
		return nil, nil
	}

	line = strings.ReplaceAll(line, `\\`, tmpBackslash)
	line = strings.ReplaceAll(line, `\|`, tmpPipe)

	parts := strings.Split(line, "|")
	if len(parts) != 4 {
		return nil, ErrorInvalidMXPassLine
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
	}, nil
}

func superCMP(fileitem, filteritem string) bool {
	if fileitem == "*" || filteritem == "*" {
		return true
	}
	return fileitem == filteritem
}

// FindPassword finds the password for the provided synapsehost, localpart, and domain. An empty
// string will be returned if no match is found.
// search rules:
//
//	search term    matches    file value
//	   value                    value, *
//	     *                      value, *, ""
//	     ""                     *, ""
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

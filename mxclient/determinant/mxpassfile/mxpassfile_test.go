// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only

package mxpassfile

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tokenComp(t *testing.T, expected string, value string) {
	if value != expected {
		t.Fatalf(`token was "%s", expected "%s"`, value, expected)
	}
}

func entryComp(t *testing.T, host, localpart, domain, token string, entry *Entry) {
	if entry.Matrixhost != host || entry.Localpart != localpart || entry.Domain != domain || entry.Token != token {
		t.Fatalf(`entry was "%+v"\nexpected "%s" "%s" "%s" "%s"`, entry, host, localpart, domain, token)
	}
}

// TestParsePassFileFailure verifies ParsePassfileValidate reports an error for a file with invalid lines.
func TestParsePassFileFailure(t *testing.T) {
	buf := bytes.NewBufferString(`# A comment
	test1:5432|larrydb|larry|whatstheidea
	and this is clearly wrong. 🫪
	# so still here?
	break it again`)

	pf, err := ParsePassfileValidate(buf)
	if len(pf.Entries) > 0 {
		t.Fatal(`ParsePassfile returned is not empty`)
	}
	if err == nil {
		t.Fatal(`ParsePassfile returned no error`)
	}
}

// TestParsePassFile verifies a typical passfile is parsed and looked up, including escapes and misses.
func TestParsePassFile(t *testing.T) {
	buf := bytes.NewBufferString(`# A comment
	test1:5432|larrydb|larry|whatstheidea
	test1:5432|moedb|moe|imbecile
	test1:5432|curlydb|curly|nyuknyuknyuk
	test2:5432|*|shemp|heymoe
	test2:5432|*|*|test\\ing\|er
	localhost|*|*|sesam
		`)

	passfile, err := ParsePassfile(buf)
	if err != nil {
		t.Fatalf(`ParsePassfile returned error: "%v"`, err)
	}

	if len(passfile.Entries) != 6 {
		t.Fatalf(`passfile.Entries is "%d", expected 6`, len(passfile.Entries))
	}

	tokenComp(t, "whatstheidea", passfile.FindPassword("test1:5432", "larrydb", "larry"))
	tokenComp(t, "imbecile", passfile.FindPassword("test1:5432", "moedb", "moe"))
	tokenComp(t, `test\ing|er`, passfile.FindPassword("test2:5432", "something", "else"))
	tokenComp(t, "sesam", passfile.FindPassword("localhost", "foo", "bare"))

	tokenComp(t, "", passfile.FindPassword("wrong:5432", "larrydb", "larry"))
	tokenComp(t, "", passfile.FindPassword("test1:wrong", "larrydb", "larry"))
	tokenComp(t, "", passfile.FindPassword("test1:5432", "wrong", "larry"))
	tokenComp(t, "", passfile.FindPassword("test1:5432", "larrydb", "wrong"))
}

// TestParseEmptyInput verifies both parsers accept an empty reader and return no entries.
func TestParseEmptyInput(t *testing.T) {
	pf, err := ParsePassfile(bytes.NewBufferString(""))
	if err != nil {
		t.Fatalf(`ParsePassfile on empty input returned error: "%v"`, err)
	}
	if len(pf.Entries) != 0 {
		t.Fatalf(`ParsePassfile on empty input returned "%d" entries, expected 0`, len(pf.Entries))
	}

	pf, err = ParsePassfileValidate(bytes.NewBufferString(""))
	if err != nil {
		t.Fatalf(`ParsePassfileValidate on empty input returned error: "%v"`, err)
	}
	if len(pf.Entries) != 0 {
		t.Fatalf(`ParsePassfileValidate on empty input returned "%d" entries, expected 0`, len(pf.Entries))
	}
}

// TestParsePassfileValidateValidOnly verifies ParsePassfileValidate succeeds and returns all entries on a fully valid file.
func TestParsePassfileValidateValidOnly(t *testing.T) {
	buf := bytes.NewBufferString(`# A comment
h|u|d|t
h2|u2|d2|t2`)

	pf, err := ParsePassfileValidate(buf)
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}
	if len(pf.Entries) != 2 {
		t.Fatalf(`passfile.Entries is "%d", expected 2`, len(pf.Entries))
	}

	tokenComp(t, "t", pf.FindPassword("h", "u", "d"))
	tokenComp(t, "t2", pf.FindPassword("h2", "u2", "d2"))
}

// TestParsePassfileValidateLineNumbers verifies the reported error mentions only the actually invalid line numbers.
func TestParsePassfileValidateLineNumbers(t *testing.T) {
	buf := bytes.NewBufferString(`# A comment
test1:5432|larrydb|larry|whatstheidea
and this is clearly wrong.
# so still here?
break it again`)

	pf, err := ParsePassfileValidate(buf)
	if len(pf.Entries) > 0 {
		t.Fatal(`ParsePassfileValidate returned a non-empty Passfile`)
	}
	if err == nil {
		t.Fatal(`ParsePassfileValidate returned no error`)
	}

	for _, bad := range []string{"invalid line 3", "invalid line 5"} {
		if !strings.Contains(err.Error(), bad) {
			t.Fatalf(`error "%v" does not mention %q`, err, bad)
		}
	}
	for _, ok := range []string{"invalid line 1", "invalid line 2", "invalid line 4"} {
		if strings.Contains(err.Error(), ok) {
			t.Fatalf(`error "%v" wrongly mentions %q`, err, ok)
		}
	}
}

// TestParsePassfileValidateDiscardsValidLines verifies ParsePassfileValidate drops valid entries when any line is invalid.
func TestParsePassfileValidateDiscardsValidLines(t *testing.T) {
	buf := bytes.NewBufferString(`# A comment
test1:5432|larrydb|larry|whatstheidea
and this is clearly wrong.`)

	pf, err := ParsePassfileValidate(buf)
	if len(pf.Entries) != 0 {
		t.Fatal(`ParsePassfileValidate kept a valid line even though another line was invalid`)
	}
	if err == nil {
		t.Fatal(`ParsePassfileValidate returned no error`)
	}
}

// TestParseBlankLines verifies blank lines are skipped by both parsers without error.
func TestParseBlankLines(t *testing.T) {
	pf, err := ParsePassfile(strings.NewReader("# A comment\n\nh|u|d|t\n"))
	if err != nil {
		t.Fatalf(`ParsePassfile returned error: "%v"`, err)
	}
	if len(pf.Entries) != 1 {
		t.Fatalf(`ParsePassfile returned "%d" entries, expected 1 (blank line skipped)`, len(pf.Entries))
	}

	pf, err = ParsePassfileValidate(strings.NewReader("# A comment\n\nh|u|d|t\n"))
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}
	if len(pf.Entries) != 1 {
		t.Fatalf(`ParsePassfileValidate returned "%d" entries, expected 1 (blank line skipped)`, len(pf.Entries))
	}
}

// TestParseAllEmptyFields verifies a line of four empty fields is a valid entry that matches an all-empty query.
func TestParseAllEmptyFields(t *testing.T) {
	buf := bytes.NewBufferString("|||tok\n")

	pf, err := ParsePassfileValidate(buf)
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}
	if len(pf.Entries) != 1 {
		t.Fatalf(`passfile.Entries is "%d", expected 1`, len(pf.Entries))
	}
	entryComp(t, "", "", "", "tok", pf.Entries[0])
	tokenComp(t, "tok", pf.FindPassword("", "", ""))
}

// TestParseEscapeTricky verifies backslash/pipe escaping edge cases, including a backslash before a real separator.
func TestParseEscapeTricky(t *testing.T) {
	cases := []struct {
		name    string
		line    string
		want    string
		wantErr bool
	}{
		{"escaped pipe alone", "h|u|d|a\\|b", "a|b", false},
		{"lone backslash kept", "h|u|d|a\\b", `a\b`, false},
		{"escaped backslash and pipe", "h|u|d|a\\\\\\|b", `a\|b`, false},
		{"backslash before separator is invalid", "h|u|d|a\\\\|b", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			entry, err := parseLine(c.line)
			if c.wantErr {
				if err == nil {
					t.Fatalf("parseLine(%q) returned no error, expected ErrorInvalidMXPassLine", c.line)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseLine(%q) returned error: %v", c.line, err)
			}
			if entry.Token != c.want {
				t.Fatalf(`parseLine(%q) token was %q, expected %q`, c.line, entry.Token, c.want)
			}
		})
	}

	// an actual carriage return character in a field is replaced by a backslash
	// during unescaping
	entry, err := parseLine("h|u|d|a\rb")
	if err != nil {
		t.Fatalf("parseLine with CR returned error: %v", err)
	}
	entryComp(t, "h", "u", "d", `a\b`, entry)
}

// TestParseHashInsideField verifies '#' is only a comment at line start and is kept as data inside fields.
func TestParseHashInsideField(t *testing.T) {
	buf := bytes.NewBufferString(`# comment # with a hash
host|user|domain|my#hash
   # indented comment
`)

	pf, err := ParsePassfile(buf)
	if err != nil {
		t.Fatalf(`ParsePassfile returned error: "%v"`, err)
	}
	if len(pf.Entries) != 1 {
		t.Fatalf(`passfile.Entries is "%d", expected 1`, len(pf.Entries))
	}

	tokenComp(t, "my#hash", pf.FindPassword("host", "user", "domain"))
}

// TestParseFieldTrimming verifies surrounding whitespace is trimmed from fields while internal spaces are kept.
func TestParseFieldTrimming(t *testing.T) {
	buf := bytes.NewBufferString("  host | user | domain | my token  \n")

	pf, err := ParsePassfileValidate(buf)
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}
	if len(pf.Entries) != 1 {
		t.Fatalf(`passfile.Entries is "%d", expected 1`, len(pf.Entries))
	}

	entryComp(t, "host", "user", "domain", "my token", pf.Entries[0])
}

// TestParseCRLF verifies files with Windows line endings are parsed correctly.
func TestParseCRLF(t *testing.T) {
	buf := bytes.NewBufferString("h|u|d|t\r\nh2|u2|d2|t2\r\n")

	pf, err := ParsePassfile(buf)
	if err != nil {
		t.Fatalf(`ParsePassfile returned error: "%v"`, err)
	}
	if len(pf.Entries) != 2 {
		t.Fatalf(`passfile.Entries is "%d", expected 2`, len(pf.Entries))
	}

	tokenComp(t, "t", pf.FindPassword("h", "u", "d"))
	tokenComp(t, "t2", pf.FindPassword("h2", "u2", "d2"))
}

// TestParserScannerLimit verifies a line longer than the scanner token limit surfaces bufio.ErrTooLong.
func TestParserScannerLimit(t *testing.T) {
	longLine := strings.Repeat("a", 70000)

	if _, err := ParsePassfile(strings.NewReader(longLine)); err == nil || !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf(`ParsePassfile returned error "%v", expected bufio.ErrTooLong`, err)
	}
	if _, err := ParsePassfileValidate(strings.NewReader(longLine)); err == nil || !errors.Is(err, bufio.ErrTooLong) {
		t.Fatalf(`ParsePassfileValidate returned error "%v", expected bufio.ErrTooLong`, err)
	}
}

// TestFindPasswordWildcards verifies '*' in any file field matches any query value.
func TestFindPasswordWildcards(t *testing.T) {
	buf := bytes.NewBufferString(`*|user|domain|hostwild
host|*|domain|localwild
host|user|*|domainwild
*|*|*|allwild`)

	pf, err := ParsePassfileValidate(buf)
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}

	tokenComp(t, "hostwild", pf.FindPassword("x", "user", "domain"))
	tokenComp(t, "localwild", pf.FindPassword("host", "y", "domain"))
	tokenComp(t, "domainwild", pf.FindPassword("host", "user", "z"))
	tokenComp(t, "allwild", pf.FindPassword("anything", "some", "other"))
}

// TestFindPasswordOrdering verifies first-match-wins when both wildcard and specific entries are present.
func TestFindPasswordOrdering(t *testing.T) {
	wildFirst := bytes.NewBufferString(`*|*|*|wild
h|u|d|specific`)
	pf, err := ParsePassfileValidate(wildFirst)
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}
	tokenComp(t, "wild", pf.FindPassword("h", "u", "d"))

	specificFirst := bytes.NewBufferString(`h|u|d|specific
*|*|*|wild`)
	pf, err = ParsePassfileValidate(specificFirst)
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}
	tokenComp(t, "specific", pf.FindPassword("h", "u", "d"))
}

// TestFindPasswordQueryStarWildcard verifies '*' as a query matches any file value, including empty ones.
func TestFindPasswordQueryStarWildcard(t *testing.T) {
	pf, err := ParsePassfileValidate(bytes.NewBufferString(`h|u|d|tok`))
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}

	tokenComp(t, "tok", pf.FindPassword("*", "u", "d"))
	tokenComp(t, "tok", pf.FindPassword("h", "*", "d"))
	tokenComp(t, "tok", pf.FindPassword("h", "u", "*"))
	tokenComp(t, "tok", pf.FindPassword("*", "*", "*"))

	// a '*' query also matches a file entry with empty fields
	pf, err = ParsePassfileValidate(bytes.NewBufferString(`h|||emptytok`))
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}
	tokenComp(t, "emptytok", pf.FindPassword("h", "*", "*"))
	tokenComp(t, "emptytok", pf.FindPassword("*", "*", "*"))
}

// TestFindPasswordEmptyFields verifies empty file fields only match empty query values.
func TestFindPasswordEmptyFields(t *testing.T) {
	buf := bytes.NewBufferString("host|||emptytok\n")

	pf, err := ParsePassfileValidate(buf)
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}

	tokenComp(t, "emptytok", pf.FindPassword("host", "", ""))
	tokenComp(t, "", pf.FindPassword("host", "u", ""))
	tokenComp(t, "", pf.FindPassword("host", "", "u"))
}

// TestFindPasswordEmptyQuery verifies an empty query matches only '*' or empty file values, never a concrete one.
func TestFindPasswordEmptyQuery(t *testing.T) {
	pf, err := ParsePassfileValidate(bytes.NewBufferString(`h|u|d|tok
h|||emptytok
*|u2|d2|wildtok`))
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}

	// empty query matches an empty file field
	tokenComp(t, "emptytok", pf.FindPassword("h", "", ""))

	// empty query matches a file wildcard
	tokenComp(t, "wildtok", pf.FindPassword("", "u2", "d2"))

	// empty query does not match a concrete file value
	tokenComp(t, "", pf.FindPassword("", "u", "d"))
	tokenComp(t, "", pf.FindPassword("h2", "u", "d"))
}

// TestSuperCMP verifies the full search rule matrix.
func TestSuperCMP(t *testing.T) {
	cases := []struct {
		fileitem, filteritem string
		want                 bool
	}{
		{"value", "value", true},
		{"value", "*", true},
		{"value", "", false},
		{"*", "value", true},
		{"*", "*", true},
		{"*", "", true},
		{"", "value", false},
		{"", "*", true},
		{"", "", true},
	}

	for _, c := range cases {
		if got := superCMP(c.fileitem, c.filteritem); got != c.want {
			t.Errorf(`superCMP(%q, %q) = %v, want %v`, c.fileitem, c.filteritem, got, c.want)
		}
	}
}

// TestFindPasswordFill verifies FindPasswordFill returns the matching entry or nil.
func TestFindPasswordFill(t *testing.T) {
	pf, err := ParsePassfileValidate(bytes.NewBufferString(`h|u|d|tok`))
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}

	entryComp(t, "h", "u", "d", "tok", pf.FindPasswordFill("h", "u", "d"))
	entryComp(t, "h", "u", "d", "tok", pf.FindPasswordFill("*", "*", "*"))
	if e := pf.FindPasswordFill("h", "z", "d"); e != nil {
		t.Fatalf("FindPasswordFill returned %+v for a non-matching query, expected nil", e)
	}

	pf, err = ParsePassfileValidate(bytes.NewBufferString(`*|*|*|wild`))
	if err != nil {
		t.Fatalf(`ParsePassfileValidate returned error: "%v"`, err)
	}
	entryComp(t, "*", "*", "*", "wild", pf.FindPasswordFill("x", "y", "z"))
}

// TestReadPassfilePermissions verifies ReadPassfile accepts 0400/0600, rejects wider permissions, and errors on a missing file.
func TestReadPassfilePermissions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "passfile")
	if err := os.WriteFile(path, []byte("h|u|d|t\n"), 0o600); err != nil {
		t.Fatalf(`os.WriteFile failed: "%v"`, err)
	}

	check := func(mode os.FileMode, wantErr bool) {
		if err := os.Chmod(path, mode); err != nil {
			t.Fatalf(`os.Chmod failed: "%v"`, err)
		}
		pf, err := ReadPassfile(path)
		if wantErr {
			if err == nil {
				t.Fatalf("ReadPassfile with %#o returned no error", mode)
			}
			if !strings.Contains(err.Error(), "Too wide") {
				t.Fatalf("ReadPassfile with %#o returned error %q without \"Too wide\"", mode, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("ReadPassfile with %#o returned error: %v", mode, err)
		}
		if len(pf.Entries) != 1 {
			t.Fatalf(`passfile.Entries is "%d", expected 1`, len(pf.Entries))
		}
	}

	check(0o600, false)
	check(0o400, false)
	check(0o644, true)
	check(0o666, true)
	check(0o444, true)
	check(0o100, true)
	check(0o700, true)

	if _, err := ReadPassfile(filepath.Join(dir, "nope")); err == nil {
		t.Fatal("ReadPassfile on a missing file returned no error")
	}
}

// TestParseLineDirect table-tests parseLine on comments, errors, and valid entries.
func TestParseLineDirect(t *testing.T) {
	cases := []struct {
		name    string
		line    string
		want    *Entry
		wantErr error
	}{
		{"comment", "# comment", nil, nil},
		{"indented comment", "  # comment", nil, nil},
		{"empty", "", nil, nil},
		{"whitespace only", "   \t", nil, ErrorInvalidMXPassLine},
		{"no pipe", "a", nil, ErrorInvalidMXPassLine},
		{"too few", "a|b|c", nil, ErrorInvalidMXPassLine},
		{"too many", "a|b|c|d|e|f", nil, ErrorInvalidMXPassLine},
		{"basic", "h|u|d|t", &Entry{"h", "u", "d", "t"}, nil},
		{"empty token", "h|u|d|", &Entry{"h", "u", "d", ""}, nil},
		{"empty middle", "h||d|t", &Entry{"h", "", "d", "t"}, nil},
		{"all empty", "|||", &Entry{"", "", "", ""}, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			entry, err := parseLine(c.line)
			if c.wantErr != nil {
				if err == nil || !errors.Is(err, c.wantErr) {
					t.Fatalf(`parseLine(%q) error was "%v", expected %v`, c.line, err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf(`parseLine(%q) returned error: "%v"`, c.line, err)
			}
			if c.want == nil {
				if entry != nil {
					t.Fatalf(`parseLine(%q) returned %+v, expected nil`, c.line, entry)
				}
				return
			}
			entryComp(t, c.want.Matrixhost, c.want.Localpart, c.want.Domain, c.want.Token, entry)
		})
	}
}

// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only

package mxpassfile

import (
	"bytes"
	"strings"
	"testing"
)

func tokenComp(t *testing.T, expected string, value string) {
	if value != expected {
		t.Fatalf(`token was "%s", expected "%s"`, value, expected)
	}
}

func unescape(s string) string {
	s = strings.Replace(s, `\:`, `:`, -1)
	s = strings.Replace(s, `\\`, `\`, -1)
	return s
}

var passfile = [][]string{
	{"test1:5432", "larrydb", "larry", "whatstheidea"},
	{"test1:5432", "moedb", "moe", "imbecile"},
	{"test1:5432", "curlydb", "curly", "nyuknyuknyuk"},
	{"test2:5432", "*", "shemp", "heymoe"},
	{"test2:5432", "*", "*", `test\\ing\|er`},
	{"localhost", "*", "*", "sesam"},
	{"test3", "", "", "swordfish"}, // user will be filled later
}

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

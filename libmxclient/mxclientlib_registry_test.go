// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"sync"
	"testing"
)

// TestClientRegistryConcurrent exercises addClient/getClient from many
// goroutines. Run with `go test -race` to verify the RWMutex guarding.
func TestClientRegistryConcurrent(t *testing.T) {
	const n = 64

	ids := make([]int, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ids[i] = addClient(&CBClient{})
		}(i)
	}
	wg.Wait()

	seen := make(map[int]bool, n)
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate client id %d", id)
		}
		seen[id] = true
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cli, err := getClient(i)
			if err != nil || cli == nil {
				t.Errorf("getClient(%d) = %v, %v", i, cli, err)
			}
		}(i)
	}
	wg.Wait()

	if _, err := getClient(-1); err == nil {
		t.Error("getClient(-1) expected error")
	}
	if _, err := getClient(n); err == nil {
		t.Error("getClient(n) expected error")
	}
}

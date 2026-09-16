// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxclient

import (
	"fmt"
	"sync"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func TestDirectMapRoundTrip(t *testing.T) {
	log.Logger = zerolog.Nop()
	mxc := &MXClient{_directMap: make(map[id.RoomID][]id.UserID)}

	uid := id.NewUserID("alice", "example.org")
	room := id.RoomID("!room:example.org")
	mxc.AddDirectRoom(uid, room)

	if !mxc.IsDirectRoom(room) {
		t.Fatalf("IsDirectRoom(%q) = false, want true", room)
	}
	if got := mxc.GetUserDM(uid.String()); len(got) != 1 || got[0] != room.String() {
		t.Fatalf("GetUserDM(%q) = %v, want [%s]", uid, got, room)
	}
	chats, ok := mxc._directMapToChats()[uid]
	if !ok || len(chats) != 1 || chats[0] != room {
		t.Fatalf("_directMapToChats() = %v, want %s -> [%s]", chats, uid, room)
	}

	mxc._replaceDirectMap(event.DirectChatsEventContent{})
	if mxc.IsDirectRoom(room) {
		t.Fatalf("IsDirectRoom(%q) = true after _replaceDirectMap(empty)", room)
	}
}

// TestDirectMapConcurrentAccess hammers every direct map helper from multiple
// goroutines. Run with `go test -race` to verify the RWMutex guarding.
func TestDirectMapConcurrentAccess(t *testing.T) {
	log.Logger = zerolog.Nop()
	mxc := &MXClient{_directMap: make(map[id.RoomID][]id.UserID)}

	const (
		writers = 8
		readers = 8
		iters   = 2000
	)

	var wg sync.WaitGroup
	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			uid := id.NewUserID(fmt.Sprintf("user%d", w), "example.org")
			for i := 0; i < iters; i++ {
				room := id.RoomID(fmt.Sprintf("!room%d-%d:example.org", w, i))
				mxc.AddDirectRoom(uid, room)
				mxc._replaceDirectMap(event.DirectChatsEventContent{uid: {room}})
			}
		}(w)
	}
	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				_ = mxc.IsDirectRoom(id.RoomID("!room0-0:example.org"))
				_ = mxc.GetUserDM("@user0:example.org")
				_ = mxc._directMapToChats()
			}
		}()
	}
	wg.Wait()
}

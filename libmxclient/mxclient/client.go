// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package mxclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mxclientlib/determinant/mxpassfile"
	"slices"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mau.fi/util/dbutil"
	_ "go.mau.fi/util/dbutil/litestream"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"maunium.net/go/mautrix/sqlstatestore"
)

type MXClient struct {
	*mautrix.Client
	OnEvent    func(string)
	OnMessage  func(string)
	OnSystem   func(string)
	_directMap map[id.RoomID][]id.UserID
}

func (mxc *MXClient) _onAccountDataDM(ctx context.Context, evt *event.Event) { // event.DirectChatsEventContent
	// TODO
	fmt.Printf("\nTODO: Got event account data dm: %#v\n", evt)
}

func (mxc *MXClient) AddDirectRoom(uid id.UserID, roomid id.RoomID) {
	room, ok := mxc._directMap[roomid]
	if ok {
		if slices.Contains(room, uid) {
			log.Error().Str("room", roomid.String()).Str("uid", uid.String()).Msg("direct mapping already exists.")
		} else {
			log.Warn().Str("room", roomid.String()).Str("uid", uid.String()).Msg("room already a dm.")
			mxc._directMap[roomid] = append(room, uid)
		}
	} else {
		mxc._directMap[roomid] = []id.UserID{uid}
	}
}

func (mxc *MXClient) IsDirectRoom(roomid id.RoomID) bool {
	_, ok := mxc._directMap[roomid]
	return ok
}

func (mxc *MXClient) _loadDirectMap() error {
	var directChats event.DirectChatsEventContent
	err := mxc.GetAccountData(context.Background(), event.AccountDataDirectChats.Type, &directChats)
	if err != nil {
		return err
	}

	new_directMap := make(map[id.RoomID][]id.UserID)
	for uid, rooms := range directChats {
		for _, room := range rooms {
			new_directMap[room] = append(new_directMap[room], uid)
		}
	}
	mxc._directMap = new_directMap
	return nil
}

func (mxc *MXClient) _storeDirectMap() error {
	directChats := make(event.DirectChatsEventContent)

	for room, uids := range mxc._directMap {
		for _, uid := range uids {
			directChats[uid] = append(directChats[uid], room)
		}
	}
	err := mxc.SetAccountData(context.Background(), event.AccountDataDirectChats.Type, &directChats)
	if err != nil {
		return err
	}
	return nil
}

func (mxc *MXClient) GetUserDM(mxid string) []string {
	var res = make([]string, 0)

	for room, uids := range mxc._directMap {
		for _, uid := range uids {
			if uid.String() == mxid {
				res = append(res, room.String())
			}
		}
	}
	return res
}

func (mxc *MXClient) _onEventMember(ctx context.Context, evt *event.Event) {
	if evt.GetStateKey() == mxc.UserID.String() && evt.Content.AsMember().Membership == event.MembershipInvite {
		if evt.Content.AsMember().IsDirect {
			mxc.AddDirectRoom(evt.Sender, evt.RoomID)
			err := mxc._storeDirectMap()
			if err != nil {
				log.Error().Err(err).Msg("failed to store direct chats account data")
			}
		}
		_, err := mxc.JoinRoomByID(ctx, evt.RoomID)
		if err == nil {
			log.Info().
				Str("room_id", evt.RoomID.String()).
				Str("inviter", evt.Sender.String()).
				Msg("Joined room after invite")
		} else {
			log.Error().Err(err).
				Str("room_id", evt.RoomID.String()).
				Str("inviter", evt.Sender.String()).
				Msg("Failed to join room after invite")
		}
	} else if evt.Content.AsMember().Membership == event.MembershipJoin {
		out, err := json.Marshal(evt)
		if err != nil {
			log.Error().Err(err).
				Str("id", evt.ID.String()).
				Str("joiner", evt.Sender.String()).
				Msg("Marshalling error")
			return
		}
		mxc.OnEvent(string(out))
	} else {
		fmt.Printf("\nGot member event: %s\n%#v\n", evt.GetStateKey(), evt)
	}
}

func (mxc *MXClient) _onMessage(ctx context.Context, evt *event.Event) {
	out, err := json.Marshal(map[string]any{"sender": evt.Sender.String(),
		"type":             evt.Type.String(),
		"server_timestamp": evt.Timestamp,
		"id":               evt.ID.String(),
		"roomid":           evt.RoomID.String(),
		"is_direct":        mxc.IsDirectRoom(evt.RoomID),
		"content":          evt.Content.Raw,
		"redacts":          evt.Redacts,
		"unsigned":         evt.Unsigned})

	if err != nil {
		log.Error().Err(err).
			Str("id", evt.ID.String()).
			Str("inviter", evt.Sender.String()).
			Msg("Marshalling error")
		return
	}
	mxc.OnMessage(string(out))

	/*
		log.Info().
			Str("sender", evt.Sender.String()).
			Str("type", evt.Type.String()).
			Str("id", evt.ID.String()).
			Str("roomid", evt.RoomID.String()).
			Str("body", evt.Content.AsMessage().Body).
			Msg("Received message")
		fmt.Printf("Got message: %#v\n", evt)
	*/
}

func (mxc *MXClient) LeaveRoomAndForget(ctx context.Context, room id.RoomID) error {
	_, err := mxc.LeaveRoom(ctx, room)
	if err != nil {
		return err
	}
	if mxc.IsDirectRoom(room) {
		delete(mxc._directMap, room)
		mxc._storeDirectMap()
	}
	_, err = mxc.ForgetRoom(ctx, room)
	if err != nil {
		return err
	}
	return nil
}

func (mxc *MXClient) CreateDM(ctx context.Context, uid id.UserID) (resp *mautrix.RespCreateRoom, err error) {
	req := mautrix.ReqCreateRoom{
		IsDirect: true,
		Preset:   "trusted_private_chat",
		Invite:   []id.UserID{uid},
		InitialState: []*event.Event{{
			Type: event.StateEncryption,
			Content: event.Content{
				Parsed: &event.EncryptionEventContent{
					Algorithm:              id.AlgorithmMegolmV1,
					RotationPeriodMillis:   604800000,
					RotationPeriodMessages: 100,
				},
			}},
		},
	}

	resp, err = mxc.CreateRoom(context.Background(), &req)
	if err != nil {
		return
	}

	mxc.AddDirectRoom(uid, resp.RoomID)
	err = mxc._storeDirectMap()
	return
}

type sendmessage_data struct {
	RoomId  id.RoomID      `json:"roomid"`
	Type    event.Type     `json:"type"`
	Content map[string]any `json:"content"`
}

func (mxc *MXClient) SendRoomMessage(ctx context.Context, data string) (*mautrix.RespSendEvent, error) {
	var smd sendmessage_data
	err := json.Unmarshal([]byte(data), &smd)
	if err != nil {
		return nil, err
	}

	resp, err := mxc.SendMessageEvent(ctx, smd.RoomId, event.EventMessage, smd.Content)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

// NewMXClient creates a new Matrix Client ready for syncing
func NewMXClient(homeserverURL string, userID id.UserID, accessToken string) (*MXClient, error) {
	client, err := mautrix.NewClient(homeserverURL, userID, accessToken)
	if err != nil {
		return nil, err
	}

	// keep this for the import
	client.Log = zerolog.Nop()
	// client.Log = zerolog.New(os.Stdout)
	// client.SyncTraceLog = true

	resp, err := client.Whoami(context.Background())
	if err != nil {
		return nil, err
	}
	client.DeviceID = resp.DeviceID

	//fmt.Printf("Device ID: %s\n", client.DeviceID)

	rawdb, err := dbutil.NewWithDialect("smalbot.db", "sqlite3")
	if err != nil {
		return nil, err
	}

	//fmt.Println("db is offen.")

	syncer, ok := client.Syncer.(*mautrix.DefaultSyncer)
	if !ok {
		return nil, errors.New("panic: syncer implementation error")
	}

	//mxclient.StateStore = mautrix.NewMemoryStateStore()

	stateStore := sqlstatestore.NewSQLStateStore(rawdb, dbutil.NoopLogger, false)
	err = stateStore.Upgrade(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade state db: %w", err)
	}
	client.StateStore = stateStore

	pickleKey := []byte("pickle")

	//cryptoStore := crypto.NewMemoryStore(nil)
	cryptoStore := crypto.NewSQLCryptoStore(rawdb, dbutil.ZeroLogger(log.With().Str("db_section", "crypto").Logger()), "", "", pickleKey)
	err = cryptoStore.DB.Upgrade(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade crypto db: %w", err)
	}

	client.Crypto, err = cryptohelper.NewCryptoHelper(client, pickleKey, cryptoStore)
	if err != nil {
		return nil, err
	}

	err = client.Crypto.Init(context.TODO())
	if err != nil {
		return nil, err
	}

	client.Store = cryptoStore

	mxclient := &MXClient{client, nil, nil, nil, make(map[id.RoomID][]id.UserID)}

	syncer.ParseEventContent = true
	syncer.OnEvent(client.StateStoreSyncHandler)

	syncer.OnEventType(event.EventMessage, mxclient._onMessage)
	syncer.OnEventType(event.StateMember, mxclient._onEventMember)
	syncer.OnEventType(event.AccountDataDirectChats, mxclient._onAccountDataDM)

	mxclient._loadDirectMap()

	return mxclient, nil
}

func CreateClient(storage_path string, url string, userID string, accessToken string) (*MXClient, error) {
	return nil, fmt.Errorf("nope.")
}

func CreateClientPass(mxpassfile_path string, storage_path string, url string, localpart string, domain string) (*MXClient, error) {
	pf, err := mxpassfile.ReadPassfile(mxpassfile_path)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("mxpass pf: '%#v'\n", pf)
	//fmt.Printf("mxpass find: '%s' '%s' '%s'\n", url, localpart, domain)
	e := pf.FindPasswordFill(url, localpart, domain)
	if e != nil {
		//fmt.Printf("mxpass: %#v\n", e)
		return NewMXClient(e.Matrixhost, id.NewUserID(e.Localpart, e.Domain), e.Token)
	}
	return nil, fmt.Errorf("nope.")
}

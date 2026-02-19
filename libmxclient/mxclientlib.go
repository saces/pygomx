// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mxclientlib/mxapi"
	"mxclientlib/mxclient"
	"mxclientlib/mxutils"
	"unsafe"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"
)

/*
#include <stdlib.h>
typedef void (*on_event_handler_ptr) (char*, void*);
typedef void (*on_message_handler_ptr) (char*, void*);
typedef void (*on_sys_handler_ptr) (char*, void*);

static inline void call_c_on_event_handler(on_event_handler_ptr ptr, char* jsonStr, void* pobj) {
    (ptr)(jsonStr, pobj);
}

static inline void call_c_on_message_handler(on_message_handler_ptr ptr, char* jsonStr, void* pobj) {
    (ptr)(jsonStr, pobj);
}

static inline void call_c_on_sys_handler(on_message_handler_ptr ptr, char* jsonStr, void* pobj) {
    (ptr)(jsonStr, pobj);
}

*/
import "C"

var apiCanceled = errors.New("canceled by api call")

/*
matrix client with c callback
*/
type CBClient struct {
	*mxclient.MXClient
	on_event_handler        C.on_event_handler_ptr
	on_event_handler_pobj   unsafe.Pointer
	on_message_handler      C.on_message_handler_ptr
	on_message_handler_pobj unsafe.Pointer
	on_sys_handler          C.on_sys_handler_ptr
	on_sys_handler_pobj     unsafe.Pointer
	syncCancelFunc          context.CancelCauseFunc
}

func (cli *CBClient) OnEvent(s string) {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	C.call_c_on_event_handler(cli.on_event_handler, cStr, cli.on_event_handler_pobj)
}

func (cli *CBClient) OnMessage(s string) {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	C.call_c_on_message_handler(cli.on_message_handler, cStr, cli.on_message_handler_pobj)
}

func (cli *CBClient) OnSystem(s string) {
	cStr := C.CString(s)
	defer C.free(unsafe.Pointer(cStr))
	C.call_c_on_sys_handler(cli.on_sys_handler, cStr, cli.on_sys_handler_pobj)
}

func (cli *CBClient) Set_on_event_handler(fn C.on_event_handler_ptr, pobj unsafe.Pointer) {
	cli.on_event_handler = fn
	cli.on_event_handler_pobj = pobj
}

func (cli *CBClient) Set_on_message_handler(fn C.on_message_handler_ptr, pobj unsafe.Pointer) {
	cli.on_message_handler = fn
	cli.on_message_handler_pobj = pobj
}

func (cli *CBClient) Set_on_sys_handler(fn C.on_sys_handler_ptr, pobj unsafe.Pointer) {
	cli.on_sys_handler = fn
	cli.on_sys_handler_pobj = pobj
}

// NewCClient creates a new Matrix Client ready for syncing
func NewCBClient(homeserverURL string, userID id.UserID, accessToken string) (*CBClient, error) {
	client, err := mxclient.NewMXClient(homeserverURL, userID, accessToken)
	if err != nil {
		return nil, err
	}
	return &CBClient{client, nil, nil, nil, nil, nil, nil, nil}, nil
}

/*
account/client management
*/
var cclients []*CBClient

func getClient(id int) (*CBClient, error) {
	if id < 0 || id >= len(cclients) {
		return nil, fmt.Errorf("index out of bounds: '%d'", id)
	}
	res := cclients[id]
	return res, nil
}

/*
helperse
*/

//export FreeCString
func FreeCString(s *C.char) {
	C.free(unsafe.Pointer(s))
}

/*
cli tools
*/

//export cli_discoverhs
func cli_discoverhs(id *C.char) *C.char {
	mxid := C.GoString(id)
	result := mxutils.DiscoverHS(mxid)
	return C.CString(result)
}

//export cli_mkmxtoken
func cli_mkmxtoken(id *C.char, pw *C.char) *C.char {
	mxid := C.GoString(id)
	mxpw := C.GoString(pw)
	result := mxutils.MkToken(mxid, mxpw)
	return C.CString(result)
}

//export cli_whoami
func cli_whoami(hs *C.char, tk *C.char) *C.char {
	_hs := C.GoString(hs)
	_tk := C.GoString(tk)
	result := mxutils.Whoami(_hs, _tk)
	return C.CString(result)
}

//export cli_accountinfo
func cli_accountinfo(hs *C.char, tk *C.char) *C.char {
	_hs := C.GoString(hs)
	_tk := C.GoString(tk)
	result := mxutils.AccountInfo(_hs, _tk)
	return C.CString(result)
}

//export cli_clearaccount
func cli_clearaccount(hs *C.char, tk *C.char) *C.char {
	_hs := C.GoString(hs)
	_tk := C.GoString(tk)
	result := mxutils.ClearAccount(_hs, _tk)
	return C.CString(result)
}

//export cli_serverinfo
func cli_serverinfo(url *C.char) *C.char {
	_url := C.GoString(url)
	result := mxutils.ServerInfo(_url)
	return C.CString(result)
}

/*
high level api, supports multiple clients
the same handler can be attached to multiple clients
*/

//export apiv0_initialize
func apiv0_initialize() C.int {
	return 0
}

//export apiv0_deinitialize
func apiv0_deinitialize() C.int {
	return 0
}

//export apiv0_discover
func apiv0_discover(mxid *C.char) *C.char {
	result, err := mxapi.Discover(C.GoString(mxid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	return C.CString(result)
}

//export apiv0_login
func apiv0_login(data *C.char) *C.char {
	result, err := mxapi.Login(C.GoString(data))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	return C.CString(result)
}

//export apiv0_createclient
func apiv0_createclient(storage_path *C.char, url *C.char, userID *C.char, accessToken *C.char) *C.char {
	mxclient, err := mxclient.CreateClient(C.GoString(storage_path), C.GoString(url), C.GoString(userID), C.GoString(accessToken))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	client := &CBClient{mxclient, nil, nil, nil, nil, nil, nil, nil}
	cclients = append(cclients, client)
	return C.CString(fmt.Sprintf("{ \"id:\"SUCESS. ID=%d\n", len(cclients)))
}

//export apiv0_createclient_pass
func apiv0_createclient_pass(mxpassfile_path *C.char, storage_path *C.char, url *C.char, localpart *C.char, domain *C.char) *C.char {
	mxclient, err := mxclient.CreateClientPass(C.GoString(mxpassfile_path), C.GoString(storage_path), C.GoString(url), C.GoString(localpart), C.GoString(domain))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	client := &CBClient{mxclient, nil, nil, nil, nil, nil, nil, nil}
	mxclient.OnEvent = client.OnEvent
	mxclient.OnMessage = client.OnMessage
	mxclient.OnSystem = client.OnSystem
	cclients = append(cclients, client)
	out, err := json.Marshal(map[string]any{"id": len(cclients) - 1, "userid": client.UserID.String(), "deviceid": client.DeviceID.String()})
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	s := string(out)
	return C.CString(s)
}

//export apiv0_set_on_event_handler
func apiv0_set_on_event_handler(cid C.int, fn C.on_event_handler_ptr, pobj unsafe.Pointer) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}

	cli.Set_on_event_handler(fn, pobj)
	return C.CString("SUCCESS.")
}

//export apiv0_set_on_message_handler
func apiv0_set_on_message_handler(cid C.int, fn C.on_message_handler_ptr, pobj unsafe.Pointer) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}

	cli.Set_on_message_handler(fn, pobj)
	return C.CString("SUCCESS.")
}

//export apiv0_set_on_sys_handler
func apiv0_set_on_sys_handler(cid C.int, fn C.on_sys_handler_ptr, pobj unsafe.Pointer) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}

	cli.Set_on_sys_handler(fn, pobj)
	return C.CString("SUCCESS.")
}

//export apiv0_startclient
func apiv0_startclient(cid C.int) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}

	ctx, cancel := context.WithCancelCause(context.Background())
	cli.syncCancelFunc = cancel

	err = cli.SyncWithContext(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			cause := context.Cause(ctx)
			if errors.Is(cause, apiCanceled) {
				return C.CString("SUCCESS.")
			} else {
				return C.CString(fmt.Sprintf("ERR: %v", cause))
			}
		}
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	return C.CString("SUCCESS.")
}

//export apiv0_stopclient
func apiv0_stopclient(cid C.int) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	cli.StopSync()
	cli.syncCancelFunc(apiCanceled)

	return C.CString("SUCCESS.")
}

//export apiv0_sendmessage
func apiv0_sendmessage(cid C.int, data *C.char) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	result, err := cli.SendRoomMessage(context.Background(), C.GoString(data))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	out, err := json.Marshal(result)
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	s := string(out)
	return C.CString(s)
}

//export apiv0_leaveroom
func apiv0_leaveroom(cid C.int, roomid *C.char) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	err = cli.LeaveRoomAndForget(context.Background(), id.RoomID(C.GoString(roomid)))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	return C.CString("SUCCESS.")
}

//export apiv0_joinedrooms
func apiv0_joinedrooms(cid C.int) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	resp, err := cli.JoinedRooms(context.Background())
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}

	type roomListItem struct {
		RoomId   id.RoomID `json:"roomid"`
		IsDirect bool      `json:"is_direct"`
	}

	var roomList []roomListItem

	for _, room := range resp.JoinedRooms {
		roomList = append(roomList, roomListItem{RoomId: room, IsDirect: cli.IsDirectRoom(room)})
	}
	out, err := json.Marshal(roomList)
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	s := string(out)
	return C.CString(s)
}

//export apiv0_createroom
func apiv0_createroom(cid C.int, data *C.char) *C.char {
	cli, err := getClient(int(cid))
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}

	var req mautrix.ReqCreateRoom
	err = json.Unmarshal([]byte(C.GoString(data)), &req)
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}

	resp, err := cli.CreateRoom(context.Background(), &req)
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return C.CString(fmt.Sprintf("ERR: %v", err))
	}
	s := string(out)
	return C.CString(s)
}

//export apiv0_removeclient
func apiv0_removeclient(cid C.int) C.int {
	return 0
}

//export apiv0_listclients
func apiv0_listclients() *C.char {
	return C.CString("{}")
}

//export apiv0_getoptions
func apiv0_getoptions(cid C.int) *C.char {
	return C.CString("{}")
}

//export apiv0_setoptions
func apiv0_setoptions(cid C.int, opts *C.char) C.int {
	return 0
}

func main() {}

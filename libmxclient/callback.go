// Copyright (C) 2026 saces@c-base.org
// SPDX-License-Identifier: AGPL-3.0-only
package main

/*
#include <stdlib.h>
typedef void (*on_event_handler_ptr) (char*);
typedef void (*on_message_handler_ptr) (char*);

static inline void call_c_on_event_handler(on_event_handler_ptr ptr, char* jsonStr) {
    (ptr)(jsonStr);
}

static inline void call_c_on_message_handler(on_message_handler_ptr ptr, char* jsonStr) {
    (ptr)(jsonStr);
}
*/
import "C"

var on_event_handler C.on_event_handler_ptr
var on_message_handler C.on_message_handler_ptr

//export register_on_event_handler
func register_on_event_handler(fn C.on_event_handler_ptr) {
	on_event_handler = fn
}

//export register_on_message_handler
func register_on_message_handler(fn C.on_message_handler_ptr) {
	on_message_handler = fn
}

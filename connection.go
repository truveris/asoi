// Copyright (c) 2015 Bertrand Janin <b@janin.com>
// Copyright 2014-2015, Truveris Inc. All Rights Reserved.
// Use of this source code is governed by the ISC license in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// Connection states
const (
	ConnStateInit            = iota
	ConnStateWaitingForHello = iota
	ConnStateLive            = iota
)

// Connection error codes (RFC 1459)
const (
	ErrCodeNoNicknameGiven  = 431
	ErrCodeErroneusNickname = 432
	ErrCodeNicknameInUse    = 433
	ErrCodeNickCollision    = 436
)

var (
	// ConnectionState is the current state of the connection state machine.
	ConnectionState = ConnStateInit

	// reServerMessage is a regexp to parse IRC server messages.
	reServerMessage = regexp.MustCompile(`^:[^ ]+ ([0-9]{2,4}) ([^ ]+) (.*)`)
)

// Send a command to the IRC server.
func sendLine(conn net.Conn, cmd string) {
	cmd = strings.TrimSpace(cmd)
	log.Printf("> %s", cmd)
	fmt.Fprintf(conn, "%s\r\n", cmd)
}

func parseServerMessageCode(line string) int16 {
	tokens := reServerMessage.FindStringSubmatch(line)
	if tokens == nil {
		return 0
	}

	code, err := strconv.ParseInt(tokens[1], 10, 16)
	if err != nil {
		log.Printf("error: invalid server message: bad code (%s) in: %s",
			err.Error(), line)
		return 0
	}

	if tokens[2] != cfg.IRCNickname {
		log.Printf("error: invalid server message: wrong nickname in: %s",
			line)
		return 0
	}

	return int16(code)
}

// Connect to the selected server and join all the specified channels.
func connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", cfg.IRCServer)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

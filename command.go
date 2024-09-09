package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDset Command = "SET"
	CMDget Command = "GET"
)

type Message struct {
	Cmd   Command
	Key   []byte
	Value []byte
	TTL   time.Duration
}

type MSGset struct {
	Key   string
	Value string
	TTL   time.Duration
}

type MSGget struct {
	Key string
}

func parseMessage(raw []byte) (*Message, error) {
	rawString := string(raw)
	parts := strings.Split(rawString, " ")

	if len(parts) < 2 {
		return nil, errors.New("invalid command")
	}

	msg := &Message{
		Cmd: Command(parts[0]),
		Key: []byte(parts[1]),
	}

	if msg.Cmd == CMDset {
		if len(parts) != 4 {
			return nil, errors.New("invalid SET command")
		}

		msg.Value = []byte(parts[2])

		ttl, err := strconv.Atoi(strings.Trim(parts[3], "\n"))
		if err != nil {
			return nil, errors.New("invalid SET command")
		}

		msg.TTL = time.Duration(ttl) * time.Second
	}

	return msg, nil
}

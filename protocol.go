package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Command byte

const (
	CmdNonce Command = iota
	CMDset
	CMDget
	CMDdel
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTl   int
}

func (cmd *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, CMDset)

	binary.Write(buf, binary.LittleEndian, int32(len(cmd.Key)))
	binary.Write(buf, binary.LittleEndian, cmd.Key)

	binary.Write(buf, binary.LittleEndian, int32(len(cmd.Value)))
	binary.Write(buf, binary.LittleEndian, cmd.Value)

	binary.Write(buf, binary.LittleEndian, int32(cmd.TTl))

	return buf.Bytes()
}

func ParseCommand(r io.Reader) any {
	var cmd Command
	binary.Read(r, binary.BigEndian, &cmd)

	switch cmd {
	case CMDset:
		return parseSetCommand(r)
	}

	return nil

}

func parseSetCommand(r io.Reader) *CommandSet {
	cmd := &CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	cmd.Key = make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &cmd.Key)

	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	cmd.Value = make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &cmd.Value)

	var ttl int32
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.TTl = int(ttl)

	return cmd
}

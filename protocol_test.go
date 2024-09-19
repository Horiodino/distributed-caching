package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cmd := &CommandSet{
		Key:   []byte("foo"),
		Value: []byte("bar"),
		TTl:   2,
	}

	r := bytes.NewReader(cmd.Bytes())

	pcmd := ParseCommand(r)

	fmt.Println("---------------------")
	fmt.Println(pcmd)
	assert.Equal(t, cmd, pcmd)
}

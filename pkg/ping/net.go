package ping

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrInvalidResponse = errors.New("invalid response")

type ServerResponse struct {
	Name     string
	Duration int
	Error    error
}

func OpenConnection(name, host string, port, timeout int, respose chan<- ServerResponse) {
	address := fmt.Sprintf("%v:%v", host, port)
	conn, err := net.DialTimeout("tcp", address, time.Millisecond*time.Duration(timeout))
	if err != nil {
		respose <- ServerResponse{
			Name:  name,
			Error: err,
		}
		return
	}
	defer conn.Close()

	buf := make([]byte, 4)
	conn.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))
	connectTime := time.Now()
	_, err = conn.Read(buf)
	duration := time.Since(connectTime)
	if err != nil && err != io.EOF {
		respose <- ServerResponse{
			Name:  name,
			Error: err,
		}
		return
	}

	var opcode uint16
	reader := bytes.NewReader(buf[2:4])
	err = binary.Read(reader, binary.LittleEndian, &opcode)
	if err != nil {
		respose <- ServerResponse{
			Name:  name,
			Error: err,
		}
		return
	}

	if opcode != SMSG_AUTH_CHALLENGE {
		respose <- ServerResponse{
			Name:  name,
			Error: ErrInvalidResponse,
		}
		return
	}

	respose <- ServerResponse{
		Name:     name,
		Duration: int(duration.Milliseconds()),
	}
}

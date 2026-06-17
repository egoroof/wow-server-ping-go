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
var ErrOSTimeout = errors.New("OS goes sleep and causes timeout")

type ServerResponse struct {
	Name     string
	Duration int
	Error    error
}

func OpenConnection(name, ip string, port int, timeout time.Duration, respose chan<- ServerResponse) {
	address := fmt.Sprintf("%v:%v", ip, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		respose <- ServerResponse{
			Name:  name,
			Error: err,
		}
		return
	}
	defer conn.Close()

	buf := make([]byte, 4)
	conn.SetDeadline(time.Now().Add(timeout))
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

	// OS can goes sleep
	if duration > timeout*2 {
		respose <- ServerResponse{
			Name:  name,
			Error: ErrOSTimeout,
		}
		return
	}

	respose <- ServerResponse{
		Name:     name,
		Duration: int(duration.Milliseconds()),
	}
}

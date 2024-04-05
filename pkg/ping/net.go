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

func OpenConnection(host string, port, timeout int) (int, error) {
	address := fmt.Sprintf("%v:%v", host, port)
	conn, err := net.DialTimeout("tcp", address, time.Millisecond*time.Duration(timeout))
	connectTime := time.Now()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	buf := make([]byte, 4)
	conn.SetDeadline(connectTime.Add(time.Millisecond * time.Duration(timeout)))
	_, err = conn.Read(buf)
	responseTime := time.Now()
	if err != nil && err != io.EOF {
		return 0, err
	}

	var opcode uint16
	reader := bytes.NewReader(buf[2:4])
	err = binary.Read(reader, binary.LittleEndian, &opcode)
	if err != nil {
		return 0, err
	}

	if opcode != SMSG_AUTH_CHALLENGE {
		return 0, ErrInvalidResponse
	}

	res := responseTime.Sub(connectTime).Milliseconds()
	return int(res), nil
}

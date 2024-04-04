package ping

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type result struct {
	Status           string
	ResponseDuration time.Duration
}

func OpenConnection(host string, port, timeout int) (result, error) {
	address := fmt.Sprintf("%v:%v", host, port)
	conn, err := net.DialTimeout("tcp", address, time.Millisecond*time.Duration(timeout))
	connectTime := time.Now()
	if err != nil {
		return result{}, err
	}
	defer conn.Close()

	buf := make([]byte, 4)
	conn.SetDeadline(connectTime.Add(time.Millisecond * time.Duration(timeout)))
	_, err = conn.Read(buf)
	responseTime := time.Now()
	if err != nil && err != io.EOF {
		return result{}, err
	}

	var opcode uint16
	reader := bytes.NewReader(buf[2:4])
	err = binary.Read(reader, binary.LittleEndian, &opcode)
	if err != nil {
		return result{}, err
	}

	responseDuration := responseTime.Sub(connectTime).Round(time.Millisecond)

	status := "fail"
	if opcode == SMSG_AUTH_CHALLENGE {
		status = "success"
	}

	return result{
		Status:           status,
		ResponseDuration: responseDuration,
	}, nil
}

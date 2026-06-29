package ping

import (
	"bytes"
	"errors"
	"net"
	"time"
)

var ErrInvalidResponse = errors.New("invalid response")
var ErrOSTimeout = errors.New("OS goes sleep and causes timeout")

type ServerResponse struct {
	Name  string
	Group string

	Duration int
	Error    error
}

func OpenConnection(
	name, group, address string,
	timeout time.Duration,
	respose chan<- ServerResponse,
) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		respose <- ServerResponse{
			Name:  name,
			Group: group,
			Error: err,
		}
		return
	}
	defer conn.Close()

	buf := make([]byte, 64)
	conn.SetDeadline(time.Now().Add(timeout))
	connectTime := time.Now()
	bytesRead, err := conn.Read(buf)
	duration := time.Since(connectTime)
	if err != nil {
		respose <- ServerResponse{
			Name:  name,
			Group: group,
			Error: err,
		}
		return
	}

	SMSG_AUTH_CHALLENGE := []byte{
		0, 42, // BE size
		236, 1, // LE opcode 0x1EC
		1, 0, 0, 0, // LE server_seed
	}
	if bytesRead != 44 || !bytes.Equal(SMSG_AUTH_CHALLENGE, buf[0:8]) {
		respose <- ServerResponse{
			Name:  name,
			Group: group,
			Error: ErrInvalidResponse,
		}
		return
	}

	// OS can goes sleep
	if duration > timeout*2 {
		respose <- ServerResponse{
			Name:  name,
			Group: group,
			Error: ErrOSTimeout,
		}
		return
	}

	respose <- ServerResponse{
		Name:     name,
		Group:    group,
		Duration: int(duration.Milliseconds()),
	}
}

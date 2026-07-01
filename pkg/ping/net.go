package ping

import (
	"bytes"
	"errors"
	"net"
	"os"
	"time"
)

var ErrInvalidResponse = errors.New("invalid response")
var ErrResponseBodyBig = errors.New("response body too big")

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

	// OS can goes sleep
	if duration > timeout*2 {
		respose <- ServerResponse{
			Name:  name,
			Group: group,
			Error: os.ErrDeadlineExceeded,
		}
		return
	}

	if bytesRead >= len(buf) {
		respose <- ServerResponse{
			Name:  name,
			Group: group,
			Error: ErrResponseBodyBig,
		}
		return
	}

	// usual response
	SMSG_AUTH_CHALLENGE := []byte{
		0, 42, // BE size
		236, 1, // LE opcode 0x1EC SMSG_AUTH_CHALLENGE
		1, 0, 0, 0, // LE unknown1
		// 4x LE server_seed
		// 32x seed
	}
	// response when our ip is blocked
	// we still can measure duration
	// can temporarily happen when trying to login with wrong username/password
	SMSG_AUTH_RESPONSE := []byte{
		0, 3, // BE size
		238, 1, // LE opcode 0x1EE SMSG_AUTH_RESPONSE
		14, // result AUTH_REJECT
	}
	if bytes.Equal(SMSG_AUTH_CHALLENGE, buf[0:8]) || bytes.Equal(SMSG_AUTH_RESPONSE, buf[0:5]) {
		respose <- ServerResponse{
			Name:     name,
			Group:    group,
			Duration: int(duration.Milliseconds()),
		}
		return
	}

	respose <- ServerResponse{
		Name:  name,
		Group: group,
		Error: ErrInvalidResponse,
	}
}

package client

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Voxelum/minecraft-launcher-core/pkg/discord-rpc/ipc"
)

var logged bool

// Login sends a handshake in the socket and returns an error or nil
func Login(clientid string) error {
	if !logged {
		payload, err := json.Marshal(Handshake{"1", clientid})
		if err != nil {
			return err
		}

		err = ipc.OpenSocket()
		if err != nil {
			return err
		}

		resp, err := ipc.Send(0, string(payload))
		if err != nil {
			return fmt.Errorf("failed to send handshake: %w", err)
		}

		var response Response
		if err := json.Unmarshal([]byte(resp), &response); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		if response.Event == "ERROR" {
			return fmt.Errorf("discord error: %v", response.Data)
		}

		var handshakeResp HandshakeResponse
		if respData, err := json.Marshal(response.Data); err == nil {
			if err := json.Unmarshal(respData, &handshakeResp); err != nil {
				return fmt.Errorf("failed to parse handshake response: %w", err)
			}
		}
	}
	logged = true

	return nil
}

func Logout() {
	logged = false

	err := ipc.CloseSocket()
	if err != nil {
		panic(err)
	}
}

func SetActivity(activity PayloadActivity) error {
	if !logged {
		return nil
	}

	payload, err := json.Marshal(Frame{
		Cmd: "SET_ACTIVITY",
		Args: Args{
			Pid:      os.Getpid(),
			Activity: &activity,
		},
		Nonce: getNonce(),
	})

	if err != nil {
		return err
	}

	resp, err := ipc.Send(1, string(payload))
	if err != nil {
		return fmt.Errorf("failed to send activity: %w", err)
	}

	var response Response
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Event == "ERROR" {
		return fmt.Errorf("discord error: %v", response.Data)
	}

	return nil
}

func getNonce() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

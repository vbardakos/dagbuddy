package rpc

import (
	"encoding/json"
	"fmt"
	"os"
)

func EncodeMessage(msg Message) ([]byte, error) {
	env := Envelope{Version: rpcVersion}
	msg.marshal(&env)

	data, err := json.Marshal(msg)

	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func DecodeMessage(data []byte) (Message, error) {
	env := Envelope{}
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "Envelope: %+v", env)

	if env.Version != rpcVersion {
		return nil, fmt.Errorf("Expected version: %s. Got: %s", rpcVersion, env.Version)
	}

	if env.ID == nil {
		return NotificationMessage{
			Version: env.Version,
			Method:  env.Method,
			Params:  env.Params,
		}, nil
	}

	id, err := newID(env.ID)
	if err != nil {
		return nil, err
	}

	// Method -> Request | Notification
	if env.Method != "" {
		return RequestMessage{
			Version: env.Version,
			ID:      id,
			Method:  env.Method,
			Params:  env.Params,
		}, nil
	}

	if env.Error != nil {
		return nil, fmt.Errorf("Response Error: %s\n", env.Error)
	}

	return ResponseMessage{
		Version: env.Version,
		ID:      id,
		Result:  env.Result,
		Error:   env.Error,
	}, nil
}

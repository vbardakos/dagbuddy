package rpc

import (
	"encoding/json"
	"fmt"
)

func EncodeMessage(msg Message) ([]byte, error) {
	env := versionedEnvelope()
	msg.marshal(&env)

	data, err := json.Marshal(env)

	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func DecodeMessage(data []byte) (Message, error) {
	env := envelope{}
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}

	if env.Version != rpcVersion {
		return nil, fmt.Errorf("Expected version: %s. Got: %s", rpcVersion, env.Version)
	}

	if env.ID == nil {
		return NotificationMessage{
			Method: env.Method,
			Params: env.Params,
		}, nil
	}

	id, err := NewID(env.ID)
	if err != nil {
		return nil, err
	}

	if env.Method != "" {
		return RequestMessage{
			ID:     id,
			Method: env.Method,
			Params: env.Params,
		}, nil
	}

	if env.Error != nil {
		return nil, env.Error
	}

	return ResponseMessage{
		ID:     id,
		Result: env.Result,
		Error:  env.Error,
	}, nil
}

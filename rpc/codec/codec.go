package codec

import (
	"encoding/json"
	p "github.com/vbardakos/dagbuddy/rpc/protocol"
)

func EncodeMessage(msg p.RPCMessage) ([]byte, error) {
	env := p.NewVersionedEnvelope()
	msg.Marshal(&env)

	data, err := json.Marshal(env)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func DecodeMessage(data []byte) (p.RPCMessage, error) {
	env := p.NewEmptyEnvelope()
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}

	if !env.VersionOk() {
		return nil, p.InternalError
	}

	if env.ID == nil {
		return p.NotificationMessage{
			Method: env.Method,
			Params: env.Params,
		}, nil
	}

	id, err := p.NewID(env.ID)
	if err != nil {
		return nil, err
	}

	if env.Method != "" {
		return p.RequestMessage{
			ID:     id,
			Method: env.Method,
			Params: env.Params,
		}, nil
	}

	if env.Error != nil {
		return nil, env.Error
	}

	if len(env.Result) == 0 {
		return nil, p.InternalError
	}

	return p.ResponseMessage{
		ID:     id,
		Result: env.Result,
		Error:  env.Error,
	}, nil
}

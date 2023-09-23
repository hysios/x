package cache

import "encoding/json"

type jsonEncoder struct {
}

func (j *jsonEncoder) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *jsonEncoder) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

package encoding

type plainUnmarshaler struct{}

func (c *plainUnmarshaler) Unmarshal(b []byte) (interface{}, error) {
	return nil, nil
}

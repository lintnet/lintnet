package encoding

type plainUnmarshaler struct{}

func (c *plainUnmarshaler) Unmarshal(_ []byte) (interface{}, error) {
	return map[string]interface{}{}, nil
}

package encoding

type plainUnmarshaler struct{}

func (c *plainUnmarshaler) Unmarshal(_ []byte) (interface{}, error) {
	return nil, nil //nolint:nilnil
}

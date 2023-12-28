package encoding

type plainUnmarshaler struct{}

func (c *plainUnmarshaler) Unmarshal(_ []byte) (any, error) {
	return nil, nil //nolint:nilnil
}

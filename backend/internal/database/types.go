package database

import "fmt"

type NullableInt int

type NullableString string

func (ni *NullableInt) Scan(value any) error {
	if value == nil {
		test := 0
		*ni = NullableInt(test)
		return nil
	}

	switch v := value.(type) {
	case int64:
		*ni = NullableInt(v)
	case int:
		*ni = NullableInt(v)
	default:
		return fmt.Errorf("cannot scan type %T into NullableInt", value)
	}

	return nil
}

func (ns *NullableString) Scan(value any) error {
	if value == nil {
		*ns = ""
		return nil
	}

	switch v := value.(type) {
	case string:
		*ns = NullableString(v)
	case []uint8:
		*ns = NullableString(v)
	default:
		return fmt.Errorf("cannot scan type %T into NullableInt", value)
	}

	return nil
}

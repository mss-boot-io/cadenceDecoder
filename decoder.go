package cadenceDecoder

import (
	"github.com/onflow/cadence"
)

// Decode decode cadence to json
func Decode(value cadence.Value) (any, error) {
	var err error
	switch value := value.(type) {
	case cadence.Array:
		data := make([]any, len(value.Values))
		for i := range value.Values {
			data[i], err = Decode(value.Values[i])
			if err != nil {
				return data, err
			}
		}
		return data, nil
	case cadence.Optional:
		return Decode(value.Value)
	case cadence.Struct:
		data := make(map[string]any)
		for i := range value.StructType.Fields {
			v, err := Decode(value.Fields[i])
			if err != nil {
				return data, err
			}
			if v != nil {
				data[value.StructType.Fields[i].Identifier] = v
			}
		}
		return data, nil
	case cadence.Dictionary:
		data := make(map[string]any)
		for i := range value.Pairs {
			v, err := Decode(value.Pairs[i].Value)
			if err != nil {
				return data, err
			}
			if v != nil {
				data[replaceBothSideMarks(value.Pairs[i].Key.String())] = v
			}
		}
		return data, nil
	case cadence.Address:
		return value.String(), nil
	case cadence.String:
		return replaceBothSideMarks(value.ToGoValue().(string)), nil
	default:
		if value == nil {
			return nil, nil
		}
		return value.ToGoValue(), nil
	}
}

func replaceBothSideMarks(data string) string {
	if data[0] == '"' {
		data = data[1:]
	}
	if data[len(data)-1] == '"' {
		data = data[:len(data)-1]
	}
	return data
}

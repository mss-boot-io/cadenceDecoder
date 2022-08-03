package cadenceDecoder

import "github.com/onflow/cadence"

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
				data[value.Pairs[i].Key.String()] = v
			}
		}
		return data, nil
	case cadence.Address:
		return value.String(), nil
	default:
		if value == nil {
			return nil, nil
		}
		return value.ToGoValue(), nil
	}
}

package assign

import "errors"

// Assign produces a command to assign a value to a variable within a go session
func Assign(symbol string, value interface{}) ([]byte, error) {
	switch value.(type) {
	case []float64:
		return assignDoubleArray(symbol, value.([]float64))
	case []int32:
		return assignIntArray(symbol, value.([]int32))
	case []string:
		return assignStrArray(symbol, value.([]string))
	case []byte:
		return assignByteArray(symbol, value.([]byte))
	case string:
		return assignStr(symbol, value.(string))
	case int32:
		return assignInt(symbol, value.(int32))
	case float64:
		return assignDouble(symbol, value.(float64))
	default:
		return nil, errors.New("session assign: type is not supported")
	}
}

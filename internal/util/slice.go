package util

func ToInt64Slice[T int | int8 | int16 | int32](input []T) []int64 {
	i64s := make([]int64, len(input))
	for i, value := range input {
		i64s[i] = int64(value)
	}
	return i64s
}

func AnySliceToInt64(input []any) (res []int64, ok bool) {
	i64s := make([]int64, len(input))
	for i, value := range input {
		switch v := value.(type) {
		case int:
			i64s[i] = int64(v)
		case int8:
			i64s[i] = int64(v)
		case int16:
			i64s[i] = int64(v)
		case int32:
			i64s[i] = int64(v)
		case int64:
			i64s[i] = v
		default:
			return nil, false
		}
	}
	return i64s, true
}

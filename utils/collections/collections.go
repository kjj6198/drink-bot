package collections

func Map(
	arr []interface{},
	mapFn func(in interface{}, idx int) interface{},
) []interface{} {
	var result []interface{}
	for i, value := range arr {
		result = append(result, mapFn(value, i))
	}

	return result
}

func ForEach(
	arr []interface{},
	fn func(in interface{}, idx int),
) {
	for i, value := range arr {
		fn(value, i)
	}
}

package grossfilter

var (
	grossfilter_array8 = grossfilter_arrayUntyped
	grossfilter_array16 = grossfilter_arrayUntyped
	grossfilter_array32 = grossfilter_arrayUntyped
	grossfilter_arrayLengthen = grossfilter_arrayLengthenUntyped
	grossfilter_arrayWiden = grossfilter_arrayWidenUntyped
)

func grossfilter_arrayUntyped(n int) ([]interface{}) {
	arr := make([]interface{}, n)
	for i, _ := range arr {
		arr[i] = 0
	}
	return arr
}

func grossfilter_arrayLengthenUntyped(array []int, length int) []int {
	if length < len(array) {
		return array
	}
	arr := make([]int, length)
	for i := len(arr); i < length; i++ {
		arr[i] = 0
	}
	return arr
}

func grossfilter_arrayWidenUntyped(array []interface{}, width int) []interface{} {
	return array;
}

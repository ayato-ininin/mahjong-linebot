package utils

/*
**

	型がfloat64であればint64に変換し、int64の型であればそのまま返す関数
	firestoreから取得したデータで、整数はint64でくるはずだが、
	float64でくることがあるので、その場合はint64に変換する

**
*/
func ConvertInt64(i interface{}) int64 {
	switch i.(type) {
	case float64:
		return int64(i.(float64))
	case int64:
		return i.(int64)
	default:
		return 0
	}
}

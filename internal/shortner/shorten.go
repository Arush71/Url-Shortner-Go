package shortner

func GetCodeFromId(id int64) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 0, 10)
	if id == 0 {
		return string(charset[0])
	}
	for id > 0 {
		remainder := id % 62
		result = append(result, charset[remainder])
		id /= 62
	}
	left := 0
	right := len(result) - 1
	for left < right {
		result[left], result[right] = result[right], result[left]
		left++
		right--
	}
	return string(result)
}

package base62

import "strings"

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Encode(num uint64) string {
	if num == 0 {
		return "0"
	}
	var b strings.Builder
	for num > 0 {
		b.WriteByte(chars[num%62])
		num /= 62
	}
	s := b.String()
	runes := []byte(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func Decode(s string) uint64 {
	var num uint64
	for _, c := range s {
		var idx int
		switch {
		case c >= '0' && c <= '9':
			idx = int(c - '0')
		case c >= 'A' && c <= 'Z':
			idx = int(c-'A') + 10
		case c >= 'a' && c <= 'z':
			idx = int(c-'a') + 36
		default:
			return 0
		}
		num = num*62 + uint64(idx)
	}
	return num
}

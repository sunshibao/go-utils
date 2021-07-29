package stringer

import "github.com/kinwyb/go/xrunes"

// String倒排算法 反转字符串
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// String倒排算法 反转字符串 支持特殊字符,性能稍微比Reverse慢一些
func Reverse2(text string) string {
	textRunes := []rune(text)
	textRunesLength := len(textRunes)
	if textRunesLength <= 1 {
		return text
	}

	i, j := 0, 0
	for i < textRunesLength && j < textRunesLength {
		j = i + 1
		for j < textRunesLength && xrunes.IsMark(textRunes[j]) {
			j++
		}

		if xrunes.IsMark(textRunes[j-1]) {
			// Reverses Combined Characters
			baseReverse(textRunes[i:j], j-i)
		}

		i = j
	}

	// Reverses the entire array
	baseReverse(textRunes, textRunesLength)

	return string(textRunes)
}

func baseReverse(runes []rune, length int) {
	for i, j := 0, length-1; i < length/2; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
}
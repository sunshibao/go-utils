package stringer

import (
	"bytes"
	"github.com/sunshibao/go-utils/xrunes"
	"strings"
)

const (
	space = " "
)

// IsEmpty returns true if the string is empty
func IsEmpty(text string) bool {
	return len(text) == 0
}

// IsNotEmpty returns true if the string is not empty
func IsNotEmpty(text string) bool {
	return !IsEmpty(text)
}

// IsBlank returns true if the string is blank (all whitespace)
func IsBlank(text string) bool {
	return len(strings.TrimSpace(text)) == 0
}

// IsNotBlank returns true if the string is not blank
func IsNotBlank(text string) bool {
	return !IsBlank(text)
}

// Left justifies the text to the left
func Left(text string, size int) string {
	spaces := size - Length(text)
	if spaces <= 0 {
		return text
	}

	var buffer bytes.Buffer
	buffer.WriteString(text)

	for i := 0; i < spaces; i++ {
		buffer.WriteString(space)
	}
	return buffer.String()
}

// Right justifies the text to the right
func Right(text string, size int) string {
	spaces := size - Length(text)
	if spaces <= 0 {
		return text
	}

	var buffer bytes.Buffer
	for i := 0; i < spaces; i++ {
		buffer.WriteString(space)
	}

	buffer.WriteString(text)
	return buffer.String()
}

// Center justifies the text in the center
func Center(text string, size int) string {
	left := Right(text, (Length(text)+size)/2)
	return Left(left, size)
}

// Length counts the input while respecting UTF8 encoding and combined characters
func Length(text string) int {
	textRunes := []rune(text)
	textRunesLength := len(textRunes)

	sum, i, j := 0, 0, 0
	for i < textRunesLength && j < textRunesLength {
		j = i + 1
		for j < textRunesLength && xrunes.IsMark(textRunes[j]) {
			j++
		}
		sum++
		i = j
	}
	return sum
}

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
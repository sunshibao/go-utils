package grand

import (
	"fmt"
	"testing"
)

func TestRandAllString(t *testing.T) {
	tString := RandAllString(56)
	fmt.Println(tString)
}

func TestRandNumString(t *testing.T) {
	tString := RandNumString(56)
	fmt.Println(tString)
}

func TestRandString(t *testing.T) {
	tString := RandString(56)
	fmt.Println(tString)
}

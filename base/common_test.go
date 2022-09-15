package base

import (
	"fmt"
	"testing"
	"time"
)

func TestUUID(t *testing.T) {
	dec := UUID()
	fmt.Println(dec)
}

func TestIDGen(t *testing.T) {
	dec := IDGen()
	fmt.Println(dec)
}

func TestID(t *testing.T) {
	dec := ID()
	fmt.Println(dec)
}

func TestRandCode(t *testing.T) {
	dec := RandCode(15, false)
	fmt.Println(dec)
}

func TestIsEmail(t *testing.T) {
	dec := IsEmail("17521009800@163.com")
	fmt.Println(dec)
}

func TestIsPhone(t *testing.T) {
	dec := IsPhone("17521009800")
	fmt.Println(dec)
}

func TestIsIdCard(t *testing.T) {
	dec := IsIdCard("371321199803022345")
	fmt.Println(dec)
}

func TestIsChineseChar(t *testing.T) {
	dec := IsChineseChar("2001-02-24 你 23:23:22")
	fmt.Println(dec)
}

func TestDecimal(t *testing.T) {
	dec := Decimal(23.455323, 4)
	fmt.Println(dec)
}

func TestDateStartTimeString(t *testing.T) {
	startTime := DateStartTimeString(time.Now())
	fmt.Println(startTime)
}

func TestDateEndTimeString(t *testing.T) {
	endTime := DateEndTimeString(time.Now())
	fmt.Println(endTime)
}

func TestIsNumeric(t *testing.T) {
	numeric := IsNumeric("3.4")
	fmt.Println(numeric)
}

func TestNumberM(t *testing.T) {
	numeric := NumberM(10, 4)
	fmt.Println(numeric)
}

func TestStructFieldAndValueToString(t *testing.T) {
	type AS struct {
		Asj string `json:"asj"`
		Bjh string `json:"bjh"`
	}
	var as AS
	as.Asj = "23"
	as.Bjh = "32"
	toString, s, b := StructFieldAndValueToString(as)

	fmt.Println(toString)
	fmt.Println(s)
	fmt.Println(b)
}

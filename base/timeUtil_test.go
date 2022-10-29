package base

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeToString(t *testing.T) {
	var tt time.Time
	tString := TimeToString(tt, "")
	fmt.Println(tString)
}

func TestTimeToTimeNum(t *testing.T) {
	var tt time.Time
	tTimeNum := TimeToTimeNum(tt)
	fmt.Println(tTimeNum)
}

func TestTimeNumToString(t *testing.T) {
	tString := TimeNumToString(0, "")
	fmt.Println(tString)
}

func TestTimeNumToTime(t *testing.T) {
	tTime := TimeNumToTime(0)
	fmt.Println(tTime)
}

func TestStringToTimeNum(t *testing.T) {
	tTimeNum := StringToTimeNum("", "")
	fmt.Println(tTimeNum)
}

func TestWhatTime(t *testing.T) {
	WhatTime()
}

func TestSubTime(t *testing.T) {
	SubTime()
}

func TestGetTimeDifference(t *testing.T) {
	difference := GetTimeDifference()
	fmt.Println(difference)
}

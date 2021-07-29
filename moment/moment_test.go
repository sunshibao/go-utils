package moment

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	date := time.Date(2019, 9, 29, 12, 34, 15, 0, time.Local)
	if got := FormatDate(date); got != "2019-09-29" {
		t.Fatalf("调用 FormatDate() 返回与预期不符，预期[2019-09-29]，实际[%s]", got)
	}
}

func TestFormatDateTerse(t *testing.T) {
	date := time.Date(2019, 9, 29, 12, 34, 15, 0, time.Local)
	if got := FormatDateTerse(date); got != "20190929" {
		t.Fatalf("调用 FormatDateTerse() 返回与预期不符，预期[20190929]，实际[%s]", got)
	}
}

func TestFormatDateTime(t *testing.T) {
	date := time.Date(2019, 9, 29, 12, 34, 15, 0, time.Local)
	if got := FormatDateTime(date); got != "2019-09-29 12:34:15" {
		t.Fatalf("调用 FormatDateTime() 返回与预期不符，预期[2019-09-29 12:34:15]，实际[%s]", got)
	}
}

func TestFormatDateTimeTerse(t *testing.T) {
	date := time.Date(2019, 9, 29, 12, 34, 15, 0, time.Local)
	if got := FormatDateTimeTerse(date); got != "20190929123415" {
		t.Fatalf("调用 FormatDateTimeTerse() 返回与预期不符，预期[20190929123415]，实际[%s]", got)
	}
}

func TestParseDate(t *testing.T) {
	s := "2019-09-29"
	date, err := ParseDate(s)
	if err != nil {
		t.Fatalf("调用 ParseDate() 返回意外的错误：%s", err)
	}
	if date.Format("20060102150405") != "20190929000000" {
		t.Fatalf("调用 ParseDate() 返回与给定日期不一致")
	}
}

func TestParseDateTerse(t *testing.T) {
	s := "20190929"
	date, err := ParseDateTerse(s)
	if err != nil {
		t.Fatalf("调用 ParseDateTerse() 返回意外的错误：%s", err)
	}
	if date.Format("20060102150405") != "20190929000000" {
		t.Fatalf("调用 ParseDateTerse() 返回与给定日期不一致")
	}
}

func TestParseDateTime(t *testing.T) {
	s := "2019-09-29 12:34:15"
	date, err := ParseDateTime(s)
	if err != nil {
		t.Fatalf("调用 ParseDateTime() 返回意外的错误：%s", err)
	}
	if date.Format("20060102150405") != "20190929123415" {
		t.Fatalf("调用 ParseDateTime() 返回与给定日期不一致")
	}
}

func TestParseDateTimeTerse(t *testing.T) {
	s := "20190929123415"
	date, err := ParseDateTimeTerse(s)
	if err != nil {
		t.Fatalf("调用 ParseDateTimeTerse() 返回意外的错误：%s", err)
	}
	if date.Format("20060102150405") != "20190929123415" {
		t.Fatalf("调用 ParseDateTimeTerse() 返回与给定日期不一致")
	}
}

func TestMustParseDate(t *testing.T) {
	s := "2019-09-29"
	date := MustParseDate(s)
	if date.Format("20060102150405") != "20190929000000" {
		t.Fatalf("调用 MustParseDate() 返回与给定日期不一致")
	}
}

func TestMustParseDateTerse(t *testing.T) {
	s := "20190929"
	date := MustParseDateTerse(s)
	if date.Format("20060102150405") != "20190929000000" {
		t.Fatalf("调用 MustParseDateTerse() 返回与给定日期不一致")
	}
}

func TestMustParseDateTime(t *testing.T) {
	s := "2019-09-29 12:34:15"
	date := MustParseDateTime(s)
	if date.Format("20060102150405") != "20190929123415" {
		t.Fatalf("调用 MustParseDateTime() 返回与给定日期不一致")
	}
}

func TestMustParseDateTimeTerse(t *testing.T) {
	s := "20190929123415"
	date := MustParseDateTimeTerse(s)
	if date.Format("20060102150405") != "20190929123415" {
		t.Fatalf("调用 MustParseDateTimeTerse() 返回与给定日期不一致")
	}
}

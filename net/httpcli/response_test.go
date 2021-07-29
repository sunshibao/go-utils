package httpcli

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

type testErrorReadCloser string

func (s testErrorReadCloser) Read([]byte) (int, error) {
	return 0, errors.New(string(s))
}

func (s testErrorReadCloser) Close() error {
	return nil
}

func TestNewResponse(t *testing.T) {
	r := &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("body"))}
	if res, err := NewResponse(r, false); err != nil {
		t.Fatalf("NewResponse() error: %s", err)
	} else {
		if res == nil {
			t.Fatal("NewResponse() return nil")
		}
	}

	r.Body = testErrorReadCloser("err")
	if _, err := NewResponse(r, false); err == nil {
		t.Fatal("NewResponse() with ErrorReadCloser body return nil error")
	}
}

func TestResponse(t *testing.T) {
	r, err := NewResponse(&http.Response{
		StatusCode: http.StatusOK,
		Status:     http.StatusText(http.StatusOK),
		Header:     http.Header{"X-Test": {"test"}},
		Body:       ioutil.NopCloser(bytes.NewBufferString("body")),
	}, false)
	if err != nil {
		t.Fatalf("NewResponse() error: %s", err)
	}

	if got := r.Headers().Get("X-Test"); got != "test" {
		t.Fatalf("Response.Headers().Get() return %s", got)
	}
	if got := r.StatusCode(); got != http.StatusOK {
		t.Fatalf("Response.StatusCode() return %d", got)
	}
	if got := r.Status(); got != http.StatusText(http.StatusOK) {
		t.Fatalf("Response.Status() return %s", got)
	}
	if got := string(r.Body()); got != "body" {
		t.Fatalf("Response.Body() return %q", got)
	}
	if got := r.Len(); got != len("body") {
		t.Fatalf("Response.Len() return %d", got)
	}
	if got := r.String(); got != "body" {
		t.Fatalf("Response.Len() return %q", got)
	}
}

func TestResponse_JSON(t *testing.T) {
	r, err := NewResponse(&http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`"body"`))}, false)
	if err != nil {
		t.Fatalf("NewResponse() error: %s", err)
	}

	var got string
	if err := r.JSON(&got); err != nil {
		t.Fatalf("Response.JSON() error: %s", err)
	}
	if got != "body" {
		t.Fatalf("Response.JSON() got: %s", got)
	}
}

func TestResponse_XML(t *testing.T) {
	r, err := NewResponse(&http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`<XML>body</XML>`))}, false)
	if err != nil {
		t.Fatalf("NewResponse() error: %s", err)
	}

	var got string
	if err := r.XML(&got); err != nil {
		t.Fatalf("Response.XML() error: %s", err)
	}
	if got != "body" {
		t.Fatalf("Response.XML() got: %s", got)
	}
}

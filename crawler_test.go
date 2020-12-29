package crawler

import (
	"strings"
	"testing"
)

// TestNormalStart calls MyCrawler's start function with normal parameters
func TestNormalStart(t *testing.T) {
	mc := &MyCrawler{
		Depth:   2,
		BaseURL: "thesaurus.com",
	}
	result, err := mc.Start()
	if result == nil || len(result) == 0 {
		t.Fatalf(`Start() : result = %q, %v, want []string, nil`, result, err)
	}
}

// TestStartWithInvalidDepth calls MyCrawler's start function with depth 1
func TestStartWithInvalidDepth(t *testing.T) {
	mc := &MyCrawler{
		Depth:   1,
		BaseURL: "naver.com",
	}
	result, err := mc.Start()
	if result != nil || err == nil {
		t.Fatalf(`Start() : result = %q, err = %v, nil, error`, result, err)
	}
	want := "Depth should be greater than 1"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf(`Start() with host("") err = %v, want %v`, err, want)
	}
}

// TestStartWithEmptyHost calls MyCrawler's start function with empty host
func TestStartWithEmptyHost(t *testing.T) {
	mc := &MyCrawler{
		Depth: 2,
	}
	result, err := mc.Start()
	want := "BaseURL is empty. Please set a base url"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf(`Start() with host("") err = %v, want %v`, err, want)
	}
	if result != nil || err == nil {
		t.Fatalf(`Start() : host("") = %q, %v, want nil, error`, result, err)
	}
}

// TestStartWithInvalidHost calls MyCrawler's start function with invalid host
func TestStartWithInvalidHost(t *testing.T) {
	mc := &MyCrawler{
		Depth:   2,
		BaseURL: "invalid host",
	}
	result, err := mc.Start()
	want := "BaseURL is invalid. Please set a valid base url"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf(`Start() with host("invalid host") err = %v, want %v`, err, want)
	}
	if result != nil || err == nil {
		t.Fatalf(`Start() : host("") = %q, %v, want nil, error`, result, err)
	}
}

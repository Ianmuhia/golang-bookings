package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)
	switch v := h.(type) {
	case http.Handler:
	//do nothing

	default:
		t.Errorf("type is not a http handler %T", v)
	}
}

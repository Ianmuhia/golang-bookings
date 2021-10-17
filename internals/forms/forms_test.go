package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	postedData := url.Values{}

	form := New(postedData)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	has := form.Has("whatever")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("Shows form does not have field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}

	form := New(postedData)

	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows minlength for non-existent field")
	}

	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error and did not get one")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "some value")
	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.Valid() {
		t.Error("shows min length if 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("another_field", "abc123")
	form = New(postedData)

	form.MinLength("another_field", 1)
	if !form.Valid() {
		t.Error("shows min length if 1 is not met when data is shorter")
	}
	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error("should not have an error and did not get one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent email")
	}

	postedData = url.Values{}
	postedData.Add("email", "me@here.com")
	form = New(postedData)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("got an invalid email when we should not have")
	}
	postedData = url.Values{}
	postedData.Add("email", "me@here")
	form = New(postedData)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("got a  valid email when we should not have")
	}

}

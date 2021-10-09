package forms

import (
	"net/http"
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
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}

	postedData.Add("ValidEmail", "john@example.com")
	postedData.Add("InvalidEmail", "joaogmail@asdasd@s.com")

	r.PostForm = postedData
	form := New(r.PostForm)

	if !form.IsEmail("ValidEmail") {
		t.Error("failed on is email test")
	}

	if form.IsEmail("InvalidEmail") {
		t.Error("failed on is email test")
	}

}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}

	postedData.Add("name", "João")
	r.PostForm = postedData
	r.Form = postedData

	form := New(r.PostForm)

	if !form.Has("name", r) {
		t.Error("Error, The field is present but it says it isn't")
	}

	if form.Has("joão", r) {
		t.Error("Error, The field is not present but it says it is")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}

	postedData.Add("Enough", "João")
	postedData.Add("NotEnough", "Mo")
	r.PostForm = postedData
	r.Form = postedData

	form := New(r.PostForm)

	if !form.MinLength("Enough", 3, r) {
		t.Error("Error with 'enough'")
	}

	if form.MinLength("NotEnough", 3, r) {
		t.Error("Error with not 'enough'")
	}
}

func TestErrors(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}

	postedData.Add("NotEnough", "Mo")
	postedData.Add("Enough", "João")
	r.PostForm = postedData
	r.Form = postedData

	form := New(r.PostForm)

	_ = form.MinLength("NotEnough", 3, r)

	out := form.Errors.Get("NotEnough")

	if out == "" {
		t.Error("Error when there is errors in the slice")
	}

	out = form.Errors.Get("Enough")

	if out != "" {
		t.Error("Error when there is no errors in the slice")
	}
}

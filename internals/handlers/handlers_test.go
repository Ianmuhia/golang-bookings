package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/ianmuhia/bookings/internals/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	//{
	//	"home",
	//	"/",
	//	"GET",
	//	[]postData{},
	//	http.StatusOK,
	//}, {
	//	"about",
	//	"/about",
	//	"GET",
	//	[]postData{},
	//	http.StatusOK,
	//}, {
	//	"gq",
	//	"/generals-quarters",
	//	"GET",
	//	[]postData{},
	//	http.StatusOK,
	//}, {
	//	"ms",
	//	"/majors-suite",
	//	"GET",
	//	[]postData{},
	//	http.StatusOK,
	//},
	//{
	//	"sa",
	//	"/search-availability",
	//	"GET",
	//	[]postData{},
	//	http.StatusOK,
	//}, {
	//	"contact",
	//	"/contact",
	//	"GET",
	//	[]postData{},
	//	http.StatusOK,
	//}, {
	//	"mr",
	//	"/make-reservation",
	//	"GET",
	//	[]postData{},
	//	http.StatusOK,
	//}, {
	//	"post-search-availability",
	//	"/search-availability",
	//	"POST",
	//	[]postData{
	//		{key: "start", value: "2020-01-20"},
	//		{key: "end", value: "2020-01-23"},
	//	},
	//	http.StatusOK,
	//}, {
	//	"post-search-availability-json",
	//	"/search-availability-json",
	//	"POST",
	//	[]postData{
	//		{key: "start", value: "2020-01-20"},
	//		{key: "end", value: "2020-01-23"},
	//	},
	//	http.StatusOK,
	//}, {
	//	"make-reservation post",
	//	"/make-reservation",
	//	"POST",
	//	[]postData{
	//		{key: "first_name", value: "ianm"},
	//		{key: "last_name", value: "smith"},
	//		{key: "email", value: "ianm@v.com"},
	//		{key: "phone", value: "112121231"},
	//	},
	//	http.StatusOK,
	//},
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)

	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("reservation handler returned wrong responce code: got %d wanted %d", rr.Code, http.StatusOK)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

func TestHandler(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			value := url.Values{}
			for _, x := range e.params {
				value.Add(x.key, x.value)
			}

			resp, err := ts.Client().PostForm(ts.URL+e.url, value)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}

}

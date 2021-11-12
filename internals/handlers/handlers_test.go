package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ianmuhia/bookings/internals/models"
)

var theTests = []struct {
	name   string
	url    string
	method string
	//params             []postData
	expectedStatusCode int
}{
	{
		"home",
		"/",
		"GET",

		http.StatusOK,
	}, {
		"about",
		"/about",
		"GET",

		http.StatusOK,
	}, {
		"gq",
		"/generals-quarters",
		"GET",

		http.StatusOK,
	}, {
		"ms",
		"/majors-suite",
		"GET",

		http.StatusOK,
	},
	{
		"sa",
		"/search-availability",
		"GET",

		http.StatusOK,
	}, {
		"contact",
		"/contact",
		"GET",

		http.StatusOK,
	},
	//{
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

	//test case where reservation is not in the session (reset everything)
	req, _ = http.NewRequest("GET", "make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong responce code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test with non-existent  room
	req, _ = http.NewRequest("GET", "make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("reservation handler returned wrong responce code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
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
		//if e.method == "GET" {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
		//}
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=ianmuhi@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=122323232")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Postreservation handler returned wrong responce code: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test for missing post body

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Postreservation handler returned wrong responce code for missing post body: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid start date
	reqBody = "start_date=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=ianmuhi@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=122323232")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Postreservation handler returned wrong responce code for missing start date: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid end date
	reqBody = "start_date=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=nil")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=ianmuhi@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=122323232")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Postreservation handler returned wrong responce code for missing end date : got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid end date
	reqBody = "start_date=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-01-01")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=ianmuhi@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=122323232")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=nil")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Postreservation handler returned wrong responce code for missing room id: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for invalid data
	reqBody = "start_date=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=ianmuhi@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=122323232")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Postreservation handler returned wrong responce code for missing room id: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test for failure to insert reservation into database
	reqBody = "start_date=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=ianmuhi@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=122323232")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Postreservation handler failed when trying ot fail insert reservation: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	//test for failure to insert restriction into database
	reqBody = "start_date=2025-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2025-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=ianmuhi@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=122323232")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Postreservation handler failed when trying ot fail insert restriction: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

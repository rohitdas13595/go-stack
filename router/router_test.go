package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterParam(t *testing.T) {
	rt := New()
	rt.Handle(http.MethodGet, "/users/:id", func(w http.ResponseWriter, r *http.Request) {
		p := ParamsFromContext(r.Context())
		if p["id"] != "42" {
			t.Fatalf("id=%q", p["id"])
		}
		w.WriteHeader(200)
	})
	s := httptest.NewServer(rt)
	defer s.Close()
	res, err := http.Get(s.URL + "/users/42")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("status %d", res.StatusCode)
	}
}

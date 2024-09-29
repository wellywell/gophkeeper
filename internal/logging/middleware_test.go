package logging

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"gotest.tools/assert"
)

type TestMiddleWare struct {
	t *testing.T
}

func (m TestMiddleWare) Handle(next http.Handler) http.Handler {

	test := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		v := reflect.ValueOf(w)

		ind := reflect.Indirect(v)

		d := ind.FieldByName("responseData")
		data := reflect.Indirect(d)

		fmt.Println(data)
		size := data.FieldByName("size")

		var expectedSize int64 = 8
		assert.Equal(m.t, expectedSize, size.Int())

		code := data.FieldByName("status")
		var expectedStatus int64 = 200
		assert.Equal(m.t, expectedStatus, code.Int())

		resp := data.FieldByName("data")
		expectedData := "response"
		assert.Equal(m.t, expectedData, string(resp.Bytes()))
	}

	result := http.HandlerFunc(test)

	lo, _ := NewLogger()
	return lo.Handle(result)
}

type Mockhandler struct {
	t *testing.T
}

func (m *Mockhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte("response"))
}

func TestLogger_Handle(t *testing.T) {

	type args struct {
		next http.Handler
	}

	r := httptest.NewRequest(http.MethodGet, "/api/", nil)

	handler := &Mockhandler{t}

	tests := []struct {
		name    string
		args    args
		request *http.Request
	}{
		{"good request", args{handler}, r},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tm := TestMiddleWare{t}
			h := tm.Handle(tt.args.next)

			w := httptest.NewRecorder()
			h.ServeHTTP(w, tt.request)

			assert.Equal(t, 200, w.Code)
			assert.Equal(t, "response", w.Body.String())
		})
	}
}

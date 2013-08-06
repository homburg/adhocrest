package adhocrest

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"bytes"
)

func testMethod(method string, callback func (w *httptest.ResponseRecorder)) {
	req, err := http.NewRequest(method, "http://example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()

	handler := NewHandler("")
	handler.ServeHTTP(w, req)

	callback(w)
}

func TestInvalidMethod(t *testing.T) {
	testMethod("fisk", func (w *httptest.ResponseRecorder) {
		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d, got %d\n", http.StatusMethodNotAllowed, w.Code)
		}
	})
}

func TestPostBadData(t *testing.T) {
	handler := NewHandler("")

	req, _ := http.NewRequest("POST", "/", strings.NewReader("[]"))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if http.StatusInternalServerError != w.Code {
		t.Errorf(
			"Expected internal server error on non-\"dict\" post, got response\n%s",
			w.Body.String(),
		)
	}

	req, _ = http.NewRequest("POST", "/", strings.NewReader("500"))
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if http.StatusInternalServerError != w.Code {
		t.Errorf(
			"Expected internal server error on non-\"dict\" post, got response\n%s",
			w.Body.String(),
		)
	}

	req, _ = http.NewRequest("POST", "/", strings.NewReader(`{"from":"reader"}`))
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if http.StatusCreated != w.Code {
		t.Errorf("Expected %d on valid \"dict\" data, got %d", http.StatusCreated, w.Code)
	}
}

func TestPostNewResource(t *testing.T) {

	handler := NewHandler("")

	// Create a resource
	data := map[string]string{"from":"httptest.Recorder"}
	jsonData, err := json.Marshal(data)
	if nil != err {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", "/fish", strings.NewReader(string(jsonData)))
	if nil != err {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Create another
	req, err = http.NewRequest("POST", "/fish", strings.NewReader(string(jsonData)))
	if nil != err {
		t.Error(err)
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if nil != err {
		t.Error(err)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	headers := w.Header()
	_, ok := headers["Location"]

	var locationHeader string
	if ok {
		locationHeader = headers.Get("location")
		expectedHeaderPrefix := "/fish/"
		if !strings.HasPrefix(locationHeader, expectedHeaderPrefix) || len(locationHeader) < len(expectedHeaderPrefix) {
			t.Errorf("Expected location to match /fish/<string:id>, got %s", locationHeader)
		}

	} else {
		t.Errorf("Expected location header on create, got:\n%s", headers)
	}

	// Get the same info again
	req, err = http.NewRequest("GET", locationHeader, nil)
	if nil != err {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	responseBodyStr := w.Body.String()

	if w.Code != http.StatusOK {
		t.Errorf("Could not get resource from %s, status: %d", locationHeader, w.Code)
	}

	// Unmarshal recieved json data, and remove id
	responseData := make(map[string]string)
	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	if nil != err {
		t.Error(err)
	}
	delete(responseData, "id")
	responseDataStr, err := json.Marshal(responseData)
	if nil != err {
		t.Error(err)
	}

	if 0 != bytes.Compare(jsonData, responseDataStr) {
		t.Errorf(
			`Expected to match posted data plus id:
-----
%s
-----
got
-----
%s
-----`,
			jsonData,
			responseBodyStr,
		)
	}

	// Get a list
	req, err = http.NewRequest("GET", "/fish", nil)
	if nil != err {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 on resource listing, got %d", w.Code)
	}

	// Delete it
	req, err = http.NewRequest("DELETE", locationHeader, nil)
	if nil != err {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected http.StatusOK on DELETE, got %d", w.Code)
	}

	// Get it again (not)
	req, err = http.NewRequest("GET", locationHeader, nil)
	if nil != err {
		t.Error(err)
	}

	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected http status 404, Not Found, on %s, got %d. Was it not deleted?",
			locationHeader, w.Code)
	}
}

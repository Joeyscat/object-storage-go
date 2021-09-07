package test

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	Obj  string
	Name string
}{
	{
		Obj:  "object1",
		Name: "name1",
	},
	{
		Obj:  "object2",
		Name: "name2",
	},
}

func TestPutObject(t *testing.T) {
	for _, testCase := range testCases {

		obj := testCase.Obj
		name := testCase.Name
		hash := hashReader(strings.NewReader(obj))
		t.Logf("hash data: %s", hash)

		// put
		resp, err := put("http://localhost:9001/objects/"+name, "", hash, strings.NewReader(obj))
		if err != nil {
			t.Logf("err: %+v", err.Error())
		}
		assert.Nil(t, err)

		t.Log(resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func TestGetObject(t *testing.T) {
	for _, testCase := range testCases {
		obj := testCase.Obj
		name := testCase.Name

		resp, err := http.Get("http://jojo:9001/objects/" + name)
		if err != nil {
			t.Logf("err: %+v", err.Error())
		}
		assert.Nil(t, err)

		result, err := ioutil.ReadAll(resp.Body)
		assert.Nil(t, err)

		t.Log(resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, obj, string(result))
	}
}

func TestDeleteObject(t *testing.T) {
	for _, testCase := range testCases {
		resp, err := delete("http://jojo:9001/objects/"+testCase.Name, "")
		assert.Nil(t, err)

		t.Log(resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func hashReader(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func put(url, contentType, hash string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Digest", "SHA-256="+hash)
	c := http.DefaultClient
	return c.Do(req)
}

func delete(url, contentType string) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	c := http.DefaultClient
	return c.Do(req)
}

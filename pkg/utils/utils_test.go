package utils

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestCalculateHash(t *testing.T) {
	type args struct {
		r io.Reader
	}
	f, err := os.Open("/home/jojo/env/objects/1/objects/aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=.0.aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=")
	assert.Nil(t, err)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"ok",
			args{f},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateHash(tt.args.r); got != tt.want {
				t.Errorf("CalculateHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAA(t *testing.T) {
	file := "/home/jojo/env/objects/1/objects/aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=.0.aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY="
	h := sha256.New()
	gzipStream(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]

	t.Logf("d: %s, hash: %s", d, hash)
	assert.Equal(t, d, hash)
}

func TestRead(t *testing.T) {
	file := "/home/jojo/env/objects/1/objects/aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=.0.aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY="
	bs, _ := ioutil.ReadFile(file)
	fmt.Println(bs)
}

func gzipStream(w io.Writer, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	defer f.Close()

	gzipStream, err := gzip.NewReader(f)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	defer gzipStream.Close()

	io.Copy(w, gzipStream)
}

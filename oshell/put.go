package oshell

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var oPutCmd = &cobra.Command{
	Use:   "put <file>",
	Short: "Put file to the object storage server",
	Args:  cobra.ExactArgs(1),
	Run:   OPut,
}

var (
	objName string
)

func init() {
	oPutCmd.Flags().StringVarP(&objName, "object-name", "n", "", "name for new object")

	RootCmd.AddCommand(oPutCmd)
}

func OPut(cmd *cobra.Command, params []string) {
	fileToUpload := params[0]

	if fileToUpload == "" {
		log.Fatalf("<filename> is empty\n")
	}
	if objName == "" {
		log.Fatalf("<object-name> is empty\n")
	}

	fp, err := os.Open(fileToUpload)
	if err != nil {
		log.Fatalf("Open file ``%s`: %v\n", fileToUpload, err)
	}
	defer fp.Close()

	hash := hashReader(fp)
	fp.Seek(0, 0)

	resp, err := put("http://localhost:9001/objects/"+objName, "", hash, fp)
	if err != nil {
		log.Fatalf("Put file `%s`: %v\n", fileToUpload, err)
		return
	}

	r, err := ioutil.ReadAll(resp.Body)
	log.Printf("Code: %d\n%s\n", resp.StatusCode, string(r))
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

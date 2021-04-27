package test

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestReadEnv(t *testing.T) {

    e1 := os.Getenv("STORAGE_ROOT")

    t.Logf("get env: %s", e1)

    println(e1)
    fmt.Println(e1)

    assert.NotEmpty(t, e1)
}

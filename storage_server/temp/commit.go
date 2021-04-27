package temp

import (
    "compress/gzip"
    "github.com/joeyscat/object-storage-go/pkg/utils"
    "io"
    "net/url"
    "os"
    "strconv"
    "strings"

    "github.com/joeyscat/object-storage-go/storage_server/locate"
)

func (t *tempInfo) hash() string {
    s := strings.Split(t.Name, ".")
    return s[0]
}

func (t *tempInfo) id() int {
    s := strings.Split(t.Name, ".")
    id, _ := strconv.Atoi(s[1])
    return id
}

func CommitTempObject(datFile string, t *tempInfo) {
    f, err := os.Open(datFile)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    d := url.PathEscape(utils.CalculateHash(f))
    f.Seek(0, io.SeekStart)

    // 压缩到正式对象文件，然后删除临时对象文件
    w, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + t.Name + "." + d)
    if err != nil {
        panic(err)
    }
    defer w.Close()
    w2 := gzip.NewWriter(w)
    defer w2.Close()
    io.Copy(w2, f)

    os.Remove(datFile)

    locate.Add(t.hash(), t.id())
}

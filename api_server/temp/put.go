package temp

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/joeyscat/object-storage-go/api_server/locate"
	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/mongo"
	"github.com/joeyscat/object-storage-go/pkg/rs"
	"github.com/joeyscat/object-storage-go/pkg/utils"
)

func put(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset, err := utils.GetOffsetFromHeader(r.Header)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if current != offset {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	bytes := make([]byte, rs.BlockSize)
	for {
		n, err := io.ReadFull(r.Body, bytes)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Warn(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Warn("resumable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != rs.BlockSize && current != stream.Size {
			return
		}
		stream.Write(bytes[:n])
		if current == stream.Size {
			stream.Flush()
			getStream, err := rs.NewRsResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			if err != nil {
				log.Warn(fmt.Sprintf("NewRsResumableGetStream error: %v", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			hash := url.PathEscape(utils.CalculateHash(getStream))
			if hash != stream.Hash {
				stream.Commit(false)
				log.Warn("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if locate.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			v, err := mongo.SearchLatestVersion(stream.Name)
			if err != nil {
				log.Warn(fmt.Sprintf("SearchLatestVersion error: %v", err))
				w.WriteHeader(http.StatusForbidden)
				return
			}
			err = mongo.AddVersion(stream.Name, stream.Hash, v.Version+1, uint64(stream.Size))
			if err != nil {
				log.Warn(fmt.Sprintf("AddVersion error: %v", err))
				w.WriteHeader(http.StatusForbidden)
				return
			}
			return
		}
	}
}

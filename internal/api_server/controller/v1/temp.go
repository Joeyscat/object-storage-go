package v1

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/joeyscat/object-storage-go/internal/pkg/object"
	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/mongo"
	"github.com/joeyscat/object-storage-go/pkg/rs"
	"github.com/joeyscat/object-storage-go/pkg/utils"
	"github.com/labstack/echo/v4"
)

type TempController struct {
}

func NewTempController() *TempController {
	return &TempController{}
}

func (t *TempController) HeadTempObject(c echo.Context) error {
	token := c.Param("token")
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		log.Warn(err.Error())
		return c.JSON(http.StatusForbidden, nil)
	}
	current := stream.CurrentSize()
	if current == -1 {
		return c.JSON(http.StatusNotFound, nil)
	}
	c.Response().Header().Set("content-length", fmt.Sprintf("%d", current))
	return nil
}

func (t *TempController) PutTempObject(c echo.Context) error {
	token := c.Param("token")
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		log.Warn(err.Error())
		return c.JSON(http.StatusForbidden, nil)
	}
	current := stream.CurrentSize()
	if current == -1 {
		return c.JSON(http.StatusNotFound, nil)
	}
	offset, err := utils.GetOffsetFromHeader(c.Request().Header)
	if err != nil {
		log.Warn(err.Error())
		return c.JSON(http.StatusInternalServerError, nil)
	}
	if current != offset {
		return c.JSON(http.StatusRequestedRangeNotSatisfiable, nil)
	}
	bytes := make([]byte, rs.BlockSize)
	for {
		n, err := io.ReadFull(c.Request().Body, bytes)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Warn(err.Error())
			return c.JSON(http.StatusInternalServerError, nil)
		}
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Warn("resumable put exceed size")
			return c.JSON(http.StatusForbidden, nil)
		}
		if n != rs.BlockSize && current != stream.Size {
			return nil
		}
		stream.Write(bytes[:n])
		if current == stream.Size {
			stream.Flush()
			getStream, err := rs.NewRsResumableGetStream(stream.Servers, stream.Uuids, stream.Size)
			if err != nil {
				log.Warn(fmt.Sprintf("NewRsResumableGetStream error: %v", err))
				return c.JSON(http.StatusInternalServerError, nil)
			}
			hash := url.PathEscape(utils.CalculateHash(getStream))
			if hash != stream.Hash {
				stream.Commit(false)
				log.Warn("resumable put done but hash mismatch")
				return c.JSON(http.StatusForbidden, nil)
			}
			if object.Exist(url.PathEscape(hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			v, err := mongo.SearchLatestVersion(stream.Name)
			if err != nil {
				log.Warn(fmt.Sprintf("SearchLatestVersion error: %v", err))
				return c.JSON(http.StatusForbidden, nil)
			}
			err = mongo.AddVersion(stream.Name, stream.Hash, v.Version+1, uint64(stream.Size))
			if err != nil {
				log.Warn(fmt.Sprintf("AddVersion error: %v", err))
				return c.JSON(http.StatusForbidden, nil)
			}
			return nil
		}
	}
}

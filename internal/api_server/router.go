package api_server

import (
	"net/http"
	"os"

	ctlv1 "github.com/joeyscat/object-storage-go/internal/api_server/controller/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitRouter() {
	e := echo.New()

	installController(e)
	installMiddleware(e)

	e.Logger.Fatal(e.Start(os.Getenv("LISTEN_ADDRESS")))
}
func installController(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	v1 := e.Group("/v1")
	{
		// object
		objectv1 := v1.Group("/objects")
		{
			objectController := ctlv1.NewObjectController(nil)

			objectv1.GET("/:name", objectController.GetObject)
			objectv1.PUT("/:name", objectController.PutObject)
			objectv1.POST("/:name", objectController.CreateObject)
			objectv1.DELETE("/:name", objectController.DeleteObject, authM)

			// TODO not for user
			objectv1.GET("/locate/:name", objectController.GetObjectLocate, authM)

			objectv1.HEAD("/version/:name", objectController.HeadObjectVersion)

			// Upload
			{
				objectv1.PUT("/:objectName/upload", objectController.PutObject)
			}

			// Multipart Upload
			{
				objectv1.PUT("/:objectName/upload/uploadId", objectController.PutObject)
			}
		}

		// temp
		tempv1 := v1.Group("/temp")
		{
			tempController := ctlv1.NewTempController()

			tempv1.HEAD("/:token", tempController.HeadTempObject)
			tempv1.PUT("/:token", tempController.PutTempObject)
		}

		// bucket
		bucketv1 := v1.Group("/buckets", authM)
		{
			bucketController := ctlv1.NewBucketController(nil)
			bucketv1.GET("/", bucketController.GetBucketList)
			bucketv1.POST("/:bucketName", bucketController.CreateBucket)
			bucketv1.DELETE("/:bucketName", bucketController.DeleteBucket)
		}
	}
}

func installMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}

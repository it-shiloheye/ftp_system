package ginserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	ftp_context "github.com/it-shiloheye/ftp_system_lib/context"
	filehandler "github.com/it-shiloheye/ftp_system_lib/file_handler/v2"
)

func RegisterRoutes(r *gin.Engine) error {

	upload_group := r.Group("/upload", func(ctx *gin.Context) {
		log.Println(ctx.Request.URL.RequestURI())
		ctx.Next()
	})

	upload_group.POST("/file/:file_hash", func(ctx *gin.Context) {
		file_hash := ctx.Param("file_hash")

		bf := bytes.NewBuffer(make([]byte, 100_000))

		bf.Reset()
		n, err1 := bf.ReadFrom(ctx.Request.Body)

		d := bf.Bytes()[:n]
		if err1 != nil {
			log.Println(err1)
			ctx.AbortWithStatusJSON(422, gin.H{
				"do": "better",
			})
			return
		}

		fh := filehandler.FileHash{}
		err2 := json.Unmarshal(d, &fh)
		if err2 != nil {
			log.Println(err2)
			ctx.AbortWithStatusJSON(500, gin.H{
				"my": "bad",
			})
			return
		}

		ctx.JSON(201, gin.H{
			"received": file_hash,
		})

		log.Println(string(d))
	})

	upload_group.POST("/stream/:file_hash", func(ctx *gin.Context) {
		hash := ctx.Param("file_hash")
		loc := fmt.Sprintf(`upload_group.POST("/stream/%s", func(ctx *gin.Context)`, hash)

		data, err1 := io.ReadAll(ctx.Request.Body)
		if err1 != nil {
			log.Println(
				&ftp_context.LogItem{
					Location:  loc,
					After:     `data, err1:= io.ReadAll(ctx.Request.Body)`,
					Message:   err1.Error(),
					CallStack: []error{err1},
				},
			)
		}

		log.Println("uploaded file: ", hash, "\ndata:\n", string(data))

		ctx.Status(404)
	})

	return nil
}

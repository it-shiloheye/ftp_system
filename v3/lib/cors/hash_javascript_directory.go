package cors

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	b "github.com/it-shiloheye/ftp_system/v3/lib/base"
)

type HashDirFs struct {
	path         string
	route        string
	algo         HashAlgorithm
	hashedFiles  map[string]*HashedFile
	replacements map[string]string
}

func NewHashDirFs(base_path string, route string, algo HashAlgorithm, replacements map[string]string) (hD *HashDirFs, err error) {
	hD = &HashDirFs{
		path:         base_path,
		algo:         algo,
		replacements: replacements,
		route:        route,
		hashedFiles:  make(map[string]*HashedFile, 4),
	}

	err = filepath.WalkDir(base_path, func(file_path string, f fs.DirEntry, err_ error) error {
		if err_ != nil {
			log.Printf("error loading file: %s\nerror: %s\n", file_path, b.PrettyPrintValue(err_))
			return err_
		}

		if !rightFile(f) {
			return nil
		}

		hF, err_ := NewHashedFile(file_path, algo)
		if err_ != nil {
			log.Printf("error hashing file: %s\nerror: %s\n", file_path, b.PrettyPrintValue(err_))
			return err_
		}

		c_ := getPath(base_path, file_path, route)

		hD.hashedFiles[c_] = hF
		return nil
	})
	if err != nil {
		log.Printf("error loading files: %s\n", b.PrettyPrintValue(err))
		hD = nil
		return
	}

	return
}

func rightFile(f fs.DirEntry) bool {
	if f.IsDir() {
		return false
	}

	a_ := strings.Split(f.Name(), "/")
	l_ := a_[len(a_)-1]
	s_ := strings.Split(l_, ".")
	ext := s_[len(s_)-1]

	return ext == "js"
}

func (hD *HashDirFs) SetRoute(r string) {
	hD.route = r
}

func (hD *HashDirFs) Route() string {
	return hD.route
}

func (hD *HashDirFs) NewServeDirFs(r *gin.Engine, cors_ *CorsObject, cache_seconds time.Duration, reload bool) (err error) {

	for route, hF := range hD.hashedFiles {
		route := route
		hF := hF
		r.GET(route, func(gtx *gin.Context) {

			if cors_ != nil {
				cors_.RenderCors(gtx)
			}

			if reload {
				hF, err = hF.Reload()
				if err != nil {
					gtx.Status(http.StatusInternalServerError)
					_, err = fmt.Fprint(gtx.Writer, "500 internal server error")
					if err != nil {
						return
					}
				}

				defer func() {
					hD.hashedFiles[gtx.Request.RequestURI] = hF
				}()
			}

			gtx.Header("Content-Type", "text/javascript")
			gtx.Header("Cache-Control", fmt.Sprintf("max-age=%d", int(cache_seconds)))
			fmt.Fprint(gtx.Writer, hF.Replace(hD.replacements).Str())
		})

	}

	return nil
}

func getPath(base_path, file_path, route string) string {
	a := strings.Split(base_path, "/")
	b := strings.Split(file_path, "\\")

	if a[0] == "." {
		a = a[1:]
	}

	s_ := strings.Join(b[len(a):], "/")
	if len(s_) < 1 {
		return route
	}
	return fmt.Sprintf("%s/%s", route, s_)
}

func (hD *HashDirFs) GetNonce() (m map[string]string) {
	m = make(map[string]string, len(hD.hashedFiles))
	for file, hashF := range hD.hashedFiles {
		file := file
		hashF := hashF

		m[file] = hashF.nonce
	}

	return
}

func (hD *HashDirFs) GetSingleNonce(name string) string {

	return hD.hashedFiles[name].nonce
}

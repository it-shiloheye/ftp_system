package cors

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type EmbedderPolicy int

const (
	EmbedderUnsafeNone EmbedderPolicy = iota + 1
	EmbedderRequireCorp
	EmbedderCredentialLess
)

func (cep EmbedderPolicy) Policy() string {

	switch cep {
	case EmbedderUnsafeNone:
		return "unsafe-none"
	case EmbedderRequireCorp:
		return "require-corp"
	case EmbedderCredentialLess:
		return "credentialless"
	}
	panic(fmt.Errorf("invalid embedder-policy"))
}

type OpenerPolicy int

const (
	OpenerUnsafeNone OpenerPolicy = iota + 1
	OpenerSameOriginAllowPopups
	OpenerSameOrigin
)

func (cep OpenerPolicy) Policy() string {

	switch cep {
	case OpenerUnsafeNone:
		return "unsafe-none"
	case OpenerSameOriginAllowPopups:
		return "same-origin-allow-popups"
	case OpenerSameOrigin:
		return "same-origin"
	}
	panic(fmt.Errorf("invalid embedder-policy"))
}

type ResourcePolicy int

const (
	ResourceSameSite ResourcePolicy = iota + 1
	ResourceSameOrigin
	ResourceCrossOrigin
)

func (rp ResourcePolicy) Policy() string {
	switch rp {
	case ResourceSameSite:
		return "same-site"
	case ResourceSameOrigin:
		return "same-origin"
	case ResourceCrossOrigin:
		return "cross-origin"

	}

	panic(fmt.Errorf("invalid resource policy"))
}

type CorsObject struct {
	FrameSrc       []string       `json:"frame-src"`
	ChildSrc       []string       `json:"child-src"`
	ConnectSrc     []string       `json:"connect-src"`
	ScriptSrc      []string       `json:"script-src"`
	OriginEmbedder EmbedderPolicy `json:"embedder-policy"`
	OriginOpener   OpenerPolicy   `json:"opener-policy"`
	OriginResource ResourcePolicy `json:"resource-policy"`
}

func (cors *CorsObject) RenderCors(gtx *gin.Context) {
	join := strings.Join
	gtx.Header("Content-Security-Policy", join([]string{
		"frame-src " + join(cors.FrameSrc, " "),
		"child-src " + join(cors.ChildSrc, " "),
		"connect-src " + join(cors.ConnectSrc, " "),
		"script-src " + join(cors.ScriptSrc, " "),
	}, "; "))

	gtx.Header("Cross-Origin-Embedder-Policy", cors.OriginEmbedder.Policy())
	gtx.Header("Cross-Origin-Opener-Policy", cors.OriginOpener.Policy())
	gtx.Header("Cross-Origin-Resource-Policy", cors.OriginResource.Policy())
}

func (cors *CorsObject) Frame(framesrc ...string) *CorsObject {

	cors.FrameSrc = append(cors.FrameSrc, framesrc...)
	return cors
}

func (cors *CorsObject) Child(childsrc ...string) *CorsObject {

	cors.ChildSrc = append(cors.ChildSrc, childsrc...)
	return cors
}

func (cors *CorsObject) Connect(connectsrc ...string) *CorsObject {

	cors.ConnectSrc = append(cors.ConnectSrc, connectsrc...)
	return cors
}

func (cors *CorsObject) Script(scriptsrc ...string) *CorsObject {

	cors.ScriptSrc = append(cors.ScriptSrc, scriptsrc...)
	return cors
}

func (cors *CorsObject) Embedder(e EmbedderPolicy) *CorsObject {

	cors.OriginEmbedder = e
	return cors
}

func (cors *CorsObject) Opener(o OpenerPolicy) *CorsObject {

	cors.OriginOpener = o
	return cors
}

func (cors *CorsObject) Resource(r ResourcePolicy) *CorsObject {

	cors.OriginResource = r
	return cors
}

func NewCors() *CorsObject {

	return &CorsObject{}
}

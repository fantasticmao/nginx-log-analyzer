package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
)

type PvAndUvHandler struct {
	pv      int32
	uv      int32
	uniqMap map[string]bool
}

func NewPvAndUvHandler() *PvAndUvHandler {
	return &PvAndUvHandler{
		pv:      0,
		uv:      0,
		uniqMap: make(map[string]bool),
	}
}

func (handler *PvAndUvHandler) Input(info *ioutil.LogInfo) {
	handler.pv++
	if _, ok := handler.uniqMap[info.RemoteAddr]; !ok {
		handler.uv++
		handler.uniqMap[info.RemoteAddr] = true
	}
}

func (handler *PvAndUvHandler) Output(limit int) {
	fmt.Printf("PV: %v\nUV: %v\n", handler.pv, handler.uv)
}

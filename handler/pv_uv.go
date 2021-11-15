package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-log-analyzer/parser"
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

func (handler *PvAndUvHandler) Input(info *parser.LogInfo) {
	handler.pv++
	if _, ok := handler.uniqMap[info.RemoteAddr]; !ok {
		handler.uv++
		handler.uniqMap[info.RemoteAddr] = true
	}
}

func (handler *PvAndUvHandler) Output(limit int) {
	fmt.Printf("PV: %v\n", handler.pv)
	fmt.Printf("UV: %v\n", handler.uv)
}

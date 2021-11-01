package main

import "fmt"

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

func (handler *PvAndUvHandler) input(info *LogInfo) {
	handler.pv++
	if _, ok := handler.uniqMap[info.RemoteAddr]; !ok {
		handler.uv++
		handler.uniqMap[info.RemoteAddr] = true
	}
}

func (handler *PvAndUvHandler) output(limit int) {
	fmt.Printf("PV: %v\nUV: %v\n", handler.pv, handler.uv)
}

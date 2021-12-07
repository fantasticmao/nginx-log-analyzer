package handler

import (
	"github.com/fantasticmao/nginx-log-analyzer/parser"
	"github.com/pterm/pterm"
	"strconv"
)

type PvAndUvHandler struct {
	pv      int
	uv      int
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
	data := pterm.TableData{
		{"PV", strconv.Itoa(handler.pv)},
		{"UV", strconv.Itoa(handler.uv)},
	}
	_ = pterm.DefaultTable.WithData(data).Render()
}

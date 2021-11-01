package main

import (
	"errors"
	"fmt"
	"sort"
)

type MostMatchFieldHandler struct {
	fieldType int
	countMap  map[string]int
}

func NewMostMatchFieldHandler(fieldType int) *MostMatchFieldHandler {
	return &MostMatchFieldHandler{
		fieldType: fieldType,
		countMap:  make(map[string]int),
	}
}

func (handler *MostMatchFieldHandler) input(info *LogInfo) {
	var field string
	switch handler.fieldType {
	case FieldIp:
		field = info.RemoteAddr
	case FieldUri:
		field = info.Request
	case FieldUserAgent:
		field = info.HttpUserAgent
	default:
		panic(errors.New("unknown file type"))
	}

	if _, ok := handler.countMap[field]; ok {
		handler.countMap[field]++
	} else {
		handler.countMap[field] = 1
	}
}

func (handler *MostMatchFieldHandler) output(limit int) {
	keys := make([]string, 0, len(handler.countMap))
	for k := range handler.countMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return handler.countMap[keys[i]] > handler.countMap[keys[j]]
	})

	for i := 0; i < limit && i < len(keys); i++ {
		fmt.Printf("\"%v\" hits: %v\n", keys[i], handler.countMap[keys[i]])
	}
}

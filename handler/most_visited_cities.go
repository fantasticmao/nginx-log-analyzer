package handler

import (
	"fmt"
	"github.com/fantasticmao/nginx-json-log-analyzer/ioutil"
	"github.com/oschwald/geoip2-golang"
	"net"
	"path"
	"sort"
	"strings"
)

const dbFile string = "City.mmdb"
const (
	LanguageEn   = "en"
	LanguageZhCn = "zh-CN"
)

type MostVisitedCities struct {
	geoLite2Db *geoip2.Reader
	// country -> city -> ip -> count
	countryCityIpCountMap map[string]map[string]map[string]int
	// country -> count
	countryCountMap map[string]int
	// city -> count
	cityCountMap map[string]int
}

func NewMostVisitedCities(configDir string) *MostVisitedCities {
	db, err := geoip2.Open(path.Join(configDir, dbFile))
	if err != nil {
		ioutil.Fatal("open %v error: %v\n", dbFile, err.Error())
	}
	return &MostVisitedCities{
		geoLite2Db:            db,
		countryCityIpCountMap: make(map[string]map[string]map[string]int),
		countryCountMap:       make(map[string]int),
		cityCountMap:          make(map[string]int),
	}
}

func (handler *MostVisitedCities) Input(info *ioutil.LogInfo) {
	ip := net.ParseIP(info.RemoteAddr)
	record, err := handler.geoLite2Db.City(ip)
	if err != nil {
		ioutil.Fatal("query from %v error: %v\n", dbFile, err.Error())
	}

	var (
		country string
		city    string
	)
	country = record.Country.Names[LanguageEn]
	if strings.EqualFold("China", country) {
		country = record.Country.Names[LanguageZhCn]
		city = record.City.Names[LanguageZhCn]
	} else {
		city = record.City.Names[LanguageEn]
	}
	if city == "" {
		city = "unknown"
	}

	if _, ok := handler.countryCityIpCountMap[country]; !ok {
		cityIpCountMap := make(map[string]map[string]int)
		handler.countryCityIpCountMap[country] = cityIpCountMap
		handler.countryCountMap[country] = 1
		handler.cityCountMap[city] = 0
	} else {
		handler.countryCountMap[country]++
	}

	if _, ok := handler.countryCityIpCountMap[country][city]; !ok {
		ipCountMap := make(map[string]int)
		handler.countryCityIpCountMap[country][city] = ipCountMap
		handler.cityCountMap[city] = 1
	} else {
		handler.cityCountMap[city]++
	}

	if _, ok := handler.countryCityIpCountMap[country][city][info.RemoteAddr]; !ok {
		handler.countryCityIpCountMap[country][city][info.RemoteAddr] = 1
	} else {
		handler.countryCityIpCountMap[country][city][info.RemoteAddr]++
	}
}

func (handler *MostVisitedCities) Output(limit int) {
	defer handler.geoLite2Db.Close()

	countryCountKeys := make([]string, len(handler.countryCountMap))
	for k := range handler.countryCityIpCountMap {
		countryCountKeys = append(countryCountKeys, k)
	}
	sort.Slice(countryCountKeys, func(i, j int) bool {
		return handler.countryCountMap[countryCountKeys[i]] > handler.countryCountMap[countryCountKeys[j]]
	})

	for i := 0; i < len(countryCountKeys); i++ {
		country := countryCountKeys[i]
		cityIpCountMap := handler.countryCityIpCountMap[country]
		fmt.Printf("[%v] hits: %v\n", country, handler.countryCountMap[country])

		cityCountKeys := make([]string, len(handler.cityCountMap))
		for k := range cityIpCountMap {
			cityCountKeys = append(cityCountKeys, k)
		}
		sort.Slice(cityCountKeys, func(i, j int) bool {
			return handler.cityCountMap[cityCountKeys[i]] > handler.cityCountMap[cityCountKeys[j]]
		})

		for j := 0; j < len(cityCountKeys); j++ {
			city := cityCountKeys[j]
			ipCountMap := cityIpCountMap[city]
			fmt.Printf("  |--[%v] hits: %v\n", city, handler.cityCountMap[city])

			ipCountKeys := make([]string, len(ipCountMap))
			for k := range ipCountMap {
				ipCountKeys = append(ipCountKeys, k)
			}

			sort.Slice(ipCountKeys, func(i, j int) bool {
				return ipCountMap[ipCountKeys[i]] > ipCountMap[ipCountKeys[j]]
			})

			for k := 0; k < limit && k < len(ipCountKeys); k++ {
				fmt.Printf("  |  |--\"%v\" hits: %v\n", ipCountKeys[k], ipCountMap[ipCountKeys[k]])
			}
		}
	}
}

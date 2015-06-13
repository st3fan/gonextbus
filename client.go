// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package main

import (
	"encoding/json"
	"launchpad.net/xmlpath"
	"log"
	"net/http"
	"os"
	"strconv"
)

const NextBusPublicXMLFeed = "http://webservices.nextbus.com/service/publicXMLFeed"

type RouteListRoute struct {
	Tag   string
	Title string
}

type RouteList struct {
	Copyright string
	Routes    []RouteListRoute
}

// Experiment with xmlpath. Not sure if I like this just yet.

var (
	routePath             = xmlpath.MustCompile("/body/route")
	tagPath               = xmlpath.MustCompile("@tag")
	titlePath             = xmlpath.MustCompile("@title")
	directionTagPath      = xmlpath.MustCompile("@tag")
	directionTitlePath    = xmlpath.MustCompile("@title")
	directionNamePath     = xmlpath.MustCompile("@name")
	directionUseForUIPath = xmlpath.MustCompile("@useForUI")
	directionBranchPath   = xmlpath.MustCompile("@branch")
	routeTagPath          = xmlpath.MustCompile("/body/route/@tag")
	routeTitlePath        = xmlpath.MustCompile("/body/route/@title")
	stopTagPath           = xmlpath.MustCompile("@tag")
	stopTitlePath         = xmlpath.MustCompile("@title")
	stopLatPath           = xmlpath.MustCompile("@lat")
	stopLonPath           = xmlpath.MustCompile("@lon")
	stopStopIdPath        = xmlpath.MustCompile("@stopId")
	routeStopPath         = xmlpath.MustCompile("/body/route/stop")
	routeDirectionPath    = xmlpath.MustCompile("/body/route/direction")
)

//
// <body copyright="All data copyright Toronto Transit Commission 2015.">
//   <route tag="501" title="501-Queen"/>
// </body>
//

func FetchRouteList(agency string) (RouteList, error) {
	resp, err := http.Get(NextBusPublicXMLFeed + "?command=routeList&a=" + agency)
	if err != nil {
		return RouteList{}, err
	}

	defer resp.Body.Close()

	root, err := xmlpath.Parse(resp.Body)
	if err != nil {
		return RouteList{}, err
	}

	var routeList RouteList

	iter := routePath.Iter(root)
	for iter.Next() {
		tag, _ := tagPath.String(iter.Node())
		title, _ := titlePath.String(iter.Node())
		routeList.Routes = append(routeList.Routes, RouteListRoute{Tag: tag, Title: title})
	}

	return routeList, nil
}

type RouteConfigStop struct {
	Tag    string
	Title  string
	Lat    float64
	Lon    float64
	StopId string
}

type RouteConfigDirection struct {
	Tag      string
	Title    string
	Name     string
	UseForUI bool
	Branch   string
	Stops    []RouteConfigStop
}

type RouteConfig struct {
	Tag           string
	Title         string
	Color         string
	OppositeColor string
	LatMin        float64
	LatMax        float64
	LonMin        float64
	LonMax        float64
	Directions    []RouteConfigDirection
}

func ParseRouteConfigDirection(node *xmlpath.Node, stopsByTag map[string]RouteConfigStop) (RouteConfigDirection, error) {
	tag, _ := directionTagPath.String(node)
	title, _ := directionTitlePath.String(node)
	name, _ := directionNamePath.String(node)
	useForUI, _ := directionUseForUIPath.String(node)
	branch, _ := directionBranchPath.String(node)

	direction := RouteConfigDirection{
		Tag:      tag,
		Title:    title,
		Name:     name,
		UseForUI: parseBool(useForUI),
		Branch:   branch,
	}

	directionStopPath := xmlpath.MustCompile("stop")
	directionStopTagPath := xmlpath.MustCompile("@tag")

	directionStopIter := directionStopPath.Iter(node)
	for directionStopIter.Next() {
		tag, _ := directionStopTagPath.String(directionStopIter.Node())
		if stop, ok := stopsByTag[tag]; ok {
			direction.Stops = append(direction.Stops, stop)
		}
	}

	return direction, nil
}

func FetchRouteConfig(agency, route string) (RouteConfig, error) {
	resp, err := http.Get(NextBusPublicXMLFeed + "?command=routeConfig&a=" + agency + "&r=" + route)
	if err != nil {
		return RouteConfig{}, err
	}

	defer resp.Body.Close()

	root, err := xmlpath.Parse(resp.Body)
	if err != nil {
		return RouteConfig{}, err
	}

	var routeConfig RouteConfig

	routeTag, _ := routeTagPath.String(root)
	routeConfig.Tag = routeTag

	routeTitle, _ := routeTitlePath.String(root)
	routeConfig.Title = routeTitle

	// Parse the stops

	var stopsByTag map[string]RouteConfigStop = map[string]RouteConfigStop{}

	iter := routeStopPath.Iter(root)
	for iter.Next() {
		tag, _ := stopTagPath.String(iter.Node())
		title, _ := stopTitlePath.String(iter.Node())
		lat, _ := stopLatPath.String(iter.Node())
		lon, _ := stopLonPath.String(iter.Node())
		stopId, _ := stopStopIdPath.String(iter.Node())

		stopsByTag[tag] = RouteConfigStop{
			Tag:    tag,
			Title:  title,
			Lat:    parseFloat(lat),
			Lon:    parseFloat(lon),
			StopId: stopId,
		}
	}

	// Parse the directions

	routeDirectionIter := routeDirectionPath.Iter(root)
	for routeDirectionIter.Next() {
		direction, err := parseRouteConfigDirection(routeDirectionIter.Node(), stopsByTag)
		if err != nil {
			return RouteConfig{}, err
		}
		routeConfig.Directions = append(routeConfig.Directions, direction)
	}

	return routeConfig, nil
}

package main

import (
	cov "Corona/entity/cov"
	"Corona/entity/covAgeGender"
	"Corona/entity/covAllVaccine"
	"Corona/entity/covHospital"
	"Corona/entity/covState"
	"github.com/zserge/lorca"
	"io/ioutil"
	"net/url"
	"strings"
)

const (
	TEMPLATE = "./templates/"
	CSS = "./static/css/"
	JS = "./static/js/"
)

var (
	indexTemplate, _ = ioutil.ReadFile(TEMPLATE + "index.html")
	mapTemplate, _ = ioutil.ReadFile(TEMPLATE + "map.html")
	ageGenderTemplate, _ = ioutil.ReadFile(TEMPLATE + "age_gender.html")
	indexCss, _ = ioutil.ReadFile(CSS + "index.css")
	indexJs, _ = ioutil.ReadFile(JS + "index.js")
	mapJs, _ = ioutil.ReadFile(JS + "map.js")
	ageGenderJs, _ = ioutil.ReadFile(JS + "ageGender.js")
	indexHtml = getHtml(&indexTemplate, &indexCss, &indexJs)
	mapHtml = getHtml(&mapTemplate, &indexCss, &mapJs)
	ageGenderHtml = getHtml(&ageGenderTemplate, &indexCss, &ageGenderJs)
)

func getHtml(template, css, js *[]byte) *string {
	html := strings.Replace(string(*template), "</title>", "</title><style>" + string(*css) + "</style>", 1)
	html = strings.Replace(html, "</body>", "<script>" + string(*js) + "</script></body>", 1)
	return &html
}

var ui lorca.UI

func main() {
	cs := covState.NewCov19State()
	c := cov.NewCov19()
	av := covAllVaccine.NewAllVacc()
	h := covHospital.NewCovHospital()
	cag := covAgeGender.NewCov19State()

	ui, _ = lorca.New("" , "", lorca.PageA4Width, lorca.PageA4Height)
	defer ui.Close()

	ui.Bind("accDefRate", c.AccDefRate)
	ui.Bind("initTitles", c.Titles)
	ui.Bind("stateItems", cs.GetItems)
	ui.Bind("movePage", movePage)
	ui.Bind("avItems", av.GetItems)
	ui.Bind("hItems", h.GetItems)
	ui.Bind("cagItems", cag.GetItems)


	ui.Load("data:text/html," + url.PathEscape(*indexHtml))

	<-ui.Done()
}

func movePage(page string) {
	if page == "index" {
		ui.Load("data:text/html," + url.PathEscape(*indexHtml))
	} else if page == "map" {
		ui.Load("data:text/html," + url.PathEscape(*mapHtml))
	} else if page == "ageGender" {
		ui.Load("data:text/html," + url.PathEscape(*ageGenderHtml))
	}
}

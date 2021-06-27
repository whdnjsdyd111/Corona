package covAllVaccine

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type (
	allVaccine struct {
		Response xml.Name `xml:"response"`
		Body body `xml:"body"`
	}

	body struct {
		XMLName xml.Name `xml:"body"`
		DataTime string `xml:"dataTime"`
		Items items `xml:"items"`
	}

	items struct {
		XMLName xml.Name `xml:"items"`
		Item []item `xml:"item"`
	}

	item struct {
		XMLName xml.Name `xml:"item"`
		Tpcd string `xml:"tpcd"`
		FirstCnt int64 `xml:"firstCnt"`
		SecondCnt int64 `xml:"secondCnt"`
	}
)


const (
	openApi = "https://nip.kdca.go.kr/irgd/cov19stats.do?list=all"
)

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func NewAllVacc() *allVaccine {
	var result allVaccine
	if xmlBytes, err := getXML(openApi); err != nil {
		log.Printf("Failed to get XML : %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &result)
	}
	return &result
}

func (av *allVaccine) GetItems() []item {
	return av.Body.Items.Item
}
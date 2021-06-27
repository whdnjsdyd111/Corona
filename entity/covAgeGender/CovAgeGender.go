package covAgeGender

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	covid19AgeGedner struct {
		Response xml.Name `xml:"response"`
		Header   header   `xml:"header"` // 헤더
		Body     body     `xml:"body"`   // 바디
	}

	header struct {
		XMLName    xml.Name `xml:"header"`
		ResultCode int      `xml:"resultCode"` // 결과 코드
		ResultMsg  string   `xml:"resultMsg"`  // 결과 메시지
	}

	body struct {
		XMLName    xml.Name `xml:"body"`
		Items      items    `xml:"items"`      // 결과 아이템들
		NumOfRows  int      `xml:"numOfRows"`  // 로우 숫자
		PageNo     int      `xml:"pageNo"`     // 페이지 숫자
		TotalCount int      `xml:"totalCount"` // 토탈 숫자
	}

	items struct {
		XMLName xml.Name `xml:"items"`
		Item    []item   `xml:"item"`
	}

	item struct {
		XMLName      xml.Name `xml:"item"`         // 아이템
		ConfCase     int      `xml:"confCase"`     // 확진자
		ConfCaseRate string    `xml:"confCaseRate"` // 확진률
		CreateDt     string   `xml:"createDt"`     // 등록일시
		CriticalRate string  `xml:"criticalRate"` // 치명률
		Death        int      `xml:"death"`        // 사망자
		DeathRate    string  `xml:"deathRate"`    // 사망률
		Gubun        string   `xml:"gubun"`        // 구분 - 성별 연령
		Seq          int      `xml:"seq"`          // 게시글 번호
		UpdateDt     string   `xml:"updateDt"`     //수정일시
	}
)

const (
	openApi    = "http://openapi.data.go.kr/openapi/service/rest/Covid19/getCovid19GenAgeCaseInfJson?"
	serviceKey = "AhuzOWhuvKcU%2BuVgEZCPWra%2BqKbvU8XR5NoHgXnXOhRkZwGfyZF9mHpnU3L%2BrNqbcWhYAAgQbgYiiR2RBgWiyQ%3D%3D"
	pageNo     = "1"
	numOfRows  = "11"
	url        = openApi + "serviceKey=" + serviceKey + "&pageNo=" + pageNo + "&numOfRows=" + numOfRows
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

func NewCov19State() *covid19AgeGedner{
	now := time.Now()
	if (now.Hour() > 0 || now.Hour() <= 9) || (now.Hour() == 9 && now.Minute() < 30) {
		now = now.AddDate(0, 0, -1)
	}
	today := now.Format("20060102")
	addUrl := url + "&startCreateDt=" + today + "&endCreateDt=" + today

	var result covid19AgeGedner
	if xmlBytes, err := getXML(addUrl); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &result)
	}
	return &result
}

func (cag *covid19AgeGedner) GetItems() []item {
	return cag.Body.Items.Item
}
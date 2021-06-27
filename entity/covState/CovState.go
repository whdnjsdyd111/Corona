package covState

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	covid19State struct {
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
		CreateDt     string   `xml:"createDt"`     // 등록일시 분초
		DeathCnt     int      `xml:"deathCnt"`     // 사망자 수
		DefCnt		 int	  `xml:"defCnt"`	   // 확진자 수
		Gubun        string   `xml:"gubun"`        // 시도명
		GubunCn      string   `xml:"gubunCn"`      // 시도명
		GubunEn      string   `xml:"gubunEn"`      // 시도명
		IncDec       int      `xml:"incDec"`       // 전일대비 증감 수
		IsolClearCnt int      `xml:"isolClearCnt"` // 격리 해제 수
		IsolIngCnt	 int	  `xml:"isolIngCnt"`   // 격리 중 환자
		LocalOccCnt  int	  `xml:"localOccCnt"`  // 지역 발생 수
		OverFlowCnt	 int	  `xml:"overFlowCnt"`  // 해외 유입 수
		QurRate      string  `xml:"qurRate"`      // 10만명당 발생률
		Seq          int      `xml:"seq"`          // 국내 시도별 발생 현황 고유값
		StdDay       string   `xml:"stdDay"`       // 기준 일시
		UpdateDt     string   `xml:"updateDt"`     // 수정일시분초
	}
)

const (
	openApi    = "http://openapi.data.go.kr/openapi/service/rest/Covid19/getCovid19SidoInfStateJson?"
	serviceKey = "AhuzOWhuvKcU%2BuVgEZCPWra%2BqKbvU8XR5NoHgXnXOhRkZwGfyZF9mHpnU3L%2BrNqbcWhYAAgQbgYiiR2RBgWiyQ%3D%3D"
	pageNo     = "1"
	numOfRows  = "38"
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

func NewCov19State() *covid19State {
	now := time.Now()
	if (now.Hour() > 0 || now.Hour() <= 9) || (now.Hour() == 9 && now.Minute() < 30) {
		now = now.AddDate(0, 0, -1)
	}
	today := now.Format("20060102")
	yesterday := now.AddDate(0, 0, -1).Format("20060102")
	addUrl := url + "&startCreateDt=" + yesterday + "&endCreateDt=" + today

	var result covid19State
	if xmlBytes, err := getXML(addUrl); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &result)
	}

	return &result
}

func (c *covid19State) GetMap() map[string][]*item {
	m := make(map[string][]*item)
	m["경기도"] = *new([]*item)
	m["강원도"] = *new([]*item)
	m["충청도"] = *new([]*item)
	m["전라도"] = *new([]*item)
	m["경상도"] = *new([]*item)
	m["제주도"] = *new([]*item)
	m["검역"] = *new([]*item)

	items := c.Body.Items.Item

	for _, v := range items {
		switch v.Gubun {
		case "서울", "경기", "인천":
			m["경기도"] = append(m["경기도"], &v)
		case "강원":
			m["강원도"] = append(m["강원도"], &v)
		case "충남", "충북", "세종", "대전":
			m["충청도"] = append(m["충청도"], &v)
		case "전북", "전남", "광주":
			m["전라도"] = append(m["전라도"], &v)
		case "경북", "대구", "경남", "울산", "부산":
			m["경상도"] = append(m["경상도"], &v)
		case "제주":
			m["제주도"] = append(m["제주도"], &v)
		case "검역":
			m["검역"] = append(m["검역"], &v)
		}
	}
	return m
}

func (c *covid19State) GetItems() []item {
	return c.Body.Items.Item
}
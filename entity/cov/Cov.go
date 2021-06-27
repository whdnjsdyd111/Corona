package cov

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"
)

type (
	covid19 struct {
		Response xml.Name `xml:"response"`
		Header   header   `xml:"header"` // 헤더
		Body     body     `xml:"body"`   // 바디
	}

	header struct {
		XMLName    xml.Name `xml:"header"`
		ResultCode int      `xml:"resultCode"`	// 결과 코드
		ResultMsg  string   `xml:"resultMsg"`	// 결과 메시지
	}

	body struct {
		XMLName    xml.Name `xml:"body"`
		Items      items    `xml:"items"`      // 결과 아이템들
		NumOfRows  int      `xml:"NumOfRows"`  // 로우 숫자
		PageNo     int      `xml:"pageNo"`     // 페이지 숫자
		TotalCount int      `xml:"totalCount"` // 토탈 숫자
	}

	items struct {
		XMLName xml.Name `xml:"items"`
		Item	[]item   `xml:"item"`
	}

	item struct {
		XMLName        xml.Name `xml:"item"`			// 아이템
		AccDefRate     float64  `xml:"accDefRate"`		// 누적 확진률
		AccExamCnt     int      `xml:"accExamCnt"`		// 누적 검사수
		AccExamCompCnt int      `xml:"accExamCompCnt"`	// 누적 검사 완료 수
		CareCnt        int      `xml:"careCnt"`			// 치료 중 환자 수
		ClearCnt       int      `xml:"clearCnt"`		// 격리 해제 수
		CreateDt	   string	`xml:"createDt"`		// 등록일시 분초
		DeathCnt       int      `xml:"deathCnt"`		// 사망자 수
		DecideCnt      int      `xml:"decideCnt"`		// 확진자 수
		ExamCnt        int      `xml:"examCnt"`			// 검사 진행 수
		ResultNegCnt   int      `xml:"resultNegCnt"`	// 결과 음성 수
		Seq            int      `xml:"seq"`				// 감염현황 고유값
		StateDt        int      `xml:"stateDt"`			// 기준일
		StateTime      string   `xml:"stateTime"`		// 기준시간
		UpdateDt       string   `xml:"updateDt"`		// 수정일시 분초
	}
)

const (
	openApi = "http://openapi.data.go.kr/openapi/service/rest/Covid19/getCovid19InfStateJson?"
	serviceKey = "AhuzOWhuvKcU%2BuVgEZCPWra%2BqKbvU8XR5NoHgXnXOhRkZwGfyZF9mHpnU3L%2BrNqbcWhYAAgQbgYiiR2RBgWiyQ%3D%3D"
	pageNo  = "1"
	numOfRows = "10"
	url = openApi + "serviceKey=" + serviceKey + "&pageNo=" + pageNo + "&numOfRows=" + numOfRows
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

func NewCov19() *covid19 {
	now := time.Now()
	if (now.Hour() > 0 || now.Hour() <= 9) || (now.Hour() == 9 && now.Minute() < 30) {
		now = now.AddDate(0, 0, -1)
	}
	today := now.Format("20060102")
	yesterday := now.AddDate(0, 0, -1).Format("20060102")
	addUrl := url + "&startCreateDt=" + yesterday + "&endCreateDt=" + today

	var result covid19
	if xmlBytes, err := getXML(addUrl); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		xml.Unmarshal(xmlBytes, &result)
	}
	return &result
}

func (c *covid19) AccDefRate() float64 {
	return math.Floor(c.Body.Items.Item[0].AccDefRate * 100) / 100
}

func (c *covid19) Titles() [3][2]int {
	today := &c.Body.Items.Item[0]
	yesterday := &c.Body.Items.Item[1]
	return [3][2]int {
		{today.DecideCnt, today.DecideCnt - yesterday.DecideCnt},
		{today.ClearCnt, today.ClearCnt - yesterday.ClearCnt},
		{today.DeathCnt, today.DeathCnt - yesterday.DeathCnt},
	}
}
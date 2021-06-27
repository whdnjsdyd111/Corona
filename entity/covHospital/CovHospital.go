package covHospital

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type (
	hospital struct {
		CurrentCount int `json:"currentCount"`
		Data []data `json:"data"`
		MatchCount int `json:"matchCount"`
		Page       int `json:"page"`
		PerPage    int `json:"perPage"`
		TotalCount int `json:"totalCount"`
	}

	data struct {
		Address      string `json:"address"`
		CenterName   string `json:"centerName"`
		CenterType   string `json:"centerType"`
		CreatedAt    string `json:"createdAt"`
		FacilityName string `json:"facilityName"`
		ID           int    `json:"id"`
		Lat          string `json:"lat"`
		Lng          string `json:"lng"`
		Org          string `json:"org"`
		PhoneNumber  string `json:"phoneNumber"`
		Sido         string `json:"sido"`
		Sigungu      string `json:"sigungu"`
		UpdatedAt    string `json:"updatedAt"`
		ZipCode      string `json:"zipCode"`
	}
)

const (
	openApi    = "https://api.odcloud.kr/api/15077586/v1/centers?"
	serviceKey = "AhuzOWhuvKcU%2BuVgEZCPWra%2BqKbvU8XR5NoHgXnXOhRkZwGfyZF9mHpnU3L%2BrNqbcWhYAAgQbgYiiR2RBgWiyQ%3D%3D"
	perPage     = "269"
	url        = openApi + "serviceKey=" + serviceKey + "&perPage=" + perPage
)

func getJSON(url string) ([]byte, error) {
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

func NewCovHospital() *hospital {
	var result hospital
	if jsonBytes, err := getJSON(url); err != nil {
		log.Printf("Failed to get XML: %v", err)
	} else {
		json.Unmarshal(jsonBytes, &result)
	}

	return &result
}

func (h *hospital) GetItems() []data {
	return h.Data
}
package mapstruct

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

type (
	SubStruct struct {
		Date      time.Time `yaml:"date"`
		Latitude  float64   `json:"lat"` //mixed json and yaml tag
		Longitude float64   `json:"lon"`
	}
	MyStruct struct {
		ID   string           //deliberately omit struct tag
		BVal *bool            `json:"b_val"`
		IVal *int             `json:"i_val"`
		UVal *uint8           `json:"u_val"`
		FVal *float64         `json:"f_val"`
		PVal *string          `json:"p_val"`
		AVal [2][]int32       `json:"a_val"`
		LVal [][2]*string     `json:"l_val"`
		MVal map[string][]int `json:"m_val"`
		SVal *SubStruct       `json:"s_val"`
		//PAVal *[2]int32   -- pointer to array, slice or map are NOT supported
	}
)

func TestParse(t *testing.T) {
	m := map[interface{}]interface{}{
		"id":    "ident123",
		"b_val": true,
		"i_val": 8848,
		"f_val": 8848, //deliberately use an integer as float
		"p_val": "pointer to string",
		"u_val": 123,
		"a_val": [3][3]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		"l_val": [][]string{{"a", "b"}, {"c", "d"}, {"e", "f", "g"}},
		"m_val": map[string][]int{"a": {1}, "b": {2, 3}, "c": {4, 5, 6}},
		"s_val": map[interface{}]interface{}{
			"date": "2018-12-21",
			"lat":  27.988056,
			"lon":  86.925278,
		},
	}
	var s MyStruct
	assert(Parse(m, &s))
	je := json.NewEncoder(os.Stdout)
	je.SetIndent("", "    ")
	je.Encode(s)
}

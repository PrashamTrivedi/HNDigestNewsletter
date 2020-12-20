package hndata

//Stories represent stories list from firebase api
type Stories []string

//APIResponse Response of API
type APIResponse struct {
	errorData error
	data      interface{}
}

//Story represent a story data
type Story struct {
	URL        string `json:"url"`
	Title      string `json:"title"`
	By         string `json:"by"`
	Text       string `json:"text"`
	Time       int64 `json:"time"`
	TimeString string `json:"time_string"`
	ID         int `json:"id"`
	Kids       []int
	KidData    []Story
}

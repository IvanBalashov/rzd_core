package rzd_gateway

type Rid struct {
	RID       int64  `json:"RID"`
	Result    string `json:"result"`
	Timestamp string `json:"timestamp"`
}

type Codes struct {
	Name string `json:"n"`
	Code int    `json:"c"`
}

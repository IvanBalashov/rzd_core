package middleware

type Message struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}

type Data struct {
	Direction string `json:"dir"`
	Target    string `json:"target"`
	Source    string `json:"source"`
	Date      string `json:"date"`
}

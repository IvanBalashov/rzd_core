package middleware

type AllTrainsRequest struct {
	Direction string `json:"dir"`
	Target    string `json:"target"`
	Source    string `json:"source"`
	Date      string `json:"date"`
}

type Trains struct {
	TrainID   string  `json:"train_id"`
	MainRoute string  `json:"main_route"`
	Segment   string  `json:"segment"`
	StartDate string  `json:"start_date"`
	Seats     []Seats `json:"seats"`
}

type Seats struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Price int    `json:"price"`
}

type SaveOneTrainRequest struct {
	MainRoute string  `json:"main_route"`
	Segment   string  `json:"segment"`
	StartDate string  `json:"start_date"`
	Seats     []Seats `json:"seats"`
}

type SaveOneTrainResponse struct {
	Status string `json:"status"`
}

package server

type Trains struct {
	TrainID   string           `json:"train_id"`
	MainRoute string           `json:"main_route"`
	Segment   string           `json:"segment"`
	StartDate string           `json:"start_date"`
	EndTime   string           `json:"travel_time"`
	Seats     map[string]Seats `json:"seats"`
}

type Seats struct {
	Count int    `json:"count"`
	Price string `json:"price"`
	Chosen bool  `json:"chosen"`
}

type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

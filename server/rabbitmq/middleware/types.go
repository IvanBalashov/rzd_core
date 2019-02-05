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
	EndTime   string  `json:"travel_time"`
	Seats     []Seats `json:"seats"`
}

type Seats struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Price int    `json:"price"`
}

type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

//TODO: ?
type SaveOneTrainRequest struct {
	User      User    `json:"user"`
	MainRoute string  `json:"main_route"`
	Segment   string  `json:"segment"`
	StartDate string  `json:"start_date"`
	Seats     []Seats `json:"seats"`
}

type Status struct {
	Status string `json:"status"`
}

type UserLength struct {
	Length int    `json:"length"`
	Status string `json:"status"`
}

type CheckUsers struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

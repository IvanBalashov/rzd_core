package middleware

// FIXME: REMOVE AFTER TESTING
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

//<------------------------------

type GetAllTrainsEvent struct {
	Event string   `json:"event"`
	Data  []Trains `json:"data"`
}

type Trains struct {
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

type SaveOneTrainEvent struct {
	Event string `json:"event"`
	Data  Train  `json:"data"`
}

type Train struct {
}

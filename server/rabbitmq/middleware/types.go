package middleware

// REQUESTS
type AllTrainsRequest struct {
	Direction string `json:"dir"`
	Target    string `json:"target"`
	Source    string `json:"source"`
	Date      string `json:"date"`
}

type SaveOneTrainRequest struct {
	Train Trains `json:"train"`
	User  User   `json:"user"`
}

type CheckUsersRequest struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type NewUserRequest struct {
	UserTelegramID string `json:"user_telegram_id"`
	UserName       string `json:"user_name"`
	Notify         bool   `json:"notify"`
}

// RESPONSES
type StatusResponse struct {
	Status string `json:"status"`
}

type UserLengthResponse struct {
	Length int    `json:"length"`
	Status string `json:"status"`
}

type NewUserResponse struct {
	Status string `json:"status"`
}

// HELPERS
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

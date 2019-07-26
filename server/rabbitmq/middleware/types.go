package middleware

import "rzd/server"

// REQUESTS
type AllTrainsRequest struct {
	Direction string `json:"dir"`
	Target    string `json:"target"`
	Source    string `json:"source"`
	Date      string `json:"date"`
}

type SaveOneTrainRequest struct {
	Train server.Trains `json:"train"`
	User  server.User   `json:"user"`
}

type CheckUsersRequest struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
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

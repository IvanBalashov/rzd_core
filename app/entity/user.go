package entity

type User struct {
	UserTelegramID string   `json:"user_telegram_id"`
	UserName       string   `json:"user_name"`
	TrainIDS       []string `json:"train_ids"`
	Notify         bool     `json:"notify"`
}

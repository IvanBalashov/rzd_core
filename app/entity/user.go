package entity

type User struct {
	UserTelegramID string
	UserName       string `json:"user_name"`
	TrainIDS       []string
	Notify         bool
}

//Rewire block for user
//dk how much needed this method, writed for postgresQL.
func (u *User) GetArgs() (string, string, []string, bool) {
	return u.UserTelegramID, u.UserName, u.TrainIDS, u.Notify
}

package entity

type User struct {
	UserTelegramID string
	FullName       string
	Nick           string
	TrainIDS       []string
	Notify         bool
}

//Rewire block for user
//dk how much needed this method, writed for postgresQL.
func (u *User) GetArgs() (string, string, string, []string, bool) {
	return u.UserTelegramID, u.FullName, u.Nick, u.TrainIDS, u.Notify
}

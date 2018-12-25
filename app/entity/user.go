package entity

type User struct {
	FullName string
	Nick     string
	TrainIDS []int64
	Notify   bool
}

func (u *User) GetArgs() (string, string, []int64, bool) {
	return u.FullName, u.Nick, u.TrainIDS, u.Notify
}

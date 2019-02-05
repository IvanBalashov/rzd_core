package middleware

func (e *EventLayer) UsersCount() (interface{}, error) {
	users, err := e.App.GetUsersList()
	if err != nil {
		return nil, err
	}

	return UserLengthResponse{
		Length: len(users),
		Status: "ok",
	}, nil
}

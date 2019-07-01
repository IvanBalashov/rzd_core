package middleware

func (e *EventLayer) UsersCount() (interface{}, error) {
	users, err := e.App.UsersCount()
	if err != nil {
		return nil, err
	}

	return UserLengthResponse{
		Length: users,
		Status: "ok",
	}, nil
}

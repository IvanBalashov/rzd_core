package middleware

func (e *EventLayer) UsersCount() (interface{}, error) {
	users, err := e.App.GetUsersList()
	if err != nil {
		return nil, err
	}

	return UserLength{
		Length: len(users),
		Status: "ok",
	}, nil
}

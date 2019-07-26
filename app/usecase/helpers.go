package usecase

import "rzd/app/entity"

func trainToArgs(t entity.Train) entity.RidArgs {
	return entity.RidArgs{
		Dir:          t.QueryArgs.Dir,
		Tfl:          t.QueryArgs.Tfl,
		CheckSeats:   t.QueryArgs.CheckSeats,
		Code0:        t.QueryArgs.Code0,
		Code1:        t.QueryArgs.Code1,
		Dt0:          t.QueryArgs.Dt0,
		WithOutSeats: t.QueryArgs.WithOutSeats,
		Version:      t.QueryArgs.Version,
	}
}

func fromRzdToCore(seatType string) entity.SeatsType {
	switch seatType {
	case "Купе":
		return entity.CSeatsType
	case "Сидячий":
		return entity.SSeatsType
	case "СВ":
		return entity.SVSeatsType
	case "Плацкартный":
		return entity.PSeatsType
	default:
		return ""
	}
}

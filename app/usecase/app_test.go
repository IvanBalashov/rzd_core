package usecase

import "testing"

var app = App{
	Trains:  nil,
	Users:   nil,
	Routes:  nil,
	Cache:   nil,
	LogChan: nil,
	Cookies: nil,
}

func init() {

}

func TestNewApp(t *testing.T) {
	app := NewApp()
}

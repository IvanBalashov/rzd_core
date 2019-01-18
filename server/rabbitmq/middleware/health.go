package middleware

import "fmt"

func (a *EventLayer) Health() {
	fmt.Printf("ok\n")
}

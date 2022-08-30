package btobet

import (
	"time"

	gohttp "github.com/ochom/go-http"
)

type impl struct {
	http            gohttp.Service
	paymentUsername string
	paymentPassword string
	paymentAPIKey   string
	btobetID        string
}

// New ...
func New(timeout time.Duration, pu, pp, pa, bi string) BtoBet {
	return &impl{
		http:            gohttp.New(timeout),
		paymentUsername: pu,
		paymentPassword: pp,
		paymentAPIKey:   pa,
		btobetID:        bi,
	}
}

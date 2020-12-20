package spl

import (
	"fmt"
	"net/http"
)

// RateLimit ..
type RateLimit struct {
	Total     int
	Remaining int
	ResetTime int
}

// RateLimitNotify ..
func RateLimitNotify(res *http.Response) {
	rates := RateLimit{
		Total:     StrToInt(res.Header.Get("X-RateLimit-Limit")),
		Remaining: StrToInt(res.Header.Get("X-RateLimit-Remaining")),
		ResetTime: StrToInt(res.Header.Get("X-RateLimit-Reset")),
	}

	if rates.Remaining < 100 {
		fmt.Println(rates)
	}
}

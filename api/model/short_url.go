package model

import "time"

type ShortenRequest struct {
	URL         string        `json:"url" validate:"required,http_url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry" validate:"required,numeric,gt=0"`
}

type ShortenResponse struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

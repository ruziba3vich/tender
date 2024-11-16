package models

type Status string

var (
	OPEN    Status = "OPEN"
	CLOSED  Status = "CLOSED"
	AWARDED Status = "AWARDED"
)

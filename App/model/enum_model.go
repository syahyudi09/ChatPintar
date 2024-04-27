package model

type StatusEnum string

const (
	Failed    StatusEnum = "failed"
	Send    StatusEnum = "send"
	Pending StatusEnum = "pending"
)

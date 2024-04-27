package model

type StatusEnum string

const (
	Failed  StatusEnum = "failed"
	Send    StatusEnum = "send"
	Pending StatusEnum = "pending"
)

type RoleEnum string

const (
	Admin  RoleEnum = "admin"
	Member RoleEnum = "member"
)

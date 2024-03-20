package types

type UserStatusEnum int

const (
	Active UserStatusEnum = iota + 1
	Inactive
)

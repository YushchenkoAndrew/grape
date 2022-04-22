package models

type BotRedisDto struct {
	Command string `json:"command" xml:"command" binding:"required"`
}

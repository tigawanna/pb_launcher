package models

type CommandAction int

const (
	ActionStop CommandAction = iota + 1
	ActionStart
	ActionRestart
)

type ServiceCommand struct {
	ID      string        `json:"id"`
	Service string        `json:"service"`
	Action  CommandAction `json:"action"`
}

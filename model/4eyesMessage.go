package model

type ActionType string

const (
	CREATE ActionType = "CREATE"
	UPDATE            = "UPDATE"
	DELETE            = "DELETE"
)

type FourEyesMessageKey struct {
	BoId        uint       `json:"boId"`
	RegistryKey any        `json:"registryKey"`
	Action      ActionType `json:"action"`
}

type FourEyesMessage struct {
	Requester string `json:"Requester"`
	Before    any    `json:"before"`
	After     any    `json:"after"`
}

func NewFourEyesMessageKey(boId uint, registryKey any, action ActionType) *FourEyesMessageKey {
	return &FourEyesMessageKey{BoId: boId, RegistryKey: registryKey, Action: action}
}

func NewFourEyesMessage(before any, after any, requester string) *FourEyesMessage {
	return &FourEyesMessage{Before: before, After: after, Requester: requester}
}

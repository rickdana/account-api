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

func NewFourEyesMessageKey(boId uint, registryKey any, action ActionType) *FourEyesMessageKey {
	return &FourEyesMessageKey{BoId: boId, RegistryKey: registryKey, Action: action}
}

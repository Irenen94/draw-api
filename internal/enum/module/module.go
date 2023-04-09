package module

import "strings"

type Module string

const (
	EdgeService  Module = "EdgeService"
	KafkaService Module = "KafkaService"
	BaseService  Module = "BaseService"
)

func (m Module) ToString() string {
	return string(m)
}

func (m Module) ToLowerString() string {
	return strings.ToLower(string(m))
}

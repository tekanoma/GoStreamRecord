package provider

import "GoRecordurbate/modules/web/provider/chaturbate"

type iProvider interface {
	IsOnline(name string) bool
	IsRoomPublic(name string) bool
}

var Web iProvider

func Init(name string) {

	switch name {
	case "chaturbate":
		Web = &chaturbate.Chaturbare{}
	}
}

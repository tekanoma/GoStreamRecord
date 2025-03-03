package provider

import (
	"GoRecordurbate/modules/web/provider/bongacams"
	"GoRecordurbate/modules/web/provider/chaturbate"
)

// --------------- Must be updated to add support for new sites.
var Known_Providers = map[string]IProvider{
	"chaturbate": &chaturbate.Chaturbate{},
	"bongacams":  &bongacams.BongaCams{},
}
//---------------

// Functions needed to call out providers
type IProvider interface {
	// Code functions
	Init(webType, username string) any
	IsOnline(name string) bool
	TrueName(name string) string

	// Extra
}

var Site IProvider

func Init(name string) IProvider {
	for k := range Known_Providers {
		if k == name {
			return Known_Providers[k]
		}
	}
	return &chaturbate.Chaturbate{} // default option
}

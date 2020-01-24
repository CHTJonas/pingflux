package hosts

import "container/list"

var Endpoints *list.List

func ResetEndpoints() {
	Endpoints = list.New()
}

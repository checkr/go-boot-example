package core

import (
	"reflect"

	"github.com/ant0ine/go-json-rest/rest"
)

var ControllerRegistry = map[string]Creator{}

type Creator func() Controller

type Controller interface {
	Routes(api *rest.Api) []*rest.Route
}

func RegisterController(creator Creator) {
	c := creator()
	ControllerRegistry[reflect.TypeOf(c).String()] = creator
}

func AvailableControllers() []Controller {
	availableControllers := make([]Controller, 0)

	for _, creator := range ControllerRegistry {
		c := creator()
		if handler, ok := reflect.ValueOf(c).Interface().(Controller); ok {
			availableControllers = append(availableControllers, handler)
		}
	}

	return availableControllers
}

// Service Scope used for dependency injection IOC.
package di_plus

import (
	"fmt"
	"reflect"
)

type ServiceScope struct {
	provider *ServiceProvider
	services map[string]*ServiceDescriptor
}

// Gets an instance of the given service or panics.
// Panics if the given service not found (has not been registered).
// Panics if the given service is not transient and not a pointer (the factory must produce a reference: &MyService{}).
// Eg. `my_service := service_scope.GetService("MyObj").(MyObj)`.
func (s ServiceScope) GetService(service_name string) interface{} {
	result, err := s.GetServiceOr(service_name)
	if err != nil {
		panic(err.Error())
	}
	return result
}

// Gets an instance of the given service or error.
// Returns an error if the given service not found (has not been registered).
// Returns an error if the given service is scoped (not allowed in root).
// Returns an error if the given service is not transient and not a pointer (the factory must produce a reference: &MyService{}).
// Eg. `my_service, err := service_scope.GetServiceOr("MyObj").(MyObj)`.
func (s ServiceScope) GetServiceOr(service_name string) (interface{}, error) {
	d, ok := s.services[service_name]
	if !ok {
		d, ok = s.provider.services[service_name]
		if !ok {
			return nil, fmt.Errorf("the service '%s' not found", service_name)
		}
	}
	switch d.serv_type {
	case ServiceTypeEnum.Singleton:
		return s.provider.GetServiceOr(service_name)
	case ServiceTypeEnum.Scoped:
		if d.instance == nil {
			inst, err := d.factory(s)
			if err != nil {
				return nil, err
			}
			v := reflect.ValueOf(inst)
			if v.Kind() != reflect.Pointer {
				return nil, fmt.Errorf("the scoped service '%s' must be a pointer (the factory must return &%s{})", service_name, service_name)
			}
			d.instance = inst
		}
		return d.instance, nil
	default:
		return d.factory(s)
	}
}

// Service Provider used for dependency injection IOC.
package di_plus

import (
	"fmt"
	"reflect"
)

type ServiceProvider struct {
	services map[string]*ServiceDescriptor
}

// Creates a new instance of service provider.
func NewServiceProvider() ServiceProvider {
	return ServiceProvider{services: make(map[string]*ServiceDescriptor)}
}

// Registers a new service with a factory constructor.
// Eg. `service_provider.Register(GetTypeOf(MyService), ServiceTypeEnum.Singleton, func(s IServiceProvider) (interface{}, error) { return &MyService{}, nil } )`.
func (p ServiceProvider) Register(inst_type reflect.Type, serv_type ServiceType, factory func(IServiceProvider) (interface{}, error)) ServiceProvider {
	p.services[inst_type.Name()] = &ServiceDescriptor{inst_type: inst_type, serv_type: serv_type, factory: factory}
	return p
}

// Registers a new service with an instance.
// Only applicable for singleton services.
// Eg. `service_provider.RegisterInst(GetTypeOf(MyService), &MyService{})`.
func (p ServiceProvider) RegisterInst(inst_type reflect.Type, instance interface{}) ServiceProvider {
	v := reflect.ValueOf(instance)
	if v.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("the given instance '%s' must be a pointer (the instance must be &%s{})", inst_type.Name(), inst_type.Name()))
	}
	p.services[inst_type.Name()] = &ServiceDescriptor{inst_type: inst_type, serv_type: ServiceTypeEnum.Singleton, instance: instance}
	return p
}

// Builds a service provider, or just simply returns a reference.
func (p ServiceProvider) Build() *ServiceProvider {
	return &p
}

// Gets an instance of the given service or panics.
// Panics if the given service not found (has not been registered).
// Panics if the given service is scoped (not allowed in root).
// Panics if the given service is not transient and not a pointer (the factory must produce a reference: &MyService{}).
// Eg. `my_service := service_provider.GetService("MyService").(*MyService)`.
func (s ServiceProvider) GetService(service_name string) interface{} {
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
// Eg. `my_service, err := service_provider.GetServiceOr("MyService").(*MyService)`.
func (s ServiceProvider) GetServiceOr(service_name string) (interface{}, error) {
	d, ok := s.services[service_name]
	if !ok {
		return nil, fmt.Errorf("the service '%s' not found", service_name)
	}
	switch d.serv_type {
	case ServiceTypeEnum.Singleton:
		if d.instance == nil {
			inst, err := d.factory(s)
			if err != nil {
				return nil, err
			}
			v := reflect.ValueOf(inst)
			if v.Kind() != reflect.Pointer {
				return nil, fmt.Errorf("the singleton service '%s' must be a pointer (the factory must return &%s{})", service_name, service_name)
			}
			d.instance = inst
		}
		return d.instance, nil
	case ServiceTypeEnum.Scoped:
		return nil, fmt.Errorf("cannot get scoped service '%s' from root", service_name)
	default:
		return d.factory(s)
	}
}

// Creates a service scope.
func (s ServiceProvider) CreateScope() *ServiceScope {
	scope := ServiceScope{provider: &s, services: make(map[string]*ServiceDescriptor)}
	for key, val := range s.services {
		if val.serv_type == ServiceTypeEnum.Scoped {
			scope.services[key] = &ServiceDescriptor{
				inst_type: val.inst_type,
				serv_type: val.serv_type,
				factory:   val.factory,
			}
		}
	}
	return &scope
}

// Returns a type of the given T.
func GetTypeOf[T any]() reflect.Type {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

// Gets an instance of the given service or panics.
// Panics if the given service not found (has not been registered).
// Returns an error if the given service is scoped (not allowed in root).
// Returns an error if the given service is not transient and not a pointer (the factory must produce a reference: &MyService{}).
// Eg. `my_service := GetService[*MyService](service_provider_or_scope)`.
func GetService[T interface{}](s IServiceProvider) T {
	result, err := GetServiceOr[T](s)
	if err != nil {
		panic(err.Error())
	}
	return result
}

// Gets an instance of the given service or error.
// Returns an error if the given service not found (has not been registered).
// Returns an error if the given service is scoped (not allowed in root).
// Returns an error if the given service is not transient and not a pointer (the factory must produce a reference: &MyService{}).
// Eg. `my_service, err := GetServiceOr[*MyService](service_provider_or_scope)`.
func GetServiceOr[T interface{}](s IServiceProvider) (T, error) {
	inst_type := GetTypeOf[T]()
	result, err := s.GetServiceOr(inst_type.Name())
	if err != nil {
		var v T
		return v, err
	}
	// Do not allow instance while the result is a pointer (eg. `GetServiceOr[MyService](s)`)
	if reflect.TypeOf(result).Kind() == reflect.Pointer && reflect.TypeOf((*T)(nil)).Elem().Kind() != reflect.Pointer {
		var v T
		return v, fmt.Errorf("cannot get an instance of the service '%s', use pointer instead", inst_type.Name())
	} else {
		return result.(T), err
	}
}

// Registers a new service with a factory constructor.
// Eg. Register[MyService](service_provider, ServiceTypeEnum.Singleton, func(s IServiceProvider) (MyService, error) { return &MyService{}, nil } ).
func Register[T interface{}](s ServiceProvider, serv_type ServiceType, factory func(IServiceProvider) (T, error)) {
	inst_type := GetTypeOf[T]()
	s.Register(inst_type, serv_type, func(sp IServiceProvider) (interface{}, error) { return factory(sp) })
}

// Registers a new service with an instance.
// Only applicable for singleton services.
// Eg. `Register[MyService](service_provider, &MyService{})`.
func RegisterInst[T interface{}](s ServiceProvider, instance T) {
	inst_type := GetTypeOf[T]()
	s.RegisterInst(inst_type, instance)
}

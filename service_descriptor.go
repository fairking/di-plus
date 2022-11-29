// Service Descriptor used for dependency injection IOC.
package di_plus

import (
	"reflect"
)

type ServiceDescriptor struct {
	inst_type reflect.Type                                // Instance Type
	serv_type ServiceType                                 // Service Type
	instance  interface{}                                 // Service Instance
	factory   func(IServiceProvider) (interface{}, error) // Factory to create an instance
}

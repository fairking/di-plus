// Service Provider Interface used for dependency injection IOC.
package di_plus

type IServiceProvider interface {
	// Gets an instance of the given service or panics
	GetService(string) interface{}
	// Gets an instance of the given service or error
	GetServiceOr(string) (interface{}, error)
}

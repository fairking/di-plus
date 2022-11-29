// Service Type Enum used for dependency injection IOC.
package di_plus

type ServiceType = uint8

// Service Type Enum.
type ServiceTypeStruct struct {
	// Every subsequent request of the singleton service implementation from the service container uses the same instance.
	Singleton ServiceType
	// Scoped lifetime indicates that services are created once per ServiceScope. Do not resolve scoped services from ServiceProvider but rather from ServiceScope.
	Scoped ServiceType
	// Transient lifetime services are created each time they're requested from the service container
	Transient ServiceType
}

var ServiceTypeEnum = &ServiceTypeStruct{
	Singleton: 0,
	Scoped:    1,
	Transient: 2,
}

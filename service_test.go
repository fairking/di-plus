package di_plus

import (
	"testing"
)

type MyService struct {
	Num int
}

type MyService2 struct {
	Obj  *MyService
	Num2 int
}

type MyService3 struct {
	Obj  *MyService
	Num3 int
}

type IMyService interface {
	get() *MyService
}

func (myService *MyService2) get() *MyService {
	return myService.Obj
}

var i int = 0

func getIncrement() int {
	i = i + 1
	return i
}

// Assert that singleton services returned as expected.
func Test_Singleton(t *testing.T) {

	// Build Service Provider
	services := NewServiceProvider().
		Register(
			GetTypeOf[MyService](),
			ServiceTypeEnum.Singleton,
			func(s IServiceProvider) (interface{}, error) {
				return &MyService{Num: getIncrement()}, nil
			},
		).
		Build()

	// Assert can get an instance of struct, the result2 is the same instance of result1.
	{
		result1 := services.GetService("MyService").(*MyService)
		result2 := services.GetService("MyService").(*MyService)

		if result1.Num == 0 || result1.Num != result2.Num {
			t.Fail()
		}
	}
}

// Assert that scoped services returned as expected.
func Test_Scoped(t *testing.T) {

	// Build Service Provider
	services := NewServiceProvider().
		Register(
			GetTypeOf[MyService](),
			ServiceTypeEnum.Scoped,
			func(s IServiceProvider) (interface{}, error) {
				return &MyService{Num: getIncrement()}, nil
			},
		).
		Build()

	// Scope 1
	{
		scope := services.CreateScope()

		result1 := scope.GetService("MyService").(*MyService)
		result2 := scope.GetService("MyService").(*MyService)

		if result1.Num == 0 || result1.Num != result2.Num {
			t.Fail()
		}
	}

	// Scope 2
	{
		scope := services.CreateScope()

		result1 := GetService[*MyService](scope)
		result2 := GetService[*MyService](scope)

		if result1.Num == 0 || result1.Num != result2.Num {
			t.Fail()
		}
	}
}

// Assert that transient services returned as expected.
func Test_Transient(t *testing.T) {

	// Build Service Provider
	services := NewServiceProvider().
		Register(
			GetTypeOf[MyService](),
			ServiceTypeEnum.Transient,
			func(s IServiceProvider) (interface{}, error) {
				// Return an instance, not a pointer
				return MyService{Num: getIncrement()}, nil
			},
		).
		Build()

	// Root
	{
		result1 := GetService[MyService](services)
		result2 := GetService[MyService](services)

		if result1.Num == 0 || result1.Num == result2.Num {
			t.Fail()
		}
	}

	// Scope
	{
		scope := services.CreateScope()

		result1 := GetService[MyService](scope)
		result2 := GetService[MyService](scope)

		if result1.Num == 0 || result1.Num == result2.Num {
			t.Fail()
		}
	}
}

// Assert that the injection works in conjunction with different types of services.
func Test_Injection(t *testing.T) {

	// Build Service Provider
	services := NewServiceProvider().
		Register(
			GetTypeOf[MyService](),
			ServiceTypeEnum.Singleton,
			func(s IServiceProvider) (interface{}, error) {
				return &MyService{
					Num: getIncrement(),
				}, nil
			},
		).
		Register(
			GetTypeOf[MyService2](),
			ServiceTypeEnum.Singleton,
			func(s IServiceProvider) (interface{}, error) {
				return &MyService2{
					Obj: GetService[*MyService](s),
				}, nil
			},
		).
		Build()

	// Root
	{
		result1 := GetService[*MyService](services)
		result2 := GetService[*MyService2](services)

		if result1.Num == 0 || result1.Num != result2.Obj.Num {
			t.Fail()
		}
	}

	// Scope
	{
		scope := services.CreateScope()

		result1 := GetService[*MyService](scope)
		result2 := GetService[*MyService2](scope)

		if result1.Num == 0 || result1.Num != result2.Obj.Num {
			t.Fail()
		}
	}
}

// Assert that the ServiceProvider returns errors.
func Test_Error(t *testing.T) {

	// Build Service Provider
	services := NewServiceProvider().
		Register(
			GetTypeOf[MyService](),
			ServiceTypeEnum.Scoped,
			func(s IServiceProvider) (interface{}, error) {
				return &MyService{
					Num: getIncrement(),
				}, nil
			},
		).
		Register(
			GetTypeOf[MyService2](),
			ServiceTypeEnum.Scoped,
			func(s IServiceProvider) (interface{}, error) {
				// We deliberately return an instance, not a reference in order to test errors
				return MyService2{
					Obj: GetService[*MyService](s),
				}, nil
			},
		).
		Build()

	// Root
	{
		result1, err1 := services.GetServiceOr("MyService")

		if result1 != nil || err1.Error() != "cannot get scoped service 'MyService' from root" {
			t.Fail()
		}

		result2, err2 := services.GetServiceOr("MyService2")

		if result2 != nil || err2.Error() != "cannot get scoped service 'MyService2' from root" {
			t.Fail()
		}

		result3, err3 := services.GetServiceOr("MyService3")

		if result3 != nil || err3.Error() != "the service 'MyService3' not found" {
			t.Fail()
		}
	}

	// Scope
	{
		scope := services.CreateScope()

		result1, err1 := GetServiceOr[MyService](scope)

		if result1.Num != 0 || err1.Error() != "cannot get a value of the service 'MyService', use pointer instead (eg. use `GetServiceOr[*MyService](s)` instead of `GetServiceOr[MyService](s)`)" {
			t.Fail()
		}

		result2, err2 := GetServiceOr[*MyService2](scope)

		if result2 != nil || err2.Error() != "the scoped service 'MyService2' must be a pointer (the factory must return &MyService2{} or interface)" {
			t.Fail()
		}

		result3, err3 := GetServiceOr[*MyService3](scope)

		if result3 != nil || err3.Error() != "the service 'MyService3' not found" {
			t.Fail()
		}
	}
}

// Assert that the service references are correct.
func Test_References(t *testing.T) {

	// Build Service Provider
	services := NewServiceProvider().
		Register(
			GetTypeOf[MyService](),
			ServiceTypeEnum.Singleton,
			func(s IServiceProvider) (interface{}, error) {
				return &MyService{
					Num: getIncrement(),
				}, nil
			},
		).
		Register(
			GetTypeOf[MyService2](),
			ServiceTypeEnum.Scoped,
			func(s IServiceProvider) (interface{}, error) {
				return &MyService2{
					Obj: GetService[*MyService](s),
				}, nil
			},
		).
		Register(
			GetTypeOf[MyService3](),
			ServiceTypeEnum.Transient,
			func(s IServiceProvider) (interface{}, error) {
				return &MyService3{
					Obj: GetService[*MyService](s),
				}, nil
			},
		).
		Build()

	// Root
	{
		result1 := services.GetService("MyService").(*MyService)
		result2 := services.GetService("MyService").(*MyService)

		result1.Num = 999

		if result1.Num != result2.Num {
			t.Fail()
		}
	}

	// Scope
	{
		scope1 := services.CreateScope()

		result1 := GetService[*MyService2](scope1)
		result2 := GetService[*MyService2](scope1)

		result1.Obj.Num = 777

		if result1.Obj.Num != result2.Obj.Num {
			t.Fail()
		}

		result1.Num2 = 333

		if result1.Num2 != result2.Num2 {
			t.Fail()
		}

		scope2 := services.CreateScope()

		result3 := GetService[*MyService2](scope2)
		result4 := GetService[*MyService2](scope2)

		result3.Obj.Num = 999

		if result3.Obj.Num != result4.Obj.Num {
			t.Fail()
		}

		result3.Num2 = 444

		if result3.Num2 != result4.Num2 {
			t.Fail()
		}

		if result1.Num2 == result3.Num2 {
			t.Fail()
		}
	}

	// Transient
	{
		result1 := services.GetService("MyService3").(*MyService3)
		result2 := services.GetService("MyService3").(*MyService3)

		result1.Obj.Num = 999

		if result1.Obj.Num != result2.Obj.Num {
			t.Fail()
		}

		result1.Num3 = 888

		if result1.Num3 == result2.Num3 {
			t.Fail()
		}
	}
}

func Test_Interfaces(t *testing.T) {
	// Build Service Provider
	services := NewServiceProvider().
		Register(
			GetTypeOf[IMyService](),
			ServiceTypeEnum.Singleton,
			func(s IServiceProvider) (interface{}, error) {
				var result IMyService
				result = &MyService2{Obj: &MyService{Num: 111}, Num2: 222}
				return result, nil
			},
		).
		Build()

	// Assert can get an instance of interface
	{
		result1 := GetService[IMyService](services)

		if result1.get().Num != 111 {
			t.Fail()
		}
	}
}

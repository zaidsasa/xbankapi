// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v5"
)

// MockDBTX is an autogenerated mock type for the DBTX type
type MockDBTX struct {
	mock.Mock
}

type MockDBTX_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDBTX) EXPECT() *MockDBTX_Expecter {
	return &MockDBTX_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDBTX) Exec(_a0 context.Context, _a1 string, _a2 ...interface{}) (pgconn.CommandTag, error) {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 pgconn.CommandTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)); ok {
		return rf(_a0, _a1, _a2...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgconn.CommandTag); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Get(0).(pgconn.CommandTag)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDBTX_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockDBTX_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
//   - _a2 ...interface{}
func (_e *MockDBTX_Expecter) Exec(_a0 interface{}, _a1 interface{}, _a2 ...interface{}) *MockDBTX_Exec_Call {
	return &MockDBTX_Exec_Call{Call: _e.mock.On("Exec",
		append([]interface{}{_a0, _a1}, _a2...)...)}
}

func (_c *MockDBTX_Exec_Call) Run(run func(_a0 context.Context, _a1 string, _a2 ...interface{})) *MockDBTX_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockDBTX_Exec_Call) Return(_a0 pgconn.CommandTag, _a1 error) *MockDBTX_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDBTX_Exec_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)) *MockDBTX_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDBTX) Query(_a0 context.Context, _a1 string, _a2 ...interface{}) (pgx.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Query")
	}

	var r0 pgx.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgx.Rows, error)); ok {
		return rf(_a0, _a1, _a2...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Rows); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDBTX_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type MockDBTX_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
//   - _a2 ...interface{}
func (_e *MockDBTX_Expecter) Query(_a0 interface{}, _a1 interface{}, _a2 ...interface{}) *MockDBTX_Query_Call {
	return &MockDBTX_Query_Call{Call: _e.mock.On("Query",
		append([]interface{}{_a0, _a1}, _a2...)...)}
}

func (_c *MockDBTX_Query_Call) Run(run func(_a0 context.Context, _a1 string, _a2 ...interface{})) *MockDBTX_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockDBTX_Query_Call) Return(_a0 pgx.Rows, _a1 error) *MockDBTX_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDBTX_Query_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgx.Rows, error)) *MockDBTX_Query_Call {
	_c.Call.Return(run)
	return _c
}

// QueryRow provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDBTX) QueryRow(_a0 context.Context, _a1 string, _a2 ...interface{}) pgx.Row {
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _a2...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for QueryRow")
	}

	var r0 pgx.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Row); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Row)
		}
	}

	return r0
}

// MockDBTX_QueryRow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryRow'
type MockDBTX_QueryRow_Call struct {
	*mock.Call
}

// QueryRow is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
//   - _a2 ...interface{}
func (_e *MockDBTX_Expecter) QueryRow(_a0 interface{}, _a1 interface{}, _a2 ...interface{}) *MockDBTX_QueryRow_Call {
	return &MockDBTX_QueryRow_Call{Call: _e.mock.On("QueryRow",
		append([]interface{}{_a0, _a1}, _a2...)...)}
}

func (_c *MockDBTX_QueryRow_Call) Run(run func(_a0 context.Context, _a1 string, _a2 ...interface{})) *MockDBTX_QueryRow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockDBTX_QueryRow_Call) Return(_a0 pgx.Row) *MockDBTX_QueryRow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDBTX_QueryRow_Call) RunAndReturn(run func(context.Context, string, ...interface{}) pgx.Row) *MockDBTX_QueryRow_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDBTX creates a new instance of MockDBTX. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDBTX(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDBTX {
	mock := &MockDBTX{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

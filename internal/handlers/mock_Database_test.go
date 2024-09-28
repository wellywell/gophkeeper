// Code generated by mockery v2.43.2. DO NOT EDIT.

package handlers

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	types "github.com/wellywell/gophkeeper/internal/types"
)

// MockDatabase is an autogenerated mock type for the Database type
type MockDatabase struct {
	mock.Mock
}

type MockDatabase_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDatabase) EXPECT() *MockDatabase_Expecter {
	return &MockDatabase_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) CreateUser(_a0 context.Context, _a1 string, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type MockDatabase_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
//   - _a2 string
func (_e *MockDatabase_Expecter) CreateUser(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_CreateUser_Call {
	return &MockDatabase_CreateUser_Call{Call: _e.mock.On("CreateUser", _a0, _a1, _a2)}
}

func (_c *MockDatabase_CreateUser_Call) Run(run func(_a0 context.Context, _a1 string, _a2 string)) *MockDatabase_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockDatabase_CreateUser_Call) Return(_a0 error) *MockDatabase_CreateUser_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_CreateUser_Call) RunAndReturn(run func(context.Context, string, string) error) *MockDatabase_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteItem provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) DeleteItem(_a0 context.Context, _a1 int, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for DeleteItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_DeleteItem_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteItem'
type MockDatabase_DeleteItem_Call struct {
	*mock.Call
}

// DeleteItem is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 string
func (_e *MockDatabase_Expecter) DeleteItem(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_DeleteItem_Call {
	return &MockDatabase_DeleteItem_Call{Call: _e.mock.On("DeleteItem", _a0, _a1, _a2)}
}

func (_c *MockDatabase_DeleteItem_Call) Run(run func(_a0 context.Context, _a1 int, _a2 string)) *MockDatabase_DeleteItem_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(string))
	})
	return _c
}

func (_c *MockDatabase_DeleteItem_Call) Return(_a0 error) *MockDatabase_DeleteItem_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_DeleteItem_Call) RunAndReturn(run func(context.Context, int, string) error) *MockDatabase_DeleteItem_Call {
	_c.Call.Return(run)
	return _c
}

// GetBinaryData provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) GetBinaryData(_a0 context.Context, _a1 int, _a2 string) ([]byte, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for GetBinaryData")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) ([]byte, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, string) []byte); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_GetBinaryData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBinaryData'
type MockDatabase_GetBinaryData_Call struct {
	*mock.Call
}

// GetBinaryData is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 string
func (_e *MockDatabase_Expecter) GetBinaryData(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_GetBinaryData_Call {
	return &MockDatabase_GetBinaryData_Call{Call: _e.mock.On("GetBinaryData", _a0, _a1, _a2)}
}

func (_c *MockDatabase_GetBinaryData_Call) Run(run func(_a0 context.Context, _a1 int, _a2 string)) *MockDatabase_GetBinaryData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(string))
	})
	return _c
}

func (_c *MockDatabase_GetBinaryData_Call) Return(_a0 []byte, _a1 error) *MockDatabase_GetBinaryData_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_GetBinaryData_Call) RunAndReturn(run func(context.Context, int, string) ([]byte, error)) *MockDatabase_GetBinaryData_Call {
	_c.Call.Return(run)
	return _c
}

// GetCreditCard provides a mock function with given fields: _a0, _a1
func (_m *MockDatabase) GetCreditCard(_a0 context.Context, _a1 int) (*types.CreditCardData, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetCreditCard")
	}

	var r0 *types.CreditCardData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*types.CreditCardData, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *types.CreditCardData); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.CreditCardData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_GetCreditCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCreditCard'
type MockDatabase_GetCreditCard_Call struct {
	*mock.Call
}

// GetCreditCard is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
func (_e *MockDatabase_Expecter) GetCreditCard(_a0 interface{}, _a1 interface{}) *MockDatabase_GetCreditCard_Call {
	return &MockDatabase_GetCreditCard_Call{Call: _e.mock.On("GetCreditCard", _a0, _a1)}
}

func (_c *MockDatabase_GetCreditCard_Call) Run(run func(_a0 context.Context, _a1 int)) *MockDatabase_GetCreditCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockDatabase_GetCreditCard_Call) Return(_a0 *types.CreditCardData, _a1 error) *MockDatabase_GetCreditCard_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_GetCreditCard_Call) RunAndReturn(run func(context.Context, int) (*types.CreditCardData, error)) *MockDatabase_GetCreditCard_Call {
	_c.Call.Return(run)
	return _c
}

// GetItem provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) GetItem(_a0 context.Context, _a1 int, _a2 string) (*types.Item, error) {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for GetItem")
	}

	var r0 *types.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) (*types.Item, error)); ok {
		return rf(_a0, _a1, _a2)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, string) *types.Item); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(_a0, _a1, _a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_GetItem_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetItem'
type MockDatabase_GetItem_Call struct {
	*mock.Call
}

// GetItem is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 string
func (_e *MockDatabase_Expecter) GetItem(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_GetItem_Call {
	return &MockDatabase_GetItem_Call{Call: _e.mock.On("GetItem", _a0, _a1, _a2)}
}

func (_c *MockDatabase_GetItem_Call) Run(run func(_a0 context.Context, _a1 int, _a2 string)) *MockDatabase_GetItem_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(string))
	})
	return _c
}

func (_c *MockDatabase_GetItem_Call) Return(_a0 *types.Item, _a1 error) *MockDatabase_GetItem_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_GetItem_Call) RunAndReturn(run func(context.Context, int, string) (*types.Item, error)) *MockDatabase_GetItem_Call {
	_c.Call.Return(run)
	return _c
}

// GetItems provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *MockDatabase) GetItems(_a0 context.Context, _a1 int, _a2 int, _a3 int) ([]types.Item, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for GetItems")
	}

	var r0 []types.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int) ([]types.Item, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int) []types.Item); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int, int) error); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_GetItems_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetItems'
type MockDatabase_GetItems_Call struct {
	*mock.Call
}

// GetItems is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 int
//   - _a3 int
func (_e *MockDatabase_Expecter) GetItems(_a0 interface{}, _a1 interface{}, _a2 interface{}, _a3 interface{}) *MockDatabase_GetItems_Call {
	return &MockDatabase_GetItems_Call{Call: _e.mock.On("GetItems", _a0, _a1, _a2, _a3)}
}

func (_c *MockDatabase_GetItems_Call) Run(run func(_a0 context.Context, _a1 int, _a2 int, _a3 int)) *MockDatabase_GetItems_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(int), args[3].(int))
	})
	return _c
}

func (_c *MockDatabase_GetItems_Call) Return(_a0 []types.Item, _a1 error) *MockDatabase_GetItems_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_GetItems_Call) RunAndReturn(run func(context.Context, int, int, int) ([]types.Item, error)) *MockDatabase_GetItems_Call {
	_c.Call.Return(run)
	return _c
}

// GetLogoPass provides a mock function with given fields: _a0, _a1
func (_m *MockDatabase) GetLogoPass(_a0 context.Context, _a1 int) (*types.LoginPassword, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetLogoPass")
	}

	var r0 *types.LoginPassword
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*types.LoginPassword, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *types.LoginPassword); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.LoginPassword)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_GetLogoPass_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetLogoPass'
type MockDatabase_GetLogoPass_Call struct {
	*mock.Call
}

// GetLogoPass is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
func (_e *MockDatabase_Expecter) GetLogoPass(_a0 interface{}, _a1 interface{}) *MockDatabase_GetLogoPass_Call {
	return &MockDatabase_GetLogoPass_Call{Call: _e.mock.On("GetLogoPass", _a0, _a1)}
}

func (_c *MockDatabase_GetLogoPass_Call) Run(run func(_a0 context.Context, _a1 int)) *MockDatabase_GetLogoPass_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockDatabase_GetLogoPass_Call) Return(_a0 *types.LoginPassword, _a1 error) *MockDatabase_GetLogoPass_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_GetLogoPass_Call) RunAndReturn(run func(context.Context, int) (*types.LoginPassword, error)) *MockDatabase_GetLogoPass_Call {
	_c.Call.Return(run)
	return _c
}

// GetText provides a mock function with given fields: _a0, _a1
func (_m *MockDatabase) GetText(_a0 context.Context, _a1 int) (*types.TextData, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetText")
	}

	var r0 *types.TextData
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (*types.TextData, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) *types.TextData); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.TextData)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_GetText_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetText'
type MockDatabase_GetText_Call struct {
	*mock.Call
}

// GetText is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
func (_e *MockDatabase_Expecter) GetText(_a0 interface{}, _a1 interface{}) *MockDatabase_GetText_Call {
	return &MockDatabase_GetText_Call{Call: _e.mock.On("GetText", _a0, _a1)}
}

func (_c *MockDatabase_GetText_Call) Run(run func(_a0 context.Context, _a1 int)) *MockDatabase_GetText_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int))
	})
	return _c
}

func (_c *MockDatabase_GetText_Call) Return(_a0 *types.TextData, _a1 error) *MockDatabase_GetText_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_GetText_Call) RunAndReturn(run func(context.Context, int) (*types.TextData, error)) *MockDatabase_GetText_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserHashedPassword provides a mock function with given fields: _a0, _a1
func (_m *MockDatabase) GetUserHashedPassword(_a0 context.Context, _a1 string) (string, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetUserHashedPassword")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_GetUserHashedPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserHashedPassword'
type MockDatabase_GetUserHashedPassword_Call struct {
	*mock.Call
}

// GetUserHashedPassword is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *MockDatabase_Expecter) GetUserHashedPassword(_a0 interface{}, _a1 interface{}) *MockDatabase_GetUserHashedPassword_Call {
	return &MockDatabase_GetUserHashedPassword_Call{Call: _e.mock.On("GetUserHashedPassword", _a0, _a1)}
}

func (_c *MockDatabase_GetUserHashedPassword_Call) Run(run func(_a0 context.Context, _a1 string)) *MockDatabase_GetUserHashedPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabase_GetUserHashedPassword_Call) Return(_a0 string, _a1 error) *MockDatabase_GetUserHashedPassword_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_GetUserHashedPassword_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockDatabase_GetUserHashedPassword_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserID provides a mock function with given fields: _a0, _a1
func (_m *MockDatabase) GetUserID(_a0 context.Context, _a1 string) (int, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetUserID")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (int, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDatabase_GetUserID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserID'
type MockDatabase_GetUserID_Call struct {
	*mock.Call
}

// GetUserID is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 string
func (_e *MockDatabase_Expecter) GetUserID(_a0 interface{}, _a1 interface{}) *MockDatabase_GetUserID_Call {
	return &MockDatabase_GetUserID_Call{Call: _e.mock.On("GetUserID", _a0, _a1)}
}

func (_c *MockDatabase_GetUserID_Call) Run(run func(_a0 context.Context, _a1 string)) *MockDatabase_GetUserID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockDatabase_GetUserID_Call) Return(_a0 int, _a1 error) *MockDatabase_GetUserID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDatabase_GetUserID_Call) RunAndReturn(run func(context.Context, string) (int, error)) *MockDatabase_GetUserID_Call {
	_c.Call.Return(run)
	return _c
}

// InsertBinaryData provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) InsertBinaryData(_a0 context.Context, _a1 int, _a2 types.BinaryItem) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for InsertBinaryData")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, types.BinaryItem) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_InsertBinaryData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertBinaryData'
type MockDatabase_InsertBinaryData_Call struct {
	*mock.Call
}

// InsertBinaryData is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 types.BinaryItem
func (_e *MockDatabase_Expecter) InsertBinaryData(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_InsertBinaryData_Call {
	return &MockDatabase_InsertBinaryData_Call{Call: _e.mock.On("InsertBinaryData", _a0, _a1, _a2)}
}

func (_c *MockDatabase_InsertBinaryData_Call) Run(run func(_a0 context.Context, _a1 int, _a2 types.BinaryItem)) *MockDatabase_InsertBinaryData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(types.BinaryItem))
	})
	return _c
}

func (_c *MockDatabase_InsertBinaryData_Call) Return(_a0 error) *MockDatabase_InsertBinaryData_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_InsertBinaryData_Call) RunAndReturn(run func(context.Context, int, types.BinaryItem) error) *MockDatabase_InsertBinaryData_Call {
	_c.Call.Return(run)
	return _c
}

// InsertCreditCard provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) InsertCreditCard(_a0 context.Context, _a1 int, _a2 types.CreditCardItem) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for InsertCreditCard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, types.CreditCardItem) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_InsertCreditCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertCreditCard'
type MockDatabase_InsertCreditCard_Call struct {
	*mock.Call
}

// InsertCreditCard is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 types.CreditCardItem
func (_e *MockDatabase_Expecter) InsertCreditCard(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_InsertCreditCard_Call {
	return &MockDatabase_InsertCreditCard_Call{Call: _e.mock.On("InsertCreditCard", _a0, _a1, _a2)}
}

func (_c *MockDatabase_InsertCreditCard_Call) Run(run func(_a0 context.Context, _a1 int, _a2 types.CreditCardItem)) *MockDatabase_InsertCreditCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(types.CreditCardItem))
	})
	return _c
}

func (_c *MockDatabase_InsertCreditCard_Call) Return(_a0 error) *MockDatabase_InsertCreditCard_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_InsertCreditCard_Call) RunAndReturn(run func(context.Context, int, types.CreditCardItem) error) *MockDatabase_InsertCreditCard_Call {
	_c.Call.Return(run)
	return _c
}

// InsertLogoPass provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) InsertLogoPass(_a0 context.Context, _a1 int, _a2 types.LoginPasswordItem) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for InsertLogoPass")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, types.LoginPasswordItem) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_InsertLogoPass_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertLogoPass'
type MockDatabase_InsertLogoPass_Call struct {
	*mock.Call
}

// InsertLogoPass is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 types.LoginPasswordItem
func (_e *MockDatabase_Expecter) InsertLogoPass(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_InsertLogoPass_Call {
	return &MockDatabase_InsertLogoPass_Call{Call: _e.mock.On("InsertLogoPass", _a0, _a1, _a2)}
}

func (_c *MockDatabase_InsertLogoPass_Call) Run(run func(_a0 context.Context, _a1 int, _a2 types.LoginPasswordItem)) *MockDatabase_InsertLogoPass_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(types.LoginPasswordItem))
	})
	return _c
}

func (_c *MockDatabase_InsertLogoPass_Call) Return(_a0 error) *MockDatabase_InsertLogoPass_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_InsertLogoPass_Call) RunAndReturn(run func(context.Context, int, types.LoginPasswordItem) error) *MockDatabase_InsertLogoPass_Call {
	_c.Call.Return(run)
	return _c
}

// InsertText provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) InsertText(_a0 context.Context, _a1 int, _a2 types.TextItem) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for InsertText")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, types.TextItem) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_InsertText_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertText'
type MockDatabase_InsertText_Call struct {
	*mock.Call
}

// InsertText is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 types.TextItem
func (_e *MockDatabase_Expecter) InsertText(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_InsertText_Call {
	return &MockDatabase_InsertText_Call{Call: _e.mock.On("InsertText", _a0, _a1, _a2)}
}

func (_c *MockDatabase_InsertText_Call) Run(run func(_a0 context.Context, _a1 int, _a2 types.TextItem)) *MockDatabase_InsertText_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(types.TextItem))
	})
	return _c
}

func (_c *MockDatabase_InsertText_Call) Return(_a0 error) *MockDatabase_InsertText_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_InsertText_Call) RunAndReturn(run func(context.Context, int, types.TextItem) error) *MockDatabase_InsertText_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateBinaryData provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) UpdateBinaryData(_a0 context.Context, _a1 int, _a2 types.BinaryItem) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for UpdateBinaryData")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, types.BinaryItem) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_UpdateBinaryData_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateBinaryData'
type MockDatabase_UpdateBinaryData_Call struct {
	*mock.Call
}

// UpdateBinaryData is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 types.BinaryItem
func (_e *MockDatabase_Expecter) UpdateBinaryData(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_UpdateBinaryData_Call {
	return &MockDatabase_UpdateBinaryData_Call{Call: _e.mock.On("UpdateBinaryData", _a0, _a1, _a2)}
}

func (_c *MockDatabase_UpdateBinaryData_Call) Run(run func(_a0 context.Context, _a1 int, _a2 types.BinaryItem)) *MockDatabase_UpdateBinaryData_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(types.BinaryItem))
	})
	return _c
}

func (_c *MockDatabase_UpdateBinaryData_Call) Return(_a0 error) *MockDatabase_UpdateBinaryData_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_UpdateBinaryData_Call) RunAndReturn(run func(context.Context, int, types.BinaryItem) error) *MockDatabase_UpdateBinaryData_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateCreditCard provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) UpdateCreditCard(_a0 context.Context, _a1 int, _a2 types.CreditCardItem) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCreditCard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, types.CreditCardItem) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_UpdateCreditCard_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateCreditCard'
type MockDatabase_UpdateCreditCard_Call struct {
	*mock.Call
}

// UpdateCreditCard is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 types.CreditCardItem
func (_e *MockDatabase_Expecter) UpdateCreditCard(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_UpdateCreditCard_Call {
	return &MockDatabase_UpdateCreditCard_Call{Call: _e.mock.On("UpdateCreditCard", _a0, _a1, _a2)}
}

func (_c *MockDatabase_UpdateCreditCard_Call) Run(run func(_a0 context.Context, _a1 int, _a2 types.CreditCardItem)) *MockDatabase_UpdateCreditCard_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(types.CreditCardItem))
	})
	return _c
}

func (_c *MockDatabase_UpdateCreditCard_Call) Return(_a0 error) *MockDatabase_UpdateCreditCard_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_UpdateCreditCard_Call) RunAndReturn(run func(context.Context, int, types.CreditCardItem) error) *MockDatabase_UpdateCreditCard_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateLogoPass provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) UpdateLogoPass(_a0 context.Context, _a1 int, _a2 types.LoginPasswordItem) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for UpdateLogoPass")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, types.LoginPasswordItem) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_UpdateLogoPass_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateLogoPass'
type MockDatabase_UpdateLogoPass_Call struct {
	*mock.Call
}

// UpdateLogoPass is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 types.LoginPasswordItem
func (_e *MockDatabase_Expecter) UpdateLogoPass(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_UpdateLogoPass_Call {
	return &MockDatabase_UpdateLogoPass_Call{Call: _e.mock.On("UpdateLogoPass", _a0, _a1, _a2)}
}

func (_c *MockDatabase_UpdateLogoPass_Call) Run(run func(_a0 context.Context, _a1 int, _a2 types.LoginPasswordItem)) *MockDatabase_UpdateLogoPass_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(types.LoginPasswordItem))
	})
	return _c
}

func (_c *MockDatabase_UpdateLogoPass_Call) Return(_a0 error) *MockDatabase_UpdateLogoPass_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_UpdateLogoPass_Call) RunAndReturn(run func(context.Context, int, types.LoginPasswordItem) error) *MockDatabase_UpdateLogoPass_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateText provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockDatabase) UpdateText(_a0 context.Context, _a1 int, _a2 types.TextItem) error {
	ret := _m.Called(_a0, _a1, _a2)

	if len(ret) == 0 {
		panic("no return value specified for UpdateText")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, types.TextItem) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockDatabase_UpdateText_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateText'
type MockDatabase_UpdateText_Call struct {
	*mock.Call
}

// UpdateText is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 int
//   - _a2 types.TextItem
func (_e *MockDatabase_Expecter) UpdateText(_a0 interface{}, _a1 interface{}, _a2 interface{}) *MockDatabase_UpdateText_Call {
	return &MockDatabase_UpdateText_Call{Call: _e.mock.On("UpdateText", _a0, _a1, _a2)}
}

func (_c *MockDatabase_UpdateText_Call) Run(run func(_a0 context.Context, _a1 int, _a2 types.TextItem)) *MockDatabase_UpdateText_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int), args[2].(types.TextItem))
	})
	return _c
}

func (_c *MockDatabase_UpdateText_Call) Return(_a0 error) *MockDatabase_UpdateText_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockDatabase_UpdateText_Call) RunAndReturn(run func(context.Context, int, types.TextItem) error) *MockDatabase_UpdateText_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDatabase creates a new instance of MockDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDatabase {
	mock := &MockDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

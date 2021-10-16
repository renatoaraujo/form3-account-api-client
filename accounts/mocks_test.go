// Code generated by mockery v2.9.4. DO NOT EDIT.

package accounts

import mock "github.com/stretchr/testify/mock"

// httpUtils is an autogenerated mock type for the httpUtils type
type mockHttpUtils struct {
	mock.Mock
}

// Delete provides a mock function with given fields: resourcePath
func (_m *mockHttpUtils) Delete(resourcePath string) error {
	ret := _m.Called(resourcePath)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(resourcePath)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: resourcePath
func (_m *mockHttpUtils) Get(resourcePath string) ([]byte, error) {
	ret := _m.Called(resourcePath)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string) []byte); ok {
		r0 = rf(resourcePath)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(resourcePath)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Post provides a mock function with given fields: resourcePath, payload
func (_m *mockHttpUtils) Post(resourcePath string, payload []byte) ([]byte, error) {
	ret := _m.Called(resourcePath, payload)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(string, []byte) []byte); ok {
		r0 = rf(resourcePath, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, []byte) error); ok {
		r1 = rf(resourcePath, payload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

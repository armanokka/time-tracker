// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/armanokka/test_task_Effective_mobile/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// UseCase is an autogenerated mock type for the UseCase type
type UseCase struct {
	mock.Mock
}

// AddMember provides a mock function with given fields: ctx, projectID, userID
func (_m *UseCase) AddMember(ctx context.Context, projectID int64, userID int64) error {
	ret := _m.Called(ctx, projectID, userID)

	if len(ret) == 0 {
		panic("no return value specified for AddMember")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, projectID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, project
func (_m *UseCase) Create(ctx context.Context, project *models.Project) (*models.Project, error) {
	ret := _m.Called(ctx, project)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *models.Project
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Project) (*models.Project, error)); ok {
		return rf(ctx, project)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Project) *models.Project); ok {
		r0 = rf(ctx, project)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Project)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Project) error); ok {
		r1 = rf(ctx, project)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, projectID
func (_m *UseCase) Delete(ctx context.Context, projectID int64) error {
	ret := _m.Called(ctx, projectID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, projectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, projectID
func (_m *UseCase) GetByID(ctx context.Context, projectID int64) (*models.Project, error) {
	ret := _m.Called(ctx, projectID)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *models.Project
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*models.Project, error)); ok {
		return rf(ctx, projectID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *models.Project); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Project)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMemberProductivity provides a mock function with given fields: ctx, projectID, userID
func (_m *UseCase) GetMemberProductivity(ctx context.Context, projectID int64, userID int64) ([]models.UserProductivity, error) {
	ret := _m.Called(ctx, projectID, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetMemberProductivity")
	}

	var r0 []models.UserProductivity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) ([]models.UserProductivity, error)); ok {
		return rf(ctx, projectID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) []models.UserProductivity); ok {
		r0 = rf(ctx, projectID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.UserProductivity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, projectID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMembers provides a mock function with given fields: ctx, projectID
func (_m *UseCase) GetMembers(ctx context.Context, projectID int64) ([]*models.User, error) {
	ret := _m.Called(ctx, projectID)

	if len(ret) == 0 {
		panic("no return value specified for GetMembers")
	}

	var r0 []*models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]*models.User, error)); ok {
		return rf(ctx, projectID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*models.User); ok {
		r0 = rf(ctx, projectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, projectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsMember provides a mock function with given fields: ctx, projectID, userID
func (_m *UseCase) IsMember(ctx context.Context, projectID int64, userID int64) (bool, error) {
	ret := _m.Called(ctx, projectID, userID)

	if len(ret) == 0 {
		panic("no return value specified for IsMember")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (bool, error)); ok {
		return rf(ctx, projectID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) bool); ok {
		r0 = rf(ctx, projectID, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, projectID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsOwner provides a mock function with given fields: ctx, projectID, userID
func (_m *UseCase) IsOwner(ctx context.Context, projectID int64, userID int64) (bool, error) {
	ret := _m.Called(ctx, projectID, userID)

	if len(ret) == 0 {
		panic("no return value specified for IsOwner")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (bool, error)); ok {
		return rf(ctx, projectID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) bool); ok {
		r0 = rf(ctx, projectID, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, projectID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveMember provides a mock function with given fields: ctx, projectID, userID
func (_m *UseCase) RemoveMember(ctx context.Context, projectID int64, userID int64) error {
	ret := _m.Called(ctx, projectID, userID)

	if len(ret) == 0 {
		panic("no return value specified for RemoveMember")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, projectID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, updates
func (_m *UseCase) Update(ctx context.Context, updates *models.Project) (*models.Project, error) {
	ret := _m.Called(ctx, updates)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *models.Project
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Project) (*models.Project, error)); ok {
		return rf(ctx, updates)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Project) *models.Project); ok {
		r0 = rf(ctx, updates)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Project)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Project) error); ok {
		r1 = rf(ctx, updates)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUseCase creates a new instance of UseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *UseCase {
	mock := &UseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

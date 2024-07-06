package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/armanokka/test_task_Effective_mobile/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProjectsRepo_Create(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()

	mock.ExpectQuery(createProjectQuery).
		WithArgs(project.Name, project.Description, project.CreatorID).
		WillReturnRows(sqlmock.NewRows(project.Columns()).AddRow(project.Fields()...))
	mock.ExpectExec(addProjectMemberQuery).
		WithArgs(project.ID, project.CreatorID).WillReturnResult(driver.ResultNoRows).WillReturnError(nil)

	gotProject, err := projectRepo.Create(context.Background(), project)
	assert.Nil(t, err)
	assert.Equal(t, project, gotProject)
}

func TestProjectsRepo_AddMember(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()
	var userID int64 = 123

	mock.ExpectExec(addProjectMemberQuery).WithArgs(project.ID, userID).
		WillReturnResult(driver.ResultNoRows).WillReturnError(nil)

	assert.Nil(t, projectRepo.AddMember(context.Background(), project.ID, userID))

	mock.ExpectExec(addProjectMemberQuery).WithArgs(project.ID, userID).
		WillReturnResult(driver.ResultNoRows).WillReturnError(fmt.Errorf("error"))
	assert.NotNil(t, projectRepo.AddMember(context.Background(), project.ID, userID))
}

func TestProjectsRepo_Delete(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()
	var userID int64 = 123

	mock.ExpectExec(removeProjectMemberQuery).WithArgs(project.ID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)

	assert.Nil(t, projectRepo.RemoveMember(context.Background(), project.ID, userID))

	mock.ExpectExec(removeProjectMemberQuery).WithArgs(project.ID, userID).
		WillReturnResult(driver.ResultNoRows).WillReturnError(sql.ErrNoRows)

	assert.NotNil(t, projectRepo.AddMember(context.Background(), project.ID, userID))
}

func TestProjectsRepo_GetByID(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()

	mock.ExpectQuery(getProjectByIDQuery).WithArgs(project.ID).WillReturnRows(
		sqlmock.NewRows(project.Columns()).AddRow(project.Fields()...))

	gotProject, err := projectRepo.GetByID(context.Background(), project.ID)
	assert.Nil(t, err)
	assert.Equal(t, project, gotProject)
}

func TestProjectsRepo_IsMember(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()
	var userID int64 = 123

	mock.ExpectExec(isProjectMemberQuery).WithArgs(project.ID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)

	assert.Nil(t, projectRepo.IsMember(context.Background(), project.ID, userID))

	mock.ExpectExec(isProjectMemberQuery).WithArgs(project.ID, userID).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows)).WillReturnError(sql.ErrNoRows)

	assert.NotNil(t, projectRepo.IsMember(context.Background(), project.ID, userID))
}

func TestProjectsRepo_GetMemberProductivity(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()
	var userID int64 = 123

	mock.ExpectQuery(getProjectMemberProductivityQuery).WithArgs(project.ID, userID).
		WillReturnRows(
			sqlmock.NewRows([]string{"task_id", "total_seconds"}).
				AddRow(1, 90*60).
				AddRow(2, 15*60),
		)
	needProductivity := []models.UserProductivity{
		{
			TaskID:       1,
			SpentHours:   1,
			SpentMinutes: 30,
		},
		{
			TaskID:       2,
			SpentHours:   0,
			SpentMinutes: 15,
		},
	}

	gotProductivity, err := projectRepo.GetMemberProductivity(context.Background(), project.ID, userID)
	assert.Nil(t, err)
	assert.Equal(t, needProductivity, gotProductivity)
}

func TestProjectsRepo_GetMembers(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()
	user1 := models.User{
		ID:             1,
		Email:          "sdfsd",
		Password:       "sdfsd",
		Name:           "sdfsd",
		Surname:        "sdfsd",
		Patronymic:     nil,
		Address:        "sdfsd",
		Admin:          false,
		PassportNumber: 3932,
		PassportSeries: 202202,
	}
	user2 := models.User{
		ID:             2,
		Email:          "ewprppd",
		Password:       "ewprppd",
		Name:           "sdfpo",
		Surname:        "owow",
		Patronymic:     nil,
		Address:        "9339",
		Admin:          false,
		PassportNumber: 9292,
		PassportSeries: 303033,
	}
	members := []*models.User{&user1, &user2}

	mock.ExpectQuery(getProjectMembersCount).WithArgs(project.ID).WillReturnRows(
		sqlmock.NewRows([]string{"result"}).AddRow(len(members)))
	mock.ExpectQuery(getProjectMembers).WithArgs(project.ID).WillReturnRows(
		sqlmock.NewRows(user1.Columns()).
			AddRow(user1.Rows()...).
			AddRow(user2.Rows()...),
	)

	gotMembers, err := projectRepo.GetMembers(context.Background(), project.ID)
	assert.Nil(t, err)
	assert.Equal(t, members, gotMembers)
}

func TestProjectsRepo_IsOwner(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()
	var userID int64 = 123

	mock.ExpectExec(isProjectOwnerQuery).WithArgs(project.ID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(nil)

	assert.Nil(t, projectRepo.IsOwner(context.Background(), project.ID, userID))

	mock.ExpectExec(isProjectOwnerQuery).WithArgs(project.ID, userID).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows)).WillReturnError(sql.ErrNoRows)

	assert.NotNil(t, projectRepo.IsOwner(context.Background(), project.ID, userID))
}

func TestProjectsRepo_RemoveMember(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()
	var userID int64 = 123

	mock.ExpectExec(removeProjectMemberQuery).WithArgs(project.ID, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	assert.Nil(t, projectRepo.RemoveMember(context.Background(), project.ID, userID))

	mock.ExpectExec(removeProjectMemberQuery).WithArgs(project.ID, userID).
		WillReturnResult(sqlmock.NewResult(0, 0))
	assert.NotNil(t, projectRepo.RemoveMember(context.Background(), project.ID, userID))
}

func TestProjectsRepo_Update(t *testing.T) {
	projectRepo, db, mock, err := newMockProjectsRepo()
	require.NoError(t, err)
	defer db.Close()

	project := getTestProject()

	mock.ExpectQuery(updateProjectQuery).WithArgs(project.Name, project.Description, project.CreatorID, project.ID).
		WillReturnRows(sqlmock.NewRows(project.Columns()).AddRow(project.Fields()...))
	gotProject, err := projectRepo.Update(context.Background(), project)
	assert.Nil(t, err)
	assert.Equal(t, project, gotProject)

	mock.ExpectQuery(updateProjectQuery).WithArgs(project.Name, project.Description, project.CreatorID, project.ID).
		WillReturnError(sql.ErrNoRows)
	gotProject, err = projectRepo.Update(context.Background(), project)
	assert.NotNil(t, gotProject)
}

package repository

// I didn't have enough time to finish this. Therefore, project is not fully covered with tests

//func TestProjectsRepo_Create(t *testing.T) {
//	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
//	require.NoError(t, err)
//	defer db.Close()
//	sqlxDB := sqlx.NewDb(db, "sqlmock")
//	defer sqlxDB.Close()
//
//	project := getTestProject()
//	projectRepo := NewProjectsRepository(sqlxDB)
//
//	mock.ExpectQuery(createProjectQuery).
//		WithArgs(project.Name, project.Description, project.CreatorID).
//		WillReturnRows(sqlmock.NewRows(project.Columns()).AddRow(project.Fields()...))
//	mock.ExpectQuery(addProjectMemberQuery).
//		WithArgs(project.ID, project.CreatorID).
//		WillReturnRows(sqlmock.NewRows([]string{})).WillReturnError(nil)
//
//	gotProject, err := projectRepo.Create(context.Background(), project)
//	assert.Nil(t, err)
//	assert.Equal(t, project, gotProject)
//}

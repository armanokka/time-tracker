package repository

const (
	createTaskQuery = `INSERT INTO tasks (name, description, project_id) 
VALUES ($1, $2, $3) RETURNING *`
	getTotalTasks     = `SELECT COUNT(id) FROM tasks WHERE project_id = $1`
	selectTasks       = `SELECT tasks.* FROM tasks WHERE project_id = $1`
	isTaskMemberQuery = `SELECT 1 FROM task_participants WHERE task_id = $1 AND user_id = $2 LIMIT 1`
	updateTaskQuery   = `UPDATE tasks SET
name = COALESCE(NULLIF($1, ''), name),
description = COALESCE(NULLIF($2, ''), description)
WHERE id = $3
RETURNING *`
	deleteTaskQuery = `DELETE FROM tasks WHERE id = $1`
	startTaskQuery  = `INSERT INTO time_entries (task_id, user_id, started_at, ended_at)
VALUES ($1, $2, now(), null)`
	endTaskQuery = `UPDATE time_entries SET ended_at = now()
WHERE ended_at IS NULL AND task_id = $1 AND user_id = $2`
	getActiveUserTasksQuery = `SELECT count(1) FROM tasks
INNER JOIN time_entries on time_entries.task_id = tasks.id
WHERE time_entries.ended_at IS NULL
AND time_entries.user_id = $1
AND tasks.project_id = (SELECT project_id FROM tasks WHERE id = $2)`

	getTotalTaskMembersQuery = `SELECT count(user_id) FROM task_participants WHERE task_id = $1`
	getTaskMembersQuery      = `SELECT users.* FROM users
INNER JOIN task_participants ON task_participants.user_id = users.id
WHERE task_id = $1`
	addTaskMemberQuery    = `INSERT INTO task_participants (task_id, user_id) VALUES ($1, $2)`
	deleteTaskMemberQuery = `DELETE FROM task_participants WHERE task_id = $1 AND user_id = $2`
	taskExistsQuery       = `SELECT EXISTS (SELECT 1 FROM tasks WHERE id = $1)`
)

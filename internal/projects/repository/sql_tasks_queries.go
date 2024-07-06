package repository

const (
	createTaskQuery = `INSERT INTO task (name, description, project_id) 
VALUES ($1, $2, $3) RETURNING *`
	getTotalTasks     = `SELECT COUNT(id) FROM task WHERE project_id = $1`
	selectTasks       = `SELECT task.* FROM task WHERE project_id = $1`
	isTaskMemberQuery = `SELECT FROM task_participant WHERE task_id = $1 AND user_id = $2 LIMIT 1`
	updateTaskQuery   = `UPDATE task SET
name = COALESCE(NULLIF($1, ''), name),
description = COALESCE(NULLIF($2, ''), description)
WHERE id = $3
RETURNING *`
	deleteTaskQuery = `DELETE FROM task WHERE id = $1`
	startTaskQuery  = `INSERT INTO time_entry (task_id, user_id, started_at, ended_at)
VALUES ($1, $2, now(), null)`
	endTaskQuery = `UPDATE time_entry SET ended_at = now()
WHERE ended_at IS NULL AND task_id = $1 AND user_id = $2`
	getActiveUserTasksQuery = `SELECT count(1) FROM task
INNER JOIN time_entry on time_entry.task_id = task.id
WHERE time_entry.ended_at IS NULL
AND time_entry.user_id = $1
AND task.project_id = (SELECT project_id FROM task WHERE id = $2)`
	getTotalTaskMembersQuery = `SELECT count(user_id) FROM task_participant WHERE task_id = $1`
	getTaskMembersQuery      = `SELECT user.* FROM "user"
INNER JOIN task_participant ON task_participant.user_id = "user".id
WHERE task_id = $1`
	addTaskMemberQuery    = `INSERT INTO task_participant (task_id, user_id) VALUES ($1, $2)`
	deleteTaskMemberQuery = `DELETE FROM task_participant WHERE task_id = $1 AND user_id = $2`
)

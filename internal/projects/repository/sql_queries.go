package repository

const (
	getProjectByIDQuery = `SELECT * FROM project WHERE id = $1`
	createProjectQuery  = `INSERT INTO project (name, description, creator_id) VALUES
	 ($1, $2, $3) RETURNING *`
	deleteProjectQuery = `DELETE FROM project WHERE id = $1`

	updateProjectQuery = `UPDATE project SET
name = COALESCE(NULLIF($1, ''), name),
description = COALESCE(NULLIF($2, ''), description),
creator_id = COALESCE(NULLIF($3, 0), creator_id)
WHERE id = $4
RETURNING *`

	isProjectMemberQuery = `SELECT FROM project_participant WHERE project_id = $1 AND user_id = $2`
	isProjectOwnerQuery  = `SELECT EXISTS (SELECT 1 FROM project WHERE id = $1 AND creator_id = $2)`

	getProjectMembersCount = `SELECT COUNT(user_id) FROM project_participant WHERE project_id = $1`
	getProjectMembers      = `SELECT * FROM "user" 
INNER JOIN project_participant ON "user".id = project_participant.user_id
WHERE project_id = $1`
	addProjectMemberQuery             = `INSERT INTO project_participant (project_id, user_id) VALUES ($1, $2)`
	removeProjectMemberQuery          = `DELETE FROM project_participant  WHERE project_id = $1 AND user_id = $2`
	getProjectMemberProductivityQuery = `SELECT task_id, SUM(EXTRACT(EPOCH FROM (time_entry.ended_at - started_at))) 
AS total_seconds FROM time_entry
INNER JOIN task ON task.id = time_entry.task_id
WHERE task.project_id = $1
  AND time_entry.user_id = $2
GROUP BY task_id
ORDER BY total_seconds DESC`
)

package repository

const (
	selectUserByIDQuery    = `SELECT * FROM "user" WHERE id = $1`
	selectUserByEmailQuery = `SELECT * FROM "user" WHERE email = $1`
	deleteUserByIDQuery    = `DELETE FROM "user" WHERE id = $1`

	createUserQuery = `INSERT INTO "user" (email, password, name, surname, patronymic, address)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`

	updateUserQuery = `UPDATE "user" SET
email = COALESCE(NULLIF($1, ''), email),
password = COALESCE(NULLIF($2, ''), password),
name = COALESCE(NULLIF($3, ''), name),
surname = COALESCE(NULLIF($4, ''), surname),
patronymic = COALESCE(NULLIF($5, ''), patronymic),
address = COALESCE(NULLIF($6, ''), address)
WHERE id = $7
RETURNING *`

	searchUsersCountQuery = `SELECT count(1)
FROM "user"
WHERE
    id >= COALESCE(NULLIF($1, 0), id) AND
    id <= COALESCE(NULLIF($2, 0), id) AND
    email ILIKE '%' || $3 || '%' AND
    name ILIKE '%' || $4 || '%' AND
    surname ILIKE '%' || $5 || '%' AND
    patronymic ILIKE '%' || $6 || '%' AND
    address ILIKE '%' || $7 || '%'`

	searchUsersQuery = `SELECT *
FROM "user"
WHERE
    id >= COALESCE(NULLIF($1, 0), id) AND
    id <= COALESCE(NULLIF($2, 0), id) AND
    email ILIKE '%' || $3 || '%' AND
    name ILIKE '%' || $4 || '%' AND
    surname ILIKE '%' || $5 || '%' AND
    patronymic ILIKE '%' || $6 || '%' AND
    address ILIKE '%' || $7 || '%'
OFFSET $8 LIMIT $9`

	getTotalUsers = `SELECT COUNT(id) FROM "user"`
)

package daos

import (
	"database/sql"
	"errors"
	"github.com/bheemeshkammak/intelops_compge_testing/intelops_compge_testing/pkg/rest/server/daos/clients/sqls"
	"github.com/bheemeshkammak/intelops_compge_testing/intelops_compge_testing/pkg/rest/server/models"
	log "github.com/sirupsen/logrus"
)

type NameDao struct {
	sqlClient *sqls.SQLiteClient
}

func migrateNames(r *sqls.SQLiteClient) error {
	query := `
	CREATE TABLE IF NOT EXISTS names(
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
        
		Firstname TEXT NOT NULL,
		Lastname TEXT NOT NULL,
		Name TEXT NOT NULL,
        CONSTRAINT id_unique_key UNIQUE (Id)
	)
	`
	_, err1 := r.DB.Exec(query)
	return err1
}

func NewNameDao() (*NameDao, error) {
	sqlClient, err := sqls.InitSqliteDB()
	if err != nil {
		return nil, err
	}
	err = migrateNames(sqlClient)
	if err != nil {
		return nil, err
	}
	return &NameDao{
		sqlClient,
	}, nil
}

func (nameDao *NameDao) CreateName(m *models.Name) (*models.Name, error) {
	insertQuery := "INSERT INTO names(Firstname, Lastname, Name)values(?, ?, ?)"
	res, err := nameDao.sqlClient.DB.Exec(insertQuery, m.Firstname, m.Lastname, m.Name)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	m.Id = id

	log.Debugf("name created")
	return m, nil
}

func (nameDao *NameDao) UpdateName(id int64, m *models.Name) (*models.Name, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}
	if id != m.Id {
		return nil, errors.New("id and payload don't match")
	}

	name, err := nameDao.GetName(id)
	if err != nil {
		return nil, err
	}
	if name == nil {
		return nil, sql.ErrNoRows
	}

	updateQuery := "UPDATE names SET Firstname = ?, Lastname = ?, Name = ? WHERE Id = ?"
	res, err := nameDao.sqlClient.DB.Exec(updateQuery, m.Firstname, m.Lastname, m.Name, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sqls.ErrUpdateFailed
	}

	log.Debugf("name updated")
	return m, nil
}

func (nameDao *NameDao) DeleteName(id int64) error {
	deleteQuery := "DELETE FROM names WHERE Id = ?"
	res, err := nameDao.sqlClient.DB.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sqls.ErrDeleteFailed
	}

	log.Debugf("name deleted")
	return nil
}

func (nameDao *NameDao) ListNames() ([]*models.Name, error) {
	selectQuery := "SELECT * FROM names"
	rows, err := nameDao.sqlClient.DB.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var names []*models.Name
	for rows.Next() {
		m := models.Name{}
		if err = rows.Scan(&m.Id, &m.Firstname, &m.Lastname, &m.Name); err != nil {
			return nil, err
		}
		names = append(names, &m)
	}
	if names == nil {
		names = []*models.Name{}
	}

	log.Debugf("name listed")
	return names, nil
}

func (nameDao *NameDao) GetName(id int64) (*models.Name, error) {
	selectQuery := "SELECT * FROM names WHERE Id = ?"
	row := nameDao.sqlClient.DB.QueryRow(selectQuery, id)
	m := models.Name{}
	if err := row.Scan(&m.Id, &m.Firstname, &m.Lastname, &m.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sqls.ErrNotExists
		}
		return nil, err
	}

	log.Debugf("name retrieved")
	return &m, nil
}

package mysql

import (
	"github.com/herrjulz/jefe/pkg/models"
	"github.com/jmoiron/sqlx"
)

const (
	getEnvsStmt    = `SELECT * FROM environments`
	updateEnvStmt  = `UPDATE environments SET claimed = ?, claimer = ?, date = UTC_TIMESTAMP() WHERE id = ?`
	unclaimAllStmt = `UPDATE environments SET claimed = false`
	createEnvStmt  = `INSERT INTO environments (name, image, about, date) VALUES (
				?,
				?,
				?,
				UTC_TIMESTAMP()
	                  )`

	deleteEnvStmt  = `DELETE FROM environments WHERE id = ?`
	hasClaimedStmt = `SELECT EXISTS(SELECT * FROM environments WHERE claimer = ? AND claimed = true)`
	isClaimerStmt  = `SELECT EXISTS(SELECT * FROM environments WHERE claimed = true AND claimer = ? AND id = ?)`
)

type Environments struct {
	DB *sqlx.DB
}

func (e *Environments) List() ([]*models.Environment, error) {
	envs := []*models.Environment{}

	err := e.DB.Select(&envs, getEnvsStmt)
	if err != nil {
		return nil, err
	}

	return envs, nil

}

func (e *Environments) Claim(id int, user string) error {
	if hasClaimed, err := e.hasClaimed(user); hasClaimed {
		return models.ErrClaimed
	} else if err != nil {
		return err
	}

	_, err := e.DB.Exec(updateEnvStmt, true, user, id)
	if err != nil {
		return err
	}

	return nil
}

func (e *Environments) Unclaim(id int, user string) error {
	if isClaimer, err := e.isClaimer(user, id); !isClaimer {
		return models.ErrUnclaimed
	} else if err != nil {
		return err
	}

	_, err := e.DB.Exec(updateEnvStmt, false, user, id)
	if err != nil {
		return err
	}

	return nil
}

func (e *Environments) UnclaimAll() error {
	_, err := e.DB.Exec(unclaimAllStmt)
	if err != nil {
		return err
	}

	return nil
}

func (e *Environments) Create(name, image, about string) error {
	_, err := e.DB.Exec(createEnvStmt, name, image, about)
	if err != nil {
		return err
	}

	return nil
}

func (e *Environments) Delete(id int) error {
	_, err := e.DB.Exec(deleteEnvStmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (e *Environments) hasClaimed(claimer string) (bool, error) {
	var exists bool
	err := e.DB.QueryRow(hasClaimedStmt, claimer).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (e *Environments) isClaimer(claimer string, id int) (bool, error) {
	var exists bool
	err := e.DB.QueryRow(isClaimerStmt, claimer, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

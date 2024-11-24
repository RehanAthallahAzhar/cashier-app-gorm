package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SessionsRepository struct {
	db *gorm.DB
}

func NewSessionsRepository(db *gorm.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	if err := u.db.Table("sessions").Create(&session).Error; err != nil {
		// return any error will rollback
		return err
	}
	return nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	results := []model.Session{}
	rows, err := u.db.Table("sessions").Select("*").Where("deleted_at is NULL").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() { // Next akan menyiapkan hasil baris berikutnya untuk dibaca dengan metode Scan.
		u.db.ScanRows(rows, &results)
	}

	err = u.db.Table("sessions").Delete(&results).Where("token = ?", tokenTarget).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) UpdateSessions(session model.Session) error {
	// UPDATE sessions SET (token = {token}, expiry = {expiry}) where username = {username}
	result := u.db.Table("sessions").Where("username = ?", session.Username).Update("token", session.Token).Update("expiry", session.Expiry)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *SessionsRepository) TokenValidity(token string) (model.Session, error) {
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, fmt.Errorf("asd")
	}

	if u.TokenExpired(session) {
		err := u.DeleteSessions(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, fmt.Errorf("Token is Expired!")
	}
	return session, nil
}

func (u *SessionsRepository) SessionAvailName(name string) (model.Session, error) {
	var data model.Session
	result := u.db.Table("sessions").Select("*").Where("username = ?", name).Scan(&data)
	if result.Error != nil {
		return model.Session{}, result.Error
	}
	if data == (model.Session{}) {
		return model.Session{}, fmt.Errorf("Session token tidak ditemukan")
	}
	return data, nil
}

func (u *SessionsRepository) SessionAvailToken(token string) (model.Session, error) {
	var data model.Session
	result := u.db.Table("sessions").Select("*").Where("token = ?", token).Where("deleted_at is null").Scan(&data)
	if result.Error != nil {
		return model.Session{}, result.Error
	}
	if data == (model.Session{}) {
		return model.Session{}, fmt.Errorf("Session token tidak ditemukan")
	}

	return data, nil
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}

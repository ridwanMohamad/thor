package sessions

import (
	"errors"
	"gorm.io/gorm"
	"thor/src/constanta"
	"thor/src/domain/tbl_session"
	"thor/src/server/database"
)

type sessions struct {
	db *database.Database
}

func NewSessionsRepository(DB *database.Database) ISessionRepository {

	return &sessions{db: DB}
}

func (s sessions) Save(data *tbl_session.Session) (resp *tbl_session.Session, err error) {
	//TODO implement me
	tx := s.db

	if err = tx.Create(data).Error; err != nil {
		return nil, constanta.DbFailedToInsertData
	}

	if err = tx.First(&resp).Where("session_id = ?", data.Pk).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (s sessions) FindByAccessToken(token string) (resp *tbl_session.Session, err error) {
	tx := s.db

	if err = tx.Where("access_token = ?", token).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, constanta.DbFailedToExecuteQuery
	}
	return
}

func (s sessions) FindByUserId(id int64) (resp *tbl_session.Session, err error) {
	tx := s.db

	if err = tx.Where("fk_user = ?", id).First(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		tx.Rollback()
		return nil, constanta.DbFailedToExecuteQuery
	}

	return
}

func (s sessions) Update(data *tbl_session.Session) (err error) {
	//TODO implement me
	tx := s.db.Begin()

	if err = tx.Model(tbl_session.Session{}).
		Where("pk = ? and fk_user = ?", data.Pk, data.FkUser).
		Select("access_token", "expired_at", "updated_at").
		Updates(tbl_session.Session{
			AccessToken: data.AccessToken,
			ExpiredAt:   data.ExpiredAt,
			UpdatedAt:   data.UpdatedAt,
		}).Error; err != nil {

		tx.Rollback()
		return constanta.DbFailedToUpdateData
	}
	tx.Commit()
	return
}

//func (s sessions) UpdateTokenAndExpired(data *tbl_session.Session) (err error) {
//	//TODO implement me
//	tx := s.db.Begin()
//
//	if err = tx.Model(tbl_session.Session{}).
//		Where("session_id = ? and user_id = ?", data.SessionId, data.UserId).
//		Updates(tbl_session.Session{
//			AccessToken: data.AccessToken,
//			ExpiredAt:   data.ExpiredAt,
//			UpdatedAt:   data.UpdatedAt,
//		}).Error; err != nil {
//		return constanta.DbFailedToUpdateData
//		tx.Rollback()
//	}
//	tx.Commit()
//	return
//}

func (s sessions) Delete(username string) (resp bool, err error) {
	//TODO implement me
	panic("implement me")
}

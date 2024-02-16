package sessions

import "thor/src/domain/tbl_session"

type ISessionRepository interface {
	Save(data *tbl_session.Session) (resp *tbl_session.Session, err error)
	FindByAccessToken(token string) (resp *tbl_session.Session, err error)
	FindByUserId(id int64) (resp *tbl_session.Session, err error)
	Update(data *tbl_session.Session) (err error)
	Delete(username string) (resp bool, err error)
}

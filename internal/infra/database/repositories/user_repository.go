package repositories

import (
	"database/sql"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/logger"
)

type user struct {
	Id                     uint64  `db:"id, omitempty"`
	Name                   string  `db:"name"`
	Email                  string  `db:"email"`
	Avatar                 *string `db:"avatar"`
	Password               string  `db:"password"`
	EmailConfirmed         bool    `db:"email_confirmed"`
	EmailConfirmationToken string  `db:"email_confirmation_token"`
}

type UserRepository interface {
	FindById(id uint64) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	FindByEmailConfirmationToken(confToken string) (domain.User, error)
	Save(user domain.User) (domain.User, error)
	UpdateUserAvatar(user domain.User) (domain.User, error)
	ConfirmUserEmail(user domain.User) error
	Delete(id uint64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur userRepository) FindById(id uint64) (domain.User, error) {
	userModel := user{}
	sqlCommand := `SELECT * FROM users WHERE id=$1`
	err := ur.db.QueryRow(sqlCommand, id).Scan(
		&userModel.Id,
		&userModel.Name,
		&userModel.Email,
		&userModel.Avatar,
		&userModel.Password,
		&userModel.EmailConfirmed,
		&userModel.EmailConfirmationToken,
	)
	if err != nil {
		logger.Logger.Error(err)
		return domain.User{}, err
	}

	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) FindByEmail(email string) (domain.User, error) {
	userModel := user{}
	sqlCommand := `SELECT * FROM users WHERE email=$1`

	err := ur.db.QueryRow(sqlCommand, email).Scan(
		&userModel.Id,
		&userModel.Name,
		&userModel.Email,
		&userModel.Avatar,
		&userModel.Password,
		&userModel.EmailConfirmed,
		&userModel.EmailConfirmationToken,
	)
	if err != nil {
		logger.Logger.Error(err)
		return domain.User{}, err
	}

	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) FindByEmailConfirmationToken(confToken string) (domain.User, error) {
	sqlCommand := `SELECT * FROM users WHERE email_confirmation_token=$1 AND email_confirmed=false`
	userModel := user{}

	err := ur.db.QueryRow(sqlCommand, confToken).Scan(
		&userModel.Id,
		&userModel.Name,
		&userModel.Email,
		&userModel.Avatar,
		&userModel.Password,
		&userModel.EmailConfirmed,
		&userModel.EmailConfirmationToken,
	)
	if err != nil {
		logger.Logger.Error(err)
		return domain.User{}, err
	}

	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) Save(user domain.User) (domain.User, error) {
	userModel := ur.domainToModel(user)
	sqlCommand := `INSERT INTO users (name, email, password, email_confirmation_token) VALUES ($1, $2, $3, $4) RETURNING id`

	err := ur.db.QueryRow(sqlCommand, userModel.Name, userModel.Email, userModel.Password, userModel.EmailConfirmationToken).Scan(&userModel.Id)
	if err != nil {
		logger.Logger.Error(err)
		return domain.User{}, err
	}
	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) UpdateUserAvatar(user domain.User) (domain.User, error) {
	userModel := ur.domainToModel(user)
	sqlCommand := `UPDATE users SET avatar=$1 WHERE id=$2`

	_, err := ur.db.Exec(sqlCommand, userModel.Avatar, userModel.Id)
	if err != nil {
		logger.Logger.Error(err)
		return domain.User{}, err
	}

	return ur.modelToDomain(userModel), nil
}

func (ur userRepository) ConfirmUserEmail(user domain.User) error {
	userModel := ur.domainToModel(user)
	sqlCommand := `UPDATE users SET email_confirmed=true WHERE id=$1 AND email_confirmation_token=$2`

	_, err := ur.db.Exec(sqlCommand, userModel.Id, userModel.EmailConfirmationToken)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	return nil
}

func (ur userRepository) Delete(id uint64) error {
	sqlCommand := `DELETE FROM users WHERE id=$1`
	_, err := ur.db.Exec(sqlCommand, id)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}
	return nil
}

func (ur userRepository) modelToDomain(u user) domain.User {
	return domain.User{
		Id:                     u.Id,
		Name:                   u.Name,
		Email:                  u.Email,
		Avatar:                 u.Avatar,
		Password:               u.Password,
		EmailConfirmed:         u.EmailConfirmed,
		EmailConfirmationToken: u.EmailConfirmationToken,
	}
}

func (ur userRepository) domainToModel(u domain.User) user {
	return user{
		Id:                     u.Id,
		Name:                   u.Name,
		Email:                  u.Email,
		Avatar:                 u.Avatar,
		Password:               u.Password,
		EmailConfirmed:         u.EmailConfirmed,
		EmailConfirmationToken: u.EmailConfirmationToken,
	}
}

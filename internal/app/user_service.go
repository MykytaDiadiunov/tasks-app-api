package app

import (
	"fmt"
	"go-rest-api/config"
	"go-rest-api/internal/domain"
	"go-rest-api/internal/infra/database/repositories"
	"go-rest-api/internal/infra/filesystem"
	"go-rest-api/internal/infra/logger"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

type UserService interface {
	FindById(id uint64) (domain.User, error)
	FindByEmail(email string) (domain.User, error)
	FindByEmailConfirmationToken(confToken string) (domain.User, error)
	Save(user domain.User) (domain.User, error)
	UpdateUserAvatar(user domain.User) (domain.User, error)
	ConfirmUserEmail(user domain.User) error
	Delete(id uint64) error
}

type userService struct {
	userRepo          repositories.UserRepository
	configuration     config.Configuration
	cloudinaryService *filesystem.CloudinaryService
}

func NewUserService(userRepository repositories.UserRepository,
	cfg config.Configuration, cloudinaryService *filesystem.CloudinaryService) UserService {

	return userService{
		userRepo:          userRepository,
		configuration:     cfg,
		cloudinaryService: cloudinaryService,
	}
}

func (u userService) FindById(id uint64) (domain.User, error) {
	user, err := u.userRepo.FindById(id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u userService) FindByEmail(email string) (domain.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (u userService) FindByEmailConfirmationToken(confToken string) (domain.User, error) {
	user, err := u.userRepo.FindByEmailConfirmationToken(confToken)
	if err != nil {
		logger.Logger.Error(err)
		return domain.User{}, err
	}

	return user, err
}

func (u userService) Save(user domain.User) (domain.User, error) {
	var err error
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return domain.User{}, err
	}

	user.EmailConfirmationToken = generateEmailConfirmationToken()

	// err = u.sendEmail(user)
	// if err != nil {
	// 	logger.Logger.Error(err)
	// 	return domain.User{}, err
	// }

	user, err = u.userRepo.Save(user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (u userService) UpdateUserAvatar(user domain.User) (domain.User, error) {
	imageFileName := "file_" + strconv.FormatInt(time.Now().UnixNano(), 32)

	imageUrl, err := u.cloudinaryService.SaveImageToCloudinary(*user.Avatar, imageFileName)
	if err != nil {
		logger.Logger.Error(err)
		return domain.User{}, err
	}

	*user.Avatar = imageUrl

	updatedUser, err := u.userRepo.UpdateUserAvatar(user)
	if err != nil {
		logger.Logger.Error(err)
		return domain.User{}, err
	}
	return updatedUser, nil
}

func (u userService) ConfirmUserEmail(user domain.User) error {
	currentUser, err := u.FindByEmailConfirmationToken(user.EmailConfirmationToken)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}

	err = u.userRepo.ConfirmUserEmail(currentUser)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}

	return nil
}

func (u userService) Delete(id uint64) error {
	deletedUser, err := u.FindById(id)
	if err != nil {
		logger.Logger.Error(err)
		return err
	}

	if deletedUser.Avatar != nil {
		err = u.cloudinaryService.DeleteImage(*deletedUser.Avatar)
		if err != nil {
			logger.Logger.Error(err)
			return err
		}
	}

	err = u.userRepo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

// func (u userService) sendEmail(user domain.User) error {
// 	emailBody := fmt.Sprintf("Your confirmation code: %s", user.EmailConfirmationToken)

// 	auth := smtp.PlainAuth(
// 		"",
// 		u.configuration.WorkGmail,
// 		u.configuration.WorkGmailPassword,
// 		u.configuration.SmtpHost,
// 	)

// 	err := smtp.SendMail(
// 		u.configuration.SmtpHost+":"+u.configuration.SmtpPort, auth, u.configuration.WorkGmail, []string{user.Email}, []byte(emailBody),
// 	)
// 	if err != nil {
// 		logger.Logger.Error(err)
// 		return err
// 	}

// 	return nil
// }

func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func generateEmailConfirmationToken() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

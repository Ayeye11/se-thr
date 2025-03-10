package repository

import (
	"github.com/Ayeye11/AuthCache/internal/common/errs"
	"github.com/Ayeye11/AuthCache/internal/common/types"
	"github.com/Ayeye11/AuthCache/internal/database/models"
	"gorm.io/gorm"
)

type userStore struct {
	db *gorm.DB
}

// Create
func (s *userStore) CreateUser(u *types.User) error {
	if !u.IsPasswordHashed() {
		return errs.ErrValidation_PasswordNotHash
	}

	model := models.UserModel{
		Email:        u.Email,
		HashPassword: u.Password,
		Firstname:    u.Firstname,
		Lastname:     u.Lastname,
		Age:          u.Age,
		RoleID:       u.Role.ID,
	}

	if err := s.db.Create(&model).Error; err != nil {
		return errs.IsErrDoX(err, gorm.ErrDuplicatedKey, errs.ErrRepoUser_DuplicatedEmail)
	}

	return nil
}

// Read
func (s *userStore) GetUserByID(id int) (*types.User, error) {
	if id < 1 {
		return nil, errs.ErrRepoUser_InvalidUserID
	}

	model := models.UserModel{}
	if err := s.db.First(&model, id).Error; err != nil {
		return nil, errs.IsErrDoX(err, gorm.ErrRecordNotFound, errs.ErrRepoUser_NotFound)
	}

	modelRole := models.AcRole{}
	if err := s.db.First(&modelRole, model.RoleID).Error; err != nil {
		return nil, errs.UnknownError(err)
	}

	user := &types.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.HashPassword,
		Firstname: model.Firstname,
		Lastname:  model.Lastname,
		Age:       model.Age,

		Role: &types.Role{
			ID:   modelRole.ID,
			Name: modelRole.Role,
		},
	}

	return user, nil
}

func (s *userStore) GetUserByEmail(email string) (*types.User, error) {
	model := models.UserModel{}
	if err := s.db.First(&model, "email = ?", email).Error; err != nil {
		return nil, errs.IsErrDoX(err, gorm.ErrRecordNotFound, errs.ErrRepoUser_NotFound)
	}

	modelRole := models.AcRole{}
	if err := s.db.First(&modelRole, model.RoleID).Error; err != nil {
		return nil, errs.UnknownError(err)
	}

	user := &types.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.HashPassword,
		Firstname: model.Firstname,
		Lastname:  model.Lastname,
		Age:       model.Age,

		Role: &types.Role{
			ID:   modelRole.ID,
			Name: modelRole.Role,
		},
	}

	return user, nil
}

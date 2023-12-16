package models

import (
	"context"
	"errors"
	"gingonic-search-server/common"
	"gingonic-search-server/component/objectstoragepvd"
	"gingonic-search-server/utils"
	"mime/multipart"
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	UserEntityName = "User"
	UserTableName  = "users"
	UserTableID    = 1
)

type User struct {
	common.BaseModel `json:",inline"`
	FirstName        string `json:"first_name" gorm:"column:first_name;type:varchar(128);not null"`
	LastName         string `json:"last_name" gorm:"column:last_name;type:varchar(128);not null"`
	Avatar           string `json:"avatar" gorm:"column:avatar;type:text;not null"`
	Phone            string `json:"phone" gorm:"column:phone;type:varchar(10);not null"`
	Email            string `json:"email" gorm:"column:email;type:varchar(256);uniqueIndex:idx_email;not null"`
	Password         string `json:"password,omitempty" gorm:"column:password;type:text;not null"`
	IsLocked         bool   `json:"is_locked" gorm:"column:is_locked;not null;default:false"`
	VerificationCode string `json:"verification_code,omitempty" gorm:"column:verification_code;type:text;not null"`
	Verified         bool   `json:"verified" gorm:"column:verified;not null;default:false"`
	Role             string `json:"role" gorm:"column:role;type:varchar(32);not null;default:'user'"`
	IsDeleted        bool   `json:"is_deleted,omitempty" gorm:"column:is_deleted;index:idx_is_deleted;not null;default:false"`
}

type UserLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UserCreate struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
}

type UserUpdate struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Password  string `json:"password,omitempty"`
}

type UserAvatar struct {
	Header *multipart.FileHeader
	File   *multipart.File
	Avatar string `json:"avatar"`
}

func (User) TableName() string { return UserTableName }

func (u *User) Mask(isAdmin bool) {
	u.FakeID = utils.EncodeUID(u.ID, UserTableID)
	if !isAdmin {
		u.Password = *new(string)
		u.VerificationCode = *new(string)
		u.IsDeleted = *new(bool)
	}
}
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u UserLogin) Validate() error {
	u.Email = strings.TrimSpace(u.Email)

	if u.Email == "" {
		return ErrEmailIsEmpty
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return ErrEmailIsInvalid
	}

	if u.Password == "" {
		return ErrPasswordIsEmpty
	}

	return nil
}

func (u UserCreate) Validate() error {
	spacesRegex := regexp.MustCompile(`\s{2,}`)

	if u.FirstName == "" {
		return ErrFirstNameIsEmpty
	}

	if spacesRegex.MatchString(u.FirstName) || strings.ContainsAny(u.FirstName, "!@#$%^&*()_+{}|:<>?0123456789") {
		return ErrFirstNameIsInvalid
	}

	if u.LastName == "" {
		return ErrLastNameIsEmpty
	}

	if spacesRegex.MatchString(u.LastName) || strings.ContainsAny(u.LastName, "!@#$%^&*()_+{}|:<>?0123456789") {
		return ErrLastNameIsInvalid
	}

	if u.Email == "" {
		return ErrEmailIsEmpty
	}

	if _, err := mail.ParseAddress(u.Email); err != nil || spacesRegex.MatchString(u.Email) {
		return ErrEmailIsInvalid
	}

	if u.Phone == "" {
		return ErrPhoneIsEmpty
	}

	if _, err := strconv.Atoi(u.Phone); err != nil || len(u.Phone) != 10 || spacesRegex.MatchString(u.Phone) {
		return ErrPhoneIsInvalid
	}

	if u.Password == "" {
		return ErrPasswordIsEmpty
	}

	if !regexp.MustCompile(`^[^ ]{8,32}$`).MatchString(u.Password) {
		return ErrPasswordIsInvalid
	}

	return nil
}

func (u UserUpdate) Validate() error {
	spacesRegex := regexp.MustCompile(`\s{2,}`)

	if u.FirstName != "" && (spacesRegex.MatchString(u.FirstName) || strings.ContainsAny(u.FirstName, "!@#$%^&*()_+{}|:<>?0123456789")) {
		return ErrFirstNameIsInvalid
	}

	if u.LastName != "" && (spacesRegex.MatchString(u.LastName) || strings.ContainsAny(u.LastName, "!@#$%^&*()_+{}|:<>?0123456789")) {
		return ErrLastNameIsInvalid
	}

	if _, err := strconv.Atoi(u.Phone); (err != nil || len(u.Phone) != 10 || spacesRegex.MatchString(u.Phone)) && u.Phone != "" {
		return ErrPhoneIsInvalid
	}

	if u.Password != "" && !regexp.MustCompile(`^[^ ]{8,32}$`).MatchString(u.Password) {
		return ErrPasswordIsInvalid
	}

	return nil
}

func (u UserAvatar) Validate() error {
	if u.Header == nil || u.File == nil {
		return ErrAvatarIsEmpty
	}

	return nil
}

func (u *UserAvatar) UploadAvatar(ctx context.Context, objectStorageProvider objectstoragepvd.ObjectStorageProvider) error {
	path, err := objectStorageProvider.PutObject(ctx, common.MinioBucketAvatars, u.File, u.Header)

	if err != nil {
		return err
	}

	u.Avatar = path

	return nil
}

var (
	ErrFirstNameIsEmpty = errors.New("first name is empty")
	ErrLastNameIsEmpty  = errors.New("last name is empty")
	ErrEmailIsEmpty     = errors.New("email is empty")
	ErrPhoneIsEmpty     = errors.New("phone is empty")
	ErrPasswordIsEmpty  = errors.New("password is empty")
	ErrAvatarIsEmpty    = errors.New("avatar is empty")

	ErrFirstNameIsInvalid = errors.New("first name is invalid")
	ErrLastNameIsInvalid  = errors.New("last name is invalid")
	ErrEmailIsInvalid     = errors.New("email is invalid")
	ErrPhoneIsInvalid     = errors.New("phone is invalid")
	ErrPasswordIsInvalid  = errors.New("password is invalid")

	ErrEmailExisted = common.NewCustomErrorResponse(
		errors.New("email has already existed"),
		"email has already existed",
		"ErrEmailExisted",
	)

	ErrUserIsLocked = common.NewCustomErrorResponse(
		errors.New("user is locked"),
		"user is locked",
		"ErrUserIsLocked",
	)

	ErrEmailOrPasswordIsIncorrect = common.NewCustomErrorResponse(
		errors.New("email or password is incorrect"),
		"email or password is incorrect",
		"ErrEmailOrPasswordIsIncorrect",
	)

	ErrCurrentUserDoesNotExist = common.NewCustomErrorResponse(
		errors.New("current user does not exist"),
		"current user does not exist",
		"ErrCurrentUserDoesNotExist",
	)
)

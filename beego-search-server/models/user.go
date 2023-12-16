package models

import (
	"beego-search-server/common"
	"beego-search-server/component/objectstoragepvd"
	"beego-search-server/utils"
	"context"
	"errors"
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
	FirstName        string `json:"first_name" orm:"column(first_name);size(128)"`
	LastName         string `json:"last_name" orm:"column(last_name);size(128)"`
	Avatar           string `json:"avatar" orm:"column(avatar)"`
	Phone            string `json:"phone" orm:"column(phone);size(10)"`
	Email            string `json:"email" orm:"column(email);size(256);unique;index(idx_email)"`
	Password         string `json:"password,omitempty" orm:"column(password)"`
	IsLocked         bool   `json:"is_locked" orm:"column(is_locked);default(false)"`
	VerificationCode string `json:"verification_code,omitempty" orm:"column(verification_code)"`
	Verified         bool   `json:"verified" orm:"column(verified);default(false)"`
	Role             string `json:"role" orm:"column(role);size(32);default(user)"`
	IsDeleted        bool   `json:"is_deleted,omitempty" orm:"column(is_deleted);index(idx_is_deleted);default(false)"`
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

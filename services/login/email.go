package login

import (
	"dankmuzikk/db"
	"dankmuzikk/entities"
	"dankmuzikk/models"
	"dankmuzikk/services/jwt"
	"dankmuzikk/services/mailer"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type EmailLoginService struct {
	accountRepo db.CRUDRepo[models.Account]
	profileRepo db.CRUDRepo[models.Profile]
	otpRepo     db.CRUDRepo[models.EmailVerificationCode]
	jwtUtil     jwt.Manager[jwt.Json]
}

func NewEmailLoginService(
	accountRepo db.CRUDRepo[models.Account],
	profileRepo db.CRUDRepo[models.Profile],
	otpRepo db.CRUDRepo[models.EmailVerificationCode],
	jwtUtil jwt.Manager[jwt.Json],
) *EmailLoginService {
	return &EmailLoginService{
		accountRepo: accountRepo,
		profileRepo: profileRepo,
		otpRepo:     otpRepo,
		jwtUtil:     jwtUtil,
	}
}

func (e *EmailLoginService) Login(user entities.LoginRequest) (string, error) {
	account, err := e.accountRepo.GetByConds("email = ? AND is_o_auth = 0", user.Email)
	if err != nil {
		return "", errors.Join(ErrAccountNotFound, err)
	}

	profile, err := e.profileRepo.GetByConds("account_id = ?", account[0].Id)
	if err != nil {
		return "", errors.Join(ErrProfileNotFound, err)
	}
	profile[0].Account = account[0]
	profile[0].AccountId = account[0].Id

	verificationToken, err := e.jwtUtil.Sign(jwt.Json{
		"name":     profile[0].Name,
		"email":    profile[0].Account.Email,
		"username": profile[0].Username,
	}, jwt.VerificationToken, time.Hour/2)
	if err != nil {
		return "", err
	}

	return verificationToken, e.sendOtp(profile[0])
}

func (e *EmailLoginService) Signup(user entities.SignupRequest) (string, error) {
	profile := models.Profile{
		Account: models.Account{
			Email: user.Email,
		},
		Name:     user.Name,
		Username: user.Email,
	}

	// creating a new account will create the account underneath it.
	err := e.profileRepo.Add(&profile)
	if errors.Is(err, db.ErrRecordExists) {
		return "", errors.Join(ErrAccountExists, err)
	}
	if err != nil {
		return "", err
	}

	verificationToken, err := e.jwtUtil.Sign(jwt.Json{
		"name":     profile.Name,
		"email":    profile.Account.Email,
		"username": profile.Username,
	}, jwt.VerificationToken, time.Hour/2)
	if err != nil {
		return "", err
	}

	return verificationToken, e.sendOtp(profile)
}

func (e *EmailLoginService) VerifyOtp(token string, otp entities.OtpRequest) (string, error) {
	tokeeeen, err := e.jwtUtil.Decode(token, jwt.VerificationToken)
	if err != nil {
		return "", err
	}

	email, emailExists := tokeeeen.Payload["email"].(string)
	// TODO: ADD THE FUCKING ERRORS SUKA
	if !emailExists {
		return "", errors.New("missing email")
	}
	name, nameExists := tokeeeen.Payload["name"].(string)
	// TODO: ADD THE FUCKING ERRORS SUKA
	if !nameExists {
		return "", errors.New("missing name")
	}
	username, usernameExists := tokeeeen.Payload["username"].(string)
	// TODO: ADD THE FUCKING ERRORS SUKA
	if !usernameExists {
		return "", errors.New("missing username")
	}

	account, err := e.accountRepo.GetByConds("email = ?", email)
	if err != nil {
		return "", err
	}

	verCodes, err := e.otpRepo.GetByConds("account_id = ?", account[0].Id)
	if err != nil {
		return "", err
	}
	verCode := verCodes[len(verCodes)-1]
	defer func() {
		_ = e.otpRepo.Delete("account_id = ?", account[0].Id)
	}()

	if verCode.CreatedAt.Add(time.Hour / 2).Before(time.Now()) {
		return "", ErrExpiredVerificationCode
	}

	err = bcrypt.CompareHashAndPassword([]byte(verCode.Code), []byte(otp.Code))
	if err != nil {
		return "", ErrInvalidVerificationCode
	}

	sessionToken, err := e.jwtUtil.Sign(jwt.Json{
		"email":    email,
		"name":     name,
		"username": username,
	}, jwt.SessionToken, time.Hour*24*30)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (e *EmailLoginService) sendOtp(profile models.Profile) error {
	otp := generateOtp()

	hashed, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = e.otpRepo.Add(&models.EmailVerificationCode{
		AccountId: profile.AccountId,
		Code:      string(hashed),
	})
	if err != nil {
		return err
	}

	err = mailer.SendOtpEmail(profile.Name, profile.Account.Email, otp)
	if err != nil {
		return err
	}

	return nil
}

func generateOtp() string {
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	n := r.Intn(1_000_000_000-100001) + 100001
	return fmt.Sprint(n)[:6]
}

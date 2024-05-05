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
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type EmailLoginService struct {
	accountRepo db.CRUDRepo[models.Account]
	profileRepo db.CRUDRepo[models.Profile]
	otpRepo     db.CRUDRepo[models.EmailVerificationCode]
	jwtUtil     jwt.Manager[any]
}

func NewEmailLoginService(
	accountRepo db.CRUDRepo[models.Account],
	profileRepo db.CRUDRepo[models.Profile],
	otpRepo db.CRUDRepo[models.EmailVerificationCode],
	jwtUtil jwt.Manager[any],
) *EmailLoginService {
	return &EmailLoginService{
		accountRepo: accountRepo,
		profileRepo: profileRepo,
		otpRepo:     otpRepo,
		jwtUtil:     jwtUtil,
	}
}

func (e *EmailLoginService) Login(user entities.LoginRequest) (string, error) {
	account, err := e.accountRepo.GetByConds("email = ?", user.Email)
	if err != nil {
		return "", err
	}

	profile, err := e.profileRepo.GetByConds("account_id = ?", account[0].Id)
	if err != nil {
		return "", err
	}
	profile[0].Account = account[0]
	profile[0].AccountId = account[0].Id

	verificationToken, err := e.jwtUtil.Sign(map[string]string{
		"name":  profile[0].Name,
		"email": profile[0].Account.Email,
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
		Username: user.Email[:strings.Index(user.Email, "@")],
	}

	// creating a new account will create the account underneeth it.
	err := e.profileRepo.Add(&profile)
	if err != nil {
		return "", err
	}

	verificationToken, err := e.jwtUtil.Sign(map[string]string{
		"name":  profile.Name,
		"email": profile.Account.Email,
	}, jwt.VerificationToken, time.Hour/2)
	if err != nil {
		return "", err
	}

	return verificationToken, e.sendOtp(profile)
}

func (e *EmailLoginService) VerifyOtp(token string, otp entities.OtpRequest) (string, error) {
	user, err := e.jwtUtil.Decode(token, jwt.VerificationToken)
	if err != nil {
		return "", err
	}

	mappedUser := user.Payload.(map[string]any)
	email, emailExists := mappedUser["email"].(string)
	// TODO: ADD THE FUCKING ERRORS SUKA
	if !emailExists {
		return "", errors.New("missing email")
	}
	name, nameExists := mappedUser["name"].(string)
	// TODO: ADD THE FUCKING ERRORS SUKA
	if !nameExists {
		return "", errors.New("missing name")
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
		_ = e.otpRepo.Delete("id = ?", verCode.Id)
	}()

	err = bcrypt.CompareHashAndPassword([]byte(verCode.Code), []byte(otp.Code))
	if err != nil {
		return "", err
	}

	sessionToken, err := e.jwtUtil.Sign(map[string]string{
		"email": email,
		"name":  name,
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

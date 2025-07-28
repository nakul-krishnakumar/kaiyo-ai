package auth

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode"

)

const (
    MaxNameLength     = 50
    MaxEmailLength    = 254
    MinPasswordLength = 8
    MaxPasswordLength = 128
)

func (ve ValidationErrors) HasErrors() bool {
    return len(ve) > 0
}

func NewValidator() Validator {
    return &validator{}
}

// Validate sign-in request
func (v *validator) ValidateSignInRequest(req SignInRequest) ValidationErrors {
    errors := make(ValidationErrors)
    
    if err := validateName(req.FirstName, "first_name"); err != nil {
        errors["first_name"] = err.Error()
    }
    
    if err := validateName(req.LastName, "last_name"); err != nil {
        errors["last_name"] = err.Error()
    }
    
    if err := validateEmail(req.Email); err != nil {
        errors["email"] = err.Error()
    }
    
    if err := validatePassword(req.Password); err != nil {
        errors["password"] = err.Error()
    }
    
    return errors	
}

// Validate login request
func (v *validator) ValidateLoginRequest(req LogInRequest) ValidationErrors {
    errors := make(ValidationErrors)
    
    if err := validateEmail(req.Email); err != nil {
        errors["email"] = err.Error()
    }
    
    if req.Password == "" {
        errors["password"] = "password is required"
    }
    
    return errors
}

// Private validation functions
func validateName(name, fieldName string) error {
    name = strings.TrimSpace(name)
    
    if name == "" {
        return fmt.Errorf("%s is required", fieldName)
    }
    
    if len(name) < 1 {
        return fmt.Errorf("%s must be at least 1 character", fieldName)
    }
    
    if len(name) > MaxNameLength {
        return fmt.Errorf("%s must be no more than %d characters", fieldName, MaxNameLength)
    }
    
    // Allow letters, spaces, hyphens, apostrophes, periods
    for _, r := range name {
        if !unicode.IsLetter(r) && r != ' ' && r != '-' && r != '\'' && r != '.' {
            return fmt.Errorf("%s contains invalid characters", fieldName)
        }
    }
    
    return nil
}

func validateEmail(email string) error {
    email = strings.TrimSpace(email)
    
    if email == "" {
        return errors.New("email is required")
    }
    
    if len(email) > MaxEmailLength {
        return fmt.Errorf("email must be no more than %d characters", MaxEmailLength)
    }
    
    if _, err := mail.ParseAddress(email); err != nil {
        return errors.New("invalid email format")
    }
    
    return nil
}

func validatePassword(password string) error {
    if password == "" {
        return errors.New("password is required")
    }
    
    if len(password) < MinPasswordLength {
        return fmt.Errorf("password must be at least %d characters long", MinPasswordLength)
    }
    
    if len(password) > MaxPasswordLength {
        return fmt.Errorf("password must be no more than %d characters long", MaxPasswordLength)
    }
    
    var hasUpper, hasLower, hasNumber, hasSpecial bool
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    var missing []string
    if !hasUpper {
        missing = append(missing, "uppercase letter")
    }
    if !hasLower {
        missing = append(missing, "lowercase letter")
    }
    if !hasNumber {
        missing = append(missing, "number")
    }
    if !hasSpecial {
        missing = append(missing, "special character")
    }
    
    if len(missing) > 0 {
        return fmt.Errorf("password must contain: %s", strings.Join(missing, ", "))
    }
    
    return nil
}
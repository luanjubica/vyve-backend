package repository

import "errors"

// Common repository errors
var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrLastAuthMethod     = errors.New("cannot remove last authentication method")
	
	// Person errors
	ErrPersonNotFound = errors.New("person not found")
	ErrPersonExists   = errors.New("person already exists")
	
	// Interaction errors
	ErrInteractionNotFound = errors.New("interaction not found")
	
	// Reflection errors
	ErrReflectionNotFound       = errors.New("reflection not found")
	ErrReflectionAlreadyExists  = errors.New("reflection already exists for today")
	
	// Nudge errors
	ErrNudgeNotFound = errors.New("nudge not found")
	ErrNudgeExpired  = errors.New("nudge has expired")
	
	// Token errors
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token has expired")
	ErrTokenRevoked  = errors.New("token has been revoked")
	ErrTokenInvalid  = errors.New("token is invalid")
	
	// Consent errors
	ErrConsentNotFound = errors.New("consent not found")
	ErrConsentRequired = errors.New("consent is required")
	
	// Export errors
	ErrExportNotFound  = errors.New("export not found")
	ErrExportPending   = errors.New("export already pending")
	ErrExportExpired   = errors.New("export has expired")
	
	// General errors
	ErrNotFound         = errors.New("record not found")
	ErrAlreadyExists    = errors.New("record already exists")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrDatabaseError    = errors.New("database error")
	ErrTransactionError = errors.New("transaction error")
)

// IsNotFound checks if error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrUserNotFound) ||
		errors.Is(err, ErrPersonNotFound) ||
		errors.Is(err, ErrInteractionNotFound) ||
		errors.Is(err, ErrReflectionNotFound) ||
		errors.Is(err, ErrNudgeNotFound) ||
		errors.Is(err, ErrTokenNotFound) ||
		errors.Is(err, ErrConsentNotFound) ||
		errors.Is(err, ErrExportNotFound)
}

// IsAlreadyExists checks if error is an already exists error
func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists) ||
		errors.Is(err, ErrUserAlreadyExists) ||
		errors.Is(err, ErrPersonExists) ||
		errors.Is(err, ErrReflectionAlreadyExists) ||
		errors.Is(err, ErrExportPending)
}

// IsUnauthorized checks if error is an unauthorized error
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrInvalidCredentials) ||
		errors.Is(err, ErrTokenInvalid) ||
		errors.Is(err, ErrTokenExpired) ||
		errors.Is(err, ErrTokenRevoked)
}

// IsForbidden checks if error is a forbidden error
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden) ||
		errors.Is(err, ErrEmailNotVerified) ||
		errors.Is(err, ErrConsentRequired)
}
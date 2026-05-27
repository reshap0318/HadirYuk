package services

import (
	"context"
	"fmt"
	"time"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthValidateToken validates a JWT token and returns the claims.
func (s *Services) AuthValidateToken(tokenString string) (*helpers.JWTClaims, error) {
	claims, err := helpers.ValidateToken(tokenString, s.JWKSManager.GetPublicKey())
	if err != nil {
		return nil, helpers.ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, helpers.ErrExpiredToken
	}

	// Try to get cached session from Redis with fallback to DB
	if s.RedisClient.IsCacheAvailable() {
		sessionKey := fmt.Sprintf("session:%d", claims.UserID)
		var cachedUserDTO dtos.UserDTO
		if err := s.RedisClient.GetJSON(sessionKey, &cachedUserDTO); err == nil {
			// Cache hit - user data is valid, return claims
			return claims, nil
		} else {
			// Cache miss or error, fallback to DB validation
			s.Logger.LogWarn("AuthValidateToken", "Cache miss/error for session:%d, falling back to DB: %v", claims.UserID, err)
		}
	}

	// Fallback: validate user from database
	_, err = s.repo.User.FindByID(s.repo.User.DB, claims.UserID)
	if err != nil {
		s.Logger.LogStep("AuthValidateToken", "User not found in DB: %d", claims.UserID)
		return nil, helpers.ErrInvalidCredential
	}

	return claims, nil
}

// AuthLogin authenticates a user and returns tokens.
func (s *Services) AuthLogin(ctx context.Context, email, password string) (*dtos.LoginResponse, error) {
	s.Logger.LogStart("AuthLogin", "User login attempt: %s", email)

	user, err := s.repo.User.FindByEmail(email)
	if err != nil {
		s.Logger.LogStep("AuthLogin", "User not found: %s", email)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - user not found")
		return nil, helpers.ErrInvalidCredential
	}
	s.Logger.LogStep("AuthLogin", "User found: %s", email)

	if !s.checkPassword(password, user.Password) {
		s.Logger.LogStep("AuthLogin", "Invalid password for user: %s", email)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - invalid password")
		return nil, helpers.ErrInvalidCredential
	}
	s.Logger.LogStep("AuthLogin", "Password validated successfully")

	token, err := s.generateTokenWithClaims(user)
	if err != nil {
		s.Logger.LogError("AuthLogin", "Failed to generate token: %v", err)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - token generation error")
		return nil, err
	}
	s.Logger.LogStep("AuthLogin", "Access token generated")

	refreshToken, err := s.generateRefreshTokenWithClaims(user)
	if err != nil {
		s.Logger.LogError("AuthLogin", "Failed to generate refresh token: %v", err)
		s.Logger.LogEndWithError("AuthLogin", "Login failed - refresh token generation error")
		return nil, err
	}
	s.Logger.LogStep("AuthLogin", "Refresh token generated")

	// Reload user with roles for response
	userWithRoles, err := s.repo.User.FindByID(s.repo.User.DB, user.ID, "Roles.Permissions")
	if err != nil {
		s.Logger.LogWarn("AuthLogin", "Failed to load user roles: %v", err)
		userWithRoles = user
	}

	// Cache UserDTO to Redis with fallback
	if s.RedisClient.IsCacheAvailable() {
		userDTO := dtos.ToUserDTO(userWithRoles)
		sessionKey := fmt.Sprintf("session:%d", userWithRoles.ID)
		if err := s.RedisClient.SetJSON(sessionKey, userDTO, s.cfg.Expiration); err != nil {
			s.Logger.LogWarn("AuthLogin", "Failed to cache session to Redis: %v", err)
		} else {
			s.Logger.LogStep("AuthLogin", "Session cached to Redis")
		}
	}

	s.Logger.LogEnd("AuthLogin", "Login successful for user: %s", email)
	return &dtos.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         dtos.ToUserDTO(userWithRoles),
	}, nil
}

// AuthRefreshToken refreshes the access token using a refresh token.
func (s *Services) AuthRefreshToken(ctx context.Context, refreshToken string) (*dtos.LoginResponse, error) {
	s.Logger.LogStart("AuthRefreshToken", "Token refresh attempt")

	claims, err := s.AuthValidateToken(refreshToken)
	if err != nil {
		s.Logger.LogStep("AuthRefreshToken", "Invalid refresh token: %v", err)
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - invalid token")
		return nil, err
	}
	s.Logger.LogStep("AuthRefreshToken", "Refresh token validated")

	user, err := s.repo.User.FindByID(s.repo.User.DB, claims.UserID)
	if err != nil {
		s.Logger.LogStep("AuthRefreshToken", "User not found: %v", err)
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - user not found")
		return nil, helpers.ErrInvalidCredential
	}
	s.Logger.LogStep("AuthRefreshToken", "User found: %s", user.Email)

	token, err := s.generateTokenWithClaims(user)
	if err != nil {
		s.Logger.LogError("AuthRefreshToken", "Failed to generate token: %v", err)
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - token generation error")
		return nil, err
	}
	s.Logger.LogStep("AuthRefreshToken", "Access token regenerated")

	newRefreshToken, err := s.generateRefreshTokenWithClaims(user)
	if err != nil {
		s.Logger.LogError("AuthRefreshToken", "Failed to generate refresh token: %v", err)
		s.Logger.LogEndWithError("AuthRefreshToken", "Token refresh failed - refresh token generation error")
		return nil, err
	}
	s.Logger.LogStep("AuthRefreshToken", "Refresh token regenerated")

	s.Logger.LogEnd("AuthRefreshToken", "Token refreshed successfully for user: %s", user.Email)
	return &dtos.LoginResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		User:         dtos.ToUserDTO(user),
	}, nil
}

// AuthLogout handles user logout (client-side token clear).
func (s *Services) AuthLogout(ctx context.Context) error {
	// For stateless JWT, logout is handled client-side
	return nil
}

// Helper functions

func (s *Services) checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// getUserRolesAndPermissions fetches role names and permission names for a user
func (s *Services) getUserRolesAndPermissions(userID uint) (roles []string, permissions []string) {
	user, err := s.repo.User.FindByID(s.repo.User.DB, userID, "Roles.Permissions")
	if err != nil {
		return []string{}, []string{}
	}

	roleSet := make(map[string]bool)
	permSet := make(map[string]bool)

	for _, r := range user.Roles {
		roleSet[r.Name] = true
		for _, p := range r.Permissions {
			permSet[p.Name] = true
		}
	}

	for name := range roleSet {
		roles = append(roles, name)
	}
	for name := range permSet {
		permissions = append(permissions, name)
	}

	if roles == nil {
		roles = []string{}
	}
	if permissions == nil {
		permissions = []string{}
	}

	return roles, permissions
}

func (s *Services) generateTokenWithClaims(user *models.User) (string, error) {
	return helpers.GenerateToken(
		user.ID,
		user.Email,
		s.JWKSManager.GetPrivateKey(),
		s.JWKSManager.GetKeyID(),
		helpers.GetEnvInt("JWT_EXPIRATION", 24),
	)
}

func (s *Services) generateRefreshTokenWithClaims(user *models.User) (string, error) {
	return helpers.GenerateRefreshToken(
		user.ID,
		user.Email,
		s.JWKSManager.GetPrivateKey(),
		s.JWKSManager.GetKeyID(),
		helpers.GetEnvInt("JWT_REFRESH_EXPIRATION", 168),
	)
}

// AuthForgetPassword generates a reset token and sends it via email.
func (s *Services) AuthForgetPassword(ctx context.Context, email string) error {
	s.Logger.LogStart("AuthForgetPassword", "Reset password request for: %s", email)

	user, err := s.repo.User.FindByEmail(email)
	if err != nil {
		s.Logger.LogStep("AuthForgetPassword", "User not found: %s", email)
		s.Logger.LogEndWithError("AuthForgetPassword", "Reset password failed - user not found")
		return helpers.ErrNotFound
	}

	s.Logger.LogStep("AuthForgetPassword", "User found: %s", email)

	frontendURL := helpers.GetEnv("APP_FE_URL", "http://localhost:3000")

	token, err := helpers.GenerateRandomString(32)
	if err != nil {
		s.Logger.LogError("AuthForgetPassword", "Failed to generate token: %v", err)
		s.Logger.LogEndWithError("AuthForgetPassword", "Reset password failed - token generation error")
		return err
	}

	s.Logger.LogStep("AuthForgetPassword", "Reset token generated")

	hashedToken, err := helpers.HashString(token)
	if err != nil {
		s.Logger.LogError("AuthForgetPassword", "Failed to hash token: %v", err)
		s.Logger.LogEndWithError("AuthForgetPassword", "Reset password failed - token hashing error")
		return err
	}

	s.Logger.LogStep("AuthForgetPassword", "Token hashed successfully")

	expiresAt := time.Now().Add(1 * time.Hour)

	passwordReset := &models.PasswordReset{
		Email:     user.Email,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
		Used:      false,
	}

	if _, err := s.repo.PasswordReset.Create(nil, passwordReset); err != nil {
		s.Logger.LogError("AuthForgetPassword", "Failed to save reset token: %v", err)
		s.Logger.LogEndWithError("AuthForgetPassword", "Reset password failed - database error")
		return err
	}

	s.Logger.LogStep("AuthForgetPassword", "Reset token saved to database")

	resetURL := frontendURL + "/reset-password?token=" + token

	// Async email sending dengan logger
	go func() {
		if err := s.EmailClient.SendResetPasswordEmail(user.Email, token, resetURL); err != nil {
			s.Logger.LogError("AuthForgetPassword", "Failed to send reset email to %s: %v", user.Email, err)
		} else {
			s.Logger.LogStep("AuthForgetPassword", "Reset email sent successfully to %s", user.Email)
		}
	}()

	s.Logger.LogEnd("AuthForgetPassword", "Reset password request processed successfully")

	return nil
}

// AuthResetPassword validates token and resets user password.
func (s *Services) AuthResetPassword(ctx context.Context, token, newPassword string) error {
	// Hash token to find in database
	hashedToken, err := helpers.HashString(token)
	if err != nil {
		return helpers.ErrTokenInvalid
	}

	_, err = s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		// Find reset record by token (hashed)
		reset, err := s.repo.PasswordReset.FindByToken(hashedToken)
		if err != nil {
			return nil, helpers.ErrTokenInvalid
		}

		// Check if token is expired
		if reset.ExpiresAt.Before(time.Now()) {
			return nil, helpers.ErrTokenExpired
		}

		// Check if token already used
		if reset.Used {
			return nil, helpers.ErrTokenUsed
		}

		// Find user by email
		if _, err := s.repo.User.FindByEmail(reset.Email); err != nil {
			return nil, helpers.ErrNotFound
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		// Update user password using generic Update
		if _, err := s.repo.User.Update(tx, &models.User{Email: reset.Email}, &models.User{Password: string(hashedPassword)}); err != nil {
			return nil, err
		}

		// Invalidate the token using generic Update
		if _, err := s.repo.PasswordReset.Update(tx, &models.PasswordReset{Token: reset.Token}, &models.PasswordReset{Used: true}); err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}

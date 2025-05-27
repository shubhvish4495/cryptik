package gootp

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/shubhvish4495/cryptik/pkg/cache"
)

var (
	ErrInvalidOTP = errors.New("invalid OTP provided")
)

// OTPService defines the interface for OTP (One-Time Password) operations.
// It provides methods for generating and validating OTPs using a secret key.
type OTPService interface {
	GenerateOTP(secret string) (string, error)
	ValidateOTP(secret, otp string) (bool, error)
}

// otpServiceInstance is the concrete implementation of the OTPService interface.
// It handles OTP generation and validation using a cache service to store and verify OTPs.
// Fields:
//   - CacheService: A cache implementation used to store generated OTPs temporarily
//   - Length: The length of generated OTPs (e.g., 6 for a 6-digit OTP)
type otpServiceInstance struct {
	CacheService cache.Cache
	Length       int
}

// GoOTPServiceConfig defines the configuration options for the OTP service.
// Fields:
//   - Cache: The cache implementation to use for storing OTPs
//   - Length: The desired length of generated OTPs (e.g., 6 for 6-digit OTPs)
type GoOTPServiceConfig struct {
	Cache  cache.Cache
	Length int
}

// GenerateOTP creates a cryptographically secure random OTP (One-Time Password) for a given key.
// It generates a random number with the specified length (e.g., 6 digits would be between 100000 and 999999).
// The generated OTP is stored in the cache with the provided key and expires after 10 minutes.
// Parameters:
//   - key: The unique identifier used to store and later validate the OTP
//
// Returns:
//   - string: The generated OTP
//   - error: An error if OTP generation or cache storage fails
func (o otpServiceInstance) GenerateOTP(key string) (string, error) {
	// Calculate min and max values for the desired length
	// For 6 digits: min=100000, max=999999
	min := int64(1)
	for i := 1; i < o.Length; i++ {
		min *= 10
	}
	max := min*10 - 1

	// Generate cryptographically secure random number in range
	diff := big.NewInt(max - min + 1)
	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return "", fmt.Errorf("failed to generate random number: %v", err)
	}

	result := n.Int64() + min
	otp := fmt.Sprintf("%0*d", o.Length, result)

	//store it in cache with key sepcified so that it can be validated later
	if err := o.CacheService.Set(key, otp, time.Now().Add(10*time.Minute).Unix()); err != nil {
		return "", fmt.Errorf("failed to store OTP in cache: %v", err)
	}

	return fmt.Sprintf("%0*d", o.Length, result), nil
}

// ValidateOTP verifies if the provided OTP matches the one stored in cache for the given key.
// The OTP is only valid if:
//   - It is not empty and matches the configured length
//   - A corresponding entry exists in the cache for the given key
//   - The cached OTP matches the provided OTP
//
// After successful validation, the OTP is deleted from cache to prevent reuse.
//
// Parameters:
//   - key: The unique identifier used when the OTP was generated
//   - otp: The OTP string to validate
//
// Returns:
//   - bool: true if the OTP is valid, false otherwise
//   - error: ErrInvalidOTP if OTP format is invalid, or other errors explaining validation failure
func (o otpServiceInstance) ValidateOTP(key, otp string) (bool, error) {
	if otp == "" || len(otp) != o.Length {
		return false, ErrInvalidOTP
	}

	cachedOTP, exists := o.CacheService.Get(key)
	if !exists {
		return false, fmt.Errorf("OTP not found in cache for secret: %s", key)
	}

	if cachedOTP == nil || cachedOTP.(string) != otp {
		return false, fmt.Errorf("OTP does not match for secret: %s", key)
	}
	// If OTP matches, delete it from cache to prevent reuse
	o.CacheService.Delete(key)
	return true, nil
}

// NewService creates and returns a new OTP service instance with the provided configuration.
// If no cache service is specified in the config, it uses a default cache implementation.
// If no length is specified (or if length < 1), it defaults to 6 digits.
//
// Parameters:
//   - conf: GoOTPServiceConfig containing the cache service and desired OTP length
//
// Returns:
//   - OTPService: An interface implementation for OTP operations
//   - error: Currently always returns nil, but maintained for future error handling
func NewService(conf GoOTPServiceConfig) (OTPService, error) {
	//if no cache service is provided, we will use default cache implementation
	if conf.Cache == nil {
		conf.Cache = cache.GetCache()
	}

	if conf.Length < 1 {
		conf.Length = 6 // Default OTP length
	}

	// Initialization logic can be added here if needed
	return otpServiceInstance{
		conf.Cache,
		conf.Length,
	}, nil
}

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// GenerateRandomBytes generates random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

// IsValidUsername validates username format
func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	return usernameRegex.MatchString(username)
}

// IsValidPassword validates password strength
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	
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
	
	return hasUpper && hasLower && hasNumber && hasSpecial
}

// SanitizeString removes potentially harmful characters
func SanitizeString(input string) string {
	// Remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	input = htmlRegex.ReplaceAllString(input, "")
	
	// Remove excessive whitespace
	input = strings.TrimSpace(input)
	spaceRegex := regexp.MustCompile(`\s+`)
	input = spaceRegex.ReplaceAllString(input, " ")
	
	return input
}

// TruncateString truncates string to specified length
func TruncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length-3] + "..."
}

// SlugifyString creates a URL-friendly slug
func SlugifyString(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	
	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	
	// Remove non-alphanumeric characters
	reg := regexp.MustCompile(`[^a-z0-9\-]`)
	s = reg.ReplaceAllString(s, "")
	
	// Remove multiple hyphens
	reg = regexp.MustCompile(`\-+`)
	s = reg.ReplaceAllString(s, "-")
	
	// Trim hyphens from start and end
	s = strings.Trim(s, "-")
	
	return s
}

// ParseDuration parses duration string (e.g., "1h", "30m", "7d")
func ParseDuration(duration string) (time.Duration, error) {
	// Handle days specially
	if strings.HasSuffix(duration, "d") {
		days := strings.TrimSuffix(duration, "d")
		var d int
		_, err := fmt.Sscanf(days, "%d", &d)
		if err != nil {
			return 0, err
		}
		return time.Duration(d) * 24 * time.Hour, nil
	}
	
	return time.ParseDuration(duration)
}

// FormatDuration formats duration in human-readable format
func FormatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	
	parts := []string{}
	
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d day(s)", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hour(s)", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d minute(s)", minutes))
	}
	
	if len(parts) == 0 {
		return "0 minutes"
	}
	
	return strings.Join(parts, " ")
}

// CalculateHealthScore calculates relationship health score
func CalculateHealthScore(interactions []map[string]interface{}) float64 {
	if len(interactions) == 0 {
		return 50.0
	}
	
	var totalScore float64
	weights := map[string]float64{
		"energizing": 100.0,
		"neutral":    50.0,
		"draining":   0.0,
	}
	
	// Recent interactions have more weight
	for i, interaction := range interactions {
		energyImpact := interaction["energy_impact"].(string)
		weight := 1.0 / float64(i+1) // Decay factor
		totalScore += weights[energyImpact] * weight
	}
	
	// Normalize score
	maxPossibleScore := 0.0
	for i := range interactions {
		weight := 1.0 / float64(i+1)
		maxPossibleScore += 100.0 * weight
	}
	
	if maxPossibleScore > 0 {
		return math.Round((totalScore/maxPossibleScore)*100*100) / 100
	}
	
	return 50.0
}

// GetTimeAgo returns human-readable time ago string
func GetTimeAgo(t time.Time) string {
	duration := time.Since(t)
	
	if duration.Seconds() < 60 {
		return "just now"
	}
	
	if duration.Minutes() < 60 {
		m := int(duration.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	}
	
	if duration.Hours() < 24 {
		h := int(duration.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	}
	
	days := int(duration.Hours() / 24)
	if days == 1 {
		return "yesterday"
	}
	if days < 7 {
		return fmt.Sprintf("%d days ago", days)
	}
	
	weeks := days / 7
	if weeks == 1 {
		return "1 week ago"
	}
	if weeks < 4 {
		return fmt.Sprintf("%d weeks ago", weeks)
	}
	
	months := days / 30
	if months == 1 {
		return "1 month ago"
	}
	if months < 12 {
		return fmt.Sprintf("%d months ago", months)
	}
	
	years := days / 365
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

// GetDateRange returns start and end dates for a period
func GetDateRange(period string) (time.Time, time.Time) {
	now := time.Now()
	var start time.Time
	
	switch period {
	case "today":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return start, now
	case "week":
		start = now.AddDate(0, 0, -7)
		return start, now
	case "month":
		start = now.AddDate(0, -1, 0)
		return start, now
	case "quarter":
		start = now.AddDate(0, -3, 0)
		return start, now
	case "year":
		start = now.AddDate(-1, 0, 0)
		return start, now
	default:
		return now.AddDate(0, 0, -7), now
	}
}

// RemoveDuplicates removes duplicate strings from slice
func RemoveDuplicates(strings []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	
	for _, str := range strings {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}
	
	return result
}

// Contains checks if slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// StructToMap converts struct to map
func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	
	if v.Kind() != reflect.Struct {
		return result
	}
	
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		
		// Use json tag if available
		tag := field.Tag.Get("json")
		if tag == "-" {
			continue
		}
		
		name := field.Name
		if tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" {
				name = parts[0]
			}
		}
		
		// Skip zero values if omitempty
		if strings.Contains(tag, "omitempty") && value.IsZero() {
			continue
		}
		
		result[name] = value.Interface()
	}
	
	return result
}

// Paginate returns paginated slice
func Paginate(slice interface{}, page, limit int) interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return slice
	}
	
	length := v.Len()
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	
	start := (page - 1) * limit
	end := start + limit
	
	if start >= length {
		return reflect.MakeSlice(v.Type(), 0, 0).Interface()
	}
	
	if end > length {
		end = length
	}
	
	return v.Slice(start, end).Interface()
}

// GetEnvOrDefault gets environment variable or returns default
func GetEnvOrDefault(key, defaultValue string) string {
	if value := strings.TrimSpace(key); value != "" {
		return value
	}
	return defaultValue
}

// Pointer returns a pointer to the value
func Pointer[T any](v T) *T {
	return &v
}

// Deref dereferences a pointer or returns default value
func Deref[T any](ptr *T, defaultValue T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// Min returns the minimum of two values
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two values
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clamp clamps a value between min and max
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
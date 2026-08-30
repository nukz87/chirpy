package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT_ValidToken(t *testing.T) {
	secret := "my-secret-key-12345"
	userID := uuid.New()
	expiresIn := time.Hour

	// Tạo token hợp lệ
	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Xác thực token
	gotUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	// Kiểm tra UserID trả về có khớp với ban đầu không
	if gotUserID != userID {
		t.Errorf("expected userID %v, got %v", userID, gotUserID)
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	secret := "my-secret-key-12345"
	userID := uuid.New()
	// Đặt thời gian hết hạn là số âm để token hết hạn ngay lập tức
	expiresIn := -time.Second

	tokenString, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Xác thực token đã hết hạn (kỳ vọng phải có lỗi)
	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Errorf("expected error for expired token, but got nil")
	}
}

func TestJWT_WrongSecret(t *testing.T) {
	correctSecret := "correct-secret-key"
	wrongSecret := "wrong-secret-key"
	userID := uuid.New()
	expiresIn := time.Hour

	tokenString, err := MakeJWT(userID, correctSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT returned error: %v", err)
	}

	// Xác thực bằng sai khóa bí mật (kỳ vọng phải có lỗi signature)
	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Errorf("expected error for invalid secret, but got nil")
	}
}

func TestGetBearerToken_ValidHeader(t *testing.T) {
	headers := http.Header{}
	expectedToken := "my-jwt-token-string-12345"
	headers.Set("Authorization", "Bearer "+expectedToken)

	gotToken, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotToken != expectedToken {
		t.Errorf("expected token %q, got %q", expectedToken, gotToken)
	}
}

func TestGetBearerToken_NoHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Errorf("expected error when header is missing, got nil")
	}
}

func TestGetBearerToken_MalformedHeader(t *testing.T) {
	testCases := []struct {
		name   string
		header string
	}{
		{
			name:   "Missing Bearer prefix",
			header: "Basic my-token-string",
		},
		{
			name:   "Missing token value",
			header: "Bearer",
		},
		{
			name:   "No prefix, raw token only",
			header: "some-raw-token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("Authorization", tc.header)

			_, err := GetBearerToken(headers)
			if err == nil {
				t.Errorf("expected error for header %q, got nil", tc.header)
			}
		})
	}
}

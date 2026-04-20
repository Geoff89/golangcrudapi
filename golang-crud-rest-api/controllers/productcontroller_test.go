package controllers

import (
	"errors"
	"golang-crud-rest-api/entities"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestParseProductID(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantID    uint
		wantError bool
	}{
		{name: "valid", path: "/api/products/12", wantID: 12, wantError: false},
		{name: "invalid format", path: "/api/products/abc", wantID: 0, wantError: true},
		{name: "zero", path: "/api/products/0", wantID: 0, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			req = mux.SetURLVars(req, map[string]string{"id": strings.TrimPrefix(tt.path, "/api/products/")})

			id, err := parseProductID(req)
			if tt.wantError && err == nil {
				t.Fatalf("expected an error but got nil")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("expected no error but got %v", err)
			}
			if id != tt.wantID {
				t.Fatalf("expected id %d, got %d", tt.wantID, id)
			}
		})
	}
}

func TestValidateProductInput(t *testing.T) {
	tests := []struct {
		name      string
		product   entities.Product
		wantError error
	}{
		{
			name:      "valid",
			product:   entities.Product{Name: "Mouse", Price: 10},
			wantError: nil,
		},
		{
			name:      "missing name",
			product:   entities.Product{Name: "   ", Price: 10},
			wantError: errors.New("name is required"),
		},
		{
			name:      "negative price",
			product:   entities.Product{Name: "Mouse", Price: -1},
			wantError: errors.New("price must be greater than or equal to 0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProductInput(tt.product)
			if tt.wantError == nil && err != nil {
				t.Fatalf("expected no error but got %v", err)
			}
			if tt.wantError != nil {
				if err == nil {
					t.Fatalf("expected error %q but got nil", tt.wantError.Error())
				}
				if err.Error() != tt.wantError.Error() {
					t.Fatalf("expected error %q but got %q", tt.wantError.Error(), err.Error())
				}
			}
		})
	}
}

func TestRespondJSONSetsStatusAndContentType(t *testing.T) {
	rec := httptest.NewRecorder()
	respondJSON(rec, http.StatusCreated, map[string]string{"message": "ok"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type application/json, got %s", got)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "message") {
		t.Fatalf("expected response body to contain message field, got %s", body)
	}
}

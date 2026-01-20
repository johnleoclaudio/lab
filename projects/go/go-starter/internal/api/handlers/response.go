package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONAPIData represents a single resource in JSON:API format
type JSONAPIData struct {
	Type       string      `json:"type"`
	ID         string      `json:"id"`
	Attributes interface{} `json:"attributes"`
}

// JSONAPIResponse represents a successful JSON:API response
type JSONAPIResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

// JSONAPIErrorSource represents the source of an error
type JSONAPIErrorSource struct {
	Pointer string `json:"pointer,omitempty"`
}

// JSONAPIError represents a single error in JSON:API format
type JSONAPIError struct {
	Status string                 `json:"status"`
	Code   string                 `json:"code"`
	Title  string                 `json:"title"`
	Detail string                 `json:"detail"`
	Source *JSONAPIErrorSource    `json:"source,omitempty"`
	Meta   map[string]interface{} `json:"meta,omitempty"`
}

// JSONAPIErrorResponse represents an error response in JSON:API format
type JSONAPIErrorResponse struct {
	Errors []JSONAPIError `json:"errors"`
}

// respondJSON writes a JSON:API success response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, log it but we've already written the header
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// respondError writes a JSON:API error response
func respondError(w http.ResponseWriter, reqID string, status int, code, detail string) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)

	response := JSONAPIErrorResponse{
		Errors: []JSONAPIError{
			{
				Status: fmt.Sprintf("%d", status),
				Code:   code,
				Title:  http.StatusText(status),
				Detail: detail,
				Meta: map[string]interface{}{
					"request_id": reqID,
				},
			},
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, there's not much we can do at this point
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

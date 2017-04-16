package scytale

const (
	ErrorCodeOK               = 0
	ErrorCodeInternal         = 1
	ErrorCodeSize             = 2
	ErrorCodeThrottle         = 3
	ErrorCodeFormat           = 4
	ErrorCodeExtensionTooLong = 5
	ErrorCodePermissionDenied = 6
)

// UploadResponse represents the structure of an upload response.
type UploadResponse struct {
	ErrorCode int    `json:"error_code"`
	Location  string `json:"location"`
}

// UploadRequest represents the structure of an upload request.
type UploadRequest struct {
	IsEncrypted bool   `json:"is_encrypted"`
	Extension   string `json:"extension"` //only set if not encrypted
	Data        string `json:"data"`
}

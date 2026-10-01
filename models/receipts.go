package models

type ReceiptVerificationResult struct {
	Valid      bool   `json:"valid"`
	Signature  string `json:"signature,omitempty"`
	VerifiedAt string `json:"verified_at,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
}
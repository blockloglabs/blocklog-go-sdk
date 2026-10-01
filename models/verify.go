package models

type VerifyLogResponse struct {
	Status        string                 `json:"status"`
	MerkleProof   map[string]interface{} `json:"merkle_proof,omitempty"`
	BatchProof    map[string]interface{} `json:"batch_proof,omitempty"`
	Details       map[string]interface{} `json:"details,omitempty"`
}

type VerifyBatchResponse struct {
	Status     string                 `json:"status"`
	Signature  string                 `json:"signature,omitempty"`
	SignedAt   string                 `json:"signed_at,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

type VerifyDecisionResponse struct {
	Status        string                 `json:"status"`
	MerkleProof   map[string]interface{} `json:"merkle_proof,omitempty"`
	Signature     map[string]interface{} `json:"signature,omitempty"`
	Details       map[string]interface{} `json:"details,omitempty"`
}
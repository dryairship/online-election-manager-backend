package models

import "github.com/dryairship/online-election-manager/config"

type (
	// Basic structure of a candidate that is stored in the database.
	Candidate struct {
		Roll       string `json:"roll"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		AuthCode   string `json:"authcode"`
		PostID     string `json:"postId"`
		Manifesto  string `json:"manifesto"`
		PublicKey  string `json:"publicKey"`
		PrivateKey string `json:"privateKey"`
		KeyState   int    `json:"keyState"`
	}

	// The object description that is returned through API calls.
	SimplifiedCandidate struct {
		Roll      string `json:"roll"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		PublicKey string `json:"publicKey"`
		//PrivateKey string `json:"privateKey"`
		Manifesto string `json:"manifesto"`
		State     string `json:"state"`
		KeyState  int    `json:"keyState"`
	}
)

// Function to return a simplified version of the data of a candidate.
func (candidate Candidate) Simplify() SimplifiedCandidate {
	return SimplifiedCandidate{
		Roll:      candidate.Roll,
		Username:  candidate.Username,
		Name:      candidate.Name,
		PublicKey: candidate.PublicKey,
		//PrivateKey: candidate.PrivateKey,
		Manifesto: candidate.Manifesto,
		State:     config.ElectionState,
		KeyState:  candidate.KeyState,
	}
}

// Function to return the mail recipient created from the candidate.
func (candidate Candidate) GetMailRecipient() MailRecipient {
	return MailRecipient{
		Name:     candidate.Name,
		EmailID:  candidate.Email,
		AuthCode: candidate.AuthCode,
	}
}

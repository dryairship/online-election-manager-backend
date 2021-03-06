package models

import (
	"github.com/dryairship/online-election-manager/config"
)

type (
	// Basic structure of a voter as stored in the database.
	Voter struct {
		Roll     string     `json:"roll"`
		Name     string     `json:"name"`
		Email    string     `json:"email"`
		Password string     `json:"password"`
		AuthCode string     `json:"authcode"`
		BallotID []BallotID `json:"ballotIds"`
		Voted    bool       `json:"voted"`
		Posts    []string   `json:"posts"`
		SHA1     string     `json:"sha1"`
		PassXac  string     `json:"passXac"`
	}

	// Struct to represent all the data required to send mail to a user.
	MailRecipient struct {
		Name     string `json:"name"`
		EmailID  string `json:"email"`
		AuthCode string `json:"authcode"`
	}

	// Basic structure of a student already present in the database.
	StudentSkeleton struct {
		Roll  string `json:"roll"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	// Voter model modified to return back to the user.
	SimplifiedVoter struct {
		Roll     string     `json:"roll"`
		Name     string     `json:"name"`
		BallotID []BallotID `json:"ballotIds"`
		Voted    bool       `json:"voted"`
		CEOKey   string     `json:"ceoKey"`
		State    string     `json:"electionState"`
		Posts    []string   `json:"posts"`
	}
)

// Function to generate a mail recipient from a voter,
func (voter Voter) GetMailRecipient() MailRecipient {
	return MailRecipient{
		Name:     voter.Name,
		EmailID:  voter.Email,
		AuthCode: voter.AuthCode,
	}
}

// Function to create a voter from the basic data of a student.
func (skeleton StudentSkeleton) CreateVoter(authcode string) Voter {
	return Voter{
		Roll:     skeleton.Roll,
		Name:     skeleton.Name,
		Email:    skeleton.Email + config.MailSuffix,
		Password: "",
		AuthCode: authcode,
		BallotID: nil,
		Voted:    false,
		Posts:    nil,
		SHA1:     "",
		PassXac:  "",
	}
}

// Function to convert a voter object to a form returnable to the user.
func (voter Voter) Simplify() SimplifiedVoter {
	return SimplifiedVoter{
		Roll:     voter.Roll,
		Name:     voter.Name,
		BallotID: voter.BallotID,
		Voted:    voter.Voted,
		CEOKey:   config.PublicKeyOfCEO,
		State:    config.ElectionState,
		Posts:    voter.Posts,
	}
}

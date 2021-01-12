package models

type (
	// Basic structure of a vote as stored in the database.
	Vote struct {
		PostID string `json:"postId"`
		Data   string `json:"data"`
	}

	// Basic structure of a parsed vote as stored in the database.
	ParsedVote struct {
		PostID      string `json:"postId"`
		BallotID    string `json:"ballotId"`
		Preference1 string `json:"preference1"`
		Preference2 string `json:"preference2"`
		Preference3 string `json:"preference3"`
	}

	// Struct to represent the received Ballot ID.
	BallotID struct {
		PostID       string `json:"postId"`
		BallotString string `json:"ballotString"`
	}

	// Struct to represent the received vote from the user.
	ReceivedVote struct {
		PostID       string `json:"postId"`
		BallotString string `json:"ballotString"`
		VoteData     string `json:"voteData"`
	}

	UsedSingleVoteBallotID struct {
		BallotString string `json:"ballotString"`
		Roll         string `json:"roll"`
		Name         string `json:"name"`
	}

	UsedBallotID struct {
		BallotString string `json:"ballotString"`
		Preference1  string `json:"preference1"`
		Preference2  string `json:"preference2"`
		Preference3  string `json:"preference3"`
	}
)

// Function to get the actual data of the vote from it.
func (receivedVote ReceivedVote) GetVote() Vote {
	return Vote{
		PostID: receivedVote.PostID,
		Data:   receivedVote.VoteData,
	}
}

// Function to get the ballot ID fron a vote.
func (receivedVote ReceivedVote) GetBallotID() BallotID {
	return BallotID{
		PostID:       receivedVote.PostID,
		BallotString: receivedVote.BallotString,
	}
}

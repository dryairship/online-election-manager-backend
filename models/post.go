package models

type (
	// Basic structure of the posts as stored in the database.
	Post struct {
		PostID     string   `json:"postId"`
		PostName   string   `json:"postName"`
		Candidates []string `json:"candidates"`
		Resolved   bool     `json:"resolved"`
		HasNota    bool     `json:"hasNota"`
	}

	// Structure of posts returned by the appropriate API call.
	VotablePost struct {
		PostID     string   `json:"postId"`
		PostName   string   `json:"postName"`
		Candidates []string `json:"candidates"`
		HasNota    bool     `json:"hasNota"`
	}
)

// Function to strip out regex from the data of the post.
func (post Post) ConvertToVotablePost() VotablePost {
	return VotablePost{
		PostID:     post.PostID,
		PostName:   post.PostName,
		Candidates: post.Candidates,
		HasNota:    post.HasNota,
	}
}

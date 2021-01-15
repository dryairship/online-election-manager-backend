package router

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/controllers"
)

// Groups various API calls and routes them to the
// respective controller functions.
func SetUpRoutes(r *gin.Engine) {

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			//users.POST("/mail/:roll", controllers.SendMailToStudent)
			users.POST("/register", controllers.EnsureBeforeVotingStopped(), controllers.RegisterNewVoter)
			users.POST("/login", controllers.CheckUserLogin)
			users.GET("/captcha", controllers.GetCaptcha)
		}

		election := api.Group("/election")
		{
			//election.GET("/getVotablePosts", controllers.EnsureVotingStarted(), controllers.EnsureLoggedIn(), controllers.GetVotablePosts)
			//election.GET("/getCandidateInfo/:username", controllers.GetCandidateInfo)
			election.GET("/getElectionState", controllers.GetElectionState)
			election.POST("/submitVote", controllers.EnsureVotingStarted(), controllers.EnsureLoggedIn(), controllers.SubmitVote)
			//election.GET("/singleVoteResults", controllers.GetSingleVoteResults)
			election.GET("/results", controllers.EnsureResultsCalculated(), controllers.GetResults)
		}

		ceo := api.Group("/ceo")
		ceo.Use(controllers.EnsureLoggedIn())
		ceo.Use(controllers.EnsureCEO())
		{
			ceo.POST("/startVoting", controllers.EnsureVotingNotYetStarted(), controllers.StartVoting)
			ceo.POST("/stopVoting", controllers.EnsureVotingStarted(), controllers.StopVoting)
			ceo.GET("/fetchPosts", controllers.EnsureVotingStopped(), controllers.FetchPosts)
			ceo.GET("/fetchVotes", controllers.EnsureVotingStopped(), controllers.FetchVotes)
			ceo.GET("/fetchCandidates", controllers.EnsureVotingStopped(), controllers.FetchCandidates)
			//ceo.POST("/submitSingleVoteResults", controllers.SubmitSingleVoteResults)
			ceo.POST("/submitResults", controllers.EnsureVotingStopped(), controllers.SubmitResults)
			//ceo.POST("/prepareForNextRound", controllers.PrepareForNextRound)
		}

		//candidate := api.Group("/candidate")
		//{
		//	candidate.POST("/confirmCandidature", controllers.ConfirmCandidature)
		//	candidate.POST("/declarePrivateKey", controllers.DeclarePrivateKey)
		//}

	}

	// To directly serve static files in the AssetsPath directory.
	r.Use(static.Serve("/api/data/", static.LocalFile(config.CandidatesOutputPath, true))) // for /data/candidates
	//r.Use(static.Serve("/ballotids/", static.LocalFile(config.BallotIDsPath, true)))
	//r.Use(static.Serve("/images/", static.LocalFile(config.ImagesPath, true)))
	// r.Use(static.Serve("/", static.LocalFile(config.AssetsPath, true)))
}

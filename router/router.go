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
			users.POST("/register", controllers.EnsureBeforeVotingStopped(), controllers.RegisterNewVoter)
			users.POST("/login", controllers.CheckUserLogin)
			users.GET("/captcha", controllers.GetCaptcha)
		}

		election := api.Group("/election")
		{
			election.GET("/getElectionState", controllers.GetElectionState)
			election.POST("/submitVote", controllers.EnsureVotingStarted(), controllers.EnsureLoggedIn(), controllers.SubmitVote)
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
			ceo.POST("/submitResults", controllers.EnsureVotingStopped(), controllers.SubmitResults)
		}
	}

	r.Use(static.Serve("/api/data/", static.LocalFile(config.CandidatesOutputPath, true))) // for /data/candidates
	r.Use(static.Serve("/ballotids/", static.LocalFile(config.BallotIDsPath, true)))
}

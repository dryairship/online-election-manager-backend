package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Possible Election States
const (
	VotingNotYetStarted = "VotingNotYetStarted"
	AcceptingVotes      = "AcceptingVotes"
	VotingStopped       = "VotingStopped"
	ResultsCalculated   = "ResultsCalculated"
)

// Possible states for a candidate
const (
	KeysNotGenerated = iota
	KeysGenerated
	KeysDeclared
)

// Global variables read from config file
var (
	ElectionState        string
	MailSuffix           string
	MongoDialURL         string
	MongoDbName          string
	MongoUsername        string
	MongoUsingAuth       bool
	MongoPassword        string
	BallotIDsPath        string
	VotersListPath       string
	CandidatesOutputPath string
	ElectionDataFilePath string
	ApplicationPort      string
	RollNumberOfCEO      string
	UsingCaptcha         bool
	CampusPassword       string
	DefaultAdminPassword string
)

// Global variable not read from config file
var PublicKeyOfCEO string

// Method to read the values of the global variables from environment variables.
func init() {
	viper.SetConfigName("backend-config")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("[WARN] Unable to locate configuration file: ", err.Error())
	}

	viper.AutomaticEnv()

	switch viper.GetString("ElectionState") {
	case VotingNotYetStarted:
		ElectionState = VotingNotYetStarted
	case AcceptingVotes:
		ElectionState = AcceptingVotes
	case VotingStopped:
		ElectionState = VotingStopped
	case ResultsCalculated:
		ElectionState = ResultsCalculated
	default:
		log.Fatal("ElectionState should be one of {VotingNotYetStarted, AcceptingVotes, VotingStopped, ResultsCalculated}")
	}

	MailSuffix = viper.GetString("MailSuffix")

	MongoDialURL = viper.GetString("MongoDialURL")
	MongoDbName = viper.GetString("MongoDbName")
	MongoUsername = viper.GetString("MongoUsername")
	MongoPassword = viper.GetString("MongoPassword")
	MongoUsingAuth = viper.GetBool("MongoUsingAuth")

	BallotIDsPath = viper.GetString("BallotIDsPath")
	VotersListPath = viper.GetString("VotersListPath")
	CandidatesOutputPath = viper.GetString("CandidatesOutputPath")
	ElectionDataFilePath = viper.GetString("ElectionDataFilePath")

	ApplicationPort = viper.GetString("ApplicationPort")
	UsingCaptcha = viper.GetBool("UsingCaptcha")

	RollNumberOfCEO = viper.GetString("RollNumberOfCEO")

	CampusPassword = viper.GetString("CampusPassword")
	DefaultAdminPassword = viper.GetString("DefaultAdminPassword")

	checkSanity()
}

func checkSanity() {
	if err := exists(BallotIDsPath); err != nil {
		log.Fatal("Error reading ballotids path")
	}
	if err := exists(CandidatesOutputPath); err != nil {
		log.Fatal("Error reading candidates output path")
	}
}

func exists(path string) error {
	_, err := os.OpenFile(path, os.O_RDONLY, 0)
	return err
}

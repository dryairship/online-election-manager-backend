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

// Global variables used by the program
var (
	ElectionState        string
	ElectionName         string
	MailSenderAuthID     string
	MailSenderEmailID    string
	MailSenderPassword   string
	MailSignature        string
	MailSMTPHost         string
	MailSMTPPort         string
	MailSuffix           string
	MongoDialURL         string
	MongoDbName          string
	MongoUsername        string
	MongoPassword        string
	AssetsPath           string
	BallotIDsPath        string
	ImagesPath           string
	DataPath             string
	ElectionDataFilePath string
	ApplicationPort      string
	SessionsKey          string
	MaxTimeDelay         int
	RollNumberOfCEO      string
	PublicKeyOfCEO       string
	PrivateKeyOfCEO      string
	ResultProgress       float64
)

// Method to read the values of the global variables from environment variables.
func init() {
	viper.SetConfigName("config-online-election-manager")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/go")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("[WARN] Unable to locate configuration file: ", err.Error())
	}

	viper.AutomaticEnv()

	switch viper.GetString("ElectionState") {
	case "VotingNotYetStarted":
		ElectionState = VotingNotYetStarted
	case "AcceptingVotes":
		ElectionState = AcceptingVotes
	case "VotingStopped":
		ElectionState = VotingStopped
	case "ResultsCalculated":
		ElectionState = ResultsCalculated
	default:
		log.Fatal("ElectionState should be one of {VotingNotYetStarted, AcceptingVotes, VotingStopped, ResultsCalculated}")
	}

	MailSenderEmailID = viper.GetString("MailSenderEmailID")
	MailSenderAuthID = viper.GetString("MailSenderAuthID")
	MailSenderPassword = viper.GetString("MailSenderPassword")
	MailSignature = viper.GetString("MailSignature")
	MailSMTPHost = viper.GetString("MailSMTPHost")
	MailSMTPPort = viper.GetString("MailSMTPPort")
	MailSuffix = viper.GetString("MailSuffix")

	MongoDialURL = viper.GetString("MongoDialURL")
	MongoDbName = viper.GetString("MongoDbName")
	MongoUsername = viper.GetString("MongoUsername")
	MongoPassword = viper.GetString("MongoPassword")

	AssetsPath = viper.GetString("AssetsPath")
	BallotIDsPath = viper.GetString("BallotIDsPath")
	ImagesPath = viper.GetString("ImagesPath")
	DataPath = viper.GetString("DataPath")
	ElectionDataFilePath = viper.GetString("ElectionDataFilePath")

	ApplicationPort = viper.GetString("ApplicationPort")
	SessionsKey = viper.GetString("SessionsKey")
	MaxTimeDelay = viper.GetInt("MaxTimeDelay")

	RollNumberOfCEO = viper.GetString("RollNumberOfCEO")
	ElectionName = viper.GetString("ElectionName")

	checkSanity()
}

func checkSanity() {
	if err := exists(AssetsPath); err != nil {
		log.Fatal("Error reading assets path:", err)
	}
	if err := exists(BallotIDsPath); err != nil {
		log.Fatal("Error reading ballotids path")
	}
	if err := exists(DataPath); err != nil {
		log.Fatal("Error reading data path")
	}
	if err := exists(ImagesPath); err != nil {
		log.Fatal("Error reading images path")
	}
}

func exists(path string) error {
	_, err := os.OpenFile(path, os.O_RDONLY, 0)
	return err
}

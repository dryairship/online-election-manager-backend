package main

import (
	"encoding/csv"
	"github.com/dryairship/online-election-manager/config"
	"github.com/dryairship/online-election-manager/db"
	"log"
	"math/rand"
	"os"
	"time"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	config.InitializeConfiguration()
	electionDb, err := db.ConnectToDatabase()
	if err != nil {
		log.Fatal("[ERROR] Could not establish database connection.")
	}

	candidates, err := electionDb.GetAllCandidates()
	if err != nil {
		log.Fatal("[ERROR] Could not fetch candidates.")
	}

	keys := readCsvFile("./keys.csv")
	perm := rand.Perm(10000)

	for i, candidate := range candidates {
		candidate.PublicKey = keys[perm[i+1]][0]
		candidate.PrivateKey = keys[perm[i+1]][1]
		err = electionDb.UpdateCandidate(candidate.Username, &candidate)
		if err != nil {
			panic(err)
		}
	}

	log.Println("Successfully initialized database for the purpose of demonstration.")
}

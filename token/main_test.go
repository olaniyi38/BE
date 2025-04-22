package token

import (
	"log"
	"os"
	"testing"

	"github.com/olaniyi38/BE/util"
)

var testJWTMaker Maker
var testPasetoMaker Maker

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../.")
	if err != nil {
		log.Fatal(err)
	}

	maker, err := NewJWTMaker(config.JWTSigningKey)

	if err != nil {
		log.Fatal(err)
	}

	testJWTMaker = maker

	maker, err = NewPasetoMaker(config.PasetoSymmetricKey)
	testPasetoMaker = maker

	os.Exit(m.Run())
}

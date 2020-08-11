package cascadeenv

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func checkOSEnv(names []string) bool {
	log.Info().Msg("Check Loaded ENV...")
	for _, name := range names {
		if os.Getenv(name) == "" {
			log.Info().Msgf("Cannot load ENV var: %s from environment...", name)

			return false
		}
	}
	return true
}

func loadAndCheckEnv(names []string, filename string) bool {
	log.Info().Msg("Check .ENV file...")
	err := godotenv.Load(filename)
	if err != nil {
		log.Info().Msg("Cannot load ENV file...")

		return false
	}
	log.Info().Msg("ENV file loaded; will check Environment Again...")
	return checkOSEnv(names)
}
func checkAWSParamStore(names []string, session *session.Session) bool {
	awsParamStore := ssm.New(session)
	for _, name := range names {
		res, err := awsParamStore.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(name),
			WithDecryption: aws.Bool(true),
		})
		if err != nil || res == nil {
			log.Info().Msg(name + " is not found in AWS ParamStore.")
			return false
		}
		os.Setenv(name, res.GoString())
	}
	return checkOSEnv(names)
}

//InitEnvVar init and checks os env for a list of names against ENV variables/ENV file/AWS ParamStore
func InitEnvVar(names []string, envFilename string, session *session.Session) error {
	// log.Info().Msg("ENV Checking. Loaded from OS.")
	if !checkOSEnv(names) {
		if !loadAndCheckEnv(names, envFilename) {
			if !checkAWSParamStore(names, session) {
				return fmt.Errorf("Neither OS nor ENV File nor AWS ParamStore has all the required names.")
			}
			log.Info().Msg("ENV Checking ok. Loaded from AWS ParamStore")
		}
		log.Info().Msg("ENV Checking ok. Loaded from ENV file.")
	}
	log.Info().Msg("ENV Checking. Loaded from OS.")
	return nil
}

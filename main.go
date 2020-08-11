package cascadeenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func checkOSEnv(names *[]string) bool {
	log.Info().Msg("Check Loaded ENV...")
	for _, name := range *names {
		if os.Getenv(name) == "" {
			log.Info().Msgf("Cannot load ENV var: %s from environment...", name)

			return false
		}
	}
	return true
}

func loadAndCheckEnv(names *[]string, filename string) bool {
	log.Info().Msg("Check .ENV file...")
	err := godotenv.Load(filename)
	if err != nil {
		log.Info().Msg("Cannot load ENV file...")

		return false
	}
	log.Info().Msg("ENV file loaded; will check Environment Again...")
	return checkOSEnv(names)
}
func checkAWSParamStore(names *[]string, session *session.Session) bool {
	if session == nil {
		return false
	}
	awsParamStore := ssm.New(session)
	for _, name := range *names {
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
	if len(names) == 0 {
		return fmt.Errorf("You cannot check 0 length of names")
	}
	if !checkOSEnv(&names) {
		if !loadAndCheckEnv(&names, envFilename) {
			if !checkAWSParamStore(&names, session) {
				return fmt.Errorf("Neither OS nor ENV File nor AWS ParamStore has all the required names.")
			}
		}
	}
	return nil
}

//ExportEnvVar init and checks os env for a list of names against ENV variables/ENV file/AWS ParamStore and exports to a map
func ExportEnvVar(names []string, types []reflect.Kind, envFilename string, session *session.Session) (map[string]interface{}, error) {
	// log.Info().Msg("ENV Checking. Loaded from OS.")
	ret := map[string]interface{}{}
	err := InitEnvVar(names, envFilename, session)
	if err != nil {
		return ret, err
	}
	for i, name := range names {
		readString := os.Getenv(name)
		if types[i] == reflect.Int64 {
			val, err := strconv.ParseInt(readString, 10, 64)
			if err != nil {
				return nil, fmt.Errorf(name + " has to be an Int64")
			}
			ret[name] = val
		} else if types[i] == reflect.Float64 {
			val, err := strconv.ParseFloat(readString, 64)
			if err != nil {
				return nil, fmt.Errorf(name + " has to be an Float64")
			}
			ret[name] = val
		} else {
			ret[name] = readString
		}
	}
	return ret, nil
}

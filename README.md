# CascadeEnv
init and checks os env for a list of names against ENV variables/ENV file/AWS ParamStore

# Usage

In the import block:
`import 	"github.com/abbychau/cascadeenv"`

In main:
```go
//NewAWSSession makes an AWS session for dynamodb
func NewAWSSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //Pretty logger for dev, can remove this line to use json-logger for production use

	log.Info().Msg("Booting...")
	err := cascadeenv.InitEnvVar(
		[]string{"HEALTH_CHECK_INTERVAL", "HEALTH_CHECK_MAX_CONSECUTIVE_ERROR", "SERVER_LISTEN", "USER_API_URL", "DB_USER", "DB_PASS", "DB_PROTOCOL", "DB_NAME"},
		".ENV",
		persistence.NewAWSSession(),
	)
}
```
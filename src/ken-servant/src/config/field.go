package config

type mFields struct {
	MASTER_HOST string "0.0.0.0"
	MASTER_PORT int    "6577"

	SERVANT_LISTEN_HOST string "0.0.0.0"
	SERVANT_LISTEN_PORT int    "6578"

	LOG_PATH  string "../logs"
	LOG_NAME  string "servant.log"
	LOG_LEVEL string "DEBUG"

	CERT_PATH string "../cert"
	CERT_PRIVATE_NAME string "private"
	CERT_PUBLIC_NAME string "public"
}
package config

/*
*	tag 为默认值
*/
type mFields struct {
	MASTER_LISTEN_HOST string "0.0.0.0"
	MASTER_LISTEN_PORT int    "6577"
	SERVANT_LISTEN_PORT int    "6578"

	LOG_PATH  string "../logs"
	LOG_LEVEL string "DEBUG"

	CERTS_PATH string "../certs"
}
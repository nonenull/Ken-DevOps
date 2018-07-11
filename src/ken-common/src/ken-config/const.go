package ken_config

const (
	LineTag        = "\r\n"
	EndTag         = LineTag + LineTag
	KeepAliveTag   = "keep-alive"
	NoKeepAliveTag = ""
	ReadBuffSize = 64
)

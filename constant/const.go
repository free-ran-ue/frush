package constant

const (
	CMD_START             = "frush > "
	CMD_HELP              = "help"
	CMD_EXIT              = "exit"
	CMD_ADD_SUBSCRIBER    = "add"
	CMD_DELETE_SUBSCRIBER = "delete"
)

const (
	OUTPUT_SUCCESS = "==> Success"
)

const (
	SYSTEM_HINT_CTRL_C_EXIT = "If you want to exit, please type 'exit'"
	SYSTEM_HINT_UNKNOWN_CMD = "Unknown command: %s"
)

const (
	HTTP_HEADER_CONTENT_TYPE      = "Content-Type"
	HTTP_HEADER_CONTENT_TYPE_JSON = "application/json"
)

const (
	CONSOLE_LOGIN_PATH          = "/api/login"
	CONSOLE_ADD_SUBSCRIBER_PATH = "/api/subscriber/imsi-%s/%s"
	CONSOLE_ACCESS_TOKEN        = "access_token"
	CONSOLE_TOKEN               = "Token"
)

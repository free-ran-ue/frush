package constant

const (
	CMD_START             = "frush > "
	CMD_HELP              = "help"
	CMD_EXIT              = "exit"
	CMD_ADD_SUBSCRIBER    = "add"
	CMD_DELETE_SUBSCRIBER = "delete"
	CMD_GNB               = "gnb"
	CMD_STATUS            = "status"
	CMD_UE_REGISTER       = "reg"
	CMD_UE_DE_REGISTER    = "dereg"
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

type ContextStatus string

const (
	CONTEXT_STATUS_GNB_RUNNING  ContextStatus = "running"
	CONTEXT_STATUS_GNB_STOPPED  ContextStatus = "stopped"
	CONTEXT_STATUS_GNB_STARTING ContextStatus = "starting"
	CONTEXT_STATUS_GNB_STOPPING ContextStatus = "stopping"
	CONTEXT_STATUS_GNB_ERROR    ContextStatus = "error"

	CONTEXT_STATUS_UE_STOPPED        ContextStatus = "stopped"
	CONTEXT_STATUS_UE_REGISTERED     ContextStatus = "registered"
	CONTEXT_STATUS_UE_REGISTERING    ContextStatus = "registering"
	CONTEXT_STATUS_UE_DE_REGISTERED  ContextStatus = "de-registered"
	CONTEXT_STATUS_UE_DE_REGISTERING ContextStatus = "de-registering"
	CONTEXT_STATUS_UE_ERROR          ContextStatus = "error"
)

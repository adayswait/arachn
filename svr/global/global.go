package global

type GroupType int
type RetCode int

//用户权限组
const (
	GROUP_ADMIN GroupType = iota
	GROUP_USER
	GROUP_VIEW
	GROUP_UNDEF
)

//api返回code枚举
const (
	RET_OK RetCode = iota
	RET_ERR
	RET_ERR_DB
	RET_ERR_HTTP_QUERY
	RET_ERR_USER_PASSWD
	RET_ERR_INPUT
	RET_ERR_URL_PARAM
	RET_ERR_BODY_PARAM
	RET_ERR_ACCESS_TOKEN
	RET_ERR_CREATE_TOKEN
	RET_ERR_NO_RIGHT
	RET_ERR_SESSION_INVALID
)

//数据库表名
const BUCKET_USR_PASSWD = "sys:usr:passwd"
const BUCKET_TOKEN_INFO = "sys:token:info"
const BUCKET_USER_TOKEN = "sys:user:token"
const BUCKET_USR_OPSLOG = "sys:usr:opslog"
const BUCKET_OPS_MACINI = "sys:ops:macini"
const BUCKET_OPS_DEPINI = "sys:ops:depini"

//cookie名
const ACCESS_TOKEN = "access_token"

//session key
const SESSION_KEY_USER = "user"
const SESSION_KEY_GROUP = "group"

const MD5_SALT = "mojo_salt"

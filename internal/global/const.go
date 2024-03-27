package global

type TaskProtocol int8
type Status int8
type TaskLevel int8

type TaskHTTPMethod int8

const (
	TaskHTTPMethodGet  TaskHTTPMethod = 1
	TaskHttpMethodPost TaskHTTPMethod = 2
)

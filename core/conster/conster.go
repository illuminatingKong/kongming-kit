package conster

const (
	Project = "kongming-kit"
)

type Level = Priority // severity

// The Priority is a combination of the syslog facility and
// severity. For example, USER | NOTICE.
type Priority int

// Level/severity
// These are the same on Linux, BSD, and OS X.
const (
	EMERGENCY Priority = iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

package wow

type loginResult = byte

// https://wowdev.wiki/Login_Packet_Results
// https://gtker.com/wow_messages/docs/loginresult.html#protocol-version-8

const (
	success loginResult = iota
	failUnknown0
	failUnknown1
	failBanned
	failUnknownAccount
	failIncorrectPassword
	failAlreadyOnline
	failNoTime
	failDbBusy
	failVersionInvalid
	loginDownloadFile
	failInvalidServer
	failSuspended
	failNoAccess
	successSurvey
	failParentalControl
	failLockedEnforced
)

var loginResultName = map[loginResult]string{
	success:               "SUCCESS",
	failUnknown0:          "FAIL_UNKNOWN0",
	failUnknown1:          "FAIL_UNKNOWN1",
	failBanned:            "FAIL_BANNED",
	failUnknownAccount:    "FAIL_UNKNOWN_ACCOUNT",
	failIncorrectPassword: "FAIL_INCORRECT_PASSWORD",
	failAlreadyOnline:     "FAIL_ALREADY_ONLINE",
	failNoTime:            "FAIL_NO_TIME",
	failDbBusy:            "FAIL_DB_BUSY",
	failVersionInvalid:    "FAIL_VERSION_INVALID",
	loginDownloadFile:     "LOGIN_DOWNLOAD_FILE",
	failInvalidServer:     "FAIL_INVALID_SERVER",
	failSuspended:         "FAIL_SUSPENDED",
	failNoAccess:          "FAIL_NO_ACCESS",
	successSurvey:         "SUCCESS_SURVEY",
	failParentalControl:   "FAIL_PARENTAL_CONTROL",
	failLockedEnforced:    "FAIL_LOCKED_ENFORCED",
}

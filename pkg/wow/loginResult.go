package wow

import "errors"

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

var ErrLoginFailed = errors.New("login failed")

var loginResultMessage = map[loginResult]string{
	success:               "success (SUCCESS)",
	failUnknown0:          "unable to connect (FAIL_UNKNOWN0)",
	failUnknown1:          "unable to connect (FAIL_UNKNOWN1)",
	failBanned:            "account is closed and no longer available for use (FAIL_BANNED)",
	failUnknownAccount:    "information you have entered is not valid (FAIL_UNKNOWN_ACCOUNT)",
	failIncorrectPassword: "information you have entered is not valid (FAIL_INCORRECT_PASSWORD)",
	failAlreadyOnline:     "account is already logged (FAIL_ALREADY_ONLINE)",
	failNoTime:            "prepaid time finished (FAIL_NO_TIME)",
	failDbBusy:            "could not log in at this time (FAIL_DB_BUSY)",
	failVersionInvalid:    "unable to validate game version (FAIL_VERSION_INVALID)",
	loginDownloadFile:     "downloading (LOGIN_DOWNLOAD_FILE)",
	failInvalidServer:     "unable to connect (FAIL_INVALID_SERVER)",
	failSuspended:         "account is temporarily suspended (FAIL_SUSPENDED)",
	failNoAccess:          "unable to connect (FAIL_NO_ACCESS)",
	successSurvey:         "success (SUCCESS_SURVEY)",
	failParentalControl:   "account is blocked by parental controls (FAIL_PARENTAL_CONTROL)",
	failLockedEnforced:    "unable to connect (FAIL_LOCKED_ENFORCED)",
}

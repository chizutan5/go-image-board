package config

import (
	"encoding/gob"
	"encoding/json"
	"html/template"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

//ConfigurationSettings contains the structure of all the settings that will be loaded at runtime.
type ConfigurationSettings struct {
	//DBName is the name of the db used for this instance
	DBName string
	//DBUser is the user name used to auth to the db
	DBUser string
	//DBPassword is the password used to auth to the db
	DBPassword string
	//DBPort the port the database is listening to
	DBPort string
	//DBHost hostname of the database server
	DBHost string
	//ImageDirectory path to where images are stored
	ImageDirectory string
	//Address hostname/port that this server should listen on
	Address string
	//ReadTimeout timeout allowed for reads
	ReadTimeout time.Duration
	//WriteTimeout timeout allowed for writes
	WriteTimeout time.Duration
	//MaxHeaderBytes maximum amount of bytes allowed in a request header
	MaxHeaderBytes int
	//SessionStoreKey stores the key to the session store
	SessionStoreKey [][]byte
	//CSRFKey stores the master key for CSRF token
	CSRFKey []byte
	//InSecureCSRF marks wether CSRF cookie should be secure or not, when developing this may be set to true, otherwise, keep false!
	InSecureCSRF bool
	//HTTPRoot directory where template and html files are kept
	HTTPRoot string
	//MaxUploadBytes maximum allowed bytes for an upload
	MaxUploadBytes int64
	//AllowAccountCreation if true, random users can create accounts, otherwise only mods can create users
	AllowAccountCreation bool
	//AccountRequiredToView if true, users must authenticate to access nearly any part of the server
	AccountRequiredToView bool
	//MaxThumbnailWidth Maximum width for automatically generated thumbnails
	MaxThumbnailWidth uint
	//MaxThumbnailHeight Maximum height for automatically generated thumbnails
	MaxThumbnailHeight uint
	//DefaultPermissions these permissions are assigned to all new users automatically
	DefaultPermissions uint64
	//UsersControlOwnObjects if this is set, permission checks are ignored for users that are trying to manage resources they contributed
	UsersControlOwnObjects bool
	//FFMPEGPath Path to the FFMPEG application
	FFMPEGPath string
	//UseFFMPEG If set, when joined with FFMPEGPath, videos that are uploaded will have a thumbnail generated using FFMPEG
	UseFFMPEG bool
	//PageStride How many images to show on one page
	PageStride uint64
	//APIThrottle How much time, in milliseconds, users using the API must wait between requests
	APIThrottle int64
	//UseTLS Enables TLS encryption on server
	UseTLS bool
	//TLSCertPath The path to the TLS/SSL cert
	TLSCertPath string
	//TLSKeyPath The path to the TLS/SSL key file for the cert
	TLSKeyPath string
	//ShowSimilarOnImages If enabled, shows similar count and link when viewing an image
	ShowSimilarOnImages bool
	//TargetLogLevel increase or decrease log verbosity
	TargetLogLevel int64
	//LoggingWhiteList regex based white-list for logging
	LoggingWhiteList string
	//LoggingBlackList regex based black-list for logging
	LoggingBlackList string
}

//SessionStore contains cookie information
var SessionStore *sessions.CookieStore

//Configuration contains all the information loaded from the config file.
var Configuration ConfigurationSettings

//ApplicationVersion Current version of application. This should be incremented every release
var ApplicationVersion = "1.0.4.4"

//SessionVariableName is used when checking cookies
var SessionVariableName = "gib-session"

//LoadConfiguration loads the specified configuration file into Configuration
func LoadConfiguration(Path string) error {
	//Open the specified file
	File, err := os.Open(Path)
	if err != nil {
		return err
	}
	defer File.Close()
	//Init a JSON Decoder
	decoder := json.NewDecoder(File)
	//Use decoder to decode into a ConfigrationSettings struct
	err = decoder.Decode(&Configuration)
	if err != nil {
		return err
	}
	return nil
}

//SaveConfiguration saves the configuration data in Configuration to the specified file path
func SaveConfiguration(Path string) error {
	//Open the specified file at Path
	File, err := os.OpenFile(Path, os.O_CREATE|os.O_RDWR, 0660)
	defer File.Close()
	if err != nil {
		return err
	}
	//Initialize an encoder to the File
	encoder := json.NewEncoder(File)
	//Encode the settings stored in configuration to File
	err = encoder.Encode(&Configuration)
	if err != nil {
		return err
	}
	return nil
}

//CreateSessionStore will create a new key store given a byte slice. If the slice is nil, a random key will be used.
func CreateSessionStore() {
	if Configuration.SessionStoreKey == nil || len(Configuration.SessionStoreKey) < 2 {
		Configuration.SessionStoreKey = [][]byte{securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)}
	} else if len(Configuration.SessionStoreKey[0]) != 64 || len(Configuration.SessionStoreKey[1]) != 32 {
		Configuration.SessionStoreKey = [][]byte{securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32)}
	}
	if Configuration.CSRFKey == nil || len(Configuration.CSRFKey) != 32 {
		Configuration.CSRFKey = securecookie.GenerateRandomKey(32)
	}
	//Register templates in gob for flash cookie usage
	gob.Register(template.HTML(""))
	SessionStore = sessions.NewCookieStore(Configuration.SessionStoreKey...)
}

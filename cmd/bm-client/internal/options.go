package internal

type options struct {
	Config   string `short:"c" long:"config" description:"Path to your configuration file"`
	Password string `short:"p" long:"password" description:"Vault password" default:""`
	Vault    string `long:"vault" description:"Custom vault file" default:""`
	Version  bool   `short:"v" long:"version" description:"Display version information"`
}

// Opts are the options set through the command line
var Opts options

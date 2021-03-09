// Copyright (c) 2021 BitMaelum Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package config

import "io"

const clientConfigTemplate string = `# BitMaelum Client Configuration Template. Edit for your own needs.
config:
    vault:
        # where are our accounts stored?
        path: "~/.bitmaelum/accounts.vault.json"
    composer:
        # Editor to use when composing messages. If not set, use the $EDITOR environment variable.
        editor: 
    server:
        # Should we be able to connect to self-signed and other insecure servers?
        allow_insecure: false
        # Display HTTP communication between client and server
        debug_http: false

    # How can we resolve public keys and accounts
    resolver:
        remote:
            # Enable remote resolving
            enabled: true
            url: "https://resolver.bitmaelum.com"
            # Allow insecure connections (to selfsigned certs)
            allow_insecure: false

`

const bridgeConfigTemplate string = `# BitMaelum Bridge Configuration Template. Edit for your own needs.
config:
    vault:
        # where are our accounts stored?
        path: "~/.bitmaelum/accounts.vault.json"

    server:
        smtp: 
            # Enable SMTP server to allow outgoing messages
            enabled: true

            # Host and port to listen for incoming connections
            host: "localhost"
            port: 1025

            # Run the SMTP server in gateway mode for the specified domain (for incoming mail)
            gateway: false
            domain: ""

            # Account from the vault to use to process incoming or outgoing mail (only for gateway mode)
            account: ""

            # Display SMTP communication between client and server
            debug: false

        imap: 
            # Enable IMAP4 server to allow outgoing messages
            enabled: true

            # Host and port to listen for incoming connections
            host: "localhost"
            port: 1143

            # Path to store a database that contains message flags (flags.db file). This
            # file should persists across reboots
            path: "~/.bitmaelum/bm-bridge"

            # Display IMAP communication between client and server
            debug: false
            
    resolver:
        # SQLite local resolver cache 
        sqlite:
            # Enable sqlite resolving
            enabled: false
            # Note: DSN currently does not support ~ homedir expansion
            dsn: "file:/tmp/keyresolve.db"
        # Remove resolver
        remote:
            # Enable remote resolving
            enabled: true
            # URL to the remote resolver
            url: "https://resolver.bitmaelum.com"
            # Allow insecure connections (to selfsigned certs)
            allow_insecure: false
    
`

const serverConfigTemplate string = `# BitMaelum Server Configuration Template. Edit for your own needs.
config:
    # Logging of information
    logging:
        # LogLevel: trace, debug, info, warn, error, crit
        log_level: trace

        # Log format. Either text or json
        log_format: text
        
        # Where to store logs. Can be one of the following:
        #   <path>,  
        #   stdout
        #   stderr 
        #   syslog
        #   syslog:<ip:port>
        log_path: stdout

        # Log apache-style HTTP combined log 
        apache_log: true
        # Path to apache logfile
        apache_log_path: "./bitmaelum.apache.log"

    work:
        # There can be multiple work-systems in the future. For now, we only accept "pow"
        pow:
            # How many bits of proof-of-work must a client/server do before a ticket will be issued
            bits: 25

    paths:
        # Path to store messages currently being processed (transient storage)
        processing: ~/.bitmaelum/processing
           
        # Path to store messages that have to be retried later (transient storage)
        retry: ~/.bitmaelum/retry

        # Path to store incoming messages (transient storage)        
        incoming: ~/.bitmaelum/incoming

        # Path to our message accounts (persistent storage)
        accounts: ~/.bitmaelum/messagedb

    acme:
        # When enabled, we can use LetsEncrypt to fetch valid SSL/TLS certificates. Note that certificate generation
        # does not happen automatically, but must be done manually with the help of the "bm-config" tool.
        enabled: false

        # Domain that we want to register (should be the same as your Server.Hostname)
        domain: localhost

        # Path to store our acme/LetsEncrypt settings and cache (persistent storage)
        path: ~/.bitmaelum/acme

        # The email address to register your LetsEncrypt account to. Important domain information
        # will be send to this address.
        email: info@example.org

        # The number of days before expiration of your certificate before we can renew with bm-config
        renew_days: 30

    server:
        # Hostname and port as this server is known on the internet in <host>:<port> setting. When 
        # no port is specified, port 2424 is used as default.
        hostname: localhost:2424

        # Address where to listen. Use 0.0.0.0 for all interfaces
        host: 127.0.0.1

        # Default port to listen on.
        port: 2424

        # Display additional version information on the server's root endpoint
        verbose_info: false

        # Should we be able to connect to self-signed and other insecure servers?
        allow_insecure: false

        # Certification and key for running on HTTPS. This should be a valid certificate 
        # via sLetEncrypt for instance or you can use self-signed certificates if you want
        certfile: "~/.bitmaelum/certs/server.cert"
        keyfile: "~/.bitmaelum/certs/server.key"

        # Routing file holds the server's keypair and routing ID.
        routingfile: "~/.bitmaelum/private/routing.json"

    organisations:
        # Temporary list of all organisations that are allowed to register on your server
        # - foo  
        # - bar
    management: 
        # When enabled, allow remote management through HTTPS instead of only local bm-config
        remote_enabled: false

    webhooks:
        # When enabled, users can add webhooks to their accounts
        enabled: false
        # Set to "default" to use the default internal worker system
        system: default
        # workers is the amount of standby workers that will deal with webhook events
        workers: 10

    bolt:
        # BoltDB database directory path to store the databases used for internal storage
        database_path: "~/.bitmaelum/database"

    redis:
        # Redis host where we store information
        host: 
        # Redis Database Number (defaults to 0)
        db: 0

    resolver:
        # SQLite local resolver cache 
        sqlite:
            # Enable sqlite resolving
            enabled: false
            # Note: DSN currently does not support ~ homedir expansion
            dsn: "file:/tmp/keyresolve.db"
        # Remove resolver
        remote:
            # Enable remote resolving
            enabled: true
            # URL to the remote resolver
            url: "https://resolver.bitmaelum.com"
            # Allow insecure connections (to selfsigned certs)
            allow_insecure: false

`

// GenerateClientConfig Generates a default client configuration
func GenerateClientConfig(w io.Writer) error {
	_, err := w.Write([]byte(clientConfigTemplate))

	return err
}

// GenerateServerConfig Generates a default server configuration
func GenerateServerConfig(w io.Writer) error {
	_, err := w.Write([]byte(serverConfigTemplate))

	return err
}

// GenerateBridgeConfig Generates a default bridge configuration
func GenerateBridgeConfig(w io.Writer) error {
	_, err := w.Write([]byte(bridgeConfigTemplate))

	return err
}

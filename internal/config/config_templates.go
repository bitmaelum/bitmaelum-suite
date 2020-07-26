package config

import "io"

const clientConfigTemplate string = `# BitMaelum Client Configuration Template. Edit for your own needs.
config:
    accounts:
        # where are our accounts stored?
        path: "~/.bitmaelum/accounts.vault.json"
        # Number of bits for proof-of-work. The higher the number, the longer it takes to generate an account
        proof_of_work: 22
    composer:
        # Editor to use when composing messages. If not set, use the $EDITOR environment variable.
        editor: 
    server:
        # Should we be able to connect to self-signed and other insecure servers?
        allow_insecure: false

    # How can we resolve public keys and accounts
    resolver:
        local:
            path: "~/.bitmaelum/.resolvecache"
        remote:
            url: "https://resolver.bitmaelum.com"
`

const serverConfigTemplate string = `# BitMaelum Server Configuration Template. Edit for your own needs.
config:
    # Logging of information
    logging:
        # LogLevel: trace, debug, info, warn, error, crit
        log_level: trace
        
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

    accounts:
        # How many bits of proof-of-work must an account have to be able to register here
        proof_of_work: 22
    
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
    redis:
        # Redis host where we store information
        host: 127.0.0.1:6379
        # Redis Database Number (defaults to 0)
        db: 0
    resolver:
        # Local resolver cache 
        local:
            path: "~/.bitmaelum/.resolvecache"
        # Remove resolver
        remote:
            url: "https://resolver.bitmaelum.com"
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

package config

import "io"

const clientConfigTemplate string = `
#
# BitMaelum Client Configuration Template
#
config:
    accounts:
        # where are our accounts stored?
        path: "~/.bitmaelum/accounts"

    # How can we resolve public keys and accounts
    resolver:
        local:
            path: "~/.bitmaelum/.resolvecache"
        remote:
            url: "https://resolver.bitmaelum.com"
`

const serverConfigTemplate string = `
#
# BitMaelum Server Configuration Template
#
config:
    # Logging of information
    logging:
        # LogLevel, trace, debug, info, warn, error, crit
        level: trace
    accounts:
        # Allow users to create an account on this server (with no limitations currently)
        registration: false
        
        # How many bits of proof-of-work must an account have to be able to register here
        proof_of_work: 22
    
        # Path to our message accounts
        path: ~/.messagedb
    server:
        # Address where to listen. Use 0.0.0.0 for all interfaces
        host: 127.0.0.1

        # Default port to listen on
        port: 2424
    tls:
        # Certification and key for running on HTTPS. This should be a valid certificate via LetEncrypt 
        # for instance.
        # You can use self-signed certificates is you want
        certfile: ~/.bitmaelum/certs/server.crt
        keyfile: ~/.bitmaelum/certs/server.key
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

// Generates a default client configuration
func GenerateClientConfig(w io.Writer) error {
    _, err := w.Write([]byte(clientConfigTemplate))

    return err
}

// Generates a default server configuration
func GenerateServerConfig(w io.Writer) error {
    _, err := w.Write([]byte(serverConfigTemplate))

    return err
}

# Tool: bm-config

Allows you to easily configure certain aspects of your mail server.

Current supported the following configuration settings:


* `invite`: Enables an account to be registered on the mailserver.
        
        ./bm-config uninvite --address [address] --days <30>
        
   Each address invite is valid for a maximum of a number of days (default 30). After this, the address
   should be registered again with invite command.
   
   This command returns a token that must be used while registring an account in the mail client.

* `uninvite`: Removes an account invite

        ./bm-config uninvite --address [address]

* `init-config`: Generates either a server or client configuration (or both) that can be your starting point.

        ./bm-config init-config [--server] [--client]

* `generate-cert`: Generates a self-signed certificate

        ./bm-config generate-cert --domain bitmaelum.mydomain.com

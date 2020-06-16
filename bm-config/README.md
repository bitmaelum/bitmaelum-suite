# Tool: bm-config

Allows you to easily configure certain aspects of your mail server.

Current supported the following configuration settings:


* `allow-registration`: Enables an account to be registered on the mailserver.
        
        ./mail-server-config allow-registration --address [address] --days <30>
        
   Each address invite is valid for a maximum of a number of days (default 30). After this, the address
   should be registered again with allow-registration command.
   
   This command returns a token that must be used while registring an account in the mail client.

* `remove-registration`: Removes an account invite

        ./mail-server-config remove-registration --address [address]

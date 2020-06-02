# What if
we could redesign email without any need for backward compatiblity. What would it look like? Probably not like this, but at least its an attempt.
We're trying to figure out how to make email a secure (end-to-end encrypted) system that will both combat spam and brings you back in charge of 
your mail again. This means you and only you can signup for mailing-lists, you can unsubscribe whenever you like (cannot be ignored by the 
mailinglist owners). It also means a vast reduction of spam emails, as sending spam is expensive for spammers.

It even provides additional functionality:

  - less spam (hopefully, spam free)
  - no more email address leaks
  - host your messages wherever you like
  - move your messages to other competitor without loosing your mail address 
  - easily detect mailing-lists for your favourite companies and organisations
  - don't get tracked (unless you want to) 

# Usage

* Think of a FQDN to host your message server. It doesn't have to exist for local development. For instance `mail.test.v2`

* Create a server certificate and key file. 

    openssl req -x509 -newkey rsa:4096 -keyout server-passphrase.key -out server.crt -days 3650

    Most important part is that you set your CN (common name)to the FQDN chosen above. You need to add a passphrase to the key.
 
* Remove the passphrase from the key:
 
    openssl rsa -in server-passphrase.key -out server.key

* Or, if you like, you can always use LetsEncrypt to fetch your own certificate (with the help of certbot). Make sure you 
  name your certificate `server.crt` and key `server.key`.

* Build the components:

    ./build.sh
    
    
# Mail server

* You'll need redis. Run either locally or via a container. Make sure you set the host info in the configuration file `./config.yml`

* Run the mail-server:

    ./release/mailv2-server -config ./config.yml

    This will run a mailserver on localhost port 2424 over a TLS connection.
    
    Configuration for the server can be found in `./server-config.example.yml`

## Paths and files

Public keys that are pushed to the mailserver are stored in the `.keydb` directory.
Mail is stored in the `.maildb` directory.


# Mail client

Some kind of client is underway. It uses rivo/tview as UI framework.

![client.png](client.png) 


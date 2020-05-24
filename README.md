# Usage

* Think of a FQDN to host your email server. It doesn't have to exist for local development. For instance `mail.test.v2`

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
    
    Configuration for the server can be found in `./config.yml`

## Paths and files

Public keys that are pushed to the mailserver are stored in the `.keydb` directory.
Mail is stored in the `.maildb` directory.


# Mail client

* There is no functional client yet.

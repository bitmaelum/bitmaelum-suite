# Current status

This document describes the current status of the project. It's in its current form completely 
useless and consists of small components that may or may not work as advertised. The current goal 
is to make sure our happy-path works the way it should be: we should be able to send mail, and 
receive it.

To make this happen, we need lots of components: message structures, encryption methods, server 
APIs, client connectivity etc etc. This is all still very much work in progress. 

## We are currently able to

 - [X] Add and manage public keys on the main resolver service: resolver.bitmaelum.com
 - [X] Run the message server with minimum capabilities
 - [X] Generate accounts and save them to the server and main resolver through bm-client.
 - [X] List accounts through bm-client
 - [X] Compose email, and uploading to the mail-server through bm-client
 - [X] Read email through utility
 - [X] Setting status flags on the mail-server (@todo: do we want this on the server, and if so, unencrypted?)
 
## Next up in our todo

 - [] Send email that is uploaded to the mail-server to an actual destination mail-server.
 - [] Read mailboxes from accounts
 

## Later

 - [] Deal with multiple recipients (multi-header uploads)
 - [] Mailing lists
 - [] Create a simple mail client UI (bm-client-ui is a start)


## Much later

 - Just about everything else...

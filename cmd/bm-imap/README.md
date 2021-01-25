# BM-Imap

The bm-imap service allows you to read bitmaelum messages directly via imap, meaning that you can use existing mail clients for 
READING bitmaelum mail.  SENDING mail however, requires SMTP, which is not implemented yet (that would be a bm-smtp).

## Usage

1. Start the bm-imap service:


     bm-imap -p password

The IMAP service is running on LOCALHOST port 1143.

2. Add account to your email client:

    Use the following credentials:

       
        email address: emailfield bitmaelum address (see below)
       
        service: IMAP
        host:    localhost (or 127.0.0.1)
        port:    1143
        security:  NONE (no STARTTLS/SSL)
       
        username: emailfied bitmaelum address (see below)  (ie: johndoe or johndoe.organisation)
        password: anything will do

If you need to specify an outgoing server, use a dummy SMTP account. 

s

### Emailfied bitmaelum addresses

One of the problems with the imap bridge is that bitmaelum addresses are not compatible with regular email address.
In order to make this work, we have to "emailfy" bitmaelum address. This is done the following way:

    johndoe!     ==>   johndoe@bitmaelum.network
    john@org!    ==>   john_org@bitmaelum.network


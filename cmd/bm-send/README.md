# bm-send

A simple tool that will let you send messages on behalf of a bitmaelum address.

  $ ./bm-send --priv "PRIVKEY" --subject "a mail" --from "from!" --to "to!" --body "...." --attachment foo.zip
  
will use the following environment variables:

    BITMAELUM_SEND_PRIVKEY
    BITMAELUM_SEND_FROM


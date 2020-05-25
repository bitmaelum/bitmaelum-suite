# Structures

## Mail database
 
    .maildb/
        1d/                                                                                 # first two chars of the sha256 hash of the email address
            9A6E493390DA308F5A082E4DC4ABA11BC0B4C0829BB3DCE222897E8B3CCE3A/                 # Remainder of the sha256 hash of the email address
                .pubkey                                                                     # Public key for this account.                                                                          
                inbox/                                                                      # mail directories
                    9fe27f14-a30b-46c9-9303-060dbcbae92d/                                   # directory per email, consisting of email ID
                        header.json                                                         # header file of email (TODO: Do we need this in our email?)
                        catalog.json                                                        # encrypted catalog.json file (see below)
                        34e941f2-5e5f-4855-94a1-73572a5f8f29.dat                            # one or more blocks or files, based on ID from the catalog (encrypted, but not base64 encoded)                
                        113c9c21-3319-4026-92d2-57181734410a.dat
                        504a2dfb-1231-4714-a5de-df6d184fbddc.dat
                outbox/
                trash/
                archive/
                ...


## Message header file
The message header file is an unencrypted file that allows mailservers to deal with incoming emails. It should contain 
as little metadata as possible.

    {
      "from": {
        "sha256": "1db14ca62f80ce8c77ae749ddf04e2d4077451db9f701f3bd37dec1dc950c76b"
        "public_key": "-----BEGIN PUBLIC KEY-----MIICCgKCAgEA1ANJCd79pIdlNNbUsNKqLHsLPzETkw/gnDFFpHgnYsUCvPO+P68c35QiTgVu+rwuLz8koxHFpbFEmMecqTbUdweNi7MeerYK07HG6L7MB9y/HzJ7Ig9eYvZXcErGh3R5fDq7aJSdp9arHbuL0PyNti5qoUrUBC5mCVdcvlg+7V19sX/2nG3DQ3gXOAIV1gMGhOww7w+sRtYU/sFZYIQY0S3HofkwhE7zsI9ARjHYW6aPcsevK/fnKHArry+SxscLxEN50Jd6QPgXREjN6Z/2HWxBIME5Mr+xaprA/YwS+MDsEY9AXMwDoAUywB7E0GWayD+QPYhxaBTzwIqKpLSmUH9/S6Q/g64n7taMSnrw4VEDH8Q1rO8p8ryutvZTuTFd4KvJOUsekpFE9UGNghm737AcgdJSXSu8tRj1putbfTgeAGm4SsiIrYUarXFUcxCPk2RAo+MD1FFj462nXD3Q4dtDZAh6+4jJhsvJKL8YxUXMHTWC7Wh69v/tLzfnGwzZSf/HDwVbI+GTYiDrAjsXlJgBqhdLBcSWYFT/weYKuzHjS+7LvfxvXQ7kH4X9fs0jUIJVIRU7Cj8YSt1JGiKUtyQY+RvSYl/HOSflx2/OeXtpZhVeLzQtGckzkUOaSJCJg7dj9JvdlnwVle3UEd1mFVvVWdfREVQ5PRVIbPxsfjkCAwEAAQ==-----END PUBLIC KEY-----"
        "proof_of_work": {
          "bits": 20,
          "proof": 396056
        }
      },
      "to": {
        "sha256": "19b14ca62f80ce8c77ae749ddf04e2d4077451db9f701f3bd37dec1dc950c76b",
      },
      "catalog": {
        "size": 1728,
        "checksum" : [
          { "hash": "sha1", "value": "8b53323e62705215db067ad8f296f490a57b89cf" },
          { "hash": "sha256", "value": "61b0476f8e249c68bbaa14222de9f29f078133ff9d0b461d69b78f5a63f5b678" },
        ],
        "crypto": "rsa+aes256"  
        "key": "hlyhaNEVfLfkHVCcexyMXGLRTGPaZElrj49ktvFFZf8eo/wvUJaG/Szt7+Bu7tJWRVnoyP2kldKl79ZrxtfmQOcXw1spJO8pjrqxy1ZM68KnawOW2jPqjNBexJER4ibkRF2ILWZplo4bumrvpfuX80HGZXeFCO7MP/SdUtyvxt/BMfJv4E/l+q/6YfP+apMvyPiEkwg4CaBRXhXhyB665L6BtKlFNI38Vg6jUaJ3SXYi0R7vGUmAfWBiuaI5G2nBw0i1IDZAu3dVundBE3PJsU48eyAt1JAEQ7I0rh3D9VfzCojRGM9pYxcQfGO0/lu0zinl+WTpSfNsOuzvg2MsSKXq5VTMJnNQv4tIYoEY68rOwa4XR//2s8FdfB2K84vWiNhrMOvPz4QqQ04yLq/4VrRB6L881JRpXEt7l7bcLPuQDyGBl8kq+k5ONpMfs/w/+WjySd8N/PGPyXBhnKWP9YnxpHbqZ+t1AyNj9bYdcyOqzUHQLk3/jAq0z7VTInYiV/vTztR2ZdKewm2Yfl1r4uH3Z7d/v6W+Vgy/DuONvF433Z9vAfuoWbtBM+qqzZuFZH6DWSFWicZkFYF50pN1rooafmycM9yemWkaOdN8tE/u26puOacH9R87stSjjNbdy9CXWy+ZJ7Jjj62o3SzfmbFeOi/24Vd0NfvMDpZkDjo=",
        "iv": "l1Pk4gkdolVrm7dpR569vg=="
      }
    }


### explanation

    from.sha256                     The SHA256 of the email address that sends the email
    from.public_key                 The public key of this user (TODO: do we need this?)
    from.proof_of_work.bits         How many bits of work has this proof done
    from.proof_of_work.proof        The actual proof, taken by from.sha256 data
    
    to.sha256                       The SHA256 of the email address that received the email. This email box should exist on the server
    
    catalog.size                    Size (in bytes) of the catalog file
    catalog.checksum[].hash         The hash method for the given checksum
    catalog.checksum[].value        Actual hash value for the whole catalog file
    catalog.crypto                  The crypto used for encrypting the catalog file. Only "rsa+aes256" is available for now.
    catalog.key                     RSA encrypted aes256 key. Encrypted with the public key of the receiver
    catalog.iv                      IV of the aes256 encryption.  


## Message Catalog file
The catalog file consists of all message meta-data and the catalog of the blocks. This file is encrypted and cannot be read
directly from the mail server.

    {
      "from": {
        "email": "joshua@noxlogic.nl",
        "name": "Joshua Thijssen",
        "organisation": "NoxLogic"
      },
      "to": {
        "email": "info@seams-cms.com",
        "name": "Seams-CMS"
      },
      "created_at": "2020-05-25T06:58:28+00:00",
      "thread_id": "bb9d16ea-b232-4280-b1cd-aead2edf0c0c",
      "subject": "Your invoice 2010-14141 from NoxLogic",
      "labels": [
        "invoice",
        "service"
      ],
      "catalog": {
        "blocks": [
          {
            "id": "113c9c21-3319-4026-92d2-57181734410a",
            "type": "html",
            "size": 149,
            "encoding": "base64",
            "compression": "gzip",
            "checksum": [
              {
                "hash": "sha1",
                "value": "c1d68280163089ceff11c54d6162a9504d9330f3"
              },
              {
                "hash": "sha256",
                "value": "f3c4de175280e55a630b11563f1d9d4a037f6f8017a548037dc75df10deb4435"
              },
              {
                "hash": "crc32",
                "value": "d4a14a92"
              }
            ]
          },
          {
            "id": "34e941f2-5e5f-4855-94a1-73572a5f8f29",
            "type": "mobile",
            "size": 149,
            "encoding": "base64",
            "compression": "gzip",
            "checksum": [
              {
                "hash": "sha1",
                "value": "804beebfc9b8078a262be51edf2d358ab1b91d85"
              },
              {
                "hash": "sha256",
                "value": "a899261a63a0cd7e2386ec8a214707e5fd3b9364cc7e90109c774f80b422bb58"
              },
              {
                "hash": "crc32",
                "value": "efa3ba3e"
              }
            ]
          },
          {
            "id": "504a2dfb-1231-4714-a5de-df6d184fbddc",
            "type": "text",
            "size": 101,
            "encoding": "base64",
            "compression": "gzip",
            "checksum": [
              {
                "hash": "sha1",
                "value": "63f69cfa77e8c03d1ee401d701a0dae901f15e9a"
              },
              {
                "hash": "sha256",
                "value": "ca24bc0c7ac3964dda00eb66dcceba016c16a3d8e8a87f2294bbadef3aaecdf6"
              },
              {
                "hash": "crc32",
                "value": "54b62f2c"
              }
            ]
          }
        ],
        "attachments": [
          {
            "id": "0f5f38cd-fdea-4dc0-a837-6d06f9c37a34",
            "mimetype": "application/binary",
            "filename": "textfile.doc",
            "size": 14914141,
            "encoding": "base64",
            "compression": "gzip",
            "checksum": [
              {
                "hash": "sha1",
                "value": "c1d68280163089ceff11c54d6162a9504d9330f3"
              },
              {
                "hash": "sha256",
                "value": "f3c4de175280e55a630b11563f1d9d4a037f6f8017a548037dc75df10deb4435"
              },
              {
                "hash": "crc32",
                "value": "d4a14a92"
              }
            ]
          }
        ]
      }
    }

### explanation

    from.email                                      Sender email
    from.name                                       Sender name
    from.organisation                               Sender organisation
    
    to.email                                        Email address of the receiver
    to.name                                         Name of the receiver
    
    created_at                                      ISO 8601 date when the message was created
    thread_id                                       Thread ID (if any)
    subject                                         Subject of the message
    labels[]                                        Additional labels (invoice, service, important etc)

    catalog.blocks[].id                             UUID of the block
    catalog.blocks[].type                           Block name (html, mobile, text)
    catalog.blocks[].size                           Size of the block in bytes
    catalog.blocks[].crypto                         Crypto used (rsa+aes supported)
    catalog.blocks[].key                            RSA encrypted aes256 key
    catalog.blocks[].iv                             IV for aes256
    catalog.blocks[].compression                    Compression method used (gzip)
    catalog.blocks[].checksum[].hash                The hash method for the given checksum
    catalog.blocks[].checksum[].value               Actual hash value for the whole catalog file

    catalog.attachments[].id                        UUID of the file
    catalog.attachments[].mimetype                  Mimetype of the file
    catalog.attachments[].filename                  Original filename
    catalog.attachments[].size                      Size of the block in bytes
    catalog.attachments[].crypto                    Crypto used (rsa+aes supported)
    catalog.attachments[].key                       RSA encrypted aes256 key
    catalog.attachments[].iv                        IV for aes256
    catalog.attachments[].compression               Compression method used (gzip)
    catalog.attachments[].checksum[].hash           The hash method for the given checksum
    catalog.attachments[].checksum[].value          Actual hash value for the whole catalog file
    
    
    
    
    

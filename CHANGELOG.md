## 0.1.0 (2021-???-??) 

Finalized "hello-world" release

### Security
- no issues found or fixed

### Changes
- Removed some obsolete tooling
- bm-client Organisations in bm-client can be suffixed with '!' so you can use either "foo" or "foo!" as organisation value.
- bm-client: Added functionality for managing organisation validations 
- bm-client: Added functionality to query info from the resolver 
- bm-client: "account create" and "organisation create" use a "step" layout for better feedback
- bm-client: reading messages is done concurrently, speeding it up a huge amount 

### Fixes
- bm-server: added whitelisting in config to only allow specific organisations to add accounts



## 0.0.2 (2021-jan-11)

Small update with some minor tweaks

### Security
-  no issues found or fixed

### Changes
- <a href="http://github.com/bitmaelum/bitmaelum-suite/commit/5ec838bca10fc0a898f76702230c29fb732719a4">view commit</a> &bull; "address" argument is now called "account" in all bm-client commands</li>
- <a href="http://github.com/bitmaelum/bitmaelum-suite/commit/945b7cfb997ac818b409d6b420e1634be0ddc0be">view commit</a> &bull; displaying mnmemonic if available in your account display</li>

### Fixes
- <a href="http://github.com/bitmaelum/bitmaelum-suite/commit/d7fd2281a96d4291d8b37c4e37bbeae9790df247">view commit</a>&bull; mandatory flag is called 'id' instead of 'message'</li>
- <a href="http://github.com/bitmaelum/bitmaelum-suite/commit/88ecab97d09aa5b912e12ea48693d1c1ccf7625d">view commit</a>&bull; #151: vault directories are now created when they do not exist</li> 



## 0.0.1 (2021-jan-10)

Initial developers-only release

package account

import (
    "github.com/bitmaelum/bitmaelum-server/bm-client/account"
    "github.com/bitmaelum/bitmaelum-server/core"
)

var Accounts []core.AccountInfo

//// Loads all accounts from the vault
//func LoadAccountsFromVault(pwd []byte) {
//    var err error
//    Accounts, err = UnLockVault(config.Client.Accounts.Path, pwd)
//    if err != nil {
//        panic("Cannot unlock vault: " + err.Error())
//    }
//}

// Returns the given account for address
func FindAccount(address core.Address) *core.AccountInfo {
    for _, account := range account.Vault.Accounts {
        if account.Address == address.String() {
            return &account
        }
    }

    return nil
}

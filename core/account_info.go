package core

type AccountInfo struct {
    Address         string          `json:"address"`                // The address of the account
    Name            string          `json:"name"`                   // Full name of the user
    Organisation    string          `json:"organisation"`           // Organisation of the user (if any)
    PrivKey         string          `json:"privKey"`                // Private key
    PubKey          string          `json:"pubKey"`                 // Public key
    Pow             ProofOfWork     `json:"pow"`                    // Proof of work
    Server          string          `json:"server"`                 // Mail server hosting this account
}

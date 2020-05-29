package account

type ProofOfWork struct {
    Bits    int     `json:"bits"`
    Proof   uint64  `json:"proof"`
}

type AccountInfo struct {
    Address         string          `json:"address"`
    Name            string          `json:"name"`
    Organisation    string          `json:"organisation"`
    PrivKey         string          `json:"privKey"`
    PubKey          string          `json:"pubKey"`
    Pow             ProofOfWork     `json:"pow"`
}
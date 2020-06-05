package core

type AccountInfo struct {
    Address         string          `json:"address"`
    Name            string          `json:"name"`
    Organisation    string          `json:"organisation"`
    PrivKey         string          `json:"privKey"`
    PubKey          string          `json:"pubKey"`
    Pow             ProofOfWork     `json:"pow"`
}

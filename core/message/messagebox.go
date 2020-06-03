package message

// Structure of the .info.json file
type MailBoxInfo struct {
    Name            string
    Description     string      `json:"description"`
    Quota           int         `json:"quota"`
}

type Flags map[string]string

type MessageInfo struct {
    Flags       map[string]string   `json:"flags"`
    Catalog     Catalog             `json:"catalog"`
}

type Pubkeys struct {
    PubKeys []string `json:"keys"`
}

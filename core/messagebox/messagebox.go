package messagebox

// Structure of the .info.json file
type MailBoxInfo struct {
    Name            string
    Description     string      `json:"description"`
    Quota           int         `json:"quota"`
    Total           int         `json:"total"`
}

type Flags struct {
    Flags []string  `json:"flags"`
}

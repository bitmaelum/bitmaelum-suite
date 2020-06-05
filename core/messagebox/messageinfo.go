package messagebox

import "github.com/jaytaph/mailv2/core/catalog"


type MessageInfo struct {
    Flags       map[string]string   `json:"flags"`
    Catalog     catalog.Catalog     `json:"catalog"`
}

package account

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "regexp"
    "strings"
)

const (
    ADDRESS_REGEX   string = "(^[a-z0-9][a-z0-9\\.\\-]{2,63})(@[a-z0-9][a-z0-9\\.\\-]{1,63})?!$"
)

func HashId(id string) string {
    sum := sha256.Sum256([]byte(strings.ToLower(id)))
    return hex.EncodeToString(sum[:])
}

func ParseId(id string) (string, string, error) {
    re := regexp.MustCompile(ADDRESS_REGEX)
    if re == nil {
        return "", "", errors.New("cannot compile regex")
    }

    matches := re.FindStringSubmatch(id)
    return matches[0], matches[1], nil
}

func ValidateId(id string) bool {
    matched, err := regexp.MatchString(ADDRESS_REGEX, id)
    if err != nil {
        return false
    }

    return matched
}



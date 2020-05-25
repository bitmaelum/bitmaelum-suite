package keys

import (
    "fmt"
    "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "path"
    "regexp"
    "strings"
)

func RemoveKey(hash string) {
    keyPath := getKeyPath(hash)

    if _, err := os.Stat(keyPath); os.IsNotExist(err) {
        return
    }

    logrus.Infof("removing key file %s", keyPath)
    _ = os.Remove(keyPath)
}

func GetKey(hash string) string {
    keyPath := getKeyPath(hash)

    if _, err := os.Stat(keyPath); os.IsNotExist(err) {
        return ""
    }

    pubKey, err := ioutil.ReadFile(keyPath)
    if err != nil {
        return ""
    }

    return string(pubKey)
}

func AddKey(hash string, key string) {
    keyPath := getKeyPath(hash)

    logrus.Tracef("storing key in path %s", keyPath)

    // Create path if needed
    dir := path.Dir(keyPath)
    err := os.MkdirAll(dir, 0700)
    if err != nil {
        logrus.Panic(err)
    }

    // Remove any newlines if found
    re := regexp.MustCompile(`\r?\n`)
    key = re.ReplaceAllString(key, "")

    // @TODO: Does not take into account we can have multiple keys
    _ = ioutil.WriteFile(keyPath, []byte(key), 0600)
}

func HasKey(hash string) bool {
    keyPath := getKeyPath(hash)

    logrus.Tracef("checking if key exists on path %s", keyPath)

    if _, err := os.Stat(keyPath); os.IsNotExist(err) {
        return false
    }

    return true
}

func getKeyPath(hash string) string {
    return fmt.Sprintf(".keydb/%s/%s", strings.ToLower(hash[:2]), strings.ToLower(hash[2:]))
}


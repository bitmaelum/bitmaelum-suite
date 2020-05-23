package keys

import (
    "fmt"
    logger "github.com/sirupsen/logrus"
    "io/ioutil"
    "os"
    "path"
)

func RemoveKey(hash string) {
    keyPath := getKeyPath(hash)

    if _, err := os.Stat(keyPath); os.IsNotExist(err) {
        return
    }

    logger.Infof("Removing key file %s", keyPath)
    _ = os.Remove(keyPath)
}

func AddKey(hash string, key string) {
    keyPath := getKeyPath(hash)

    logger.Tracef("Storing key in path %s", keyPath)

    // Create path if needed
    dir := path.Dir(keyPath)
    err := os.MkdirAll(dir, 0700)
    if err != nil {
        logger.Panic(err)
    }

    // @TODO: Does not take into account we can have multiple keys
    _ = ioutil.WriteFile(keyPath, []byte(key), 0600)
}

func HasKey(hash string) bool {
    keyPath := getKeyPath(hash)

    logger.Tracef("checking if key exists on path %s", keyPath)

    if _, err := os.Stat(keyPath); os.IsNotExist(err) {
        return false
    }

    return true
}

func getKeyPath(hash string) string {
    return fmt.Sprintf(".keydb/%s/%s", hash[:2], hash[2:])
}


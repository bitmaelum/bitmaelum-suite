package bitmaelumClient

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	vaultFile = `{
	"type": "vault",
	"version": 2,
	"data": "d7h27532BAXtz5hMgR+Pho5VBiFEm+ZIa0oA6yTfdvBVGGV1OMS5dIW/R6Kxal32ja6pYnaYoF/nlyxdhCWa6Jz8CuNh/0MYu34NPKSQnYbw7GoaZ8z+2op/ISGneHa574Zc0pmWsdvhIe+uAaegEDcy7cDUcL3AGdMLKZflGMg3jfatyS+j8TnywDyvkU5pbUr1c8Kvf6AAMX7M7o13prFj6D30CrO8X7dQwi6DyG9PjXftszJgWVl43L1/ytl9OoRFESUo6BXYMSUqED6rw9nSxJsCE3b1Qo4u2ZAlhn4mYFa8FMEas16FxXa85+YSKAAYaLhKq17FcPhJQ1rofVr1vTk9X07T4YJdzIQLN27QWSXXK8crKdGNbX/U4WjtCWMUpuSmI0oEzwHag0QEa1PZM8wwdwVUUJFNQpdVULk9FeWdJwxf0uBNgmcGAkian1uuVo9boieybQ9WLXZ+6SrdHMfoSJVnK3ENZBxP2cYOPcS5CipqmF37cRjfMVRsomVUZNsbOZPilSeSnTzkcN3wwsQHwMrmihfqq8pzrHdNpA089TY696wTENiJAuuADyzpXmxgSzajdMpKCL/smfheXecvjCOO1YUkqno9RdaNle6Wmw5idTHHfRBHW/KVqbHFDvHiv7OhcSmJome++JWFQt9awIDkLXgpLlzneQuBTTl2I+rj9XySZBPhsr2i7bJj4jFomHzIhi0WO7t3TjUlaeOpdRiZlM7QXFlq7FkOOjCXpbjLVRlhVZ7GqPWXca0yG5W1obK4nUoceIJLywVynXc6tury7JUMbNZoewYwTtGtVQa7qbQmonqcMmq+SYw4d9LDdgPUXlLy0CF9d+NskZJuvq1356seW4yggELCAjoQnWAcHLQt2K+ToRELIeyvDNkD6iSyab0zRUajgoHA6vW/EpKB8gG593Z6QQ7CW4idGKlDvPLc/l6HH2fjEdc7dbS7FvdP0qSDa4L34uXxnusvlcehq7cTdICvTNHOv02RXSajZXmPDYh1DtLV33e+nx7eUChJsJENtBL+3hwBG/gWz6efJhqEKvJoABJyYqhqnhToaT3fJnJNIqUd9kArefTSnBRc/6IHLD3lOohZ6g1It6h0I99ooKvGKZTa0oLDD42xAFgVxWj7q8yZI9HqwidATyfpWBt9czXTw+UxoVj5dppBkQUGydgMp60SufLyVYSrRhLwGLemcxwpdlqIL2OYHs+cwdqIQlWyDYQJLv1STnXEFCoSdpLgOtrIA9VXOv2Flot3eS4Y0EJxBYDLcXq1Vg/fbqnEBi5s1kgWllHMI63ppUP/015QgjAy6caVx/JJnHNSc9lVK+MJwBcMzry7YmoRAtvLSuaYxQiA5ajEwVZANWYoyl4Xio6f6BnWhc/yi5pL4225qtdRMRE+012iRsk+49ImVqmAl7bxTxF44sCo4yy+b1P+b7K+Q0pfAjiBTbuH4VSFMIynFCsT",
	"salt": "WWr5aYsVn1YjLqoh9p1f2Hp48dP2iQ2okj7syjHgjwUVULs85P13s4P+5J2dNLRkg+Iho5+whuH0WlSYTyCSOA==",
	"iv": "Is/5cNuwY3F/UoP/m4VZIw==",
	"hmac": "PQZryF4oeQHJkCnX16sWrIK51nkOO9rxOsQg5zVSSWs="
}`
	vaultPrefix   = "vault"
	vaultPassword = "pass"
)

func TestOpenVault(t *testing.T) {
	// Write vault to file
	file, err := ioutil.TempFile(os.TempDir(), vaultPrefix)
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	err = os.MkdirAll(filepath.Dir(file.Name()), 0755)
	assert.NoError(t, err)

	err = ioutil.WriteFile(file.Name(), []byte(vaultFile), 0600)
	assert.NoError(t, err)

	// Open non-existant vault
	_, err = NewBitMaelumClient().OpenVault("does-not-exist.vault", vaultPassword)
	assert.Error(t, err)

	// Open vault
	v, err := NewBitMaelumClient().OpenVault(file.Name(), vaultPassword)
	assert.NoError(t, err)
	assert.Equal(t, v.([]map[string]interface{})[0]["address"], "test30!")
}

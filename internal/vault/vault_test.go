package vault

import (
	"testing"

	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	v, err := New("", []byte{})
	assert.NoError(t, err)
	assert.NotNil(t, v)
}

func Test_FindShortRoutingId(t *testing.T) {
	var acc internal.AccountInfo

	v, _ := New("", []byte{})

	acc = internal.AccountInfo{Address: "example!", RoutingID: "123456780000"}
	v.AddAccount(acc)
	acc = internal.AccountInfo{Address: "example!", RoutingID: "123456780001"}
	v.AddAccount(acc)
	acc = internal.AccountInfo{Address: "example!", RoutingID: "123456780002"}
	v.AddAccount(acc)
	acc = internal.AccountInfo{Address: "example!", RoutingID: "154353535335"}
	v.AddAccount(acc)

	assert.Equal(t, "154353535335", v.FindShortRoutingId("154"))
	assert.Equal(t, "154353535335", v.FindShortRoutingId("15435"))
	assert.Equal(t, "", v.FindShortRoutingId("12345"))
	assert.Equal(t, "", v.FindShortRoutingId("1"))
}

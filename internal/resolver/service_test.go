package resolver

import (
	"github.com/bitmaelum/bitmaelum-suite/internal"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_generateAddressSignature(t *testing.T) {
	privKey1, pubKey1, _ := internal.ReadTestKey("../../testdata/key-1.json")
	_, pubKey2, _ := internal.ReadTestKey("../../testdata/key-2.json")

	info := AddressInfo{
		Hash:      "12345",
		PublicKey: *pubKey1,
		RoutingID: "56789",
		Pow:       proofofwork.New(5, "foo", 1234).String(),
		RoutingInfo: RoutingInfo{
			Hash:      "12345",
			PublicKey: *pubKey2,
			Routing:   "127.0.0.1",
		},
	}
	s := generateAddressSignature(&info, *privKey1)
	assert.Equal(t, "oxKOBosJTMhXxTOKSqBYkPRtjaswd0BiW/J7bVQT5BcRHKmwpsulsxK8rUtk2Grf70Xu+Ja+s2Xjla08r9297yPp06TM4xMoPk835/HfmYg6TD5BglxiTNqvJv7TXlxfyiMKPc6GMJN171w0XNCqQ/kc13rRQe/vY1MLWKqTvA/opgD2D0O5eo6d9OK429NLKbMmaadz1sB2wYY7vSY27H1bZ63uD7HxV+XBaKYWpXW0pCtP/vmVKTCTGWxMlRpWYHAdOCc9/iWT89UQTbpCb/nJdcPp0FGOIC4eESf403JQbgrkiYC6O9PXOoQwLEBqEFUM11ErMECZERqSmUVXvw==", s)
}

func Test_generateOrganisationSignature(t *testing.T) {
	privKey1, pubKey1, _ := internal.ReadTestKey("../../testdata/key-1.json")

	info := OrganisationInfo{
		Hash:        "12345",
		PublicKey:   *pubKey1,
		Pow:         proofofwork.New(5, "foo", 1234).String(),
		Validations: nil,
	}
	s := generateOrganisationSignature(&info, *privKey1)
	assert.Equal(t, "L0kbS1VZrjhz1VSVP271DrTVMLWjmb7XyOC5/jSZB8lhZY2KLO+M9Si50loxOOR7EhmoJjPzp9OONffNJYIv8rD4e86zr/sgYL8aJ//a49yOjc0C57DI8E44TF827ibzyUXOtg2IzlDLrkbIwyFDvbQU7MBD3DXNVBVLC3OG3eKILhv8siK6wJpVnpIygR8c9PTq6Zy0KfKOuOreANMGJTDt/oKYQTKLrOpcPV85B+ch1LB9A+sMqrVllAyKLomTmzyENSYy9ZPsrzUcbZyzq892EiZYB8Alg554ejmSn4ic3yXKmtrnlY4gxUHYc93fFDDjZdUE8Q2tnNvJwCDwig==", s)
}

func Test_generateRoutingSignature(t *testing.T) {
	privKey1, pubKey1, _ := internal.ReadTestKey("../../testdata/key-1.json")

	info := RoutingInfo{
		Hash:      "12345",
		PublicKey: *pubKey1,
		Routing:   "127.0.0.1",
	}
	s := generateRoutingSignature(&info, *privKey1)
	assert.Equal(t, "L0kbS1VZrjhz1VSVP271DrTVMLWjmb7XyOC5/jSZB8lhZY2KLO+M9Si50loxOOR7EhmoJjPzp9OONffNJYIv8rD4e86zr/sgYL8aJ//a49yOjc0C57DI8E44TF827ibzyUXOtg2IzlDLrkbIwyFDvbQU7MBD3DXNVBVLC3OG3eKILhv8siK6wJpVnpIygR8c9PTq6Zy0KfKOuOreANMGJTDt/oKYQTKLrOpcPV85B+ch1LB9A+sMqrVllAyKLomTmzyENSYy9ZPsrzUcbZyzq892EiZYB8Alg554ejmSn4ic3yXKmtrnlY4gxUHYc93fFDDjZdUE8Q2tnNvJwCDwig==", s)
}

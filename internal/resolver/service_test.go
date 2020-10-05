package resolver

import (
	"testing"

	bmtest "github.com/bitmaelum/bitmaelum-suite/internal/testing"
	"github.com/bitmaelum/bitmaelum-suite/pkg/proofofwork"
	"github.com/stretchr/testify/assert"
)

func Test_generateAddressSignature(t *testing.T) {
	privKey1, pubKey1, _ := bmtest.ReadTestKey("../../testdata/key-1.json")
	_, pubKey2, _ := bmtest.ReadTestKey("../../testdata/key-2.json")

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
	s := generateAddressSignature(&info, *privKey1, 1234567890)
	assert.Equal(t, "SQ//6Uxbt6YOkl7povRPgQOuIoB788iMHd2SqpXxHmgBRZhxh1CYoXuAIf1ry2jdzq6lwdrydlEFYmdfKzFNOh6oy2vGOQul24QOuHL1qnYnsVPc13H2clsl/jtw3T0Is6hh+JFOmu5tDY/3xRex4AWA9KRBsl1cTlNPdHAXanV6UVk9MBZSUE4H+kJoA7BBom0rUeZ14V/SalVbGRb62qiFjW6JMYQFXZNGEqOxGCp+maZNCta5vSjK/7oEG3a8Kmns+OVwFMGdU+1fkmi6ebaYNbMh3ta9cEv87qzoodGLBCI5Q/RjTCNnZ8ERwFatnmh4sGzEMXSuUHkqc2zo/g==", s)

	s = generateAddressSignature(&info, *privKey1, 4325262)
	assert.Equal(t, "dcMB2GoQAaQ+OmcVPPjwroIJUyE1b7bv/C21xRhgOAwN9rrPGw9n5h+VzBJ2pcZsSdieCVBV4L2gg/imzzf5iKTtHk+EED1ZirZR3CJC/HrrCCeDVIVDcAOk1NoFlBwdf3RawAoVCfCPX/y1QIg5A1FG/vkKvtxGNNYwEgXHIhjNm5cVyRvGEVbrukmc1O1/xDBeelzP32KH60MuC5Fht5PNHwKJwOD4XtmmVjufqQJw0bAJhsUch3Pkti2bkkGvWP4iloNPy76hojlcCTAokFKk+qSt9HX9CvWAL749V3nl5CnGwrOL2yjavuB608NnuuuJH3F1k6EAKVMI8aUqVQ==", s)

	s = generateAddressSignature(&info, *privKey1, 1234567890)
	assert.Equal(t, "SQ//6Uxbt6YOkl7povRPgQOuIoB788iMHd2SqpXxHmgBRZhxh1CYoXuAIf1ry2jdzq6lwdrydlEFYmdfKzFNOh6oy2vGOQul24QOuHL1qnYnsVPc13H2clsl/jtw3T0Is6hh+JFOmu5tDY/3xRex4AWA9KRBsl1cTlNPdHAXanV6UVk9MBZSUE4H+kJoA7BBom0rUeZ14V/SalVbGRb62qiFjW6JMYQFXZNGEqOxGCp+maZNCta5vSjK/7oEG3a8Kmns+OVwFMGdU+1fkmi6ebaYNbMh3ta9cEv87qzoodGLBCI5Q/RjTCNnZ8ERwFatnmh4sGzEMXSuUHkqc2zo/g==", s)
}

func Test_generateOrganisationSignature(t *testing.T) {
	privKey1, pubKey1, _ := bmtest.ReadTestKey("../../testdata/key-1.json")

	info := OrganisationInfo{
		Hash:        "12345",
		PublicKey:   *pubKey1,
		Pow:         proofofwork.New(5, "foo", 1234).String(),
		Validations: nil,
	}
	s := generateOrganisationSignature(&info, *privKey1, 1234567890)
	assert.Equal(t, "fDyQ1/pGp+CWvnzGp2XasFwK0jfUP6VbrYtYGiqVJAEjh/qWZLjfOCrXn6a/mwVmIUh6UVNz/uf7nchED8UARpdIovXFDvTvjGgwitcz5T8Bo+b0WWVymO7PGJVBTZmRk1cIjM4N9CuSH8ObUIl6McaDL7CmR6MuyYn6BYglzd5uwde8B0AfhnmgaAFu2NDSPkeRxboh6jCH60PiAV03a/TAN+t7apkEv+uQ5RXTGLYZXw0CxCmPIr8aQMyF+/d20fLirq7SDh7taHcsTZFgZ1VWEchEGpm3ZyWw3W1lCx1/tdXOgKA1/9Ott2HVzRvrrn1Uh20R7OgJxk2ou8319g==", s)

	s = generateOrganisationSignature(&info, *privKey1, 643632)
	assert.Equal(t, "pG0c9GRNt9LB9OLXBaRgGWXm6Y8ZpFUXgnm0K48C/aKhsgBSpd8w18k0hQLhUfLBJj5fby0CFHvvm6B9Dy4speYeQHQiH5jKOVNSJdvjbdDHvehVdgJ9E58GG1Ck2vTON2q+oFcs1WdTe2K5PoqP+A1ITgkEM+F4br3AGESOKyyNgMV0PjOrBDUaxc1U0t2wCITrtgM4dHee8VsIB+31ys+5LljuFSwVGo55DGBhyfdtPUj6WEG1tzMim9s9NMWE/gMTh2TDIBP7x2C7Q08FSJuQ18QPMZLs7kFWKa8tSbudSxKqtaAFEaKOc22NlWc9MnpoA+dXJJz0dA/egOC4dw==", s)

	s = generateOrganisationSignature(&info, *privKey1, 1234567890)
	assert.Equal(t, "fDyQ1/pGp+CWvnzGp2XasFwK0jfUP6VbrYtYGiqVJAEjh/qWZLjfOCrXn6a/mwVmIUh6UVNz/uf7nchED8UARpdIovXFDvTvjGgwitcz5T8Bo+b0WWVymO7PGJVBTZmRk1cIjM4N9CuSH8ObUIl6McaDL7CmR6MuyYn6BYglzd5uwde8B0AfhnmgaAFu2NDSPkeRxboh6jCH60PiAV03a/TAN+t7apkEv+uQ5RXTGLYZXw0CxCmPIr8aQMyF+/d20fLirq7SDh7taHcsTZFgZ1VWEchEGpm3ZyWw3W1lCx1/tdXOgKA1/9Ott2HVzRvrrn1Uh20R7OgJxk2ou8319g==", s)
}

func Test_generateRoutingSignature(t *testing.T) {
	privKey1, pubKey1, _ := bmtest.ReadTestKey("../../testdata/key-1.json")

	info := RoutingInfo{
		Hash:      "12345",
		PublicKey: *pubKey1,
		Routing:   "127.0.0.1",
	}

	s := generateRoutingSignature(&info, *privKey1, 1234567890)
	assert.Equal(t, "fDyQ1/pGp+CWvnzGp2XasFwK0jfUP6VbrYtYGiqVJAEjh/qWZLjfOCrXn6a/mwVmIUh6UVNz/uf7nchED8UARpdIovXFDvTvjGgwitcz5T8Bo+b0WWVymO7PGJVBTZmRk1cIjM4N9CuSH8ObUIl6McaDL7CmR6MuyYn6BYglzd5uwde8B0AfhnmgaAFu2NDSPkeRxboh6jCH60PiAV03a/TAN+t7apkEv+uQ5RXTGLYZXw0CxCmPIr8aQMyF+/d20fLirq7SDh7taHcsTZFgZ1VWEchEGpm3ZyWw3W1lCx1/tdXOgKA1/9Ott2HVzRvrrn1Uh20R7OgJxk2ou8319g==", s)

	s = generateRoutingSignature(&info, *privKey1, 643263262)
	assert.Equal(t, "tbhINAg4aASAGiRITGIuVqMTOlWJXlg+qSKVVCikE1ACtzuc1/4t72KiyzSNZIiWo3Gw0HzyGQBt980fNTCZOk6gV2VdF4o/Fg6q/B/Up7WVhTKcPjHypIOihStsnuqaGQzQZouJxpIVA8HF+Y0rxpW9+lCmK4Xesq0qbmsQLXlRWb6ttADUaz0Uqu49PDwq0iGeC2Fzs4XxK5oTqB168VL0JcuwBQyuvNBu+vs21wgBenFMOK7oTvKp6n/xa54Z2Cs5mgEcTYCxjwn4u7kSWEvJg/p9Th/e05SnWL11446UufyY2z+2C1G16FD6bv6H81SpfbbUEXQeH2SMoWfHyA==", s)

	s = generateRoutingSignature(&info, *privKey1, 1234567890)
	assert.Equal(t, "fDyQ1/pGp+CWvnzGp2XasFwK0jfUP6VbrYtYGiqVJAEjh/qWZLjfOCrXn6a/mwVmIUh6UVNz/uf7nchED8UARpdIovXFDvTvjGgwitcz5T8Bo+b0WWVymO7PGJVBTZmRk1cIjM4N9CuSH8ObUIl6McaDL7CmR6MuyYn6BYglzd5uwde8B0AfhnmgaAFu2NDSPkeRxboh6jCH60PiAV03a/TAN+t7apkEv+uQ5RXTGLYZXw0CxCmPIr8aQMyF+/d20fLirq7SDh7taHcsTZFgZ1VWEchEGpm3ZyWw3W1lCx1/tdXOgKA1/9Ott2HVzRvrrn1Uh20R7OgJxk2ou8319g==", s)

}

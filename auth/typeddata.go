package auth

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/libs4go/bcf4go/eip712"
	identity "github.com/web3zerotrust/trust-identity"
)

var typesStandard = eip712.Types{
	"EIP712Domain": {
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "version",
			Type: "string",
		},
		{
			Name: "chainId",
			Type: "uint256",
		},
		{
			Name: "verifyingContract",
			Type: "address",
		},
	},
	"Token": {
		{
			Name: "did",
			Type: "string",
		},
		{
			Name: "session",
			Type: "string",
		},
	},
}

func newDomain(chainId uint) eip712.TypedDataDomain {
	return eip712.TypedDataDomain{
		Name:              "Trust-Identity-Eth",
		Version:           "1",
		ChainId:           eip712.NewHexOrDecimal256(int64(chainId)),
		VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		Salt:              "",
	}
}

func toMessage(token *identity.Token) eip712.TypedDataMessage {
	buff, _ := json.Marshal(token)

	var msg eip712.TypedDataMessage

	json.Unmarshal(buff, &msg)

	return msg
}

func newTypedData(token *identity.Token, chainId uint) *eip712.TypedData {
	return &eip712.TypedData{
		Types:       typesStandard,
		PrimaryType: "Token",
		Domain:      newDomain(chainId),
		Message:     toMessage(token),
	}
}

func etherAddressToDID(address string) string {
	return fmt.Sprintf("did:TypedData:%s", strings.ToLower(address))
}

func verifyTypedData(token *identity.Token, chainId uint, signature []byte) (bool, error) {
	typedData := newTypedData(token, chainId)

	recoveredAddr, err := eip712.Recover(typedData, signature)

	if err != nil {
		return false, err
	}

	return strings.HasSuffix(strings.ToLower(token.DID), strings.ToLower(recoveredAddr)), nil
}

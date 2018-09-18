/*
 * Copyright (C) 2017 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package promise

import (
	"encoding/base64"
	"testing"

	"github.com/mysteriumnetwork/node/identity"
	"github.com/mysteriumnetwork/node/money"
	"github.com/stretchr/testify/assert"
)

const CurrencyToken = money.Currency("Token")

func Test_PromiseBody(t *testing.T) {

	amount := money.Money{
		Amount:   uint64(5),
		Currency: CurrencyToken,
	}

	promise := Promise{
		SerialNumber: 1,
		IssuerID:     "issuer1",
		BenefiterID:  "benefiter1",
		Amount:       amount,
	}

	assert.Equal(t, 1, promise.SerialNumber)
	assert.Equal(t, "issuer1", promise.IssuerID)
	assert.Equal(t, "benefiter1", promise.BenefiterID)
	assert.Equal(t, uint64(5), promise.Amount.Amount)
	assert.Equal(t, CurrencyToken, promise.Amount.Currency)
}

func Test_SignedPromise(t *testing.T) {

	promise := Promise{}

	signedPromise := SignedPromise{
		Promise:         promise,
		IssuerSignature: "signature",
	}

	assert.Equal(t, promise, signedPromise.Promise)
	assert.Equal(t, Signature("signature"), signedPromise.IssuerSignature)
}

func TestNewPromise(t *testing.T) {
	promise := NewPromise(
		identity.Identity{Address: "Consumer"},
		identity.Identity{Address: "Provider"},
		money.Money{Amount: 123, Currency: CurrencyToken},
	)

	assert.Equal(t, 1, promise.SerialNumber)
	assert.Equal(t, "Consumer", promise.IssuerID)
	assert.Equal(t, "Provider", promise.BenefiterID)
	assert.Equal(t, uint64(123), promise.Amount.Amount)
	assert.Equal(t, CurrencyToken, promise.Amount.Currency)
}

type fakeSigner struct{}

func (fs *fakeSigner) Sign(message []byte) (identity.Signature, error) {
	return identity.SignatureBytes([]byte("FakeSignature")), nil
}

func TestSignByIssuer(t *testing.T) {
	promise := NewPromise(
		identity.Identity{Address: "Consumer"},
		identity.Identity{Address: "Provider"},
		money.Money{Amount: 123, Currency: "TEST"},
	)

	signedPromise, err := SignByIssuer(promise, &fakeSigner{})
	assert.NoError(t, err)

	expectedSignature := base64.StdEncoding.EncodeToString([]byte("FakeSignature"))
	assert.Equal(t, expectedSignature, string(signedPromise.IssuerSignature))
}

/*
 * MIT License
 *
 * Copyright (c) 2024 Anton Stremovskyy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package go_easypay

import (
	"strconv"
)

type Merchant struct {
	// Merchant Name
	Name string
	// Merchant ID
	PartnerKey string
	// Merchant Key
	ServiceKey string
	// System Key
	SecretKey string

	// SuccessRedirect
	SuccessRedirect string

	// FailRedirect
	FailRedirect string

	PayeeID          string
	PayeeName        string
	PayeeBankAccount string
	PayeeNarative    string
	PayerName        string
}

func (m *Merchant) GetMerchantID() *int64 {
	id, err := strconv.ParseInt(m.PartnerKey, 10, 64)

	if err != nil {
		return nil
	}

	return &id
}

func (m *Merchant) getPartnerKey() string {
	return m.PartnerKey
}

func (m *Merchant) GetSecretKey() string {
	return m.SecretKey
}

func (m *Merchant) GetServiceKey() string {
	return m.ServiceKey
}

// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package native

import "testing"

// The crypto module is exercised end-to-end from corpus/bal/library/subset2/crypto-*.bal
// (hashing, HMAC, KDF, password, AES, RSA/EC sign-verify-encrypt, PKCS#8) and from the
// corpus/extern crypto-keystore test (PKCS#12 keystore/trust-store).
//
// rc2New is a verbatim copy of golang.org/x/crypto/pkcs12/internal/rc2 (an internal
// package we cannot import). Ballerina exposes no RC2 encryption — RC2 is used only to
// *decrypt* PKCS#12-PBE-protected keys via cipher.NewCBCDecrypter, which exercises
// Decrypt (covered by crypto-keys-enc-v.bal). The Encrypt method is required by the
// cipher.Block interface but never reached from Ballerina, so this is the one place it
// is covered. Encrypt cannot simply be deleted: NewCBCDecrypter takes a cipher.Block,
// which mandates an Encrypt method.
func TestRC2Cipher(t *testing.T) {
	block, err := rc2New([]byte{1, 2, 3, 4, 5}, 40)
	if err != nil {
		t.Fatalf("rc2New: %v", err)
	}
	if block.BlockSize() != rc2BlockSize {
		t.Errorf("BlockSize = %d, want %d", block.BlockSize(), rc2BlockSize)
	}
	src := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	enc := make([]byte, rc2BlockSize)
	block.Encrypt(enc, src)
	dec := make([]byte, rc2BlockSize)
	block.Decrypt(dec, enc)
	if string(dec) != string(src) {
		t.Errorf("RC2 round-trip: got %v, want %v", dec, src)
	}
}

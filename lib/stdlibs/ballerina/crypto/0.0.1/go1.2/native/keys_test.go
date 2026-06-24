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

import (
	"crypto/rsa"
	"os"
	"testing"

	"ballerina-lang-go/runtime"
	"ballerina-lang-go/runtime/extern"
	"ballerina-lang-go/semtypes"
	"ballerina-lang-go/test_util/testharness"
	"ballerina-lang-go/values"
)

// newCryptoTestEnv builds the cryptoTypes and extern.Context needed to drive the
// key-building and keystore/trust-store helpers directly in unit tests.
func newCryptoTestEnv() (cryptoTypes, *extern.Context) {
	env := semtypes.CreateTypeEnv()
	byteArrBld := semtypes.NewListDefinition()
	keyMapBld := semtypes.NewMappingDefinition()
	utcBld := semtypes.NewListDefinition()
	types := cryptoTypes{
		byteArrTy: byteArrBld.DefineListTypeWrappedWithEnvSemType(env, semtypes.BYTE),
		keyMapTy:  keyMapBld.DefineMappingTypeWrapped(env, nil, semtypes.STRING),
		utcTy:     utcBld.TupleTypeWrappedRo(env, semtypes.INT, semtypes.DECIMAL),
	}
	return types, &extern.Context{TypeCtx: semtypes.ContextFrom(env)}
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("reading fixture %s: %v", name, err)
	}
	return data
}

// testKeyStore is a PKCS#12 keystore (RSA key + self-signed cert, friendly name
// "ballerina") protected by the password "secret". Used to drive the custom
// PKCS#12 trust-store parser (decodeCertFromTrustStore and its helpers), which
// is not reachable from a corpus .bal test because it requires a binary file.
const (
	testKeyStorePath  = "testdata/keystore.p12"
	testKeyStorePass  = "secret"
	testKeyStoreAlias = "ballerina"
)

func readKeyStore(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(testKeyStorePath)
	if err != nil {
		t.Fatalf("reading keystore fixture: %v", err)
	}
	return data
}

func TestDecodeCertFromTrustStore(t *testing.T) {
	data := readKeyStore(t)

	// Correct password + alias yields the embedded certificate.
	cert, err := decodeCertFromTrustStore(data, testKeyStorePass, testKeyStoreAlias)
	if err != nil {
		t.Fatalf("decodeCertFromTrustStore: %v", err)
	}
	if cert == nil {
		t.Fatal("expected a certificate, got nil")
	}
	if cert.PublicKey == nil {
		t.Error("certificate has no public key")
	}
}

func TestDecodeCertFromTrustStoreErrors(t *testing.T) {
	data := readKeyStore(t)

	// A wrong password must fail to decrypt the store.
	if _, err := decodeCertFromTrustStore(data, "wrong-password", testKeyStoreAlias); err == nil {
		t.Error("expected error for wrong password")
	}

	// Malformed PKCS#12 data must be rejected.
	if _, err := decodeCertFromTrustStore([]byte("not a pkcs12 file"), testKeyStorePass, testKeyStoreAlias); err == nil {
		t.Error("expected error for malformed PKCS#12 data")
	}
}

func TestDecodeBMPString(t *testing.T) {
	// "AB" encoded as a big-endian UTF-16 (BMPString): 0x0041 0x0042.
	got, err := decodeBMPString([]byte{0x00, 0x41, 0x00, 0x42})
	if err != nil {
		t.Fatalf("decodeBMPString: %v", err)
	}
	if got != "AB" {
		t.Errorf("decodeBMPString = %q, want %q", got, "AB")
	}

	// Odd-length input is invalid.
	if _, err := decodeBMPString([]byte{0x00}); err == nil {
		t.Error("expected error for odd-length BMPString")
	}
}

func TestParsePrivateKeyPEM(t *testing.T) {
	// Unencrypted PKCS#8 RSA key.
	if _, err := parsePrivateKeyPEM(readFixture(t, "rsa_key.pem"), ""); err != nil {
		t.Errorf("PKCS#8 RSA: %v", err)
	}
	// EC key.
	if _, err := parsePrivateKeyPEM(readFixture(t, "ec_key.pem"), ""); err != nil {
		t.Errorf("EC: %v", err)
	}
	// Encrypted PKCS#8 key with the correct password.
	if _, err := parsePrivateKeyPEM(readFixture(t, "rsa_pbes2.pem"), "secret"); err != nil {
		t.Errorf("encrypted with password: %v", err)
	}
	// Encrypted key with the wrong password fails.
	if _, err := parsePrivateKeyPEM(readFixture(t, "rsa_pbes2.pem"), "wrong"); err == nil {
		t.Error("expected error for wrong password")
	}
	// Non-PEM input fails.
	if _, err := parsePrivateKeyPEM([]byte("not a pem"), ""); err == nil {
		t.Error("expected error for non-PEM input")
	}
}

func TestParseCertificatePEM(t *testing.T) {
	if _, err := parseCertificatePEM(readFixture(t, "rsa_cert.pem")); err != nil {
		t.Errorf("valid cert: %v", err)
	}
	if _, err := parseCertificatePEM([]byte("not a pem")); err == nil {
		t.Error("expected error for non-PEM input")
	}
}

func TestBuildKeyMaps(t *testing.T) {
	types, ctx := newCryptoTestEnv()

	cert, err := parseCertificatePEM(readFixture(t, "rsa_cert.pem"))
	if err != nil {
		t.Fatalf("parseCertificatePEM: %v", err)
	}

	// buildCertMap + buildPublicKeyMap (with embedded certificate).
	rsaPub := cert.PublicKey.(*rsa.PublicKey)
	pubMap := buildPublicKeyMap(types, ctx, rsaPub, "RSA", cert)
	if pubMap == nil {
		t.Fatal("buildPublicKeyMap returned nil")
	}
	if _, ok := pubMap.Get("certificate"); !ok {
		t.Error("public key map missing certificate")
	}

	// buildPrivateKeyMap.
	key, err := parsePrivateKeyPEM(readFixture(t, "rsa_key.pem"), "")
	if err != nil {
		t.Fatalf("parsePrivateKeyPEM: %v", err)
	}
	if buildPrivateKeyMap(types, ctx, key.(*rsa.PrivateKey), "RSA") == nil {
		t.Error("buildPrivateKeyMap returned nil")
	}
}

func TestDecodeKeyStorePrivateKey(t *testing.T) {
	types, ctx := newCryptoTestEnv()
	data := readKeyStore(t)

	// RSA key recovered from the keystore.
	if isErr(decodeKeyStorePrivateKey(types, ctx, data, testKeyStorePass, "", testKeyStoreAlias, "RSA")) {
		t.Error("RSA keystore decode should succeed")
	}
	// Requesting an EC key from an RSA keystore is a type mismatch.
	if !isErr(decodeKeyStorePrivateKey(types, ctx, data, testKeyStorePass, "", testKeyStoreAlias, "EC")) {
		t.Error("EC request against RSA keystore should error")
	}
	// Wrong password fails recovery.
	if !isErr(decodeKeyStorePrivateKey(types, ctx, data, "bad", "bad", testKeyStoreAlias, "RSA")) {
		t.Error("wrong password should error")
	}
}

func TestDecodeTrustStorePublicKey(t *testing.T) {
	types, ctx := newCryptoTestEnv()
	data := readKeyStore(t)

	// RSA public key recovered from the trust store.
	if isErr(decodeTrustStorePublicKey(types, ctx, data, testKeyStoreAlias, testKeyStorePass, "RSA")) {
		t.Error("RSA trust-store decode should succeed")
	}
	// Requesting an EC key from an RSA cert is a type mismatch.
	if !isErr(decodeTrustStorePublicKey(types, ctx, data, testKeyStoreAlias, testKeyStorePass, "EC")) {
		t.Error("EC request against RSA cert should error")
	}
	// Wrong password fails.
	if !isErr(decodeTrustStorePublicKey(types, ctx, data, testKeyStoreAlias, "bad", "RSA")) {
		t.Error("wrong password should error")
	}
}

func isErr(v values.BalValue) bool {
	_, ok := v.(*values.Error)
	return ok
}

func storeMap(ctx *extern.Context, path, password string) *values.Map {
	atomic := semtypes.ToMappingAtomicType(ctx.TypeCtx, semtypes.MAPPING)
	return values.NewMap(semtypes.MAPPING, atomic, false, []values.MapEntry{
		{Key: "path", Value: path},
		{Key: "password", Value: password},
	})
}

func TestDecodeKeyStoreFromPath(t *testing.T) {
	types, ctx := newCryptoTestEnv()
	pal := testharness.NewTestPal().Platform()
	pal.FS.ReadFile = os.ReadFile
	rt := runtime.NewRuntime(pal, semtypes.CreateTypeEnv())

	// Keystore key recovered from the fixture path.
	ks := storeMap(ctx, testKeyStorePath, testKeyStorePass)
	if isErr(decodeKeyStoreKeyFromPath(rt, types, ctx, ks, testKeyStoreAlias, "", "RSA")) {
		t.Error("keystore decode from path should succeed")
	}
	// Trust store public key from the fixture path.
	ts := storeMap(ctx, testKeyStorePath, testKeyStorePass)
	if isErr(decodeTrustStoreKeyFromPath(rt, types, ctx, ts, testKeyStoreAlias, "RSA")) {
		t.Error("trust-store decode from path should succeed")
	}
	// A missing file yields an error.
	missing := storeMap(ctx, "testdata/does-not-exist.p12", testKeyStorePass)
	if !isErr(decodeKeyStoreKeyFromPath(rt, types, ctx, missing, testKeyStoreAlias, "", "RSA")) {
		t.Error("missing keystore should error")
	}
	if !isErr(decodeTrustStoreKeyFromPath(rt, types, ctx, missing, testKeyStoreAlias, "RSA")) {
		t.Error("missing trust store should error")
	}
}

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

func TestOaepHashForPadding(t *testing.T) {
	for _, padding := range []string{
		"OAEPwithMD5andMGF1",
		"OAEPWithSHA1AndMGF1",
		"OAEPWithSHA256AndMGF1",
		"OAEPwithSHA384andMGF1",
		"OAEPwithSHA512andMGF1",
	} {
		if oaepHashForPadding(padding) == nil {
			t.Errorf("oaepHashForPadding(%q) = nil", padding)
		}
	}
}

func TestHashValueToBytes(t *testing.T) {
	if got := hashValueToBytes("abc"); string(got) != "abc" {
		t.Errorf("string HashValue = %q, want %q", got, "abc")
	}
	_, ctx := newCryptoTestEnv()
	list := values.NewList(semtypes.LIST, semtypes.ToListAtomicType(ctx.TypeCtx, semtypes.LIST), false, nil, 0,
		[]values.BalValue{int64(1), int64(2), int64(3)})
	if got := hashValueToBytes(list); len(got) != 3 {
		t.Errorf("byte[] HashValue length = %d, want 3", len(got))
	}
}

func TestPbkdf2Params(t *testing.T) {
	for _, alg := range []string{"SHA1", "SHA256", "SHA512"} {
		if _, _, _, err := pbkdf2Params(alg); err != nil {
			t.Errorf("pbkdf2Params(%q): %v", alg, err)
		}
	}
	// SHA384 is not a supported PBKDF2 PRF, and unknown algorithms error too.
	if _, _, _, err := pbkdf2Params("SHA384"); err == nil {
		t.Error("expected error for unsupported SHA384")
	}
	if _, _, _, err := pbkdf2Params("UNKNOWN"); err == nil {
		t.Error("expected error for unknown algorithm")
	}
}

func TestPkcs7Unpad(t *testing.T) {
	// "AB" padded to an 8-byte block with 0x06 repeated six times.
	padded := []byte{'A', 'B', 6, 6, 6, 6, 6, 6}
	got, err := pkcs7Unpad(padded)
	if err != nil {
		t.Fatalf("pkcs7Unpad: %v", err)
	}
	if string(got) != "AB" {
		t.Errorf("pkcs7Unpad = %q, want %q", got, "AB")
	}
	// A pad length exceeding the block is invalid.
	if _, err := pkcs7Unpad([]byte{1, 2, 9}); err == nil {
		t.Error("expected error for invalid padding length")
	}
	// Empty input is invalid.
	if _, err := pkcs7Unpad(nil); err == nil {
		t.Error("expected error for empty input")
	}
}

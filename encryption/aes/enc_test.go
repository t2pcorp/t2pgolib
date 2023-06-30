//
/*
 * Created by Visual Studio Code.
 * User: HM
 * Date: 06/07/2018
 * Time: 12:36
 */

package aes

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {

	sharedKey, _ := base64.StdEncoding.DecodeString("7Zcf2SnEoMn16DTbV3JbFqGcgj2CcV8U+snPg87BZAw=")

	//Encrypt
	encrypted, err := EncryptAES256GCMHex(`Hello ทดสอบ 1234 `, sharedKey)
	if err != nil {
		t.Errorf(fmt.Sprintf("Cannot Encrypt : %s.", err))
	}
	fmt.Printf("Encrypted: %s\n", encrypted)

	//Decrypt
	decrypted, err := DecryptAES256GCMHex(encrypted, sharedKey)
	if err != nil {
		t.Errorf(fmt.Sprintf("Cannot Decrypt : %s.", err))
	}
	fmt.Printf("Decrypted: %s\n", decrypted)
}

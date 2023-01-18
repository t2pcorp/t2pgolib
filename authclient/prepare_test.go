package authclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateTokenTypeH(t *testing.T) {

	// Prepare Request Token (Client)
	locat, _ := time.LoadLocation("Asia/Bangkok")
	currentunix := time.Now().In(locat)
	bodyStr := `{"refcode": "CLIENIT00001002"}`

	// Generate token type H
	header, requestBody, err := GenerateTokenTypeH(bodyStr, clientKeyStr003, currentunix, true)
	fmt.Println("Type H Header:\n", header)
	fmt.Println()
	fmt.Println("Type H requestBody:\n", requestBody)
	fmt.Println()
	fmt.Println("err:", err)

	// Test Decrypt
	fmt.Println()
	decrypted, err := Decrypt(requestBody, serverKeyStr003)
	fmt.Println("decrypted:", err, decrypted)
}

func Test_GenerateTokenTypeC(t *testing.T) {

	locat, _ := time.LoadLocation("Asia/Bangkok")
	currentunix := time.Now().In(locat)
	bodyStr := `{"refcodeOfToken": "CLIENIT00001002"}`
	tokenUrl := `https://test-api-authen.t2p.co.th/authen/v1/clientToken/generate`

	tokenTypeC, err := RequestTokenTypeC(bodyStr, clientKeyStr003, currentunix, tokenUrl)

	if err != nil {
		fmt.Println("request token error:", err)
	}

	///////////////////////////////////////////////////////////////////////////////////////////////
	// Test Use Token (Use Token on T2P resource server)
	tokenUrl = `https://test-api-authen.t2p.co.th/authen/v1/clientToken/testClientToken`

	jsonStr := []byte(`{"data":"Your REQUEST DATA From Mobile Client"}`)
	req, _ := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %v", tokenTypeC))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Type C token is:", tokenTypeC)
	fmt.Println()
	fmt.Println("Test verify response Status:", resp.Status)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Test verify response Body:\n", string(body))

}

func Test_Hmac(t *testing.T) {

	plainText := `This is plain Text {"Test":"Encrypt"} ทดสอบ`
	hmac, err := GenerateHMac(plainText, clientKeyStr003)
	fmt.Println("Hmac: ", hmac, err)

	verify, err := VerifyHMac(plainText, hmac, clientKeyStr003)
	fmt.Println("Verify: ", verify, err)

	if !assert.Equal(t, verify, true) {
		panic("verify Failed !!!!")
	}
}

func Test_Encrypt_Decrypt(t *testing.T) {
	plainText := `This is plain Text {"Test":"Encrypt"} ทดสอบ`
	encryptedText, err := Encrypt(plainText, clientKeyStr003)
	if err != nil {
		fmt.Println(err)
		return
	}
	decryptedText, err := Decrypt(encryptedText, serverKeyStr003)
	if err != nil {
		panic(err)
	}
	fmt.Println("Original:", plainText)
	fmt.Println("Encrypted:", encryptedText)
	fmt.Println("Decrypted:", decryptedText)
	if !assert.Equal(t, plainText, decryptedText) {
		panic("Encrypt / Decrypt Failed !!!!")
	}

	fromStringPHPEncrypted := `DZAKYhl9Q0f2s3ezsDAlnl1HTT9wMhZbbUo4fqQYxCPNXM7RoD1hkzniVJV/f5T6Y126A9qNgZypKdSCd/6M0d8pxMP+isuYUUCFaMH3IH+8KicVpec7GIPi3da3bw7yLfSG8OM5ZKYKDRVwI1Ze9dzexPGgn+M6Ksrfu4uZqmcnKCGoO8YdXn00JkksiAD5/UuAyCTD8IyU8Skh6sVEeCfvMB/n0mC5GLclfMDvi0bXeoGt3X+T6EBmDW2DGCXzsTRuWf/2CzoSNXz+ANRRM4UHH4nfy5XS/j1JnMipYSP1vOHYvBb4G3BlgK9947JcK2gUYba4/kKvohkLA1KVyg==:aXOmWWRZaJ/1NGOVLHivaIndtTyfBxuS/zDAFOfHEGjs4/318ocYmBBmsjbna+ZNz6HRVXJZfTT8F7/jHSHXMw==`

	decryptedTextFromPHP, err := Decrypt(fromStringPHPEncrypted, serverKeyStr003)
	if err != nil {
		panic(err)
	}

	fmt.Println("DecryptedFromPHP:", decryptedTextFromPHP)
	if !assert.Equal(t, plainText, decryptedTextFromPHP) {
		panic("Decrypt From PHP Failed !!!!")
	}

	fromStringNodeJsEncrypted := `fYOyL++Dx9L0oQNaumchdmFeyVdyGiGEOpHUmA0vY7wllVC3e438rJXc3pUTV3lRvMNhjm/pX2sXSvTxXNk+m63BHbMQlL/YcPAAEwCXPcefWqDKSotTiOPLSUf34QFBHVpwgjkQYZ85OyPO/GaPw3vHjDGYmNcDTFrXix8VL4qOpd2zGJPi0nUVTeE7bVUEuF+BN1suBlTse4Znh+yXfGy5jOJvTUw4FbLw5dhN0BOFtzgFHEzu2arOGJLIEZ6qtOTLikbjEo3TGlaAQYwIOT6MpLs1yrYn3NF8i7F7Vc/faGcc/luy5KgLsUUykHhfF1tP6aA+wqYkpKHY++WA2g==:ztV1JwFogHj/4g7pfL0kBT5czD5sCEA0VT5+a1PyTSCYlTNbISnRyEejNCBpr6d2/UXazKxEj0/zcvUkvWQnWA==`

	decryptedTextFromNodeJS, err := Decrypt(fromStringNodeJsEncrypted, serverKeyStr003)
	if err != nil {
		panic(err)
	}

	fmt.Println("DecryptedFromNodeJs:", decryptedTextFromNodeJS)
	if !assert.Equal(t, plainText, decryptedTextFromNodeJS) {
		panic("Decrypt From NodeJs Failed !!!!")
	}

}

func Test_Decrypt_Cipher(t *testing.T) {

}

func Test_Prepare_Request(t *testing.T) {
	testCases := []struct {
		name          string
		hashInfo      map[string]string
		body          string
		key           string
		isEncryptBody bool
		expected      string
	}{
		{
			name:          "should error when hashInfo have no method, url or timestamp",
			hashInfo:      map[string]string{},
			body:          `{"title":"xxx"}`,
			isEncryptBody: false,
			expected:      `{"meta":{"responseCode":1666,"responseMessage":"invalid header info data","language":"en_EN","version":"1.0.0"},"data":{"header":"","body":""}}`,
		},
		{
			name:          "should error when hashInfo['timestamp'] is invalid",
			hashInfo:      map[string]string{"method": "POST", "uri": "/server-authen.php", "timestamp": "2019060812545"},
			body:          `{"title":"xxx"}`,
			isEncryptBody: false,
			expected:      `{"meta":{"responseCode":1666,"responseMessage":"invalid header info timestamp","language":"en_EN","version":"1.0.0"},"data":{"header":"","body":""}}`,
		},
		{
			name:          "should success when is not encrypt body",
			hashInfo:      map[string]string{"method": "POST", "uri": "/server-authen.php", "timestamp": "20190608125456"},
			body:          `{"title":"xxx"}`,
			isEncryptBody: false,
			expected:      `{"meta":{"responseCode":1000,"responseMessage":"Success","language":"en_EN","version":"1.0.0"},"data":{"header":"ZXlKamJHbGxiblJEYjJSbElqb2lRMHhKTVRBd01EQXpJaXdpWTJ4cFpXNTBUR2xpVm1WeWMybHZiaUk2SWpFdU1DNHdJaXdpYTJWNVEyOWtaU0k2SWt0RE1UQXdNREF5TUNJc0ltMWxkR2h2WkNJNklsQlBVMVFpTENKMGFXMWxjM1JoYlhBaU9pSXlNREU1TURZd09ERXlOVFExTmlJc0luVnlhU0k2SWk5elpYSjJaWEl0WVhWMGFHVnVMbkJvY0NKOTphZmNhYTU4ZmU2ZWI1Mjc5ZDRlMWU2Mjc4NjMxYTJiZWNlZGExYzE2NGJlMWQ4MTZhYTE3Y2YxMDg1ZTIzNmI5YjJjZjU3ZjQxNjc2NWEyZGQxYjM0MWIzMzExZDhiZWQ0NjhlMzdkOTY3YjY0MWIzZmQ2ODVkMGJjNGZhNGZjMQ==","body":"{\"title\":\"xxx\"}"}}`,
		},
		{
			name:          "should success when is encrypt body and hashInfo['tokenType'] == 'C'",
			hashInfo:      map[string]string{"method": "POST", "uri": "/server-authen.php", "timestamp": "20190608125456", "tokenType": "C"},
			body:          `{"title":"xxx"}`,
			isEncryptBody: true,
			expected:      `{"meta":{"responseCode":1000,"responseMessage":"Success","language":"en_EN","version":"1.0.0"},"data":{"header":"ZXlKamJHbGxiblJEYjJSbElqb2lRMHhKTVRBd01EQXpJaXdpWTJ4cFpXNTBUR2xpVm1WeWMybHZiaUk2SWpFdU1DNHdJaXdpYTJWNVEyOWtaU0k2SWt0RE1UQXdNREF5TUNJc0ltMWxkR2h2WkNJNklsQlBVMVFpTENKMGFXMWxjM1JoYlhBaU9pSXlNREU1TURZd09ERXlOVFExTmlJc0luUnZhMlZ1Vkhsd1pTSTZJa01pTENKMWNta2lPaUl2YzJWeWRtVnlMV0YxZEdobGJpNXdhSEFpZlE9PTozYjkxNDI1NzljNGNiMDNlNjhlMzBkMzNmYmY0MGJjOWMwMzFmM2JmOGRmZGJkNzFlZGY2NTlhN2U2MDM4ZGYwZDZiMDcyNmE0ZGFjODhjODFhN2M4ZmZmMmYyM2Q3Y2UwMTg0YWE0N2M1NGY5MzViNzhlZTQ5YTljZTYwYWNmNg==","body":"{\"title\":\"xxx\"}"}}`,
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			r := PrepareRequest(v.hashInfo, v.body, clientKeyStr003, v.isEncryptBody)
			fmt.Println(r)
			assert.Equal(t, v.expected, r)
		})
	}
}

func Test_Prepare_Request_Should_Return_Correct_Header_And_Body_When_Not_Encrypt_Body(t *testing.T) {
	hashInfo := map[string]string{"method": "POST", "uri": "/server-authen.php", "timestamp": "20190608125456", "token": "H"}
	r := PrepareRequest(hashInfo, `{"title":"Testing Client-Server Verification"}`, clientKeyStr003, false)

	expected := `{"meta":{"responseCode":1000,"responseMessage":"Success","language":"en_EN","version":"1.0.0"},"data":{"header":"ZXlKamJHbGxiblJEYjJSbElqb2lRMHhKTVRBd01EQXpJaXdpWTJ4cFpXNTBUR2xpVm1WeWMybHZiaUk2SWpFdU1DNHdJaXdpYTJWNVEyOWtaU0k2SWt0RE1UQXdNREF5TUNJc0ltMWxkR2h2WkNJNklsQlBVMVFpTENKMGFXMWxjM1JoYlhBaU9pSXlNREU1TURZd09ERXlOVFExTmlJc0luUnZhMlZ1SWpvaVNDSXNJblZ5YVNJNklpOXpaWEoyWlhJdFlYVjBhR1Z1TG5Cb2NDSjk6MWQ2NDA3MmQ4NTU3NjEzY2I1M2Y0MjQwNzg0NzMyNTRiNTAxZmJmMTM2MWRiOGQ2MjdmODc2NGEyZTVhMGQyM2ZmMzM0NDcyNmI0Y2MzYzVjMTg3NmE1OGRkYjg5YzlmYzk4OGJlOTIwMjVmNzlhNTMwOWM1MzA2ZjRlMjU4Y2U=","body":"{\"title\":\"Testing Client-Server Verification\"}"}}`
	assert.Equal(t, expected, r)
}

func Test_Function_HSGetToken(t *testing.T) {
	///////////////////////////////////////////////////////////////////////////////////////////////
	//Client Prepare Request
	// Prepare Request Token (Client)
	locat, _ := time.LoadLocation("Asia/Bangkok")
	currentunix := time.Now().In(locat)

	hashMap := make(map[string]map[string]string)
	hashMapServerInfo := make(map[string]string)
	// json.Unmarshal(c.PostBody(), &hashMap)
	hashMapServerInfo["timestamp"] = currentunix.Format(`20060102150405`)
	hashMapServerInfo["tokenType"] = "HS"
	hashMapServerInfo["method"] = "POST"
	hashMapServerInfo["uri"] = "/authen/v1/clientToken/generateNew"
	hashMap["ServerInfo"] = hashMapServerInfo
	bodyStr := "CLIENIT00001002"
	requestInfo := PrepareRequest(hashMap["ServerInfo"], bodyStr, clientKeyStr003, false)
	fmt.Println(requestInfo)
	/*
		#Provide Partner
			- clientCode
			- keyCode
			- hashmackey

		#PREPARE LOGIC
		hashInfo["uri"] = "/authen/v1/clientToken/generateNew"
		hashInfo["timestamp"] = "Current Time UTC+7"
		hashInfo["clientCode"] = clientCode
		hashInfo["keyCode"] = keyCode
		hashInfo["clientLibVersion"] = "1.0.0"

		#sort Hash by Key ASC
		hashInfo=sortKey(hashInfo)

		#Concat Value of Hash and body String
		rawString = concatString(valueOf(hashInfo)+bodyString)

		#create Hash Signature
		hashValue = hashMac512(rawString, byte(hashmackey))

		#build Header Token
		headerToken = base64_encode(json_encode(hashInfo))+':'+hashValue
		FinalHeaderToken = base64_encode(headerToken)

		#use
		- FinalHeaderToken
		- bodyString (same value using with create header token)
	*/
	reqObj := make(map[string]map[string]string)
	json.Unmarshal([]byte(requestInfo), &reqObj)

	///////////////////////////////////////////////////////////////////////////////////////////////
	// Gernerate Client Token (Call API)
	tokenUrl := `https://test-api-authen.t2p.co.th/authen/v1/clientToken/generate`

	var jsonStr = []byte(reqObj["data"]["body"])
	req, _ := http.NewRequest("POST", tokenUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %v", reqObj["data"]["header"]))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	///////////////////////////////////////////////////////////////////////////////////////////////
	// Verify Token (Use Token)
	tokenUrl = `https://test-api-authen.t2p.co.th/authen/v1/clientToken/testClientToken`

	reqObj = make(map[string]map[string]string)
	json.Unmarshal(body, &reqObj)

	jsonStr = []byte(`{"data":"ANY REQUEST DATA"}`)
	req, _ = http.NewRequest("POST", tokenUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %v", reqObj["data"]["authToken"]))
	req.Header.Set("Content-Type", "application/json")

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ = io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

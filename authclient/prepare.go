package authclient

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type meta struct {
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Language        string `json:"language"`
	Version         string `json:"version"`
}

type clientKey struct {
	clientCode        string
	keyCode           string
	hmacKey           string
	publicKey         string
	privateKey        string
	keyAuthorize      map[string][]string
	requireInfo       map[string]string
	tokenTimeout      map[string]int
	hostTokenTimezone string
}

type header struct {
	info map[string]string
	// method     string
	// timestamp  string
	// clientCode string
	// keyCode    string
	hmac string
}

type response struct {
	Meta meta `json:"meta"`
	Data data `json:"data"`
}

type data struct {
	Header     *string             `json:"header,omitempty"`
	AuthenInfo map[string]string   `json:"authenInfo,omitempty"`
	ClientInfo *string             `json:"clientInfo,omitempty"`
	Body       string              `json:"body"`
	Authorize  map[string][]string `json:"authorize,omitempty"`
}

// KeyInfor key infor to redis
type KeyInfor struct {
	KeyContents       string              `json:"keyContents"`
	KeyAuthorize      map[string][]string `json:"keyAuthorize"`
	RequireInfo       map[string]string   `json:"requireInfo"`
	TokenTimeout      map[string]int      `json:"tokenTimeout"`
	HostTokenTimeZone string              `json:"hostTokenTimeZone"`
}

type responseEncrypt struct {
	Meta meta        `json:"meta"`
	Data interface{} `json:"data"`
}

// LIBVERSION is version of library
const (
	LIBVERSION               = "1.0.0"
	ErrFatal                 = "fatal error"
	ErrInvalidHeaderSignatur = "invalid header signature"
)

// GenerateTokenTypeH use for Host to Host request short life time 60 secconds
// params: requestData type: string description: string or json string to generate token type H with request data
// params: key type: string description: key content get from T2P
// params: dtNowWithTimezone type: time.Time description: current date time with server time zone eg. "Asia/Bangkok"
// params: isRequestDataEncrypt type: bool description: will requestData be encrypted or plaintext for secure use encrypt
func GenerateTokenTypeH(requestData string, key string, dtNowWithTimezone time.Time, isRequestDataEncrypt bool) (string, string, error) {
	hashMapServerInfo := make(map[string]string)
	hashMapServerInfo["timestamp"] = dtNowWithTimezone.Format(`20060102150405`)
	hashMapServerInfo["tokenType"] = "H"
	hashMapServerInfo["method"] = "POST"
	hashMapServerInfo["uri"] = "/authen/v1/clientToken/generate"
	hashMapServerInfo["clientLibVersion"] = "1.0.0"

	out := PrepareRequest(hashMapServerInfo, requestData, key, isRequestDataEncrypt)
	reqObj := make(map[string]map[string]interface{})
	json.Unmarshal([]byte(out), &reqObj)

	if fmt.Sprintf("%v", reqObj["meta"]["responseCode"]) != "1000" {
		return ``, ``, fmt.Errorf("%v", reqObj["meta"]["responseMessage"])
	} else {
		return fmt.Sprintf("%v", reqObj["data"]["header"]), fmt.Sprintf("%v", reqObj["data"]["body"]), nil
	}
}

// RequestTokenTypeC use for Host to Host request short life time 60 secconds
// params: tokenReferenceData type: string description: short reference for tracing if needed
// params: key type: string description: key content get from T2P
// params: dtNowWithTimezone type: time.Time description: current date time with server time zone eg. "Asia/Bangkok"
// params: authUrl type: string description: T2P authen server URL (UAT) `https://test-api-authen.t2p.co.th/authen/v1/clientToken/generate`, (PROD) `https://api-authen.t2p.co.th/authen/v1/clientToken/generate`
func RequestTokenTypeC(tokenReferenceData string, key string, dtNowWithTimezone time.Time, authUrl string) (string, error) {
	header, requestBody, err := GenerateTokenTypeH(tokenReferenceData, key, dtNowWithTimezone, true)
	if err != nil {
		return ``, err
	}

	client := resty.New()
	resp, err := client.SetTimeout(5*time.Second).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Basic %v", header)).
		SetBody([]byte(requestBody)).
		Post(authUrl)

	if resp.StatusCode() != 200 {
		return ``, fmt.Errorf("request for client token failed: %v", err.Error())
	}
	body := string(resp.Body())

	reqObj := make(map[string]map[string]interface{})
	json.Unmarshal([]byte(body), &reqObj)

	if fmt.Sprintf("%v", reqObj["meta"]["responseCode"]) != "1000" {
		return ``, fmt.Errorf("%v", reqObj["meta"]["responseMessage"])
	} else {
		return fmt.Sprintf("%v", reqObj["data"]["authToken"]), nil
	}
}

// PrepareRequest is a function for preparing requests to call APIs
func PrepareRequest(hashInfo map[string]string, body, key string, isEncryptBody bool) string {
	ck, err := extractKey(key)
	if err != nil {
		return makeResponseErrorPrepare(1888, err.Error())
	}

	body, err = makeRequestBody(hashInfo, body, isEncryptBody, ck.publicKey)
	if err != nil {
		return makeResponseErrorPrepare(1888, ErrFatal+":"+err.Error())
	}

	h, err := createHeaderSignature(hashInfo, body, ck)
	if err != nil {
		return makeResponseErrorPrepare(1666, err.Error())
	}
	return makeResponse(1000, "Success", &h, nil, nil, body, nil)
}

func makeResponseErrorPrepare(code int, err string) string {
	var header = ""
	return makeResponse(code, err, &header, nil, nil, "", nil)
}

func makeRequestBody(hashInfo map[string]string, body string, isEncryptBody bool, publicKey string) (string, error) {
	if hashInfo["tokenType"] == "C" || !isEncryptBody {
		return body, nil
	}

	body, err := encryptMessage(body, publicKey)
	if err != nil {
		return "", err
	}

	return body, nil
}

// EncryptData encrypt text with key input
func EncryptData(text string, key string) string {
	ck, err := extractKey(key)
	if err != nil {
		return makeResponseErrorPrepare(1888, err.Error())
	}
	encText, err := encryptMessage(text, ck.publicKey)
	if err != nil {
		return makeResponseErrorPrepare(1888, ErrFatal+":"+err.Error())
	}
	return makeResponseEncryptData(1000, "Success", encText)
}

// DecryptData encrypt text with key input
func DecryptData(text string, key string) string {
	ck, err := extractKey(key)
	if err != nil {
		return makeResponseErrorPrepare(1888, err.Error())
	}
	decText, err := decryptMessage(ck, text)
	if err != nil {
		return makeResponseErrorPrepare(1888, ErrFatal+":"+err.Error())
	}
	return makeResponseEncryptData(1000, "Success", decText)
}

func decryptMessage(ck clientKey, encMessage string) (string, error) {
	ek, i, err := getDecryptMessageInfo(encMessage)
	if err != nil {
		return "", errors.New("GetDecryptInfo:" + err.Error())
	}

	key, iv, err := decryptAESKey(ck.privateKey, ek)
	if err != nil {
		return "", errors.New("DecryptKeys:" + err.Error())
	}

	decoded, err := decryptAES(key, iv, i)
	if err != nil {
		return "", errors.New("DecryptMessage:" + err.Error())
	}

	return string(decoded), nil
}

func encryptMessage(body, publicKey string) (string, error) {
	output, err := encryptAES([]byte(body))
	if err != nil {
		return "", err
	}

	keyEncrypted, err := encryptAESKey(publicKey, output.Key, output.IV)
	if err != nil {
		return "", err
	}

	keyEncryptedEnc := base64.StdEncoding.EncodeToString(keyEncrypted)

	return fmt.Sprintf("%s:%s", keyEncryptedEnc, output.Output), nil
}

func createHeaderSignature(hashInfo map[string]string, body string, ck clientKey) (string, error) {

	_, ok1 := hashInfo["method"]
	timestamp, ok3 := hashInfo["timestamp"]

	if !(ok1 && ok3) {
		return "", errors.New("invalid header info data")
	}

	if len(timestamp) < 14 {
		return "", errors.New("invalid header info timestamp")
	}

	hashInfo["clientCode"] = strings.TrimSpace(ck.clientCode)
	hashInfo["keyCode"] = strings.TrimSpace(ck.keyCode)
	hashInfo["clientLibVersion"] = LIBVERSION

	raw := sortAndImplodeInfo(hashInfo)

	// raw := joinMap(hashInfo)
	hashMac, err := hashHMAC(raw+body, ck.hmacKey)
	if err != nil {
		return "", errors.New("can not hash hmac")
	}

	hashInfoB, err := json.Marshal(hashInfo)
	if err != nil {
		return "", err
	}

	hashInfoStr := base64.StdEncoding.EncodeToString(hashInfoB)

	header := hashInfoStr + ":" + hashMac

	header = base64.StdEncoding.EncodeToString([]byte(header))
	return header, nil
}

func hashHMAC(raw, key string) (string, error) {
	k, err := base64.StdEncoding.DecodeString(strings.Trim(key, " "))
	if err != nil {
		return "", err
	}

	h := hmac.New(sha512.New, k)
	h.Write([]byte(raw))

	return hex.EncodeToString(h.Sum(nil)), nil
}

func sortAndImplodeInfo(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	info := ""
	for _, k := range keys {
		info += m[k]
	}
	return info
}

func sortMapByKey(m map[string]string) map[string]string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	newMap := make(map[string]string, len(m))
	for _, k := range keys {
		newMap[k] = m[k]
	}
	return newMap
}

func joinMap(m map[string]string) string {
	strs := make([]string, 0, len(m))

	for _, v := range m {
		strs = append(strs, v)
	}

	return strings.Join(strs, "")
}

func extractKey(k string) (clientKey, error) {
	ks := strings.Split(k, ":")
	if len(ks) != 5 {
		return clientKey{}, errors.New("invalid key")
	}
	return clientKey{
		clientCode: ks[0],
		keyCode:    ks[1],
		hmacKey:    ks[2],
		publicKey:  ks[3],
		privateKey: ks[4],
	}, nil
}

func makeResponse(code int, message string, header *string, authInfo map[string]string, clientInfo *string, body string, authorize map[string][]string) string {
	b, _ := json.Marshal(response{
		Meta: meta{
			ResponseCode:    code,
			ResponseMessage: message,
			Language:        "en_EN",
			Version:         LIBVERSION,
		},
		Data: data{
			Header:     header,
			AuthenInfo: authInfo,
			ClientInfo: clientInfo,
			Body:       body,
			Authorize:  authorize,
		},
	})

	return string(b)
}

func makeResponseEncryptData(code int, message string, data interface{}) string {
	b, _ := json.Marshal(responseEncrypt{
		Meta: meta{
			ResponseCode:    code,
			ResponseMessage: message,
			Language:        "en_EN",
			Version:         LIBVERSION,
		},
		Data: data,
	})

	return string(b)
}

func GenerateHMac(message string, key string) (string, error) {
	ck, err := extractKey(key)
	if err != nil {
		return "", err
	}
	return hashHMAC(message, ck.hmacKey)
}

func VerifyHMac(message string, hMac string, key string) (bool, error) {
	ck, err := extractKey(key)
	if err != nil {
		return false, err
	}
	newHMac, err := hashHMAC(message, ck.hmacKey)
	if err != nil {
		return false, err
	}
	return newHMac == hMac, nil
}

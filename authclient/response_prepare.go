package authclient

//ResponsePreparePayload main payload
type ResponsePreparePayload struct {
	Meta MetaPrepare `json:"meta"`
	Data DataPrepare `json:"data"`
}

//MetaPrepare meta part
type MetaPrepare struct {
	Language        string `json:"language"`
	ResponseCode    int    `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Version         string `json:"version"`
}

//DataPrepare data part
type DataPrepare struct {
	Body   string `json:"body"`
	Header string `json:"header"`
}

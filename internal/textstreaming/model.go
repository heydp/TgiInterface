package textstreaming

type Request struct {
	Text string `json:"text"`
}

type Response struct {
	Result string `json:"result"`
}

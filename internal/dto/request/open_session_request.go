package request

type OpensessionRequest struct {
	SendId    string `json:"send_id"`
	ReceiveId string `json:"receive_id"`
}

package plc_mitsubishi

type ResponseBinder struct {
	CommonMessageFormat
	ResDataLen int
	Payload string
}

func parseResponse(msg []byte) (*ResponseBinder,error){
	return &ResponseBinder{

	},nil
}

func GetMsg(msg []byte) (string,error) {
	res,err := parseResponse(msg)
	if err != nil {
		return "",err
	}
	return res.Payload,nil
}
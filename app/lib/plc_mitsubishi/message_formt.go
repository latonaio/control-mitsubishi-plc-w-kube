package plc_mitsubishi

type CommonMessageFormat struct {
	Subheader   string
	AccessRoute *AccessRoute
}

type AccessRoute struct {
	NetworkNo string
	PcNo      string
	IONo      string
	StationNo string
}

func newAccessRoute() *AccessRoute {
	return &AccessRoute{
		NetworkNo: NETWORK_BO_TO_HOST_STATION,
		PcNo:      PC_NO_TO_HOST_STATION,
		IONo:      TARGET_IO_UNIT_NO,
		StationNo: "00",
	}
}

func newFmt() *CommonMessageFormat {
	return &CommonMessageFormat{
		Subheader: REQUEST_SUBHEADER,
		AccessRoute: newAccessRoute(),
	}
}

package plc_mitsubishi

const (
	REQUEST_SUBHEADER  = "5000"
	RESPONSE_SUBHEADER = "D000"

	END_CODE = "0000"
)

const (
	HOST_STATION_NO_           = "00"
	NETWORK_BO_TO_HOST_STATION = "00"
	PC_NO_TO_HOST_STATION      = "FF"
	TARGET_IO_UNIT_NO          = "03FF"
)

const (
	BULK_READ_CMD  = "0401"
	BULK_WRITE_CMD = "1401"
	SUB_CMD        = "0000"
)

const (
	INPUT = iota
	OUTPUT
)


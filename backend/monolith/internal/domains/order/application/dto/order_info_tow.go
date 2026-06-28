package dto

type OrderParams struct {
	CarParams CarParams
}

type CarParams struct {
	CarWeight           uint8
	WheelsBlockedNumber uint8
}

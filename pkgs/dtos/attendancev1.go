package dtos

type AttendanceV1Response struct {
	AttendanceID int64  `json:"attendance_id"`
	PositionID   int64  `json:"position_id"`
	ClockInTime  string `json:"clock_in_time"`
	ClockOutTime string `json:"clock_out_time"`
}

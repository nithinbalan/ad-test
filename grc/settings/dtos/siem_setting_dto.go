package dtos

type SIEMSettingDTO struct {
	SyslogServer         string `json:"syslog_server" binding:"required"`   // e.g., 192.168.1.100
	SyslogPort           int    `json:"syslog_port" binding:"required"`     // e.g., 514
	SyslogProtocol       string `json:"syslog_protocol" binding:"required"` // e.g., udp, tcp
	SyslogTimezone       string `json:"syslog_timezone"`                    // Optional timezone info
	SyslogDatetimeFormat string `json:"syslog_datetime_format"`             // Optional datetime format

}

package log_monitoring

type Log_monitoring struct {
	ID_AUDIT     int64  `json:"id_audit"`
	MENU         string `json:"menu"`
	ACTION       string `json:"action"`
	OLD_VALUE    string `json:"old_value"`
	NEW_VALUE    string `json:"new_value"`
	RESPONSE_MSG string `json:"response_msg"`
	CHANGED_BY   string `json:"changed_by"`
	CHANGED_DATE string `json:"changed_date"`
	IP_CLIENT    string `json:"ip_client"`
	USER_AGENT   string `json:"user_agent"`
}

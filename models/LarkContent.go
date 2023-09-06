package models

type AlertContent struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Alerts   []struct {
		Status string `json:"status"`
		Labels struct {
			Alertname  string `json:"alertname"`
			Alertype   string `json:"alertype"`
			Device     string `json:"device"`
			Exp        string `json:"exp"`
			Fstype     string `json:"fstype"`
			Group      string `json:"group"`
			Iid        string `json:"iid"`
			Instance   string `json:"instance"`
			Job        string `json:"job"`
			Level      string `json:"level"`
			Mountpoint string `json:"mountpoint"`
			Name       string `json:"name"`
			Region     string `json:"region"`
			Vendor     string `json:"vendor"`
			Severity   string `json:"severity"`
			AlertLevel string `json:"alert_level"`
			Id         int    `json:"id"`
		} `json:"labels"`
		Annotations struct {
			Description string `json:"description"`
			Summary     string `json:"summary"`
		} `json:"annotations"`
		StartsAt     string `json:"startsAt"`
		EndsAt       string `json:"endsAt"`
		Duration     string `json:"duration"`
		GeneratorURL string `json:"generatorURL"`
		Fingerprint  string `json:"fingerprint"`
	} `json:"alerts"`
	GroupLabels struct {
		Alertname string `json:"alertname"`
		Group     string `json:"group"`
		Instance  string `json:"instance"`
		Level     string `json:"level"`
	} `json:"groupLabels"`
	CommonAnnotations struct {
		Description string `json:"description"`
	} `json:"commonAnnotations"`
	ExternalURL     string `json:"externalURL"`
	Version         string `json:"version"`
	GroupKey        string `json:"groupKey"`
	TruncatedAlerts int    `json:"truncatedAlerts"`
}

type SSLContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
type SSLContentText struct {
	Domain     string `json:"domain"`
	RootDomain string `json:"root_domain"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	ExpireTime int    `json:"expire_time"`
}

type MarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

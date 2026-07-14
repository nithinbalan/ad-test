package dtos

type QualysKBDetailsDTO struct {
	QID         string               `json:"qid"`
	Title       string               `json:"title"`
	Diagnosis   string               `json:"diagnosis"`
	Category    string               `json:"category"`
	Consequence string               `json:"consequence"`
	Solution    string               `json:"solution"`
	CVSS        string               `json:"cvss"`
	CVEList     []string             `json:"cve_list"`
	Detections  []QualysDetectionDTO `json:"detections"`
}

type QualysKBResponse struct {
	Response struct {
		Vulns []struct {
			QID         string `xml:"QID"`
			Title       string `xml:"TITLE"`
			Diagnosis   string `xml:"DIAGNOSIS"`
			Category    string `xml:"CATEGORY"`
			Consequence string `xml:"CONSEQUENCE"`
			Solution    string `xml:"SOLUTION"`
			CVSS3       string `xml:"CVSS_V3"`
			CVEList     []struct {
				ID  string `xml:"ID"`
				URL string `xml:"URL"`
			} `xml:"CVE_LIST>CVE"`
		} `xml:"VULN_LIST>VULN"`
	} `xml:"RESPONSE"`
}

type QualysDetectionDTO struct {
	HostID           string   `json:"host_id"`
	IP               string   `json:"ip"`
	IPv6             string   `json:"ipv6"`
	DNS              string   `json:"dns"`
	QID              string   `json:"qid"`
	Severity         string   `json:"severity"`
	Status           string   `json:"status"`
	QDS              string   `json:"qds"`
	CriticalityScore string   `json:"criticality_score"`
	Port             string   `json:"port"`
	Protocol         string   `json:"protocol"`
	LastFound        string   `json:"last_found"`
	Tags             []string `json:"tags"`
}

type QualysDetectionResponse struct {
	Response struct {
		HostList struct {
			Hosts []struct {
				ID         string `xml:"ID"`
				IP         string `xml:"IP"`
				DNS        string `xml:"DNS"`
				IPv6       string `xml:"IPV6"`
				Detections []struct {
					QID      string `xml:"QID"`
					Severity string `xml:"SEVERITY"`
					Status   string `xml:"STATUS"`
					QDS      struct {
						Value    string `xml:",chardata"`
						Severity string `xml:"severity,attr"`
					} `xml:"QDS"`
					CriticalityScore string   `xml:"CRITICALITY_SCORE"`
					Port             string   `xml:"PORT"`
					Protocol         string   `xml:"PROTOCOL"`
					LastFound        string   `xml:"LAST_FOUND_DATETIME"`
					Tags             []string `xml:"TAG_LIST>TAG"`
				} `xml:"DETECTION_LIST>DETECTION"`
			} `xml:"HOST"`
		} `xml:"HOST_LIST"`
	} `xml:"RESPONSE"`
}

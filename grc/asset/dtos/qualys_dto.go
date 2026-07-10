package dtos

type QualysKBDetailsDTO struct {
	QID         string   `json:"qid"`
	Title       string   `json:"title"`
	Diagnosis   string   `json:"diagnosis"`
	Category    string   `json:"category"`
	Consequence string   `json:"consequence"`
	Solution    string   `json:"solution"`
	CVSS        string   `json:"cvss"`
	CVEList     []string `json:"cve_list"`
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

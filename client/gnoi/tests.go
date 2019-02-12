package gnoi

var RebootTests = []struct {
	Desc    string
	Message string
}{{
	Desc:    "reboot system",
	Message: "Testing OVS Reboot Functionality",
}}

var GetCertificatesTests = []struct {
	Desc        string
	ExpCertID   string
	ExpCertPath string
}{{
	Desc:        "get certificates",
	ExpCertID:   "c5e5a1cb-8e1f-43c1-be4a-ab8e513fc667",
	ExpCertPath: "certs/target.crt",
}}

var RotateCertificatesTests = []struct {
	Desc   string
	CertID string
}{{
	Desc:   "rotate certificates",
	CertID: "d7f58600-4b8e-4260-be3d-ff1641e1c8e9",
}}

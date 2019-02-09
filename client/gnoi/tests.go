package gnoi

var RebootTests = []struct {
	Desc    string
	Message string
}{{
	Desc:    "reboot system",
	Message: "Testing OVS Reboot Functionality",
}}

var GetCertificatesTests = []struct {
	Desc string
}{{
	Desc: "get certificates",
}}

var RotateCertificatesTests = []struct {
	Desc   string
	CertID string
}{{
	Desc:   "rotate certificates",
	CertID: "d7f58600-4b8e-4260-be3d-ff1641e1c8e9",
}}

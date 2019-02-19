/* Copyright 2019 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

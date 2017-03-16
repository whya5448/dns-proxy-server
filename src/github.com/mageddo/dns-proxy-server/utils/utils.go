package utils

import (
	"path/filepath"
	"os"
	"encoding/json"
	"io"
)

var QTypeCodes = map[uint16] string {
0 : "TypeNone",
1 : "TypeA",
2 : "TypeNS",
3 : "TypeMD",
4 : "TypeMF",
5 : "TypeCNAME",
6 : "TypeSOA",
7 : "TypeMB",
8 : "TypeMG",
9 : "TypeMR",
10 : "TypeNULL",
12 : "TypePTR",
13 : "TypeHINFO",
14 : "TypeMINFO",
15 : "TypeMX",
16 : "TypeTXT",
17 : "TypeRP",
18 : "TypeAFSDB",
19 : "TypeX25",
20 : "TypeISDN",
21 : "TypeRT",
23 : "TypeNSAPPTR",
24 : "TypeSIG",
25 : "TypeKEY",
26 : "TypePX",
27 : "TypeGPOS",
28 : "TypeAAAA",
29 : "TypeLOC",
30 : "TypeNXT",
31 : "TypeEID",
32 : "TypeNIMLOC",
33 : "TypeSRV",
34 : "TypeATMA",
35 : "TypeNAPTR",
36 : "TypeKX",
37 : "TypeCERT",
39 : "TypeDNAME",
41 : "TypeOPT", // EDNS
43 : "TypeDS",
44 : "TypeSSHFP",
46 : "TypeRRSIG",
47 : "TypeNSEC",
48 : "TypeDNSKEY",
49 : "TypeDHCID",
50 : "TypeNSEC3",
51 : "TypeNSEC3PARAM",
52 : "TypeTLSA",
53 : "TypeSMIMEA",
55 : "TypeHIP",
56 : "TypeNINFO",
57 : "TypeRKEY",
58 : "TypeTALINK",
59 : "TypeCDS",
60 : "TypeCDNSKEY",
61 : "TypeOPENPGPKEY",
99 : "TypeSPF",
100 : "TypeUINFO",
101 : "TypeUID",
102 : "TypeGID",
103 : "TypeUNSPEC",
104 : "TypeNID",
105 : "TypeL32",
106 : "TypeL64",
107 : "TypeLP",
108 : "TypeEUI48",
109 : "TypeEUI64",
256 : "TypeURI",
257 : "TypeCAA",

249 : "TypeTKEY",
250 : "TypeTSIG",

// valid Question.Qtype only,
251 : "TypeIXFR",
252 : "TypeAXFR",
253 : "TypeMAILB",
254 : "TypeMAILA",
255 : "TypeANY",

32768 : "TypeTA",
32769 : "TypeDLV",
65535 : "TypeReserved",

}

var QClassCodes = map[uint16] string {

	// valid Question.Qclass,
	1 : "ClassINET",
	2 : "ClassCSNET",
	3 : "ClassCHAOS",
	4 : "ClassHESIOD",
	254 : "ClassNONE",
	255 : "ClassANY",
}

var RCodes  = map[uint16] string {
	// Message Response Codes.,
	0 : "RcodeSuccess",
	1 : "RcodeFormatError",
	2 : "RcodeServerFailure",
	3 : "RcodeNameError",
	4 : "RcodeNotImplemented",
	5 : "RcodeRefused",
	6 : "RcodeYXDomain",
	7 : "RcodeYXRrset",
	8 : "RcodeNXRrset",
	9 : "RcodeNotAuth",
	10 : "RcodeNotZone",
	16 : "RcodeBadSig", // TSIG
	//16 : "RcodeBadVers", // EDNS0
	17 : "RcodeBadKey",
	18 : "RcodeBadTime",
	19 : "RcodeBadMode", // TKEY
	20 : "RcodeBadName",
	21 : "RcodeBadAlg",
	22 : "RcodeBadTrunc", // TSIG
	23 : "RcodeBadCookie", // DNS Cookies
}

var opCodes  = map[uint16] string {
	// Message Opcodes. There is no 3.,
	0 : "OpcodeQuery",
	1 : "OpcodeIQuery",
	2 : "OpcodeStatus",
	4 : "OpcodeNotify",
	5 : "OpcodeUpdate",
}

func DnsQTypeCodeToName(code uint16) string {
	return QTypeCodes[code]
}

func GetCurrentPath() string {

	currDIr := os.Getenv("MG_WORK_DIR")
	if len(currDIr) != 0 {
		return currDIr
	}
	currentPath, _ := filepath.Abs(filepath.Dir("."))
	return currentPath

}

func GetPath(path string) string {
	if path[:1] != "/" {
		path = "/" + path
	}
	return GetCurrentPath() + path
}

func GetJsonEncoder(w io.Writer) *json.Encoder {
	decoder := json.NewEncoder(w)
	decoder.SetIndent("", "\t")
	return decoder
}

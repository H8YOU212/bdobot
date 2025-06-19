package logfdb

func Logapi(reqdata string, resdata string) logapi {
	if reqdata != "" && resdata != "" {
		return logapi{
			is_success: true,
			res_body: resdata,
			req_body: reqdata,
		}

	} 
	return logapi{
		is_success: true,
		req_body: "",
		res_body: "",
	}
}

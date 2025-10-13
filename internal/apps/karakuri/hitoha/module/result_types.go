package karakuri_mod

type ResponseEnableModule struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type ResponseDisableModule struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type ResponseModuleList struct {
	Result string  `json:"result"`
	List   ModList `json:"module_list"`
}

func CreateResponseEnableModule(result string, message string) ResponseEnableModule {
	resp := ResponseEnableModule{
		Result:  result,
		Message: message,
	}
	return resp
}

func CreateResponseDisableModule(result string, message string) ResponseDisableModule {
	resp := ResponseDisableModule{
		Result:  result,
		Message: message,
	}
	return resp
}

func CreateResponseModuleList(result string, module_list ModList) ResponseModuleList {
	resp := ResponseModuleList{
		Result: result,
		List:   module_list,
	}
	return resp
}

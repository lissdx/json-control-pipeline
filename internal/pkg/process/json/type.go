package json

type addConfiguration struct {
	FieldName  string `json:"fieldName"`
	FieldValue string `json:"fieldValue"`
}

type countNumOfFieldConfiguration struct {
	TargetFieldName string `json:"targetFieldName"`
}

type removeFieldConfiguration struct {
	FieldName string `json:"fieldName"`
}

type ProcessStep struct {
	ProcessName string
	ProcessConf string
}

type UnmarshalT struct {
	Steps []struct {
		Processor     string            `json:"processor"`
		Configuration map[string]string `json:"configuration"`
	} `json:"steps"`
}

//type AmpTagsRequest struct {
//	*go_entities.BaseRequest
//	UserIP    net.IP `json:"-"`
//	Component string `json:"component"`
//	Geo       string `json:"-"`
//}
//string `json:"-"`

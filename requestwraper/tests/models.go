//go:generate easyjson models.go

package tests

type flag struct {
	Type  string `json:"type"`
	Event string `json:"event"`
}

// device model
// easyjson:json
// swagger:model createDevice
type device struct {
	UserID string `json:"user_id" validate:"required"`
	Locale string `json:"locale,omitempty"`
	Num    int    `json:"num" validate:"min=4,max=15"`
	// required: true
	Type            string  `json:"type" validate:"required"`
	NativePushToken string  `json:"native_push_token" validate:"required_without_all=NativeVoIPToken"`
	NativeVoIPToken string  `json:"native_voip_token" validate:"required_without_all=NativePushToken"`
	Carrier         string  `json:"carrier,omitempty"`
	Mcc             string  `json:"mcc,omitempty"`
	Mnc             string  `json:"mnc,omitempty"`
	OsVersion       string  `json:"os_version,omitempty"`
	BuildNumber     string  `json:"build_number,omitempty"`
	AppVersion      string  `json:"app_version,omitempty"`
	CountryCode     string  `json:"country_code,omitempty"`
	PhoneNumber     string  `json:"phone_number,omitempty"`
	Mode            string  `json:"mode,omitempty"`
	Flags           []*flag `json:"flags"`
	MyData          *data   `validate:"required"`
}

type data struct {
	N string `validate:"required"`
}

// A deviceCreateRequestParams model.
//
// This is used for operations that want an Order as body of the request
// swagger:parameters createDevice
type deviceCreateRequestParams struct {
	// The device data to submit.
	//
	// in: body
	// required: true
	Body *device `json:"device"`

	Header
}

type Header struct {
	// Header
	//
	// in: header
	// required: true
	// schema:
	//   type: string
	//	 format: mongoid
	LnAuthorId string `json:"ln-author-id"`
}

//go:generate easyjson models.go

package tests

type flag struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// item model
// easyjson:json
// swagger:model createDevice
type item struct {
	UserID string  `json:"user_id" validate:"required"`
	Locale string  `json:"locale,omitempty"`
	Num    int     `json:"num" validate:"min=4,max=15"`
	Type   string  `json:"type" validate:"required"`
	Assn   string  `json:"assn" validate:"required_without_all=Assn1"`
	Assn1  string  `json:"assn1" validate:"required_without_all=Assn"`
	Flags  []*flag `json:"flags"`
	MyData *data   `validate:"required"`
}

type data struct {
	N string `validate:"required"`
}

// A deviceCreateRequestParams model.
//
// This is used for operations that want an Order as body of the request
// swagger:parameters createDevice
type deviceCreateRequestParams struct {
	// The item data to submit.
	//
	// in: body
	// required: true
	Body *item `json:"item"`

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

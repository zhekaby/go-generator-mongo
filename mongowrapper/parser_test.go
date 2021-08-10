package main

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zhekaby/go-generator-mongo/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

type VerificationStatus byte

var (
	VerificationStatusApproved VerificationStatus = 1
)

func TestParseUser(t *testing.T) {
	Convey(t.Name(), t, func() {
		p := common.NewParser("parser_test.go")
		So(p.Parse(), ShouldBeNil)
		So(p.Collections, ShouldHaveLength, 1)
		Convey("Checking big result", func() {
			s, _ := json.Marshal(p.Collections[0].Fields)
			So(string(s), ShouldEqual, testExpected)
		})
	})
}

var testExpected = `[{"Prop":"Var","Type":"string","JsonProp":"var","JsonPath":"Data.var","BsonProp":"Var","BsonPath":"Data.Var","GoPath":"DataVar","Ns":".Data.Var","NsShort":"Data.Var","NsCompact":"DataVar","Validations":{}},{"Prop":"ID","Type":"primitive.ObjectID","JsonProp":"ID","JsonPath":"ID","BsonProp":"_id","BsonPath":"_id","GoPath":"ID","Ns":".ID","NsShort":"ID","NsCompact":"ID","Validations":{}},{"Prop":"VerificationStatus","Type":"VerificationStatus","JsonProp":"VerificationStatus","JsonPath":"VerificationStatus","BsonProp":"verification_status","BsonPath":"verification_status","GoPath":"VerificationStatus","Ns":".VerificationStatus","NsShort":"VerificationStatus","NsCompact":"VerificationStatus","Validations":{}},{"Prop":"VerificationRequestedAt","Type":"time.Time","JsonProp":"VerificationRequestedAt","JsonPath":"VerificationRequestedAt","BsonProp":"verification_requested_at","BsonPath":"verification_requested_at","GoPath":"VerificationRequestedAt","Ns":".VerificationRequestedAt","NsShort":"VerificationRequestedAt","NsCompact":"VerificationRequestedAt","Validations":{}},{"Prop":"Email","Type":"string","JsonProp":"Email","JsonPath":"Email","BsonProp":"email","BsonPath":"email","GoPath":"Email","Ns":".Email","NsShort":"Email","NsCompact":"Email","Validations":{}},{"Prop":"Phone","Type":"string","JsonProp":"Phone","JsonPath":"Phone","BsonProp":"phone","BsonPath":"phone","GoPath":"Phone","Ns":".Phone","NsShort":"Phone","NsCompact":"Phone","Validations":{}},{"Prop":"Password","Type":"string","JsonProp":"Password","JsonPath":"Password","BsonProp":"password","BsonPath":"password","GoPath":"Password","Ns":".Password","NsShort":"Password","NsCompact":"Password","Validations":{}},{"Prop":"Pwd","Type":"string","JsonProp":"Pwd","JsonPath":"Pwd","BsonProp":"pwd","BsonPath":"pwd","GoPath":"Pwd","Ns":".Pwd","NsShort":"Pwd","NsCompact":"Pwd","Validations":{}},{"Prop":"Enabled","Type":"bool","JsonProp":"Enabled","JsonPath":"Enabled","BsonProp":"enabled","BsonPath":"enabled","GoPath":"Enabled","Ns":".Enabled","NsShort":"Enabled","NsCompact":"Enabled","Validations":{}},{"Prop":"FirstName","Type":"string","JsonProp":"FirstName","JsonPath":"Profile.FirstName","BsonProp":"first_name","BsonPath":"profile.first_name","GoPath":"ProfileFirstName","Ns":".FirstName","NsShort":"FirstName","NsCompact":"FirstName","Validations":{}},{"Prop":"LastName","Type":"string","JsonProp":"LastName","JsonPath":"Profile.LastName","BsonProp":"last_name","BsonPath":"profile.last_name","GoPath":"ProfileLastName","Ns":".LastName","NsShort":"LastName","NsCompact":"LastName","Validations":{}},{"Prop":"NickName","Type":"string","JsonProp":"NickName","JsonPath":"Profile.NickName","BsonProp":"nick_name","BsonPath":"profile.nick_name","GoPath":"ProfileNickName","Ns":".NickName","NsShort":"NickName","NsCompact":"NickName","Validations":{}},{"Prop":"ZipCode","Type":"string","JsonProp":"ZipCode","JsonPath":"Profile.Address.ZipCode","BsonProp":"zip_code","BsonPath":"profile.address.zip_code","GoPath":"ProfileAddressZipCode","Ns":".ZipCode","NsShort":"ZipCode","NsCompact":"ZipCode","Validations":{}},{"Prop":"Country","Type":"string","JsonProp":"Country","JsonPath":"Profile.Address.Country","BsonProp":"country","BsonPath":"profile.address.country","GoPath":"ProfileAddressCountry","Ns":".Country","NsShort":"Country","NsCompact":"Country","Validations":{}},{"Prop":"City","Type":"string","JsonProp":"City","JsonPath":"Profile.Address.City","BsonProp":"city","BsonPath":"profile.address.city","GoPath":"ProfileAddressCity","Ns":".City","NsShort":"City","NsCompact":"City","Validations":{}},{"Prop":"Address","Type":"string","JsonProp":"Address","JsonPath":"Profile.Address.Address","BsonProp":"address","BsonPath":"profile.address.address","GoPath":"ProfileAddressAddress","Ns":".Address","NsShort":"Address","NsCompact":"Address","Validations":{}},{"Prop":"Lang","Type":"string","JsonProp":"Lang","JsonPath":"Profile.Lang","BsonProp":"lang","BsonPath":"profile.lang","GoPath":"ProfileLang","Ns":".Lang","NsShort":"Lang","NsCompact":"Lang","Validations":{}},{"Prop":"Target","Type":"string","JsonProp":"Target","JsonPath":"TwoFA.Target","BsonProp":"target","BsonPath":"_2fa.target","GoPath":"TwoFATarget","Ns":".Target","NsShort":"Target","NsCompact":"Target","Validations":{}},{"Prop":"Secret","Type":"string","JsonProp":"Secret","JsonPath":"TwoFA.Secret","BsonProp":"secret","BsonPath":"_2fa.secret","GoPath":"TwoFASecret","Ns":".Secret","NsShort":"Secret","NsCompact":"Secret","Validations":{}},{"Prop":"AffiliateId","Type":"primitive.ObjectID","JsonProp":"AffiliateId","JsonPath":"AffiliateId","BsonProp":"affiliate_id","BsonPath":"affiliate_id","GoPath":"AffiliateId","Ns":".AffiliateId","NsShort":"AffiliateId","NsCompact":"AffiliateId","Validations":{}},{"Prop":"PartnerCode","Type":"string","JsonProp":"PartnerCode","JsonPath":"PartnerCode","BsonProp":"partner_code","BsonPath":"partner_code","GoPath":"PartnerCode","Ns":".PartnerCode","NsShort":"PartnerCode","NsCompact":"PartnerCode","Validations":{}},{"Prop":"PartnerRate","Type":"byte","JsonProp":"PartnerRate","JsonPath":"PartnerRate","BsonProp":"partner_rate","BsonPath":"partner_rate","GoPath":"PartnerRate","Ns":".PartnerRate","NsShort":"PartnerRate","NsCompact":"PartnerRate","Validations":{}},{"Prop":"PartnerCount","Type":"int","JsonProp":"PartnerCount","JsonPath":"PartnerCount","BsonProp":"partner_count","BsonPath":"partner_count","GoPath":"PartnerCount","Ns":".PartnerCount","NsShort":"PartnerCount","NsCompact":"PartnerCount","Validations":{}}]`

//mongowrapper:collection users
type User struct {
	Data                    *Data
	ID                      primitive.ObjectID `bson:"_id,omitempty"`
	VerificationStatus      VerificationStatus `bson:"verification_status"`
	VerificationRequestedAt time.Time          `bson:"verification_requested_at,omitempty"`
	Email                   string             `bson:"email"`
	Phone                   string             `bson:"phone,omitempty"`
	Password                *string            `bson:"password,omitempty"`
	Pwd                     string             `bson:"pwd,omitempty"`
	Enabled                 bool               `bson:"enabled"`
	Profile                 struct {
		FirstName string `bson:"first_name,omitempty"`
		LastName  string `bson:"last_name,omitempty"`
		NickName  string `bson:"nick_name,omitempty"`
		Address   struct {
			ZipCode string `bson:"zip_code,omitempty"`
			Country string `bson:"country,omitempty"`
			City    string `bson:"city,omitempty"`
			Address string `bson:"address,omitempty"`
		} `bson:"address"`
		Lang     string     `bson:"lang"`
		Birthday *time.Time `bson:"birthday"`
	} `bson:"profile"`
	Permissions map[string]interface{} `bson:"permissions"`
	TwoFA       struct {
		Target string   `bson:"target,omitempty"`
		Secret string   `bson:"secret,omitempty"`
		Codes  []string `bson:"codes,omitempty"`
	} `bson:"_2fa,omitempty"`
	Subscription map[string]string `bson:"subscription,omitempty"`

	AffiliateId  primitive.ObjectID `bson:"affiliate_id,omitempty"`
	PartnerCode  string             `bson:"partner_code,omitempty"`
	PartnerRate  *byte              `bson:"partner_rate,omitempty"`
	PartnerCount int                `bson:"partner_count,omitempty"`
}

type Data struct {
	Var *string `json:"var"`
}

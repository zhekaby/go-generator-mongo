package main

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zhekaby/go-generator-mongo/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go/parser"
	"go/token"
	"testing"
	"time"
)

type VerificationStatus byte

var (
	VerificationStatusApproved VerificationStatus = 1
)

func TestParseUser(t *testing.T) {
	Convey(t.Name(), t, func() {
		fset := token.NewFileSet()
		f, _ := parser.ParseFile(fset, "parser_test.go", nil, parser.ParseComments)
		fields := make([]field, 0, 100)
		for _, d := range f.Decls {
			if common.HasComment(d, "test_user") {
				deep(d, "", "", "", "", &fields)
			}
		}
		Convey("Checking big result", func() {
			s, _ := json.Marshal(fields)
			So(string(s), ShouldEqual, testExpected)
		})
	})
}

var testExpected = `[{"Prop":"Var","Type":"string","BsonProp":"Var","BsonPath":"Data.Var","GoPath":"DataVar"},{"Prop":"ID","Type":"primitive.ObjectID","BsonProp":"_id","BsonPath":"_id","GoPath":"ID"},{"Prop":"VerificationStatus","Type":"VerificationStatus","BsonProp":"verification_status","BsonPath":"verification_status","GoPath":"VerificationStatus"},{"Prop":"VerificationRequestedAt","Type":"time.Time","BsonProp":"verification_requested_at","BsonPath":"verification_requested_at","GoPath":"VerificationRequestedAt"},{"Prop":"Email","Type":"string","BsonProp":"email","BsonPath":"email","GoPath":"Email"},{"Prop":"Phone","Type":"string","BsonProp":"phone","BsonPath":"phone","GoPath":"Phone"},{"Prop":"Password","Type":"string","BsonProp":"password","BsonPath":"password","GoPath":"Password"},{"Prop":"Pwd","Type":"string","BsonProp":"pwd","BsonPath":"pwd","GoPath":"Pwd"},{"Prop":"Enabled","Type":"bool","BsonProp":"enabled","BsonPath":"enabled","GoPath":"Enabled"},{"Prop":"FirstName","Type":"string","BsonProp":"first_name","BsonPath":"profile.first_name","GoPath":"ProfileFirstName"},{"Prop":"LastName","Type":"string","BsonProp":"last_name","BsonPath":"profile.last_name","GoPath":"ProfileLastName"},{"Prop":"NickName","Type":"string","BsonProp":"nick_name","BsonPath":"profile.nick_name","GoPath":"ProfileNickName"},{"Prop":"ZipCode","Type":"string","BsonProp":"zip_code","BsonPath":"profile.address.zip_code","GoPath":"ProfileAddressZipCode"},{"Prop":"Country","Type":"string","BsonProp":"country","BsonPath":"profile.address.country","GoPath":"ProfileAddressCountry"},{"Prop":"City","Type":"string","BsonProp":"city","BsonPath":"profile.address.city","GoPath":"ProfileAddressCity"},{"Prop":"Address","Type":"string","BsonProp":"address","BsonPath":"profile.address.address","GoPath":"ProfileAddressAddress"},{"Prop":"Lang","Type":"string","BsonProp":"lang","BsonPath":"profile.lang","GoPath":"ProfileLang"},{"Prop":"Target","Type":"string","BsonProp":"target","BsonPath":"_2fa.target","GoPath":"TwoFATarget"},{"Prop":"Secret","Type":"string","BsonProp":"secret","BsonPath":"_2fa.secret","GoPath":"TwoFASecret"},{"Prop":"AffiliateId","Type":"primitive.ObjectID","BsonProp":"affiliate_id","BsonPath":"affiliate_id","GoPath":"AffiliateId"},{"Prop":"PartnerCode","Type":"string","BsonProp":"partner_code","BsonPath":"partner_code","GoPath":"PartnerCode"},{"Prop":"PartnerRate","Type":"byte","BsonProp":"partner_rate","BsonPath":"partner_rate","GoPath":"PartnerRate"},{"Prop":"PartnerCount","Type":"int","BsonProp":"partner_count","BsonPath":"partner_count","GoPath":"PartnerCount"}]`

//test_user
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

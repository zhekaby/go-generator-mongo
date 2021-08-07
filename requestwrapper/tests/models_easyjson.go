// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package tests

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonD2b7633eDecodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests(in *jlexer.Lexer, out *item) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "user_id":
			out.UserID = string(in.String())
		case "locale":
			out.Locale = string(in.String())
		case "num":
			out.Num = int(in.Int())
		case "type":
			out.Type = string(in.String())
		case "assn":
			out.Assn = string(in.String())
		case "assn1":
			out.Assn1 = string(in.String())
		case "flags":
			if in.IsNull() {
				in.Skip()
				out.Flags = nil
			} else {
				in.Delim('[')
				if out.Flags == nil {
					if !in.IsDelim(']') {
						out.Flags = make([]*flag, 0, 8)
					} else {
						out.Flags = []*flag{}
					}
				} else {
					out.Flags = (out.Flags)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *flag
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(flag)
						}
						easyjsonD2b7633eDecodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests1(in, v1)
					}
					out.Flags = append(out.Flags, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "MyData":
			if in.IsNull() {
				in.Skip()
				out.MyData = nil
			} else {
				if out.MyData == nil {
					out.MyData = new(data)
				}
				easyjsonD2b7633eDecodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests2(in, out.MyData)
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests(out *jwriter.Writer, in item) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix[1:])
		out.String(string(in.UserID))
	}
	if in.Locale != "" {
		const prefix string = ",\"locale\":"
		out.RawString(prefix)
		out.String(string(in.Locale))
	}
	{
		const prefix string = ",\"num\":"
		out.RawString(prefix)
		out.Int(int(in.Num))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"assn\":"
		out.RawString(prefix)
		out.String(string(in.Assn))
	}
	{
		const prefix string = ",\"assn1\":"
		out.RawString(prefix)
		out.String(string(in.Assn1))
	}
	{
		const prefix string = ",\"flags\":"
		out.RawString(prefix)
		if in.Flags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Flags {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					easyjsonD2b7633eEncodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests1(out, *v3)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"MyData\":"
		out.RawString(prefix)
		if in.MyData == nil {
			out.RawString("null")
		} else {
			easyjsonD2b7633eEncodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests2(out, *in.MyData)
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v item) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v item) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *item) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *item) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests(l, v)
}
func easyjsonD2b7633eDecodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests2(in *jlexer.Lexer, out *data) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "N":
			out.N = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests2(out *jwriter.Writer, in data) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"N\":"
		out.RawString(prefix[1:])
		out.String(string(in.N))
	}
	out.RawByte('}')
}
func easyjsonD2b7633eDecodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests1(in *jlexer.Lexer, out *flag) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			out.Type = string(in.String())
		case "value":
			out.Value = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonD2b7633eEncodeGithubComZhekabyGoGeneratorMongoRequestwrapperTests1(out *jwriter.Writer, in flag) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix[1:])
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.String(string(in.Value))
	}
	out.RawByte('}')
}
package HttpContext

import (
	"encoding/json"
	"log"
	"reflect"
)

type ContentEncoder interface {
	Encode(content any) []byte
	CanAccept(contentType string) bool
	CanAcceptRequest(request *Request) bool
	CanSendResponse(request *Request) bool
	EncodeError(err error) []byte
	HeadersForType() map[string]string
	// SetHeaders(headers map[string][]string)
}

type ContentEncoders struct {
	PlainText ContentEncoder
	Json      ContentEncoder
}

type BaseEncoder struct {
	Name        string
	ContentType string
	Headers     map[string][]string
}

func createBaseEncoder(name, contentType string) BaseEncoder {
	return BaseEncoder{
		Name:        name,
		ContentType: contentType,
		Headers:     make(map[string][]string),
	}
}

type EncoderFor struct {
	Request  ContentEncoder
	Response ContentEncoder
}

type EncoderService struct {
	Encoders ContentEncoders
}

var Encoder = &EncoderService{
	Encoders: ContentEncoders{
		PlainText: NewPlainTextEncoder(),
		Json:      NewJsonEncoder(),
	},
}

func (encoder *EncoderService) ContentTypeEncoderForRequest(request *Request) EncoderFor {
	encodersVal := reflect.ValueOf(encoder.Encoders)
	fieldCount := encodersVal.NumField()

	encoderFor := EncoderFor{
		Request:  nil,
		Response: nil,
	}

	for i := 0; i < fieldCount; i++ {
		field := encodersVal.Field(i)
		enc := field.Interface().(ContentEncoder)

		if enc.CanAcceptRequest(request) {
			encoderFor.Request = enc
		}
		if enc.CanSendResponse(request) {
			encoderFor.Response = enc
		}
	}

	if encoderFor.Request == nil {
		encoderFor.Request = encoder.Encoders.Json
	}
	if encoderFor.Response == nil {
		encoderFor.Response = encoder.Encoders.Json
	}

	return encoderFor
}

//
//
//  JSON ENCODER
//
//

type JsonEncoder struct {
	BaseEncoder
}

func NewJsonEncoder() ContentEncoder {
	return &JsonEncoder{
		BaseEncoder: createBaseEncoder("Json", HEADER_CONTENT_TYPE_JSON),
	}
}

func (p *JsonEncoder) HeadersForType() map[string]string {
	return map[string]string{
		"content-type": p.ContentType,
	}
}

func (p *JsonEncoder) CanAccept(contentType string) bool {
	return false
}

func (p *JsonEncoder) CanAcceptRequest(request *Request) bool {
	return request.Headers().Has("content-type", "application/json")
}

func (p *JsonEncoder) CanSendResponse(request *Request) bool {
	return request.Headers().Has("accept", "application/json")
}

func (p *JsonEncoder) Encode(content any) []byte {
	data, err := json.Marshal(content)

	if err != nil {
		log.Printf("Failed to encode json response: %s", err)
		panic("Failed to encode json")
	}

	return data
}

func (p *JsonEncoder) EncodeError(err error) []byte {
	type ErrorRes struct {
		Message string `json:"message"`
	}
	return p.Encode(ErrorRes{Message: err.Error()})
}

//
//
//  PLAIN TEXT ENCODER
//
//

type PlainTextEncoder struct {
	BaseEncoder
}

func NewPlainTextEncoder() ContentEncoder {
	return &PlainTextEncoder{
		BaseEncoder: createBaseEncoder("PlainText", HEADER_CONTENT_TYPE_TEXT_PLAIN),
	}
}

func (p *PlainTextEncoder) CanAccept(contentType string) bool {
	return false
}

func (p *PlainTextEncoder) HeadersForType() map[string]string {
	return map[string]string{
		"content-type": p.ContentType,
	}
}

func (p *PlainTextEncoder) CanAcceptRequest(request *Request) bool {
	return request.Headers().Has("content-type", "text/plain")
}

func (p *PlainTextEncoder) CanSendResponse(request *Request) bool {
	return request.Headers().Has("accept", "text/plain")
}

func (p *PlainTextEncoder) Encode(content any) []byte {

	if data, ok := content.(string); ok {
		return []byte(data)
	}

	log.Printf("Data: %v", content)
	panic("Data cannot be converted to string?")
}

func (p *PlainTextEncoder) EncodeError(err error) []byte {
	return p.Encode("Woops. Something went wrong!\n" + err.Error())
}

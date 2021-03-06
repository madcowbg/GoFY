// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

package generated

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type InstrumentType int32

const (
	InstrumentType_Option InstrumentType = 1
)

var InstrumentType_name = map[int32]string{
	1: "Option",
}

var InstrumentType_value = map[string]int32{
	"Option": 1,
}

func (x InstrumentType) Enum() *InstrumentType {
	p := new(InstrumentType)
	*p = x
	return p
}

func (x InstrumentType) String() string {
	return proto.EnumName(InstrumentType_name, int32(x))
}

func (x *InstrumentType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(InstrumentType_value, data, "InstrumentType")
	if err != nil {
		return err
	}
	*x = InstrumentType(value)
	return nil
}

func (InstrumentType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

type OptionType int32

const (
	OptionType_American OptionType = 1
	OptionType_European OptionType = 2
)

var OptionType_name = map[int32]string{
	1: "American",
	2: "European",
}

var OptionType_value = map[string]int32{
	"American": 1,
	"European": 2,
}

func (x OptionType) Enum() *OptionType {
	p := new(OptionType)
	*p = x
	return p
}

func (x OptionType) String() string {
	return proto.EnumName(OptionType_name, int32(x))
}

func (x *OptionType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(OptionType_value, data, "OptionType")
	if err != nil {
		return err
	}
	*x = OptionType(value)
	return nil
}

func (OptionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

type OptionParity int32

const (
	OptionParity_Call OptionParity = 1
	OptionParity_Put  OptionParity = 2
)

var OptionParity_name = map[int32]string{
	1: "Call",
	2: "Put",
}

var OptionParity_value = map[string]int32{
	"Call": 1,
	"Put":  2,
}

func (x OptionParity) Enum() *OptionParity {
	p := new(OptionParity)
	*p = x
	return p
}

func (x OptionParity) String() string {
	return proto.EnumName(OptionParity_name, int32(x))
}

func (x *OptionParity) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(OptionParity_value, data, "OptionParity")
	if err != nil {
		return err
	}
	*x = OptionParity(value)
	return nil
}

func (OptionParity) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

type BootstrapMethod int32

const (
	BootstrapMethod_Naive          BootstrapMethod = 1
	BootstrapMethod_MonotoneConvex BootstrapMethod = 2
)

var BootstrapMethod_name = map[int32]string{
	1: "Naive",
	2: "MonotoneConvex",
}

var BootstrapMethod_value = map[string]int32{
	"Naive":          1,
	"MonotoneConvex": 2,
}

func (x BootstrapMethod) Enum() *BootstrapMethod {
	p := new(BootstrapMethod)
	*p = x
	return p
}

func (x BootstrapMethod) String() string {
	return proto.EnumName(BootstrapMethod_name, int32(x))
}

func (x *BootstrapMethod) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(BootstrapMethod_value, data, "BootstrapMethod")
	if err != nil {
		return err
	}
	*x = BootstrapMethod(value)
	return nil
}

func (BootstrapMethod) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

type RequestCalculateOptionAnalytics struct {
	TermsAndConditions   *OptionTermsAndConditions `protobuf:"bytes,1,req,name=termsAndConditions" json:"termsAndConditions,omitempty"`
	StateOfWorld         *StateOfWorld             `protobuf:"bytes,2,req,name=stateOfWorld" json:"stateOfWorld,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *RequestCalculateOptionAnalytics) Reset()         { *m = RequestCalculateOptionAnalytics{} }
func (m *RequestCalculateOptionAnalytics) String() string { return proto.CompactTextString(m) }
func (*RequestCalculateOptionAnalytics) ProtoMessage()    {}
func (*RequestCalculateOptionAnalytics) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{0}
}

func (m *RequestCalculateOptionAnalytics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestCalculateOptionAnalytics.Unmarshal(m, b)
}
func (m *RequestCalculateOptionAnalytics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestCalculateOptionAnalytics.Marshal(b, m, deterministic)
}
func (m *RequestCalculateOptionAnalytics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestCalculateOptionAnalytics.Merge(m, src)
}
func (m *RequestCalculateOptionAnalytics) XXX_Size() int {
	return xxx_messageInfo_RequestCalculateOptionAnalytics.Size(m)
}
func (m *RequestCalculateOptionAnalytics) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestCalculateOptionAnalytics.DiscardUnknown(m)
}

var xxx_messageInfo_RequestCalculateOptionAnalytics proto.InternalMessageInfo

func (m *RequestCalculateOptionAnalytics) GetTermsAndConditions() *OptionTermsAndConditions {
	if m != nil {
		return m.TermsAndConditions
	}
	return nil
}

func (m *RequestCalculateOptionAnalytics) GetStateOfWorld() *StateOfWorld {
	if m != nil {
		return m.StateOfWorld
	}
	return nil
}

type ResponseCalculateOptionAnalytics struct {
	Price                *float32 `protobuf:"fixed32,1,req,name=Price" json:"Price,omitempty"`
	Delta                *float32 `protobuf:"fixed32,2,req,name=Delta" json:"Delta,omitempty"`
	Gamma                *float32 `protobuf:"fixed32,3,req,name=Gamma" json:"Gamma,omitempty"`
	Theta                *float32 `protobuf:"fixed32,4,req,name=Theta" json:"Theta,omitempty"`
	Rho                  *float32 `protobuf:"fixed32,5,req,name=Rho" json:"Rho,omitempty"`
	Intrinsic            *float32 `protobuf:"fixed32,6,req,name=Intrinsic" json:"Intrinsic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseCalculateOptionAnalytics) Reset()         { *m = ResponseCalculateOptionAnalytics{} }
func (m *ResponseCalculateOptionAnalytics) String() string { return proto.CompactTextString(m) }
func (*ResponseCalculateOptionAnalytics) ProtoMessage()    {}
func (*ResponseCalculateOptionAnalytics) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{1}
}

func (m *ResponseCalculateOptionAnalytics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseCalculateOptionAnalytics.Unmarshal(m, b)
}
func (m *ResponseCalculateOptionAnalytics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseCalculateOptionAnalytics.Marshal(b, m, deterministic)
}
func (m *ResponseCalculateOptionAnalytics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseCalculateOptionAnalytics.Merge(m, src)
}
func (m *ResponseCalculateOptionAnalytics) XXX_Size() int {
	return xxx_messageInfo_ResponseCalculateOptionAnalytics.Size(m)
}
func (m *ResponseCalculateOptionAnalytics) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseCalculateOptionAnalytics.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseCalculateOptionAnalytics proto.InternalMessageInfo

func (m *ResponseCalculateOptionAnalytics) GetPrice() float32 {
	if m != nil && m.Price != nil {
		return *m.Price
	}
	return 0
}

func (m *ResponseCalculateOptionAnalytics) GetDelta() float32 {
	if m != nil && m.Delta != nil {
		return *m.Delta
	}
	return 0
}

func (m *ResponseCalculateOptionAnalytics) GetGamma() float32 {
	if m != nil && m.Gamma != nil {
		return *m.Gamma
	}
	return 0
}

func (m *ResponseCalculateOptionAnalytics) GetTheta() float32 {
	if m != nil && m.Theta != nil {
		return *m.Theta
	}
	return 0
}

func (m *ResponseCalculateOptionAnalytics) GetRho() float32 {
	if m != nil && m.Rho != nil {
		return *m.Rho
	}
	return 0
}

func (m *ResponseCalculateOptionAnalytics) GetIntrinsic() float32 {
	if m != nil && m.Intrinsic != nil {
		return *m.Intrinsic
	}
	return 0
}

type OptionTermsAndConditions struct {
	S                    *float32      `protobuf:"fixed32,1,req,name=S" json:"S,omitempty"`
	T                    *float32      `protobuf:"fixed32,2,req,name=T" json:"T,omitempty"`
	Type                 *OptionType   `protobuf:"varint,3,req,name=Type,enum=proto.generated.OptionType" json:"Type,omitempty"`
	Parity               *OptionParity `protobuf:"varint,4,req,name=Parity,enum=proto.generated.OptionParity" json:"Parity,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *OptionTermsAndConditions) Reset()         { *m = OptionTermsAndConditions{} }
func (m *OptionTermsAndConditions) String() string { return proto.CompactTextString(m) }
func (*OptionTermsAndConditions) ProtoMessage()    {}
func (*OptionTermsAndConditions) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{2}
}

func (m *OptionTermsAndConditions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OptionTermsAndConditions.Unmarshal(m, b)
}
func (m *OptionTermsAndConditions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OptionTermsAndConditions.Marshal(b, m, deterministic)
}
func (m *OptionTermsAndConditions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OptionTermsAndConditions.Merge(m, src)
}
func (m *OptionTermsAndConditions) XXX_Size() int {
	return xxx_messageInfo_OptionTermsAndConditions.Size(m)
}
func (m *OptionTermsAndConditions) XXX_DiscardUnknown() {
	xxx_messageInfo_OptionTermsAndConditions.DiscardUnknown(m)
}

var xxx_messageInfo_OptionTermsAndConditions proto.InternalMessageInfo

func (m *OptionTermsAndConditions) GetS() float32 {
	if m != nil && m.S != nil {
		return *m.S
	}
	return 0
}

func (m *OptionTermsAndConditions) GetT() float32 {
	if m != nil && m.T != nil {
		return *m.T
	}
	return 0
}

func (m *OptionTermsAndConditions) GetType() OptionType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return OptionType_American
}

func (m *OptionTermsAndConditions) GetParity() OptionParity {
	if m != nil && m.Parity != nil {
		return *m.Parity
	}
	return OptionParity_Call
}

type StateOfWorld struct {
	Parameters           *PricingParameters `protobuf:"bytes,1,req,name=parameters" json:"parameters,omitempty"`
	Spot                 *float32           `protobuf:"fixed32,2,req,name=Spot" json:"Spot,omitempty"`
	Time                 *float32           `protobuf:"fixed32,3,req,name=Time" json:"Time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *StateOfWorld) Reset()         { *m = StateOfWorld{} }
func (m *StateOfWorld) String() string { return proto.CompactTextString(m) }
func (*StateOfWorld) ProtoMessage()    {}
func (*StateOfWorld) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{3}
}

func (m *StateOfWorld) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateOfWorld.Unmarshal(m, b)
}
func (m *StateOfWorld) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateOfWorld.Marshal(b, m, deterministic)
}
func (m *StateOfWorld) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateOfWorld.Merge(m, src)
}
func (m *StateOfWorld) XXX_Size() int {
	return xxx_messageInfo_StateOfWorld.Size(m)
}
func (m *StateOfWorld) XXX_DiscardUnknown() {
	xxx_messageInfo_StateOfWorld.DiscardUnknown(m)
}

var xxx_messageInfo_StateOfWorld proto.InternalMessageInfo

func (m *StateOfWorld) GetParameters() *PricingParameters {
	if m != nil {
		return m.Parameters
	}
	return nil
}

func (m *StateOfWorld) GetSpot() float32 {
	if m != nil && m.Spot != nil {
		return *m.Spot
	}
	return 0
}

func (m *StateOfWorld) GetTime() float32 {
	if m != nil && m.Time != nil {
		return *m.Time
	}
	return 0
}

type PricingParameters struct {
	Sigma                *float32 `protobuf:"fixed32,1,req,name=Sigma" json:"Sigma,omitempty"`
	R                    *float32 `protobuf:"fixed32,2,req,name=R" json:"R,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PricingParameters) Reset()         { *m = PricingParameters{} }
func (m *PricingParameters) String() string { return proto.CompactTextString(m) }
func (*PricingParameters) ProtoMessage()    {}
func (*PricingParameters) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{4}
}

func (m *PricingParameters) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PricingParameters.Unmarshal(m, b)
}
func (m *PricingParameters) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PricingParameters.Marshal(b, m, deterministic)
}
func (m *PricingParameters) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PricingParameters.Merge(m, src)
}
func (m *PricingParameters) XXX_Size() int {
	return xxx_messageInfo_PricingParameters.Size(m)
}
func (m *PricingParameters) XXX_DiscardUnknown() {
	xxx_messageInfo_PricingParameters.DiscardUnknown(m)
}

var xxx_messageInfo_PricingParameters proto.InternalMessageInfo

func (m *PricingParameters) GetSigma() float32 {
	if m != nil && m.Sigma != nil {
		return *m.Sigma
	}
	return 0
}

func (m *PricingParameters) GetR() float32 {
	if m != nil && m.R != nil {
		return *m.R
	}
	return 0
}

type RequestBootstrapCurve struct {
	Method               *BootstrapMethod    `protobuf:"varint,1,req,name=method,enum=proto.generated.BootstrapMethod" json:"method,omitempty"`
	Lambda               *float64            `protobuf:"fixed64,2,req,name=lambda" json:"lambda,omitempty"`
	T0                   *float32            `protobuf:"fixed32,3,req,name=t0" json:"t0,omitempty"`
	BootstrapData        *CurveBootstrapData `protobuf:"bytes,4,req,name=bootstrapData" json:"bootstrapData,omitempty"`
	TenorData            *TenorDefs          `protobuf:"bytes,5,req,name=tenorData" json:"tenorData,omitempty"`
	OutputTenors         *TenorDefs          `protobuf:"bytes,6,req,name=outputTenors" json:"outputTenors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *RequestBootstrapCurve) Reset()         { *m = RequestBootstrapCurve{} }
func (m *RequestBootstrapCurve) String() string { return proto.CompactTextString(m) }
func (*RequestBootstrapCurve) ProtoMessage()    {}
func (*RequestBootstrapCurve) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{5}
}

func (m *RequestBootstrapCurve) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RequestBootstrapCurve.Unmarshal(m, b)
}
func (m *RequestBootstrapCurve) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RequestBootstrapCurve.Marshal(b, m, deterministic)
}
func (m *RequestBootstrapCurve) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RequestBootstrapCurve.Merge(m, src)
}
func (m *RequestBootstrapCurve) XXX_Size() int {
	return xxx_messageInfo_RequestBootstrapCurve.Size(m)
}
func (m *RequestBootstrapCurve) XXX_DiscardUnknown() {
	xxx_messageInfo_RequestBootstrapCurve.DiscardUnknown(m)
}

var xxx_messageInfo_RequestBootstrapCurve proto.InternalMessageInfo

func (m *RequestBootstrapCurve) GetMethod() BootstrapMethod {
	if m != nil && m.Method != nil {
		return *m.Method
	}
	return BootstrapMethod_Naive
}

func (m *RequestBootstrapCurve) GetLambda() float64 {
	if m != nil && m.Lambda != nil {
		return *m.Lambda
	}
	return 0
}

func (m *RequestBootstrapCurve) GetT0() float32 {
	if m != nil && m.T0 != nil {
		return *m.T0
	}
	return 0
}

func (m *RequestBootstrapCurve) GetBootstrapData() *CurveBootstrapData {
	if m != nil {
		return m.BootstrapData
	}
	return nil
}

func (m *RequestBootstrapCurve) GetTenorData() *TenorDefs {
	if m != nil {
		return m.TenorData
	}
	return nil
}

func (m *RequestBootstrapCurve) GetOutputTenors() *TenorDefs {
	if m != nil {
		return m.OutputTenors
	}
	return nil
}

type ResponseBootstrapCurve struct {
	SpotCurve                *Curve   `protobuf:"bytes,1,req,name=SpotCurve" json:"SpotCurve,omitempty"`
	InterpolatedSpotCurve    *Curve   `protobuf:"bytes,2,opt,name=InterpolatedSpotCurve" json:"InterpolatedSpotCurve,omitempty"`
	InterpolatedForwardCurve *Curve   `protobuf:"bytes,3,opt,name=InterpolatedForwardCurve" json:"InterpolatedForwardCurve,omitempty"`
	XXX_NoUnkeyedLiteral     struct{} `json:"-"`
	XXX_unrecognized         []byte   `json:"-"`
	XXX_sizecache            int32    `json:"-"`
}

func (m *ResponseBootstrapCurve) Reset()         { *m = ResponseBootstrapCurve{} }
func (m *ResponseBootstrapCurve) String() string { return proto.CompactTextString(m) }
func (*ResponseBootstrapCurve) ProtoMessage()    {}
func (*ResponseBootstrapCurve) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{6}
}

func (m *ResponseBootstrapCurve) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseBootstrapCurve.Unmarshal(m, b)
}
func (m *ResponseBootstrapCurve) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseBootstrapCurve.Marshal(b, m, deterministic)
}
func (m *ResponseBootstrapCurve) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseBootstrapCurve.Merge(m, src)
}
func (m *ResponseBootstrapCurve) XXX_Size() int {
	return xxx_messageInfo_ResponseBootstrapCurve.Size(m)
}
func (m *ResponseBootstrapCurve) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseBootstrapCurve.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseBootstrapCurve proto.InternalMessageInfo

func (m *ResponseBootstrapCurve) GetSpotCurve() *Curve {
	if m != nil {
		return m.SpotCurve
	}
	return nil
}

func (m *ResponseBootstrapCurve) GetInterpolatedSpotCurve() *Curve {
	if m != nil {
		return m.InterpolatedSpotCurve
	}
	return nil
}

func (m *ResponseBootstrapCurve) GetInterpolatedForwardCurve() *Curve {
	if m != nil {
		return m.InterpolatedForwardCurve
	}
	return nil
}

type TenorDefs struct {
	Tenors               []float32 `protobuf:"fixed32,1,rep,name=tenors" json:"tenors,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *TenorDefs) Reset()         { *m = TenorDefs{} }
func (m *TenorDefs) String() string { return proto.CompactTextString(m) }
func (*TenorDefs) ProtoMessage()    {}
func (*TenorDefs) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{7}
}

func (m *TenorDefs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TenorDefs.Unmarshal(m, b)
}
func (m *TenorDefs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TenorDefs.Marshal(b, m, deterministic)
}
func (m *TenorDefs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TenorDefs.Merge(m, src)
}
func (m *TenorDefs) XXX_Size() int {
	return xxx_messageInfo_TenorDefs.Size(m)
}
func (m *TenorDefs) XXX_DiscardUnknown() {
	xxx_messageInfo_TenorDefs.DiscardUnknown(m)
}

var xxx_messageInfo_TenorDefs proto.InternalMessageInfo

func (m *TenorDefs) GetTenors() []float32 {
	if m != nil {
		return m.Tenors
	}
	return nil
}

type CurveBootstrapData struct {
	BondDefinitions      []*CouponBondDef `protobuf:"bytes,1,rep,name=bondDefinitions" json:"bondDefinitions,omitempty"`
	Yields               []float32        `protobuf:"fixed32,2,rep,name=Yields" json:"Yields,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *CurveBootstrapData) Reset()         { *m = CurveBootstrapData{} }
func (m *CurveBootstrapData) String() string { return proto.CompactTextString(m) }
func (*CurveBootstrapData) ProtoMessage()    {}
func (*CurveBootstrapData) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{8}
}

func (m *CurveBootstrapData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CurveBootstrapData.Unmarshal(m, b)
}
func (m *CurveBootstrapData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CurveBootstrapData.Marshal(b, m, deterministic)
}
func (m *CurveBootstrapData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CurveBootstrapData.Merge(m, src)
}
func (m *CurveBootstrapData) XXX_Size() int {
	return xxx_messageInfo_CurveBootstrapData.Size(m)
}
func (m *CurveBootstrapData) XXX_DiscardUnknown() {
	xxx_messageInfo_CurveBootstrapData.DiscardUnknown(m)
}

var xxx_messageInfo_CurveBootstrapData proto.InternalMessageInfo

func (m *CurveBootstrapData) GetBondDefinitions() []*CouponBondDef {
	if m != nil {
		return m.BondDefinitions
	}
	return nil
}

func (m *CurveBootstrapData) GetYields() []float32 {
	if m != nil {
		return m.Yields
	}
	return nil
}

type CouponBondDef struct {
	IssueTime            *float32 `protobuf:"fixed32,1,req,name=IssueTime" json:"IssueTime,omitempty"`
	Maturity             *float32 `protobuf:"fixed32,2,req,name=Maturity" json:"Maturity,omitempty"`
	CouponFrequency      *float32 `protobuf:"fixed32,3,req,name=CouponFrequency" json:"CouponFrequency,omitempty"`
	Coupon               *float32 `protobuf:"fixed32,4,req,name=Coupon" json:"Coupon,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CouponBondDef) Reset()         { *m = CouponBondDef{} }
func (m *CouponBondDef) String() string { return proto.CompactTextString(m) }
func (*CouponBondDef) ProtoMessage()    {}
func (*CouponBondDef) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{9}
}

func (m *CouponBondDef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CouponBondDef.Unmarshal(m, b)
}
func (m *CouponBondDef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CouponBondDef.Marshal(b, m, deterministic)
}
func (m *CouponBondDef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CouponBondDef.Merge(m, src)
}
func (m *CouponBondDef) XXX_Size() int {
	return xxx_messageInfo_CouponBondDef.Size(m)
}
func (m *CouponBondDef) XXX_DiscardUnknown() {
	xxx_messageInfo_CouponBondDef.DiscardUnknown(m)
}

var xxx_messageInfo_CouponBondDef proto.InternalMessageInfo

func (m *CouponBondDef) GetIssueTime() float32 {
	if m != nil && m.IssueTime != nil {
		return *m.IssueTime
	}
	return 0
}

func (m *CouponBondDef) GetMaturity() float32 {
	if m != nil && m.Maturity != nil {
		return *m.Maturity
	}
	return 0
}

func (m *CouponBondDef) GetCouponFrequency() float32 {
	if m != nil && m.CouponFrequency != nil {
		return *m.CouponFrequency
	}
	return 0
}

func (m *CouponBondDef) GetCoupon() float32 {
	if m != nil && m.Coupon != nil {
		return *m.Coupon
	}
	return 0
}

type Curve struct {
	Tenors               []float32 `protobuf:"fixed32,1,rep,name=Tenors" json:"Tenors,omitempty"`
	Rates                []float32 `protobuf:"fixed32,2,rep,name=Rates" json:"Rates,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Curve) Reset()         { *m = Curve{} }
func (m *Curve) String() string { return proto.CompactTextString(m) }
func (*Curve) ProtoMessage()    {}
func (*Curve) Descriptor() ([]byte, []int) {
	return fileDescriptor_00212fb1f9d3bf1c, []int{10}
}

func (m *Curve) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Curve.Unmarshal(m, b)
}
func (m *Curve) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Curve.Marshal(b, m, deterministic)
}
func (m *Curve) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Curve.Merge(m, src)
}
func (m *Curve) XXX_Size() int {
	return xxx_messageInfo_Curve.Size(m)
}
func (m *Curve) XXX_DiscardUnknown() {
	xxx_messageInfo_Curve.DiscardUnknown(m)
}

var xxx_messageInfo_Curve proto.InternalMessageInfo

func (m *Curve) GetTenors() []float32 {
	if m != nil {
		return m.Tenors
	}
	return nil
}

func (m *Curve) GetRates() []float32 {
	if m != nil {
		return m.Rates
	}
	return nil
}

func init() {
	proto.RegisterEnum("proto.generated.InstrumentType", InstrumentType_name, InstrumentType_value)
	proto.RegisterEnum("proto.generated.OptionType", OptionType_name, OptionType_value)
	proto.RegisterEnum("proto.generated.OptionParity", OptionParity_name, OptionParity_value)
	proto.RegisterEnum("proto.generated.BootstrapMethod", BootstrapMethod_name, BootstrapMethod_value)
	proto.RegisterType((*RequestCalculateOptionAnalytics)(nil), "proto.generated.RequestCalculateOptionAnalytics")
	proto.RegisterType((*ResponseCalculateOptionAnalytics)(nil), "proto.generated.ResponseCalculateOptionAnalytics")
	proto.RegisterType((*OptionTermsAndConditions)(nil), "proto.generated.OptionTermsAndConditions")
	proto.RegisterType((*StateOfWorld)(nil), "proto.generated.StateOfWorld")
	proto.RegisterType((*PricingParameters)(nil), "proto.generated.PricingParameters")
	proto.RegisterType((*RequestBootstrapCurve)(nil), "proto.generated.RequestBootstrapCurve")
	proto.RegisterType((*ResponseBootstrapCurve)(nil), "proto.generated.ResponseBootstrapCurve")
	proto.RegisterType((*TenorDefs)(nil), "proto.generated.TenorDefs")
	proto.RegisterType((*CurveBootstrapData)(nil), "proto.generated.CurveBootstrapData")
	proto.RegisterType((*CouponBondDef)(nil), "proto.generated.CouponBondDef")
	proto.RegisterType((*Curve)(nil), "proto.generated.Curve")
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_00212fb1f9d3bf1c) }

var fileDescriptor_00212fb1f9d3bf1c = []byte{
	// 805 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0xd6, 0x38, 0x4d, 0x68, 0x4e, 0xb3, 0xa9, 0x77, 0xc4, 0x56, 0xd6, 0xb2, 0x40, 0xf0, 0xde,
	0x84, 0x5e, 0x74, 0xab, 0x88, 0x15, 0x7b, 0x85, 0xd4, 0xa4, 0x2c, 0x54, 0xa2, 0x50, 0x4d, 0x2c,
	0xa1, 0xe5, 0x6e, 0x1a, 0x9f, 0xb6, 0x96, 0xec, 0x19, 0x33, 0x1e, 0x17, 0xc2, 0x2b, 0xf0, 0x0e,
	0x3c, 0x00, 0x0f, 0x01, 0xaf, 0xc4, 0x23, 0xa0, 0xf9, 0xc9, 0xaf, 0x93, 0x5e, 0xc5, 0xdf, 0xf1,
	0xf9, 0xbe, 0x39, 0xfe, 0xce, 0x9c, 0x13, 0xe8, 0xf2, 0x32, 0x3b, 0x2b, 0x95, 0xd4, 0x92, 0x1e,
	0xdb, 0x9f, 0xb3, 0x7b, 0x14, 0xa8, 0xb8, 0xc6, 0x34, 0xfe, 0x97, 0xc0, 0xe7, 0x0c, 0x7f, 0xad,
	0xb1, 0xd2, 0x13, 0x9e, 0xcf, 0xea, 0x9c, 0x6b, 0xfc, 0xa9, 0xd4, 0x99, 0x14, 0x17, 0x82, 0xe7,
	0x73, 0x9d, 0xcd, 0x2a, 0xfa, 0x01, 0xa8, 0x46, 0x55, 0x54, 0x17, 0x22, 0x9d, 0x48, 0x91, 0x66,
	0xe6, 0x6d, 0x15, 0x91, 0x41, 0x30, 0x3c, 0x1a, 0x7d, 0x79, 0xb6, 0xa5, 0x78, 0xe6, 0xd8, 0x49,
	0x83, 0xc0, 0x76, 0x88, 0xd0, 0x0b, 0xe8, 0x55, 0xda, 0x1c, 0x79, 0xf7, 0xb3, 0x54, 0x79, 0x1a,
	0x05, 0x56, 0xf4, 0xd3, 0x86, 0xe8, 0x74, 0x2d, 0x89, 0x6d, 0x50, 0xe2, 0xbf, 0x09, 0x0c, 0x18,
	0x56, 0xa5, 0x14, 0x15, 0xee, 0xfd, 0x84, 0x8f, 0xa1, 0x7d, 0xa3, 0xb2, 0x19, 0xda, 0xaa, 0x03,
	0xe6, 0x80, 0x89, 0x5e, 0x62, 0xae, 0xb9, 0x3d, 0x36, 0x60, 0x0e, 0x98, 0xe8, 0x77, 0xbc, 0x28,
	0x78, 0xd4, 0x72, 0x51, 0x0b, 0x4c, 0x34, 0x79, 0x40, 0xcd, 0xa3, 0x03, 0x17, 0xb5, 0x80, 0x86,
	0xd0, 0x62, 0x0f, 0x32, 0x6a, 0xdb, 0x98, 0x79, 0xa4, 0xaf, 0xa0, 0x7b, 0x25, 0xb4, 0xca, 0x44,
	0x95, 0xcd, 0xa2, 0x8e, 0x8d, 0xaf, 0x02, 0xf1, 0x5f, 0x04, 0xa2, 0x7d, 0x06, 0xd1, 0x1e, 0x90,
	0xa9, 0x2f, 0x90, 0x4c, 0x0d, 0x4a, 0x7c, 0x61, 0x24, 0xa1, 0x6f, 0xe0, 0x20, 0x99, 0x97, 0x68,
	0x6b, 0xea, 0x8f, 0x3e, 0xd9, 0xe7, 0xfa, 0xbc, 0x44, 0x66, 0x13, 0xe9, 0x5b, 0xe8, 0xdc, 0x70,
	0x95, 0xe9, 0xb9, 0x2d, 0xb8, 0xbf, 0xc3, 0x53, 0x47, 0x71, 0x49, 0xcc, 0x27, 0xc7, 0x7f, 0x40,
	0x6f, 0xdd, 0x6b, 0x3a, 0x06, 0x28, 0xb9, 0xe2, 0x05, 0x6a, 0x54, 0x8b, 0x9e, 0xc7, 0x0d, 0x29,
	0x63, 0x67, 0x26, 0xee, 0x6f, 0x96, 0x99, 0x6c, 0x8d, 0x45, 0x29, 0x1c, 0x4c, 0x4b, 0xa9, 0xfd,
	0xc7, 0xd8, 0x67, 0x13, 0x4b, 0xb2, 0x02, 0xbd, 0xc7, 0xf6, 0x39, 0xfe, 0x1a, 0x9e, 0x37, 0x84,
	0x8c, 0xef, 0xd3, 0xec, 0xbe, 0xe0, 0x8b, 0xce, 0x59, 0x60, 0xcc, 0x61, 0x0b, 0x73, 0x58, 0xfc,
	0x4f, 0x00, 0x2f, 0xfc, 0x25, 0x1e, 0x4b, 0xa9, 0x2b, 0xad, 0x78, 0x39, 0xa9, 0xd5, 0x23, 0xd2,
	0x77, 0xd0, 0x29, 0x50, 0x3f, 0xc8, 0xd4, 0xd2, 0xfb, 0xa3, 0x41, 0xa3, 0xf4, 0x25, 0xe1, 0xda,
	0xe6, 0x31, 0x9f, 0x4f, 0x4f, 0xa0, 0x93, 0xf3, 0xe2, 0x36, 0x75, 0x97, 0x83, 0x30, 0x8f, 0x68,
	0x1f, 0x02, 0x7d, 0xee, 0xcb, 0x0e, 0xf4, 0x39, 0xbd, 0x82, 0x67, 0xb7, 0x0b, 0x89, 0x4b, 0xee,
	0xef, 0xc7, 0xd1, 0xe8, 0x75, 0xe3, 0x20, 0x5b, 0xd0, 0x78, 0x3d, 0x95, 0x6d, 0x32, 0xe9, 0x3b,
	0xe8, 0x6a, 0x14, 0x52, 0x59, 0x99, 0xb6, 0x95, 0x79, 0xd9, 0x90, 0x49, 0x6c, 0x06, 0xde, 0x55,
	0x6c, 0x95, 0x4c, 0xbf, 0x81, 0x9e, 0xac, 0x75, 0x59, 0x6b, 0xfb, 0xb6, 0xb2, 0xf7, 0xee, 0x69,
	0xf2, 0x46, 0x7e, 0xfc, 0x1f, 0x81, 0x93, 0xc5, 0x0c, 0x6d, 0x39, 0xf8, 0x15, 0x74, 0x4d, 0xc3,
	0x2c, 0xf0, 0xfd, 0x3f, 0xd9, 0xfd, 0x6d, 0x6c, 0x95, 0x48, 0x7f, 0x80, 0x17, 0x57, 0x42, 0xa3,
	0x2a, 0xa5, 0x99, 0xc6, 0x74, 0xa5, 0x10, 0x0c, 0xc8, 0x13, 0x0a, 0xbb, 0x49, 0x94, 0x41, 0xb4,
	0xfe, 0xe2, 0xbd, 0x54, 0xbf, 0x71, 0x95, 0x3a, 0xc1, 0xd6, 0x93, 0x82, 0x7b, 0x79, 0xf1, 0x6b,
	0xe8, 0x2e, 0xdd, 0x30, 0xcd, 0xd6, 0xce, 0x39, 0x32, 0x68, 0x0d, 0x03, 0xe6, 0x51, 0xfc, 0x08,
	0xb4, 0xd9, 0x36, 0xfa, 0x3d, 0x1c, 0xdf, 0x4a, 0x91, 0x5e, 0xe2, 0x5d, 0x26, 0x96, 0xcb, 0xb0,
	0x35, 0x3c, 0x1a, 0x7d, 0xd6, 0xac, 0x42, 0xd6, 0xa5, 0x14, 0x63, 0x97, 0xcd, 0xb6, 0x69, 0xe6,
	0xdc, 0x0f, 0x19, 0xe6, 0x69, 0x15, 0x05, 0xee, 0x5c, 0x87, 0xe2, 0x3f, 0x09, 0x3c, 0xdb, 0xa0,
	0xda, 0xb5, 0x52, 0x55, 0x35, 0xda, 0xa1, 0x21, 0x7e, 0xad, 0x2c, 0x02, 0xf4, 0x25, 0x1c, 0x5e,
	0x73, 0x5d, 0xdb, 0x71, 0x77, 0x53, 0xb1, 0xc4, 0x74, 0x08, 0xc7, 0x4e, 0xea, 0xbd, 0x32, 0x23,
	0x22, 0x66, 0x73, 0x7f, 0x7b, 0xb7, 0xc3, 0xa6, 0x1a, 0x17, 0xf2, 0x3b, 0xce, 0xa3, 0xf8, 0x2d,
	0xb4, 0x5d, 0x1f, 0x4e, 0xa0, 0x93, 0x6c, 0xd8, 0xe4, 0x90, 0x99, 0x51, 0xc6, 0x35, 0x2e, 0xbe,
	0xc2, 0x81, 0xd3, 0x57, 0xd0, 0xbf, 0x12, 0x95, 0x56, 0x75, 0x81, 0x42, 0xdb, 0x9d, 0x04, 0xd0,
	0x71, 0x4b, 0x27, 0x24, 0xa7, 0x43, 0x80, 0xd5, 0xce, 0xa2, 0x3d, 0x38, 0xbc, 0x28, 0x50, 0x65,
	0x33, 0x2e, 0x42, 0x62, 0xd0, 0xb7, 0xb5, 0x92, 0x25, 0x72, 0x11, 0x06, 0xa7, 0x5f, 0x40, 0x6f,
	0x7d, 0x55, 0xd1, 0x43, 0x38, 0x98, 0xf0, 0x3c, 0x0f, 0x09, 0xfd, 0x08, 0x5a, 0x37, 0xb5, 0x0e,
	0x83, 0xd3, 0x73, 0x38, 0xde, 0x9a, 0x63, 0xda, 0x85, 0xf6, 0x8f, 0x3c, 0x7b, 0xc4, 0x90, 0x50,
	0x0a, 0xfd, 0x6b, 0x29, 0xa4, 0x96, 0x02, 0x27, 0x52, 0x3c, 0xe2, 0xef, 0x61, 0x30, 0x7e, 0xfe,
	0x8b, 0xfb, 0x2b, 0x7c, 0xb3, 0xec, 0xd5, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0xda, 0x7b, 0x21,
	0x25, 0x27, 0x07, 0x00, 0x00,
}

// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/house/house.proto

package house

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for House service

type HouseService interface {
	// 查询房屋信息
	GetHouseInfo(ctx context.Context, in *InfoReq, opts ...client.CallOption) (*InfoResp, error)
	// 发布房屋信息
	PubHouse(ctx context.Context, in *PubReq, opts ...client.CallOption) (*PubResp, error)
	// 上传房屋图片
	UploadHouseImg(ctx context.Context, in *ImgReq, opts ...client.CallOption) (*ImgResp, error)
	// 查询房屋具体信息
	GetHouseDetail(ctx context.Context, in *DetailReq, opts ...client.CallOption) (*DetailResp, error)
	// 首页轮播图
	GetHouseIndex(ctx context.Context, in *IndexReq, opts ...client.CallOption) (*GetResp, error)
	// 搜索房屋
	SearchHouse(ctx context.Context, in *SearchReq, opts ...client.CallOption) (*GetResp, error)
}

type houseService struct {
	c    client.Client
	name string
}

func NewHouseService(name string, c client.Client) HouseService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "go.micro.srv.house"
	}
	return &houseService{
		c:    c,
		name: name,
	}
}

func (c *houseService) GetHouseInfo(ctx context.Context, in *InfoReq, opts ...client.CallOption) (*InfoResp, error) {
	req := c.c.NewRequest(c.name, "House.GetHouseInfo", in)
	out := new(InfoResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *houseService) PubHouse(ctx context.Context, in *PubReq, opts ...client.CallOption) (*PubResp, error) {
	req := c.c.NewRequest(c.name, "House.PubHouse", in)
	out := new(PubResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *houseService) UploadHouseImg(ctx context.Context, in *ImgReq, opts ...client.CallOption) (*ImgResp, error) {
	req := c.c.NewRequest(c.name, "House.UploadHouseImg", in)
	out := new(ImgResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *houseService) GetHouseDetail(ctx context.Context, in *DetailReq, opts ...client.CallOption) (*DetailResp, error) {
	req := c.c.NewRequest(c.name, "House.GetHouseDetail", in)
	out := new(DetailResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *houseService) GetHouseIndex(ctx context.Context, in *IndexReq, opts ...client.CallOption) (*GetResp, error) {
	req := c.c.NewRequest(c.name, "House.GetHouseIndex", in)
	out := new(GetResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *houseService) SearchHouse(ctx context.Context, in *SearchReq, opts ...client.CallOption) (*GetResp, error) {
	req := c.c.NewRequest(c.name, "House.SearchHouse", in)
	out := new(GetResp)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for House service

type HouseHandler interface {
	// 查询房屋信息
	GetHouseInfo(context.Context, *InfoReq, *InfoResp) error
	// 发布房屋信息
	PubHouse(context.Context, *PubReq, *PubResp) error
	// 上传房屋图片
	UploadHouseImg(context.Context, *ImgReq, *ImgResp) error
	// 查询房屋具体信息
	GetHouseDetail(context.Context, *DetailReq, *DetailResp) error
	// 首页轮播图
	GetHouseIndex(context.Context, *IndexReq, *GetResp) error
	// 搜索房屋
	SearchHouse(context.Context, *SearchReq, *GetResp) error
}

func RegisterHouseHandler(s server.Server, hdlr HouseHandler, opts ...server.HandlerOption) error {
	type house interface {
		GetHouseInfo(ctx context.Context, in *InfoReq, out *InfoResp) error
		PubHouse(ctx context.Context, in *PubReq, out *PubResp) error
		UploadHouseImg(ctx context.Context, in *ImgReq, out *ImgResp) error
		GetHouseDetail(ctx context.Context, in *DetailReq, out *DetailResp) error
		GetHouseIndex(ctx context.Context, in *IndexReq, out *GetResp) error
		SearchHouse(ctx context.Context, in *SearchReq, out *GetResp) error
	}
	type House struct {
		house
	}
	h := &houseHandler{hdlr}
	return s.Handle(s.NewHandler(&House{h}, opts...))
}

type houseHandler struct {
	HouseHandler
}

func (h *houseHandler) GetHouseInfo(ctx context.Context, in *InfoReq, out *InfoResp) error {
	return h.HouseHandler.GetHouseInfo(ctx, in, out)
}

func (h *houseHandler) PubHouse(ctx context.Context, in *PubReq, out *PubResp) error {
	return h.HouseHandler.PubHouse(ctx, in, out)
}

func (h *houseHandler) UploadHouseImg(ctx context.Context, in *ImgReq, out *ImgResp) error {
	return h.HouseHandler.UploadHouseImg(ctx, in, out)
}

func (h *houseHandler) GetHouseDetail(ctx context.Context, in *DetailReq, out *DetailResp) error {
	return h.HouseHandler.GetHouseDetail(ctx, in, out)
}

func (h *houseHandler) GetHouseIndex(ctx context.Context, in *IndexReq, out *GetResp) error {
	return h.HouseHandler.GetHouseIndex(ctx, in, out)
}

func (h *houseHandler) SearchHouse(ctx context.Context, in *SearchReq, out *GetResp) error {
	return h.HouseHandler.SearchHouse(ctx, in, out)
}

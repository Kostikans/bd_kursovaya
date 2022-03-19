// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: api/forum/forum.proto

package forum

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ForumClient is the client API for Forum service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ForumClient interface {
	CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error)
	DeleteAccount(ctx context.Context, in *DeleteAccountRequest, opts ...grpc.CallOption) (*DeleteAccountResponse, error)
	CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostResponse, error)
	CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error)
	CreateTag(ctx context.Context, in *CreateTagRequest, opts ...grpc.CallOption) (*CreateTagResponse, error)
	AssignTagsToPost(ctx context.Context, in *AssignTagsToPostRequest, opts ...grpc.CallOption) (*AssignTagsToPostResponse, error)
	CreatePostVote(ctx context.Context, in *CreatePostVoteRequest, opts ...grpc.CallOption) (*CreatePostVoteResponse, error)
	CreateCommentVote(ctx context.Context, in *CreateCommentVoteRequest, opts ...grpc.CallOption) (*CreateCommentVoteResponse, error)
	GetPosts(ctx context.Context, in *GetPostListRequest, opts ...grpc.CallOption) (*GetPostListResponse, error)
	GetComments(ctx context.Context, in *GetCommentListRequest, opts ...grpc.CallOption) (*GetCommentListResponse, error)
	Truncate(ctx context.Context, in *TruncateRequest, opts ...grpc.CallOption) (*TruncateResponse, error)
}

type forumClient struct {
	cc grpc.ClientConnInterface
}

func NewForumClient(cc grpc.ClientConnInterface) ForumClient {
	return &forumClient{cc}
}

func (c *forumClient) CreateAccount(ctx context.Context, in *CreateAccountRequest, opts ...grpc.CallOption) (*CreateAccountResponse, error) {
	out := new(CreateAccountResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/CreateAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) DeleteAccount(ctx context.Context, in *DeleteAccountRequest, opts ...grpc.CallOption) (*DeleteAccountResponse, error) {
	out := new(DeleteAccountResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/DeleteAccount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) CreatePost(ctx context.Context, in *CreatePostRequest, opts ...grpc.CallOption) (*CreatePostResponse, error) {
	out := new(CreatePostResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/CreatePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) CreateComment(ctx context.Context, in *CreateCommentRequest, opts ...grpc.CallOption) (*CreateCommentResponse, error) {
	out := new(CreateCommentResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/CreateComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) CreateTag(ctx context.Context, in *CreateTagRequest, opts ...grpc.CallOption) (*CreateTagResponse, error) {
	out := new(CreateTagResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/CreateTag", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) AssignTagsToPost(ctx context.Context, in *AssignTagsToPostRequest, opts ...grpc.CallOption) (*AssignTagsToPostResponse, error) {
	out := new(AssignTagsToPostResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/AssignTagsToPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) CreatePostVote(ctx context.Context, in *CreatePostVoteRequest, opts ...grpc.CallOption) (*CreatePostVoteResponse, error) {
	out := new(CreatePostVoteResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/CreatePostVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) CreateCommentVote(ctx context.Context, in *CreateCommentVoteRequest, opts ...grpc.CallOption) (*CreateCommentVoteResponse, error) {
	out := new(CreateCommentVoteResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/CreateCommentVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) GetPosts(ctx context.Context, in *GetPostListRequest, opts ...grpc.CallOption) (*GetPostListResponse, error) {
	out := new(GetPostListResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/GetPosts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) GetComments(ctx context.Context, in *GetCommentListRequest, opts ...grpc.CallOption) (*GetCommentListResponse, error) {
	out := new(GetCommentListResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/GetComments", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *forumClient) Truncate(ctx context.Context, in *TruncateRequest, opts ...grpc.CallOption) (*TruncateResponse, error) {
	out := new(TruncateResponse)
	err := c.cc.Invoke(ctx, "/forum.Forum/Truncate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ForumServer is the server API for Forum service.
// All implementations must embed UnimplementedForumServer
// for forward compatibility
type ForumServer interface {
	CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error)
	DeleteAccount(context.Context, *DeleteAccountRequest) (*DeleteAccountResponse, error)
	CreatePost(context.Context, *CreatePostRequest) (*CreatePostResponse, error)
	CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error)
	CreateTag(context.Context, *CreateTagRequest) (*CreateTagResponse, error)
	AssignTagsToPost(context.Context, *AssignTagsToPostRequest) (*AssignTagsToPostResponse, error)
	CreatePostVote(context.Context, *CreatePostVoteRequest) (*CreatePostVoteResponse, error)
	CreateCommentVote(context.Context, *CreateCommentVoteRequest) (*CreateCommentVoteResponse, error)
	GetPosts(context.Context, *GetPostListRequest) (*GetPostListResponse, error)
	GetComments(context.Context, *GetCommentListRequest) (*GetCommentListResponse, error)
	Truncate(context.Context, *TruncateRequest) (*TruncateResponse, error)
	mustEmbedUnimplementedForumServer()
}

// UnimplementedForumServer must be embedded to have forward compatible implementations.
type UnimplementedForumServer struct {
}

func (UnimplementedForumServer) CreateAccount(context.Context, *CreateAccountRequest) (*CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (UnimplementedForumServer) DeleteAccount(context.Context, *DeleteAccountRequest) (*DeleteAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccount not implemented")
}
func (UnimplementedForumServer) CreatePost(context.Context, *CreatePostRequest) (*CreatePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}
func (UnimplementedForumServer) CreateComment(context.Context, *CreateCommentRequest) (*CreateCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateComment not implemented")
}
func (UnimplementedForumServer) CreateTag(context.Context, *CreateTagRequest) (*CreateTagResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTag not implemented")
}
func (UnimplementedForumServer) AssignTagsToPost(context.Context, *AssignTagsToPostRequest) (*AssignTagsToPostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignTagsToPost not implemented")
}
func (UnimplementedForumServer) CreatePostVote(context.Context, *CreatePostVoteRequest) (*CreatePostVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePostVote not implemented")
}
func (UnimplementedForumServer) CreateCommentVote(context.Context, *CreateCommentVoteRequest) (*CreateCommentVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCommentVote not implemented")
}
func (UnimplementedForumServer) GetPosts(context.Context, *GetPostListRequest) (*GetPostListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPosts not implemented")
}
func (UnimplementedForumServer) GetComments(context.Context, *GetCommentListRequest) (*GetCommentListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetComments not implemented")
}
func (UnimplementedForumServer) Truncate(context.Context, *TruncateRequest) (*TruncateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Truncate not implemented")
}
func (UnimplementedForumServer) mustEmbedUnimplementedForumServer() {}

// UnsafeForumServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ForumServer will
// result in compilation errors.
type UnsafeForumServer interface {
	mustEmbedUnimplementedForumServer()
}

func RegisterForumServer(s grpc.ServiceRegistrar, srv ForumServer) {
	s.RegisterService(&Forum_ServiceDesc, srv)
}

func _Forum_CreateAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).CreateAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/CreateAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).CreateAccount(ctx, req.(*CreateAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_DeleteAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAccountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).DeleteAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/DeleteAccount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).DeleteAccount(ctx, req.(*DeleteAccountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_CreatePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).CreatePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/CreatePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).CreatePost(ctx, req.(*CreatePostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_CreateComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).CreateComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/CreateComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).CreateComment(ctx, req.(*CreateCommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_CreateTag_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTagRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).CreateTag(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/CreateTag",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).CreateTag(ctx, req.(*CreateTagRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_AssignTagsToPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignTagsToPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).AssignTagsToPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/AssignTagsToPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).AssignTagsToPost(ctx, req.(*AssignTagsToPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_CreatePostVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePostVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).CreatePostVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/CreatePostVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).CreatePostVote(ctx, req.(*CreatePostVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_CreateCommentVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCommentVoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).CreateCommentVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/CreateCommentVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).CreateCommentVote(ctx, req.(*CreateCommentVoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_GetPosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPostListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).GetPosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/GetPosts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).GetPosts(ctx, req.(*GetPostListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_GetComments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCommentListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).GetComments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/GetComments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).GetComments(ctx, req.(*GetCommentListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Forum_Truncate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TruncateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ForumServer).Truncate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/forum.Forum/Truncate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ForumServer).Truncate(ctx, req.(*TruncateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Forum_ServiceDesc is the grpc.ServiceDesc for Forum service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Forum_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "forum.Forum",
	HandlerType: (*ForumServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAccount",
			Handler:    _Forum_CreateAccount_Handler,
		},
		{
			MethodName: "DeleteAccount",
			Handler:    _Forum_DeleteAccount_Handler,
		},
		{
			MethodName: "CreatePost",
			Handler:    _Forum_CreatePost_Handler,
		},
		{
			MethodName: "CreateComment",
			Handler:    _Forum_CreateComment_Handler,
		},
		{
			MethodName: "CreateTag",
			Handler:    _Forum_CreateTag_Handler,
		},
		{
			MethodName: "AssignTagsToPost",
			Handler:    _Forum_AssignTagsToPost_Handler,
		},
		{
			MethodName: "CreatePostVote",
			Handler:    _Forum_CreatePostVote_Handler,
		},
		{
			MethodName: "CreateCommentVote",
			Handler:    _Forum_CreateCommentVote_Handler,
		},
		{
			MethodName: "GetPosts",
			Handler:    _Forum_GetPosts_Handler,
		},
		{
			MethodName: "GetComments",
			Handler:    _Forum_GetComments_Handler,
		},
		{
			MethodName: "Truncate",
			Handler:    _Forum_Truncate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/forum/forum.proto",
}

// Code generated from quiz/v1/quiz_result.proto. DO NOT EDIT.

package quizresultv1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	QuizResultService_SubmitQuiz_FullMethodName      = "/quiz.v1.QuizResultService/SubmitQuiz"
	QuizResultService_GetQuizResults_FullMethodName  = "/quiz.v1.QuizResultService/GetQuizResults"
)

// QuizResultServiceClient is the client API for QuizResultService.
type QuizResultServiceClient interface {
	SubmitQuiz(ctx context.Context, in *SubmitQuizRequest, opts ...grpc.CallOption) (*SubmitQuizResponse, error)
	GetQuizResults(ctx context.Context, in *GetQuizResultsRequest, opts ...grpc.CallOption) (*GetQuizResultsResponse, error)
}

type quizResultServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQuizResultServiceClient(cc grpc.ClientConnInterface) QuizResultServiceClient {
	return &quizResultServiceClient{cc}
}

func (c *quizResultServiceClient) SubmitQuiz(ctx context.Context, in *SubmitQuizRequest, opts ...grpc.CallOption) (*SubmitQuizResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SubmitQuizResponse)
	if err := c.cc.Invoke(ctx, QuizResultService_SubmitQuiz_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *quizResultServiceClient) GetQuizResults(ctx context.Context, in *GetQuizResultsRequest, opts ...grpc.CallOption) (*GetQuizResultsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetQuizResultsResponse)
	if err := c.cc.Invoke(ctx, QuizResultService_GetQuizResults_FullMethodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

// QuizResultServiceServer is the server API for QuizResultService.
// All implementations must embed UnimplementedQuizResultServiceServer
// for forward compatibility.
type QuizResultServiceServer interface {
	SubmitQuiz(context.Context, *SubmitQuizRequest) (*SubmitQuizResponse, error)
	GetQuizResults(context.Context, *GetQuizResultsRequest) (*GetQuizResultsResponse, error)
	mustEmbedUnimplementedQuizResultServiceServer()
}

// UnimplementedQuizResultServiceServer must be embedded to have forward compatible implementations.
type UnimplementedQuizResultServiceServer struct{}

func (UnimplementedQuizResultServiceServer) SubmitQuiz(context.Context, *SubmitQuizRequest) (*SubmitQuizResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method SubmitQuiz not implemented")
}
func (UnimplementedQuizResultServiceServer) GetQuizResults(context.Context, *GetQuizResultsRequest) (*GetQuizResultsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetQuizResults not implemented")
}
func (UnimplementedQuizResultServiceServer) mustEmbedUnimplementedQuizResultServiceServer() {}
func (UnimplementedQuizResultServiceServer) testEmbeddedByValue()                           {}

// UnsafeQuizResultServiceServer may be embedded to opt out of forward compatibility.
type UnsafeQuizResultServiceServer interface {
	mustEmbedUnimplementedQuizResultServiceServer()
}

func RegisterQuizResultServiceServer(s grpc.ServiceRegistrar, srv QuizResultServiceServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&QuizResultService_ServiceDesc, srv)
}

func _QuizResultService_SubmitQuiz_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(SubmitQuizRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuizResultServiceServer).SubmitQuiz(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: QuizResultService_SubmitQuiz_FullMethodName}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(QuizResultServiceServer).SubmitQuiz(ctx, req.(*SubmitQuizRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuizResultService_GetQuizResults_Handler(srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor) (any, error) {
	in := new(GetQuizResultsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuizResultServiceServer).GetQuizResults(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: QuizResultService_GetQuizResults_FullMethodName}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(QuizResultServiceServer).GetQuizResults(ctx, req.(*GetQuizResultsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var QuizResultService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "quiz.v1.QuizResultService",
	HandlerType: (*QuizResultServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "SubmitQuiz", Handler: _QuizResultService_SubmitQuiz_Handler},
		{MethodName: "GetQuizResults", Handler: _QuizResultService_GetQuizResults_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "quiz/v1/quiz_result.proto",
}

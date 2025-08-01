/*
 *
 * Copyright 2016 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package stats is for collecting and reporting various network and RPC stats.
// This package is for monitoring purpose only. All fields are read-only.
// All APIs are experimental.
package stats // import "google.golang.org/grpc/stats"

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc/metadata"
)

// RPCStats contains stats information about RPCs.
type RPCStats interface {
	isRPCStats()
	// IsClient returns true if this RPCStats is from client side.
	IsClient() bool
}

// Begin contains stats for the start of an RPC attempt.
//
//   - Server-side: Triggered after `InHeader`, as headers are processed
//     before the RPC lifecycle begins.
//   - Client-side: The first stats event recorded.
//
// FailFast is only valid if this Begin is from client side.
type Begin struct {
	// Client is true if this Begin is from client side.
	Client bool
	// BeginTime is the time when the RPC attempt begins.
	BeginTime time.Time
	// FailFast indicates if this RPC is failfast.
	FailFast bool
	// IsClientStream indicates whether the RPC is a client streaming RPC.
	IsClientStream bool
	// IsServerStream indicates whether the RPC is a server streaming RPC.
	IsServerStream bool
	// IsTransparentRetryAttempt indicates whether this attempt was initiated
	// due to transparently retrying a previous attempt.
	IsTransparentRetryAttempt bool
}

// IsClient indicates if the stats information is from client side.
func (s *Begin) IsClient() bool { return s.Client }

func (s *Begin) isRPCStats() {}

// DelayedPickComplete indicates that the RPC is unblocked following a delay in
// selecting a connection for the call.
type DelayedPickComplete struct{}

// IsClient indicates DelayedPickComplete is available on the client.
func (*DelayedPickComplete) IsClient() bool { return true }

func (*DelayedPickComplete) isRPCStats() {}

// PickerUpdated indicates that the RPC is unblocked following a delay in
// selecting a connection for the call.
//
// Deprecated: will be removed in a future release; use DelayedPickComplete
// instead.
type PickerUpdated = DelayedPickComplete

// InPayload contains stats about an incoming payload.
type InPayload struct {
	// Client is true if this InPayload is from client side.
	Client bool
	// Payload is the payload with original type.  This may be modified after
	// the call to HandleRPC which provides the InPayload returns and must be
	// copied if needed later.
	Payload any

	// Length is the size of the uncompressed payload data. Does not include any
	// framing (gRPC or HTTP/2).
	Length int
	// CompressedLength is the size of the compressed payload data. Does not
	// include any framing (gRPC or HTTP/2). Same as Length if compression not
	// enabled.
	CompressedLength int
	// WireLength is the size of the compressed payload data plus gRPC framing.
	// Does not include HTTP/2 framing.
	WireLength int

	// RecvTime is the time when the payload is received.
	RecvTime time.Time
}

// IsClient indicates if the stats information is from client side.
func (s *InPayload) IsClient() bool { return s.Client }

func (s *InPayload) isRPCStats() {}

// InHeader contains stats about header reception.
//
// - Server-side: The first stats event after the RPC request is received.
type InHeader struct {
	// Client is true if this InHeader is from client side.
	Client bool
	// WireLength is the wire length of header.
	WireLength int
	// Compression is the compression algorithm used for the RPC.
	Compression string
	// Header contains the header metadata received.
	Header metadata.MD

	// The following fields are valid only if Client is false.
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
	// RemoteAddr is the remote address of the corresponding connection.
	RemoteAddr net.Addr
	// LocalAddr is the local address of the corresponding connection.
	LocalAddr net.Addr
}

// IsClient indicates if the stats information is from client side.
func (s *InHeader) IsClient() bool { return s.Client }

func (s *InHeader) isRPCStats() {}

// InTrailer contains stats about trailer reception.
type InTrailer struct {
	// Client is true if this InTrailer is from client side.
	Client bool
	// WireLength is the wire length of trailer.
	WireLength int
	// Trailer contains the trailer metadata received from the server. This
	// field is only valid if this InTrailer is from the client side.
	Trailer metadata.MD
}

// IsClient indicates if the stats information is from client side.
func (s *InTrailer) IsClient() bool { return s.Client }

func (s *InTrailer) isRPCStats() {}

// OutPayload contains stats about an outgoing payload.
type OutPayload struct {
	// Client is true if this OutPayload is from client side.
	Client bool
	// Payload is the payload with original type.  This may be modified after
	// the call to HandleRPC which provides the OutPayload returns and must be
	// copied if needed later.
	Payload any
	// Length is the size of the uncompressed payload data. Does not include any
	// framing (gRPC or HTTP/2).
	Length int
	// CompressedLength is the size of the compressed payload data. Does not
	// include any framing (gRPC or HTTP/2). Same as Length if compression not
	// enabled.
	CompressedLength int
	// WireLength is the size of the compressed payload data plus gRPC framing.
	// Does not include HTTP/2 framing.
	WireLength int
	// SentTime is the time when the payload is sent.
	SentTime time.Time
}

// IsClient indicates if this stats information is from client side.
func (s *OutPayload) IsClient() bool { return s.Client }

func (s *OutPayload) isRPCStats() {}

// OutHeader contains stats about header transmission.
//
//   - Client-side: Only occurs after 'Begin', as headers are always the first
//     thing sent on a stream.
type OutHeader struct {
	// Client is true if this OutHeader is from client side.
	Client bool
	// Compression is the compression algorithm used for the RPC.
	Compression string
	// Header contains the header metadata sent.
	Header metadata.MD

	// The following fields are valid only if Client is true.
	// FullMethod is the full RPC method string, i.e., /package.service/method.
	FullMethod string
	// RemoteAddr is the remote address of the corresponding connection.
	RemoteAddr net.Addr
	// LocalAddr is the local address of the corresponding connection.
	LocalAddr net.Addr
}

// IsClient indicates if this stats information is from client side.
func (s *OutHeader) IsClient() bool { return s.Client }

func (s *OutHeader) isRPCStats() {}

// OutTrailer contains stats about trailer transmission.
type OutTrailer struct {
	// Client is true if this OutTrailer is from client side.
	Client bool
	// WireLength is the wire length of trailer.
	//
	// Deprecated: This field is never set. The length is not known when this
	// message is emitted because the trailer fields are compressed with hpack
	// after that.
	WireLength int
	// Trailer contains the trailer metadata sent to the client. This
	// field is only valid if this OutTrailer is from the server side.
	Trailer metadata.MD
}

// IsClient indicates if this stats information is from client side.
func (s *OutTrailer) IsClient() bool { return s.Client }

func (s *OutTrailer) isRPCStats() {}

// End contains stats about RPC completion.
type End struct {
	// Client is true if this End is from client side.
	Client bool
	// BeginTime is the time when the RPC began.
	BeginTime time.Time
	// EndTime is the time when the RPC ends.
	EndTime time.Time
	// Trailer contains the trailer metadata received from the server. This
	// field is only valid if this End is from the client side.
	// Deprecated: use Trailer in InTrailer instead.
	Trailer metadata.MD
	// Error is the error the RPC ended with. It is an error generated from
	// status.Status and can be converted back to status.Status using
	// status.FromError if non-nil.
	Error error
}

// IsClient indicates if this is from client side.
func (s *End) IsClient() bool { return s.Client }

func (s *End) isRPCStats() {}

// ConnStats contains stats information about connections.
type ConnStats interface {
	isConnStats()
	// IsClient returns true if this ConnStats is from client side.
	IsClient() bool
}

// ConnBegin contains stats about connection establishment.
type ConnBegin struct {
	// Client is true if this ConnBegin is from client side.
	Client bool
}

// IsClient indicates if this is from client side.
func (s *ConnBegin) IsClient() bool { return s.Client }

func (s *ConnBegin) isConnStats() {}

// ConnEnd contains stats about connection termination.
type ConnEnd struct {
	// Client is true if this ConnEnd is from client side.
	Client bool
}

// IsClient indicates if this is from client side.
func (s *ConnEnd) IsClient() bool { return s.Client }

func (s *ConnEnd) isConnStats() {}

// SetTags attaches stats tagging data to the context, which will be sent in
// the outgoing RPC with the header grpc-tags-bin.  Subsequent calls to
// SetTags will overwrite the values from earlier calls.
//
// Deprecated: set the `grpc-tags-bin` header in the metadata instead.
func SetTags(ctx context.Context, b []byte) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "grpc-tags-bin", string(b))
}

// Tags returns the tags from the context for the inbound RPC.
//
// Deprecated: obtain the `grpc-tags-bin` header from metadata instead.
func Tags(ctx context.Context) []byte {
	traceValues := metadata.ValueFromIncomingContext(ctx, "grpc-tags-bin")
	if len(traceValues) == 0 {
		return nil
	}
	return []byte(traceValues[len(traceValues)-1])
}

// SetTrace attaches stats tagging data to the context, which will be sent in
// the outgoing RPC with the header grpc-trace-bin.  Subsequent calls to
// SetTrace will overwrite the values from earlier calls.
//
// Deprecated: set the `grpc-trace-bin` header in the metadata instead.
func SetTrace(ctx context.Context, b []byte) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "grpc-trace-bin", string(b))
}

// Trace returns the trace from the context for the inbound RPC.
//
// Deprecated: obtain the `grpc-trace-bin` header from metadata instead.
func Trace(ctx context.Context) []byte {
	traceValues := metadata.ValueFromIncomingContext(ctx, "grpc-trace-bin")
	if len(traceValues) == 0 {
		return nil
	}
	return []byte(traceValues[len(traceValues)-1])
}

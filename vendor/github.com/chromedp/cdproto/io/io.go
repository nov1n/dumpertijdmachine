// Package io provides the Chrome Debugging Protocol
// commands, types, and events for the IO domain.
//
// Input/Output operations for streams produced by DevTools.
//
// Generated by the cdproto-gen command.
package io

// Code generated by cdproto-gen. DO NOT EDIT.

import (
	"context"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/runtime"
)

// CloseParams close the stream, discard any temporary backing storage.
type CloseParams struct {
	Handle StreamHandle `json:"handle"` // Handle of the stream to close.
}

// Close close the stream, discard any temporary backing storage.
//
// parameters:
//   handle - Handle of the stream to close.
func Close(handle StreamHandle) *CloseParams {
	return &CloseParams{
		Handle: handle,
	}
}

// Do executes IO.close against the provided context.
func (p *CloseParams) Do(ctxt context.Context, h cdp.Executor) (err error) {
	return h.Execute(ctxt, CommandClose, p, nil)
}

// ReadParams read a chunk of the stream.
type ReadParams struct {
	Handle StreamHandle `json:"handle"`           // Handle of the stream to read.
	Offset int64        `json:"offset,omitempty"` // Seek to the specified offset before reading (if not specificed, proceed with offset following the last read).
	Size   int64        `json:"size,omitempty"`   // Maximum number of bytes to read (left upon the agent discretion if not specified).
}

// Read read a chunk of the stream.
//
// parameters:
//   handle - Handle of the stream to read.
func Read(handle StreamHandle) *ReadParams {
	return &ReadParams{
		Handle: handle,
	}
}

// WithOffset seek to the specified offset before reading (if not specificed,
// proceed with offset following the last read).
func (p ReadParams) WithOffset(offset int64) *ReadParams {
	p.Offset = offset
	return &p
}

// WithSize maximum number of bytes to read (left upon the agent discretion
// if not specified).
func (p ReadParams) WithSize(size int64) *ReadParams {
	p.Size = size
	return &p
}

// ReadReturns return values.
type ReadReturns struct {
	Base64encoded bool   `json:"base64Encoded,omitempty"` // Set if the data is base64-encoded
	Data          string `json:"data,omitempty"`          // Data that were read.
	EOF           bool   `json:"eof,omitempty"`           // Set if the end-of-file condition occurred while reading.
}

// Do executes IO.read against the provided context.
//
// returns:
//   data - Data that were read.
//   eof - Set if the end-of-file condition occurred while reading.
func (p *ReadParams) Do(ctxt context.Context, h cdp.Executor) (data string, eof bool, err error) {
	// execute
	var res ReadReturns
	err = h.Execute(ctxt, CommandRead, p, &res)
	if err != nil {
		return "", false, err
	}

	return res.Data, res.EOF, nil
}

// ResolveBlobParams return UUID of Blob object specified by a remote object
// id.
type ResolveBlobParams struct {
	ObjectID runtime.RemoteObjectID `json:"objectId"` // Object id of a Blob object wrapper.
}

// ResolveBlob return UUID of Blob object specified by a remote object id.
//
// parameters:
//   objectID - Object id of a Blob object wrapper.
func ResolveBlob(objectID runtime.RemoteObjectID) *ResolveBlobParams {
	return &ResolveBlobParams{
		ObjectID: objectID,
	}
}

// ResolveBlobReturns return values.
type ResolveBlobReturns struct {
	UUID string `json:"uuid,omitempty"` // UUID of the specified Blob.
}

// Do executes IO.resolveBlob against the provided context.
//
// returns:
//   uuid - UUID of the specified Blob.
func (p *ResolveBlobParams) Do(ctxt context.Context, h cdp.Executor) (uuid string, err error) {
	// execute
	var res ResolveBlobReturns
	err = h.Execute(ctxt, CommandResolveBlob, p, &res)
	if err != nil {
		return "", err
	}

	return res.UUID, nil
}

// Command names.
const (
	CommandClose       = "IO.close"
	CommandRead        = "IO.read"
	CommandResolveBlob = "IO.resolveBlob"
)

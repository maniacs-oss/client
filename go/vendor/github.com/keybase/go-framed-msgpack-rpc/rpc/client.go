package rpc

import (
	"errors"

	"golang.org/x/net/context"
)

// Client allows calls and notifies on the given transporter, or any protocol
// type. All will share the same ErrorUnwrapper hook for unwrapping incoming
// msgpack objects and converting to possible Go-native `Error` types
type Client struct {
	xp             Transporter
	errorUnwrapper ErrorUnwrapper
	tagsFunc       LogTagsFromContext
}

// NewClient constructs a new client from the given RPC Transporter and the
// ErrorUnwrapper.
func NewClient(xp Transporter, u ErrorUnwrapper,
	tagsFunc LogTagsFromContext) *Client {
	return &Client{xp, u, tagsFunc}
}

// Call makes an msgpack RPC call over the transports that's bound to this
// client. The name of the method, and the argument are given. On reply,
// the result field will be populated (if applicable). It returns an Error
// on error, where the error might have been unwrapped from Msgpack via the
// UnwrapErrorFunc in this client.
func (c *Client) Call(ctx context.Context, method string, arg interface{}, res interface{}) (err error) {
	if ctx == nil {
		return errors.New("No Context provided for this call")
	}

	if c.tagsFunc != nil {
		tags, ok := c.tagsFunc(ctx)
		if ok {
			rpcTags := make(CtxRpcTags)
			for key, tagName := range tags {
				if v := ctx.Value(key); v != nil {
					rpcTags[tagName] = v
				}
			}
			ctx = AddRpcTagsToContext(ctx, rpcTags)
		}
	}

	c.xp.receiveFrames()
	d, err := c.xp.getDispatcher()
	if err != nil {
		return err
	}
	return d.Call(ctx, method, arg, res, c.errorUnwrapper)
}

// Notify notifies the server, with the given method and argument. It does not
// wait to hear back for an error. An error might happen in sending the call, in
// which case a native Go Error is returned. The UnwrapErrorFunc in the underlying
// client isn't relevant in this case.
func (c *Client) Notify(ctx context.Context, method string, arg interface{}) (err error) {
	if ctx == nil {
		return errors.New("No Context provided for this notification")
	}
	d, err := c.xp.getDispatcher()
	if err != nil {
		return err
	}
	return d.Notify(ctx, method, arg)
}

// GenericClient is the interface that is exported to autogenerated RPC stubs
// from AVDL files.
type GenericClient interface {
	Call(ctx context.Context, method string, arg interface{}, res interface{}) error
	Notify(ctx context.Context, method string, arg interface{}) error
}

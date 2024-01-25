package gRPC

import (
	"context"
	"fmt"

	"github.com/Sskrill/gRpc-log/proto/audit"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn        *grpc.ClientConn
	auditClient audit.AuditLogClient
}

func NewAuditClient(port int) (*Client, error) {
	addr := fmt.Sprint(":", port)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, auditClient: audit.NewAuditLogClient(conn)}, nil
}

func (c *Client) SenLogReq(ctx context.Context, req audit.LogItem) error {
	entity, err := audit.ToPbEntity(req.Entity)
	if err != nil {
		return err
	}
	action, err := audit.ToPbAction(req.Action)
	if err != nil {
		return err
	}
	_, err = c.auditClient.Log(ctx, &audit.LogRequest{Action: action, Entity: entity, EntityId: req.EntityID, Timestamp: timestamppb.New(req.Timestamp)})

	return err
}

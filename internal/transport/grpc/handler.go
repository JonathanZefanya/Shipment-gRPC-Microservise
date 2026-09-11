package grpc

import (
	"context"
	"errors"
	"shipment/internal/application"
	"shipment/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "shipment/gen"
)

type ShipmentHandler struct {
	pb.UnimplementedShipmentServiceServer
	service *application.ShipmentService
}

func NewShipmentHandler(service *application.ShipmentService) *ShipmentHandler {
	return &ShipmentHandler{service: service}
}

func (h *ShipmentHandler) CreateShipment(
	ctx context.Context,
	req *pb.CreateShipmentRequest,
) (*pb.CreateShipmentResponse, error) {
	shipment, err := h.service.CreateShipment(
		ctx,
		req.ReferenceNumber,
		req.Origin,
		req.Destination,
		req.DriverName,
		req.DriverUnit,
		req.Amount,
		req.DriverRevenue,
	)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &pb.CreateShipmentResponse{
		Shipment: shipmentToProto(shipment),
	}, nil
}

func (h *ShipmentHandler) GetShipment(
	ctx context.Context,
	req *pb.GetShipmentRequest,
) (*pb.GetShipmentResponse, error) {
	shipment, err := h.service.GetShipment(ctx, req.Id)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &pb.GetShipmentResponse{
		Shipment: shipmentToProto(shipment),
	}, nil
}

func (h *ShipmentHandler) AddShipmentEvent(
	ctx context.Context,
	req *pb.AddShipmentEventRequest,
) (*pb.AddShipmentEventResponse, error) {
	event, err := h.service.AddShipmentEvent(
		ctx,
		req.ShipmentId,
		protoStatusToDomain(req.Status),
		req.Note,
	)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &pb.AddShipmentEventResponse{
		Event: eventToProto(event),
	}, nil
}

func (h *ShipmentHandler) GetShipmentHistory(
	ctx context.Context,
	req *pb.GetShipmentHistoryRequest,
) (*pb.GetShipmentHistoryResponse, error) {
	events, err := h.service.GetShipmentHistory(ctx, req.ShipmentId)
	if err != nil {
		return nil, toGRPCError(err)
	}

	protoEvents := make([]*pb.ShipmentEvent, len(events))
	for i, e := range events {
		protoEvents[i] = eventToProto(e)
	}

	return &pb.GetShipmentHistoryResponse{
		Events: protoEvents,
	}, nil
}

func toGRPCError(err error) error {
	switch {
	case errors.Is(err, domain.ErrShipmentNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrInvalidShipment):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}

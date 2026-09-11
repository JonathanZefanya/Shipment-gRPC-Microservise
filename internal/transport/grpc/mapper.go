package grpc

import (
	"shipment/internal/domain"

	pb "shipment/gen"
)

func domainStatusToProto(s domain.ShipmentStatus) pb.ShipmentStatus {
	switch s {
	case domain.StatusPending:
		return pb.ShipmentStatus_SHIPMENT_STATUS_PENDING
	case domain.StatusPickedUp:
		return pb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP
	case domain.StatusInTransit:
		return pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT
	case domain.StatusDelivered:
		return pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED
	case domain.StatusCancelled:
		return pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED
	default:
		return pb.ShipmentStatus_SHIPMENT_STATUS_UNSPECIFIED
	}
}

func protoStatusToDomain(s pb.ShipmentStatus) domain.ShipmentStatus {
	switch s {
	case pb.ShipmentStatus_SHIPMENT_STATUS_PENDING:
		return domain.StatusPending
	case pb.ShipmentStatus_SHIPMENT_STATUS_PICKED_UP:
		return domain.StatusPickedUp
	case pb.ShipmentStatus_SHIPMENT_STATUS_IN_TRANSIT:
		return domain.StatusInTransit
	case pb.ShipmentStatus_SHIPMENT_STATUS_DELIVERED:
		return domain.StatusDelivered
	case pb.ShipmentStatus_SHIPMENT_STATUS_CANCELLED:
		return domain.StatusCancelled
	default:
		return domain.StatusPending
	}
}

func shipmentToProto(s *domain.Shipment) *pb.Shipment {
	return &pb.Shipment{
		Id:              s.ID,
		ReferenceNumber: s.ReferenceNumber,
		Origin:          s.Origin,
		Destination:     s.Destination,
		Status:          domainStatusToProto(s.Status),
		DriverName:      s.DriverName,
		DriverUnit:      s.DriverUnit,
		Amount:          s.Amount,
		DriverRevenue:   s.DriverRevenue,
		CreatedAt:       s.CreatedAt.String(),
		UpdatedAt:       s.UpdatedAt.String(),
	}
}

func eventToProto(e *domain.ShipmentEvent) *pb.ShipmentEvent {
	return &pb.ShipmentEvent{
		Id:         e.ID,
		ShipmentId: e.ShipmentID,
		Status:     domainStatusToProto(e.Status),
		Note:       e.Note,
		CreatedAt:  e.CreatedAt.String(),
	}
}

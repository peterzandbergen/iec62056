package server

import (
	"context"

	"github.com/peterzandbergen/iec62056/service/proto"
)

type protoServer struct {
	proto.UnimplementedMeasurementServiceServer
}

func (ps *protoServer) GetMeasurements(proto.MeasurementService_GetMeasurementsServer) error {
	return nil
}

func (ps *protoServer) GetLastMeasurement(context.Context, *proto.VoidRequest) (*proto.GetMeasurementResponse, error) {
	
	return nil, nil
}

func (ps *protoServer) GetFirstMeasurement(context.Context, *proto.VoidRequest) (*proto.GetMeasurementResponse, error) {
	return nil, nil
}



package grpc

import (
	"context"
	"github.com/dhikaroofi/stock.git/internal/core/use-case/ohlc"
	"github.com/dhikaroofi/stock.git/internal/entity"
	"github.com/dhikaroofi/stock.git/pkg/protobuf/pb"
)

type OhlcHandler struct {
	pb.UnimplementedOhlcServiceServer
	usecase ohlc.Task
}

func NewOhlcHandler(usecase ohlc.Task) *OhlcHandler {
	return &OhlcHandler{
		usecase: usecase,
	}
}

func (s *OhlcHandler) GetSummaryStock(ctx context.Context, req *pb.Request) (*pb.Summary, error) {
	resp, err := s.usecase.GetSummary(ctx, req.GetStockCode(), req.GetDate())
	if err != nil {
		return &pb.Summary{}, err
	}
	return ohlcSummaryToProto(resp), nil
}

func ohlcSummaryToProto(summary entity.OHLCSummary) *pb.Summary {
	return &pb.Summary{
		Average:       summary.Average,
		AverageRound:  summary.AverageRound,
		Volume:        summary.Volume,
		Value:         summary.Value,
		PreviousPrice: int32(summary.PreviousPrice),
		HighestPrice:  int32(summary.HighestPrice),
		LowestPrice:   int32(summary.LowestPrice),
		ClosePrice:    int32(summary.ClosePrice),
		TotalTrans:    int32(summary.TotalTrans),
	}
}

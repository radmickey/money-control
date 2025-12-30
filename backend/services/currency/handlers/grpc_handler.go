package handlers

import (
	"context"

	pb "github.com/radmickey/money-control/backend/proto/currency"
	"github.com/radmickey/money-control/backend/services/currency/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler implements the CurrencyServiceServer interface
type GRPCHandler struct {
	pb.UnimplementedCurrencyServiceServer
	currencyService *service.CurrencyService
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(currencyService *service.CurrencyService) *GRPCHandler {
	return &GRPCHandler{
		currencyService: currencyService,
	}
}

// GetExchangeRate gets the exchange rate between two currencies
func (h *GRPCHandler) GetExchangeRate(ctx context.Context, req *pb.GetExchangeRateRequest) (*pb.ExchangeRateResponse, error) {
	rate, err := h.currencyService.GetExchangeRate(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exchange rate: %v", err)
	}

	return &pb.ExchangeRateResponse{
		Rate: &pb.ExchangeRate{
			FromCurrency: rate.FromCurrency,
			ToCurrency:   rate.ToCurrency,
			Rate:         rate.Rate,
			UpdatedAt:    timestamppb.New(rate.UpdatedAt),
		},
	}, nil
}

// GetMultipleExchangeRates gets exchange rates for multiple currencies
func (h *GRPCHandler) GetMultipleExchangeRates(ctx context.Context, req *pb.GetMultipleExchangeRatesRequest) (*pb.GetMultipleExchangeRatesResponse, error) {
	rates, updatedAt, err := h.currencyService.GetMultipleExchangeRates(ctx, req.BaseCurrency, req.TargetCurrencies)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get exchange rates: %v", err)
	}

	return &pb.GetMultipleExchangeRatesResponse{
		BaseCurrency: req.BaseCurrency,
		Rates:        rates,
		UpdatedAt:    timestamppb.New(updatedAt),
	}, nil
}

// ConvertAmount converts an amount from one currency to another
func (h *GRPCHandler) ConvertAmount(ctx context.Context, req *pb.ConvertAmountRequest) (*pb.ConvertAmountResponse, error) {
	convertedAmount, rate, err := h.currencyService.ConvertAmount(ctx, req.Amount, req.FromCurrency, req.ToCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert currency: %v", err)
	}

	return &pb.ConvertAmountResponse{
		OriginalAmount:  req.Amount,
		FromCurrency:    req.FromCurrency,
		ConvertedAmount: convertedAmount,
		ToCurrency:      req.ToCurrency,
		RateUsed:        rate,
		RateTimestamp:   timestamppb.Now(),
	}, nil
}

// ConvertMultipleAmounts converts multiple amounts
func (h *GRPCHandler) ConvertMultipleAmounts(ctx context.Context, req *pb.ConvertMultipleAmountsRequest) (*pb.ConvertMultipleAmountsResponse, error) {
	var converted []*pb.ConvertedAmount
	var total float64

	for _, a := range req.Amounts {
		convertedAmount, rate, err := h.currencyService.ConvertAmount(ctx, a.Amount, a.FromCurrency, req.ToCurrency)
		if err != nil {
			continue
		}
		total += convertedAmount
		converted = append(converted, &pb.ConvertedAmount{
			Id:              a.Id,
			OriginalAmount:  a.Amount,
			FromCurrency:    a.FromCurrency,
			ConvertedAmount: convertedAmount,
			RateUsed:        rate,
		})
	}

	return &pb.ConvertMultipleAmountsResponse{
		Converted:  converted,
		Total:      total,
		ToCurrency: req.ToCurrency,
	}, nil
}

// ListSupportedCurrencies lists all supported currencies
func (h *GRPCHandler) ListSupportedCurrencies(ctx context.Context, req *pb.ListSupportedCurrenciesRequest) (*pb.ListSupportedCurrenciesResponse, error) {
	currencies, err := h.currencyService.ListSupportedCurrencies(ctx, req.IncludeCrypto)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list currencies: %v", err)
	}

	pbCurrencies := make([]*pb.Currency, len(currencies))
	for i, c := range currencies {
		pbCurrencies[i] = &pb.Currency{
			Code:          c.Code,
			Name:          c.Name,
			Symbol:        c.Symbol,
			DecimalPlaces: int32(c.DecimalPlaces),
			IsCrypto:      c.IsCrypto,
		}
	}

	return &pb.ListSupportedCurrenciesResponse{
		Currencies: pbCurrencies,
	}, nil
}

// GetRateHistory gets rate history
func (h *GRPCHandler) GetRateHistory(ctx context.Context, req *pb.GetRateHistoryRequest) (*pb.RateHistoryResponse, error) {
	history, err := h.currencyService.GetRateHistory(ctx, req.FromCurrency, req.ToCurrency, req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get rate history: %v", err)
	}

	pbHistory := make([]*pb.RatePoint, len(history))
	for i, h := range history {
		pbHistory[i] = &pb.RatePoint{
			Date: timestamppb.New(h.Date),
			Rate: h.Rate,
		}
	}

	return &pb.RateHistoryResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		History:      pbHistory,
	}, nil
}

// RefreshRates refreshes exchange rates
func (h *GRPCHandler) RefreshRates(ctx context.Context, req *pb.RefreshRatesRequest) (*pb.RefreshRatesResponse, error) {
	count, err := h.currencyService.RefreshRates(ctx, req.BaseCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to refresh rates: %v", err)
	}

	return &pb.RefreshRatesResponse{
		Success:      true,
		RatesUpdated: int32(count),
		UpdatedAt:    timestamppb.Now(),
	}, nil
}

// GetLastUpdateTime gets the last time rates were updated
func (h *GRPCHandler) GetLastUpdateTime(ctx context.Context, req *pb.GetLastUpdateTimeRequest) (*pb.GetLastUpdateTimeResponse, error) {
	baseCurrency := "USD" // default
	lastUpdate, err := h.currencyService.GetLastUpdateTime(ctx, baseCurrency)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get last update time: %v", err)
	}

	return &pb.GetLastUpdateTimeResponse{
		LastUpdate:   timestamppb.New(lastUpdate),
		BaseCurrency: baseCurrency,
	}, nil
}

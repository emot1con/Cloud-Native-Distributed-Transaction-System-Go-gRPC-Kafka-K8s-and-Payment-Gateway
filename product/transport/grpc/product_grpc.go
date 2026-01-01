package grpc

import (
	"context"
	"net"
	"product_service/cmd/db"
	"product_service/proto"
	"product_service/repository"
	"product_service/service"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ProductGRPCServer struct {
	service *service.ProductService
	proto.UnimplementedProductServiceServer
}

func NewProductGRPCServer(service *service.ProductService) *ProductGRPCServer {
	return &ProductGRPCServer{
		service: service,
	}
}

func (u *ProductGRPCServer) CreateProduct(ctx context.Context, req *proto.ProductRequest) (*proto.Empty, error) {
	logrus.Info(req.Payload.Name)
	logrus.Info("create product handler")
	if err := u.service.Create(req); err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *ProductGRPCServer) GetProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.Product, error) {
	product, err := u.service.GetProductByID(int(req.Id))
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (u *ProductGRPCServer) ListProducts(ctx context.Context, req *proto.Offset) (*proto.ProductList, error) {
	products, totalProducts, totalPages, page, err := u.service.GetAll(int(req.Id))
	if err != nil {
		return nil, err
	}
	logrus.Info("proto Listing products")

	return &proto.ProductList{
		Products:  products,
		TotalPage: int32(totalPages),
		Total:     int64(totalProducts),
		Page:      int32(page),
	}, nil
}

func (u *ProductGRPCServer) UpdateProduct(ctx context.Context, req *proto.Product) (*proto.Product, error) {
	productResult, err := u.service.Update(req)
	if err != nil {
		return nil, err
	}

	return productResult, nil
}

func (u *ProductGRPCServer) DeleteProduct(ctx context.Context, req *proto.GetProductRequest) (*proto.Empty, error) {
	if err := u.service.Delete(int(req.Id)); err != nil {
		return nil, err
	}

	return nil, nil
}

func GRPCListen() {
	DB, err := db.Connect()
	if err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.InitRedis(); err != nil {
		logrus.Warnf("Redis initialization failed: %v. Continuing without cache.", err)
	}

	ctx := context.Background()
	repo := repository.NewProductRepositoryImpl()
	service := service.NewProductService(repo, DB, ctx)
	connection := NewProductGRPCServer(service)

	lis, err := net.Listen("tcp", ":40001")
	if err != nil {
		logrus.Fatalf("Failed to listen for gRPC: %v", err)
	}

	srv := grpc.NewServer()
	proto.RegisterProductServiceServer(srv, connection)
	logrus.Infof("gRPC Server started on port 40001")

	if err := srv.Serve(lis); err != nil {
		logrus.Fatalf("error when connect to gRPC Server  %v", err)
	}
}

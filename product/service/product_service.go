package service

import (
	"context"
	"database/sql"
	"fmt"
	"product_service/cmd/db"
	"product_service/helper"
	"product_service/proto"
	"product_service/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	productCacheTTL     = 10 * time.Minute
	productListCacheTTL = 5 * time.Minute
)

type ProductService struct {
	repo repository.ProductRepository
	DB   *sql.DB
	ctx  context.Context
}

func NewProductService(ProductRepository repository.ProductRepository, DB *sql.DB, ctx context.Context) *ProductService {
	return &ProductService{
		repo: ProductRepository,
		DB:   DB,
		ctx:  ctx,
	}
}

func (u *ProductService) Create(payload *proto.ProductRequest) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	logrus.Info("calling create product service")
	if err := u.repo.Create(u.ctx, tx, payload.Payload); err != nil {
		return err
	}

	// Invalidate product list cache after create
	if db.RedisClient != nil {
		if err := db.DeleteCachePattern(u.ctx, "products:list:*"); err != nil {
			logrus.Warnf("Failed to invalidate product list cache: %v", err)
		} else {
			logrus.Info("Product list cache invalidated after create")
		}
	}

	return nil
}

func (u *ProductService) GetProductByID(ID int) (*proto.Product, error) {
	cacheKey := fmt.Sprintf("product:%d", ID)

	cachedProduct, err := db.GetCacheProduct(u.ctx, cacheKey)
	if err == nil {
		logrus.Infof("Cache HIT for product ID: %d", ID)
		return cachedProduct, nil
	}
	logrus.Debugf("Cache MISS for product ID: %d", ID)

	// Cache miss or Redis not available, query from database
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	logrus.Info("getting product from database")
	product, err := u.repo.GetProductByID(u.ctx, tx, ID)
	if err != nil {
		return nil, err
	}

	if err := db.SetCache(u.ctx, cacheKey, product, productCacheTTL); err != nil {
		logrus.Warnf("Failed to cache product ID %d: %v", ID, err)
	} else {
		logrus.Debugf("Product ID %d cached successfully", ID)
	}

	return product, nil
}

func (u *ProductService) Update(payload *proto.Product) (*proto.Product, error) {
	tx, err := u.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer helper.CommitOrRollback(tx)

	logrus.Info("getting product")
	productResult, err := u.repo.GetProductByID(u.ctx, tx, int(payload.Id))
	if err != nil {
		return nil, err
	}

	productResult.Name = payload.Name
	productResult.Description = payload.Description
	productResult.Price = payload.Price
	productResult.Stock = payload.Stock

	if err := u.repo.UpdateProduct(u.ctx, tx, payload); err != nil {
		return nil, err
	}
	logrus.Info("product updated")

	// Invalidate cache for this product and product list
	productKey := fmt.Sprintf("product:%d", payload.Id)
	if err := db.DeleteCache(u.ctx, productKey); err != nil {
		logrus.Warnf("Failed to invalidate product cache for ID %d: %v", payload.Id, err)
	}

	if err := db.DeleteCachePattern(u.ctx, "products:list:*"); err != nil {
		logrus.Warnf("Failed to invalidate product list cache: %v", err)
	}

	logrus.Infof("Cache invalidated for product ID: %d", payload.Id)

	return productResult, nil
}

func (u *ProductService) GetAll(page int) ([]*proto.Product, int, int, int, error) {
	offset := (page - 1) * 10

	cacheKey := fmt.Sprintf("products:list:page:%d", page)
	cachedList, err := db.GetCachedProductList(u.ctx, cacheKey)
	if err == nil && cachedList != nil {
		logrus.Infof("Cache HIT for product list page: %d", page)
		return cachedList.Products, cachedList.TotalProduct, cachedList.TotalPage, page, nil
	}
	logrus.Infof("Cache MISS for product list page: %d", page)

	logrus.Info("getting all products from database")
	products, totalProducts, totalPage, err := u.repo.GetAllProduct(u.ctx, u.DB, offset)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// Store in cache for next time
	setCachedList := &db.CachedProductList{
		Products:     products,
		TotalProduct: totalProducts,
		TotalPage:    totalPage,
	}
	if err := db.SetCache(u.ctx, cacheKey, setCachedList, productListCacheTTL); err != nil {
		logrus.Warnf("Failed to cache product list: %v", err)
	} else {
		logrus.Debugf("Product list page %d cached successfully", page)
	}

	return products, totalProducts, totalPage, page, nil
}

func (u *ProductService) Delete(ID int) error {
	tx, err := u.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.CommitOrRollback(tx)

	if err := u.repo.DeleteProduct(u.ctx, tx, ID); err != nil {
		return err
	}

	// Invalidate cache for this product and product list
	productKey := fmt.Sprintf("product:%d", ID)
	if err := db.DeleteCache(u.ctx, productKey); err != nil {
		logrus.Warnf("Failed to invalidate product cache for ID %d: %v", ID, err)
	}

	if err := db.DeleteCachePattern(u.ctx, "products:list:*"); err != nil {
		logrus.Warnf("Failed to invalidate product list cache: %v", err)
	}

	logrus.Infof("Cache invalidated for deleted product ID: %d", ID)

	return nil
}

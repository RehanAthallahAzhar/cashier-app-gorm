package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db}
}

func (p *ProductRepository) AddProduct(product model.Product) error {
	if err := p.db.Create(&product).Error; err != nil {
		// return any error will rollback
		return err
	}
	return nil // TODO: replace this
}

func (p *ProductRepository) ReadProducts() ([]model.Product, error) {
	results := []model.Product{}
	rows, err := p.db.Table("products").Select("*").Where("deleted_at is null").Rows()
	if err != nil {
		return []model.Product{}, err
	}
	defer rows.Close()
	for rows.Next() { // Next akan menyiapkan hasil baris berikutnya untuk dibaca dengan metode Scan.
		p.db.ScanRows(rows, &results)
	}

	return results, nil // TODO: replace this
}

func (p *ProductRepository) DeleteProduct(id uint) error {
	data, err := p.ReadProducts()
	if err != nil {
		return err
	}

	// DELETE FROM session where id = {id}
	err = p.db.Delete(&data, id).Error
	if err != nil {
		return err
	}

	p.db.Save(&data)
	return nil // TODO: replace this
}

func (p *ProductRepository) UpdateProduct(id uint, product model.Product) error {
	// err := p.db.Table("products").Where("id = ?", id).Updates("name", product.Name).Update("price", product.Price).Update("stock", product.Stock).Update("discount", product.Discount).Update("type", product.Type).Update("price", product.Price).Update("stock", product.Stock).Update("discount", product.Discount).Update("type", product.Type).Error
	err := p.db.Table("products").Where("id = ?", id).Updates(&product).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

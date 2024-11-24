package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return CartRepository{db}
}

func (c *CartRepository) ReadCart() ([]model.JoinCart, error) {
	// results := []model.JoinCart{}
	// rows, err := c.db.Table("carts").Select("*").Where("carts.deleted_at is NULL").Joins("left join products on products.id = carts.product_id").Rows()
	// if err != nil {
	// 	return []model.JoinCart{}, err
	// }
	// defer rows.Close()
	// for rows.Next() { // Next akan menyiapkan hasil baris berikutnya untuk dibaca dengan metode Scan.
	// 	c.db.ScanRows(rows, &results)
	// }

	// return results, nil // TODO: replace this

	var listCart []model.JoinCart
	err := c.db.Table("carts").Select("carts.id, carts.product_id, products.name, carts.quantity, carts.total_price").Joins("left join products on products.id = carts.product_id").Where("carts.deleted_at is NULL").Scan(&listCart).Error
	return listCart, err
}

func (c *CartRepository) AddCart(product model.Product) error {
	// // add to cart
	// data := model.Cart{
	// 	ProductID:  product.ID,
	// 	Quantity:   float64(1),
	// 	TotalPrice: float64(product.Price) - (float64(product.Price) * float64(product.Discount/100)),
	// }

	// var findData model.Cart
	// result := c.db.Table("carts").Select("*").Where("product_id = ?", product.ID).Scan(&findData)
	// if result.Error != nil {
	// 	return result.Error
	// }

	// if findData == (model.Cart{}) {
	// 	// Transaction add cart if product doesnt exist in cart table
	// 	c.db.Transaction(func(tx *gorm.DB) error {
	// 		// insert cart
	// 		if err := tx.Table("carts").Create(&data).Error; err != nil {
	// 			// return any error will rollback
	// 			return err
	// 		}

	// 		// get stock product
	// 		var slctedProduct model.Product
	// 		stokProduct := tx.Table("products").Select("*").Where("id = ?", product.ID).Scan(&slctedProduct)
	// 		if stokProduct.Error != nil {
	// 			// return any error will rollback
	// 			return stokProduct.Error
	// 		}

	// 		// mengurangi stok product
	// 		result := tx.Table("products").Where("id = ?", product.ID).Update("stock", (float64(slctedProduct.Stock) - 1.0))
	// 		if result.Error != nil {
	// 			// return any error will rollback
	// 			return result.Error
	// 		}
	// 		return nil
	// 	})
	// } else {
	// 	c.db.Transaction(func(tx *gorm.DB) error {

	// 		// get stock cart table
	// 		var slctedCart model.Cart
	// 		err := tx.Table("carts").Select("*").Where("product_id = ?", product.ID).Scan(&slctedCart).Error
	// 		if err != nil {
	// 			// return any error will rollback
	// 			return err
	// 		}

	// 		// add stock in cart table
	// 		updateData := model.Cart{
	// 			ProductID:  data.ProductID,
	// 			Quantity:   (float64(slctedCart.Quantity) + 1.0),
	// 			TotalPrice: (float64(slctedCart.TotalPrice) + (data.TotalPrice)),
	// 		}
	// 		// err = tx.Table("carts").Where("product_id = ?", product.ID).Updates(updateData{Quantity: (float64(slctedCart.Quantity) + 1.0), TotalPrice: (float64(slctedCart.TotalPrice) + (data.TotalPrice))}).Error
	// 		err = tx.Table("carts").Where("product_id = ?", data.ProductID).Updates(&updateData).Error
	// 		if err != nil {
	// 			// return any error will rollback
	// 			return err
	// 		}

	// 		// get stock product table
	// 		var slctedProduct model.Product
	// 		err = tx.Table("products").Select("*").Where("id = ?", product.ID).Scan(&slctedProduct).Error
	// 		if err != nil {
	// 			// return any error will rollback
	// 			return err
	// 		}

	// 		// reduce stock in product table
	// 		err = tx.Table("products").Where("id = ?", product.ID).Update("stock", (float64(slctedProduct.Stock) - 1.0)).Error
	// 		if err != nil {
	// 			// return any error will rollback
	// 			return err
	// 		}
	// 		return nil
	// 	})
	// }

	// // c.UpdateCart(product.ID, data)

	// return nil // TODO: replace this

	var cart model.Cart
	cartExist := c.db.First(&cart, "product_id = ?", product.ID).Error
	if cartExist == gorm.ErrRecordNotFound {
		return c.db.Transaction(func(tx *gorm.DB) error {
			totalPrice := product.Price - (product.Price * (product.Discount / 100))
			var newCart = &model.Cart{
				ProductID:  product.ID,
				Quantity:   1,
				TotalPrice: totalPrice,
			}

			err := c.db.Create(newCart).Error

			if err != nil {
				return err
			}

			err = c.db.Model(&model.Product{}).Where("id = ?", product.ID).Update("stock", product.Stock-1).Error
			if err != nil {
				return err
			}

			return nil
		})
	} else if cartExist != nil {
		return cartExist
	}

	return c.db.Transaction(func(tx *gorm.DB) error {
		totalPrice := product.Price - (product.Price * (product.Discount / 100))
		err := c.db.Model(&model.Cart{}).Where("product_id = ?", product.ID).Updates(model.Cart{Quantity: cart.Quantity + 1, TotalPrice: cart.TotalPrice + totalPrice}).Error

		if err != nil {
			return err
		}

		err = c.db.Model(&model.Product{}).Where("id = ?", product.ID).Update("stock", product.Stock-1).Error

		if err != nil {
			return err
		}

		return nil
	})
}

func (c *CartRepository) DeleteCart(id uint, productID uint) error {
	c.db.Transaction(func(tx *gorm.DB) error {

		// get stock cart table
		var slctedCart model.Cart
		err := tx.Table("carts").Select("*").Where("product_id = ?", productID).Scan(&slctedCart).Error
		if err != nil {
			// return any error will rollback
			return err
		}

		// get stock product table
		var slctedProduct model.Product
		err = tx.Table("products").Select("*").Where("id = ?", productID).Scan(&slctedProduct).Error
		if err != nil {
			// return any error will rollback
			return err
		}

		// add stock in product table
		err = tx.Table("products").Where("id = ?", productID).Update("stock", (float64(slctedProduct.Stock) + float64(slctedCart.Quantity+1))).Error
		if err != nil {
			// return any error will rollback
			return err
		}

		data := model.Cart{}
		// DELETE FROM carts where id = {id} AND where product_id = {productID}
		err = c.db.Table("carts").Delete(&data, id).Error
		if err != nil {
			return err
		}
		return nil
	})

	return nil // TODO: replace this
}

func (c *CartRepository) UpdateCart(id uint, cart model.Cart) error {
	// update data in carts table
	err := c.db.Table("carts").Where("id = ?", id).Updates(&cart).Error
	if err != nil {
		return err
	}
	return nil // TODO: replace this
}

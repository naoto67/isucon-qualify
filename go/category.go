package main

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

var memCategories []Category

func initializeCategories() error {
	if len(memCategories) > 0 {
		return nil
	}
	var categories []Category
	err := dbx.Select(&categories, "SELECT * FROM `categories`")
	if err != nil {
		return err
	}
	for _, category := range categories {
		if category.ParentID != 0 {
			for _, v := range categories {
				if v.ID == category.ParentID {
					category.ParentCategoryName = v.CategoryName
				}
			}
		}
		memCategories = append(memCategories, category)
	}

	fmt.Println("DEBUG: ", memCategories)

	return nil
}

func getCategoryByID(q sqlx.Queryer, categoryID int) (category Category, err error) {
	for _, category := range memCategories {
		if category.ID == categoryID {
			return category, nil
		}
	}

	return Category{}, errors.New("NotFound")
}

package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang-crud-rest-api/database"
	"golang-crud-rest-api/entities"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func parseProductID(r *http.Request) (uint, error) {
	idParam := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil || id == 0 {
		return 0, fmt.Errorf("invalid product id: %s", idParam)
	}
	return uint(id), nil
}

func validateProductInput(product entities.Product) error {
	if strings.TrimSpace(product.Name) == "" {
		return errors.New("name is required")
	}
	if product.Price < 0 {
		return errors.New("price must be greater than or equal to 0")
	}
	return nil
}

func getProductByID(id uint) (*entities.Product, error) {
	var product entities.Product
	err := database.Instance.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product entities.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := validateProductInput(product); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := database.Instance.Create(&product).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create product"})
		return
	}
	respondJSON(w, http.StatusCreated, product)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	productID, err := parseProductID(r)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	product, err := getProductByID(productID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch product"})
		return
	}

	respondJSON(w, http.StatusOK, product)
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	var products []entities.Product
	if err := database.Instance.Find(&products).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch products"})
		return
	}
	respondJSON(w, http.StatusOK, products)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := parseProductID(r)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	existingProduct, err := getProductByID(productID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch product"})
		return
	}

	var input entities.Product
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	input.ID = existingProduct.ID
	if err := validateProductInput(input); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := database.Instance.Save(&input).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not update product"})
		return
	}

	respondJSON(w, http.StatusOK, input)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := parseProductID(r)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_, err = getProductByID(productID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondJSON(w, http.StatusNotFound, map[string]string{"error": "product not found"})
		return
	}
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not fetch product"})
		return
	}

	if err := database.Instance.Delete(&entities.Product{}, productID).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not delete product"})
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "product deleted successfully"})
}

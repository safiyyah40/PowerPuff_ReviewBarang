package reviewcontroller

import (
	"PowerPuff_ReviewBarang/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Fungsi untuk membuat ulasan baru
func Create(c *gin.Context) {
	var review models.Review

	if err := c.ShouldBindJSON(&review); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Data tidak valid", "error": err.Error()})
		return
	}

	models.DB.Create(&review)
	c.JSON(http.StatusOK, gin.H{"message": "Review berhasil dibuat", "review": review})
}

type Stack []models.Review

func (s *Stack) Push(review models.Review) {
	*s = append(*s, review)
}

// Variabel global untuk stack
var reviewStack = make(Stack, 0)

func PushToStack(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Data tidak valid", "error": err.Error()})
		return
	}
	reviewStack.Push(review)
	models.DB.Create(&review)
	c.JSON(http.StatusOK, gin.H{"message": "Review berhasil dibuat", "review": review})
}

func (s *Stack) Peek() (models.Review, bool) {
	if len(*s) == 0 {
		return models.Review{}, false
	}
	return (*s)[len(*s)-1], true
}

// Fungsi untuk melihat data terakhir di stack
func PeekStack(c *gin.Context) {
    models.DB.Find(&reviewStack)
	review, ok := reviewStack.Peek()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "Stack kosong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ulasan teratas di stack", "review": review})
}

// Fungsi untuk menampilkan semua ulasan
func Index(c *gin.Context) {
	var reviews []models.Review
	models.DB.Find(&reviews)
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// Fungsi untuk menampilkan semua ulasan berdasarkan stack (Last in FIrst Out)
func GetAllFromStack(c *gin.Context) {
    var reviewStack Stack
    var reversedStack Stack
    models.DB.Find(&reviewStack)

    if len(reviewStack) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "Data Review kosong"})
        return
    }

    for i := len(reviewStack) - 1; i >= 0; i-- {
        reversedStack = append(reversedStack, reviewStack[i])
    }

    c.JSON(http.StatusOK, gin.H{
        "stack": reversedStack,
    })
}

// Fungsi untuk menampilkan ulasan berdasarkan nama produk
func Show(c *gin.Context) {
	var reviews []models.Review
	productName := c.Param("ProductName")

	if err := models.DB.Where("product_name = ?", productName).Find(&reviews).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saat mengambil data", "error": err.Error()})
		return
	}

	if len(reviews) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// Fungsi untuk mencari ulasan berdasarkan nama produk dan rating
func SearchByProductAndRating(c *gin.Context) {
	namaProduk := c.DefaultQuery("product_name", "")
	rating := c.DefaultQuery("rating", "")

	if namaProduk == "" || rating == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nama produk dan rating wajib diisi"})
		return
	}

	ratingInt, err := strconv.Atoi(rating)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Format rating tidak valid"})
		return
	}

	var ulasan []models.Review
	result := models.DB.Where("product_name = ? AND rating = ?", namaProduk, ratingInt).Find(&ulasan)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
		return
	}

	if len(ulasan) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Data tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": ulasan})
}

// Fungsi untuk memperbarui data ulasan
func Update(c *gin.Context) {
	var review models.Review
	id := c.Param("ID")

	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data tidak valid", "error": err.Error()})
		return
	}

	result := models.DB.Model(&models.Review{}).Where("id = ?", id).Updates(review)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saat memperbarui data", "error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Tidak ada data yang diperbarui"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diperbarui"})
}
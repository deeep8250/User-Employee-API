package controllers

import (
	"context"
	"fmt"
	"net/http"
	"user-management-system/config"
	"user-management-system/models"
	"user-management-system/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func GetFilteredData(c *gin.Context) {
	coll := config.GetCollection()
	var filterUse models.User

	id := c.Query("id")
	intID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error main": err.Error(), "id": id})

		return
	}

	filter := bson.M{
		"_id": intID,
	}

	err = coll.FindOne(context.TODO(), filter).Decode(&filterUse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong while searching for data"})
		}
		return
	}

	c.JSON(http.StatusOK, filterUse)

}
func Signin(c *gin.Context) {

	var signingUser models.User
	var storedUser models.User

	if err := c.ShouldBindJSON(&signingUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "invalid request"})
		return
	}

	var coll = config.GetCollection()

	err := coll.FindOne(context.TODO(), bson.M{"email": signingUser.Email}).Decode(&storedUser)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data already exist"})
		return
	}
	if signingUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is empty"})
		return
	}
	hashpass, err := bcrypt.GenerateFromPassword([]byte(signingUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password cant be hashed"})
		return
	}

	signingUser.Password = string(hashpass)

	fmt.Print(hashpass)

	user := bson.M{
		// "_id":       signingUser.ID,
		"name":     signingUser.Name,
		"role":     signingUser.Role,
		"email":    signingUser.Email,
		"password": signingUser.Password,
	}

	_, err = coll.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "data inserted"})

}
func Login(c *gin.Context) {

	type UserLoginStruct struct {
		Email    string `json:"email" bson:"email"`
		Password string `json:"password" bson:"password"`
	}

	var coll = config.GetCollection()
	var LoginUser UserLoginStruct
	var LoginUserDecode UserLoginStruct

	if err := c.ShouldBindJSON(&LoginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant take input"})
		return
	}

	err := coll.FindOne(context.TODO(), bson.M{"email": LoginUser.Email}).Decode(&LoginUserDecode)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "no data found"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "something went wrong while finding our provided email"})
			return
		}

	}

	err = bcrypt.CompareHashAndPassword([]byte(LoginUserDecode.Password), []byte(LoginUser.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password doesn't match"})
		return
	}

	token, err := utils.CreateJWT(LoginUser.Email, LoginUser.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error to crate jwt tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": " login successful",
		"token ":  token,
	})

}
func Update(c *gin.Context) {
	var coll = config.GetCollection()

	var UpdateData models.User
	var DecodeData models.User

	err := c.ShouldBindJSON(&UpdateData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant take input"})
		return
	}

	//check if email is exist or not
	err = coll.FindOne(context.TODO(), bson.M{"_id": UpdateData.ID}).Decode(&DecodeData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = coll.UpdateOne(context.TODO(), bson.M{"_id": UpdateData.ID}, bson.M{"$set": UpdateData})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant insert the data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "data updated successfully"})

}
func Delete(c *gin.Context) {

	var coll = config.GetCollection()
	type Request struct {
		ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	}

	var id Request

	err := c.ShouldBindJSON(&id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant take input"})
		return
	}

	_, err = coll.DeleteOne(context.TODO(), bson.M{"_id": id.ID})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cant delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "data successfully deleted"})

}

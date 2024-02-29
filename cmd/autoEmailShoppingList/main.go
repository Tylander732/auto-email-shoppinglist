package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/magiconair/properties"

	"net/smtp"
	"os"
	"slices"
	"strings"

	"github.com/tylander732/autoEmailShoppingList/internal/model"
    // "github.com/tylander732/autoEmailShoppingList/internal/consts"
	// "github.com/tylander732/autoEmailShoppingList/internal/projectpath"
)


var propertiesFile = "./resources/app.properties"
var username string
var password string
var distributionList string


var vegetables = []string {"potatoes"}
var fruits = []string {"apples"}
var proteins = []string {"apples"}

func main() {
    readProperties()
    userJsonFile, errReadingJson := readJsonFile()
    if errReadingJson != nil {
        fmt.Println("Error gathering file from disk: ", errReadingJson)
        return
    }

	userArray := userJsonFile.UserJArray
    var emailReceivers []string

	for i := 0; i < len(userArray); i++ {
		meals, err := selectMeals(userArray[i])

		if err != nil {
			fmt.Println("Was unable to succesfully select meal for users", err)
		}
		emailString := makeMealEmailString(meals)
        //This is a bug. We don't want to append to the list of receivers and then resend another email
        //They will already have received an email the first time around
        //Send 1 email per loop for each receiver, or batch all the emails to be sent at once?
		emailReceivers = append(emailReceivers, userArray[i].Email)
        fmt.Println(emailReceivers)

		sendEmail(emailString, emailReceivers)
	}

    fmt.Println("after")

}

func readJsonFile() (model.JsonFile, error) {
	contents, err := os.ReadFile("./resources/userList.json")
	if err != nil {
		return model.JsonFile{}, err
	}

	data := model.JsonFile{}

	err = json.Unmarshal(contents, &data)

	if err != nil {
		fmt.Println(err)
	}

	return data, nil

}

func selectMeals(usersData model.User) ([]model.Meal, error) {
	numOfUsersMeal := len(usersData.MealJArray)
	numOfmealsToSelect := usersData.NumOfMealsToSelect

	mealsToSend := []model.Meal{}

	randomMealsToBeSelected, err := generateUniqueRandomIntegers(numOfUsersMeal, numOfmealsToSelect)
	if err != nil {
		return []model.Meal{}, err
	}

	for i := 0; i < len(randomMealsToBeSelected); i++ {
		mealsToSend = append(mealsToSend, usersData.MealJArray[randomMealsToBeSelected[i]])
	}

	//get the meal objects
	return mealsToSend, nil
}

func generateUniqueRandomIntegers(numberRange int, amountToGenerate int) ([]int, error) {
	if amountToGenerate > numberRange {
		return []int{}, errors.New("Amount asked to generate, is higher than the number of user meals added")
	}
	uniqueInts := []int{}

	i := 0
	// I might have an off by one error in this loop, we'll see
	for i < amountToGenerate {
		currentInt := rand.Intn(numberRange)
		if !slices.Contains(uniqueInts, currentInt) {
			uniqueInts = append(uniqueInts, currentInt)
			i++
		}
	}

	return uniqueInts, nil
}

func makeMealEmailString(meal []model.Meal) string {
	var emailString strings.Builder
	for i := 0; i < len(meal); i++ {
		currentMeal := meal[i]
		emailString.WriteString(currentMeal.Name + "\n")
		emailString.WriteString("Ingredients: ")
		emailString.WriteString(strings.Join(currentMeal.IngredientsJArray, ", "))
		emailString.WriteString("\n \n")
	}

	return emailString.String()
}

func sendEmail(emailString string, receivers []string) {
	// smtp server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", username, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, username, receivers, []byte(emailString))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email sent successfully")
}

func readProperties() {
	p := properties.MustLoadFile(propertiesFile, properties.UTF8)
	username, _ = p.Get("username")
	password, _ = p.Get("password")
	distributionList, _ = p.Get("distributionList")
}

func seperateIngredients(ingredients []string) ([]string, []string, []string) {

    vegetables := []string{}
    fruits := []string{}
    proteins := []string{}

    for i := 0; i < len(ingredients); i++ {
        currentIngredient := ingredients[i]
        if(isVegetable(currentIngredient)) {
            vegetables = append(vegetables, currentIngredient)
        }
        if(isFruit(currentIngredient)) {
            fruits = append(fruits, currentIngredient)
        }
        if(isProtein(currentIngredient)) {
            proteins = append(proteins, currentIngredient)
        }
    }

    return vegetables, fruits, proteins
}

func isVegetable(ingredient string) bool {
    if(!slices.Contains(vegetables, ingredient)){
        return false
    }

    return true
}

func isFruit(ingredient string) bool {
    if(!slices.Contains(fruits, ingredient)){
        return false
    }

    return true
}

func isProtein(ingredient string) bool {
    if(!slices.Contains(proteins, ingredient)){
        return false
    }

    return true
}
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magiconair/properties"
	"math/rand"
	"net/smtp"
	"os"
	"slices"
	"strings"
)

var propertiesFile = "./resources/app.properties"
var username string
var password string
var distributionList string

type JsonFile struct {
	UserJArray []User `json:"users"`
}

type User struct {
	Email              string `json:"email"`
	NumOfMealsToSelect int    `json:"numOfmealsToSelect"`
	MealJArray         []Meal `json:"meals"`
}

type Meal struct {
	Name              string   `json:"name"`
	IngrediantsJArray []string `json:"ingrediants"`
}

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
		emailReceivers = append(emailReceivers, userArray[i].Email)
        // __AUTO_GENERATED_PRINT_VAR_START__
        fmt.Println(fmt.Sprintf("main emailString: %v", emailString)) // __AUTO_GENERATED_PRINT_VAR_END__

		sendEmail(emailString, emailReceivers)
	}

}

func readJsonFile() (JsonFile, error) {
	contents, err := os.ReadFile("./resources/userList.json")
	if err != nil {
		return JsonFile{}, err
	}

	data := JsonFile{}

	err = json.Unmarshal(contents, &data)

	if err != nil {
		fmt.Println(err)
	}

	return data, nil

}

func selectMeals(usersData User) ([]Meal, error) {
	numOfUsersMeal := len(usersData.MealJArray)
	numOfmealsToSelect := usersData.NumOfMealsToSelect

	mealsToSend := []Meal{}

	randomMealsToBeSelected, err := generateUniqueRandomIntegers(numOfUsersMeal, numOfmealsToSelect)
	if err != nil {
		return []Meal{}, err
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

func makeMealEmailString(meal []Meal) string {
	var emailString strings.Builder
	for i := 0; i < len(meal); i++ {
		currentMeal := meal[i]
		emailString.WriteString(currentMeal.Name + "\n")
		emailString.WriteString("Ingrediants: ")
		emailString.WriteString(strings.Join(currentMeal.IngrediantsJArray, ", "))
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

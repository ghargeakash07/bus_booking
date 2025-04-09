package utils

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"math/rand"

	// "errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	//"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestBodyLogger(c *gin.Context) string {
	requestBody, _ := c.GetRawData()
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
	return string(requestBody)
}

func SubtractOneYearFromDate(inputDate string) (string, error) {
	// Parse the input date string
	date, err := time.Parse("2006-01-02", inputDate)
	if err != nil {
		return "", err
	}

	// Subtract one year
	oneYearAgo := date.AddDate(-1, 0, 1)

	// Format the result as a string
	result := oneYearAgo.Format("2006-01-02")

	return result, nil
}

func SubtractFiveYearsFromDate(inputDate string) (string, error) {
	// Parse the input date string
	date, err := time.Parse("2006-01-02", inputDate)
	if err != nil {
		return "", err
	}
	// Subtract five years
	fiveYearsAgo := date.AddDate(-5, 0, 0)
	// Check for leap years between input date and five years ago
	for i := date.Year(); i > fiveYearsAgo.Year(); i-- {
		// If the current year is a leap year, adjust the result
		if isLeapYear(i) {
			fiveYearsAgo = fiveYearsAgo.AddDate(0, 0, -1) // Subtract one day for a leap year
		}
	}
	// Format the result as a string
	result := fiveYearsAgo.Format("2006-01-02")
	return result, nil
}

// Function to check if a year is a leap year
func isLeapYear(year int) bool {
	// Leap year if divisible by 4, but not divisible by 100 unless divisible by 400
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
func InArray(item string, array []string, caseSensitive bool) bool {
	for _, element := range array {
		if (caseSensitive && element == item) || (!caseSensitive && strings.EqualFold(element, item)) {
			return true
		}
	}
	return false
}

func AddFourYearsToDate(inputDate string) (string, error) {
	// Parse the input date string
	date, err := time.Parse("2006-01-02", inputDate)
	if err != nil {
		return "", err
	}

	// Add four years
	fourYearsLater := date.AddDate(4, 0, 0)

	// Check for leap years between input date and four years later
	for i := date.Year(); i < fourYearsLater.Year(); i++ {
		// If the current year is a leap year, adjust the result
		if isLeapYear(i) {
			fourYearsLater = fourYearsLater.AddDate(0, 0, 1) // Add one day for a leap year
		}
	}

	// Format the result as a string
	result := fourYearsLater.Format("2006-01-02")

	return result, nil
}

func ConvertMapToString(data map[string]interface{}) (string, error) {
	// Marshal the map into JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	// Convert JSON bytes to string
	jsonString := string(jsonBytes)

	return jsonString, nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func CalculateAgeInYear(birthdate string) (int, error) {
	// Parse the birthdate string
	layout := "2006-01-02"
	birth, err := time.Parse(layout, birthdate)
	if err != nil {
		return 0, err
	}

	// Calculate the age
	now := time.Now()
	age := now.Year() - birth.Year()

	// Adjust age if the birthday hasn't occurred yet this year
	if now.YearDay() < birth.YearDay() {
		age--
	}

	return age, nil
}

func ParseDateDDMMYYYY(dateString string) string {
	if dateString != "" {
		inputDate, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			logrus.Error("ParseDateDDMMYYYY@Error parsing input date:", err)
		}

		// Format the time in the desired output format
		outputDateString := inputDate.Format("02/01/2006")
		return outputDateString
	}
	return ""
}

// CreateZip creates a ZIP file from all files in the source directory
func CreateZip(sourceDir, zipFilePath string) error {
	logrus.Infof("Creating ZIP file from %s to %s", sourceDir, zipFilePath)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		log.Fatalf("Source directory %s does not exist", sourceDir)
	}
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Error("Error walking file:", err)
			return err
		}

		if !info.IsDir() {
			relativePath, err := filepath.Rel(sourceDir, filePath)
			if err != nil {
				return err
			}

			zipEntry, err := archive.Create(relativePath)
			if err != nil {
				return err
			}

			fileContent, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			// Log or print file content for debugging
			fmt.Printf("File Content for %s:\n%s\n", filePath, fileContent)
			_, err = zipEntry.Write(fileContent)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func MakeTimestampMilli() string {
	unixMilli := unixMilli(time.Now())
	return strconv.FormatInt(unixMilli, 10)
}

func unixMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func InArrayContain(val string, array []string) (ok bool) {
	var i int
	for i = range array {
		if ok = array[i] == val; ok {
			return true
		}
	}
	return false
}

func Contain(val int, array []int) (ok bool) {
	var i int
	for i = range array {
		if ok = array[i] == val; ok {
			return true
		}
	}
	return false
}

type Config struct {
	TravelRequestConfig []ApprovalConfig `json:"travel_request_config"`
}

type ApprovalConfig struct {
	Type string `json:"type"`
	ID   []int  `json:"id,omitempty"`
}

type ApprovalUser struct {
	UserID int
	Type   string
}

type TravelDeskUsers struct {
	ID []int `json:"id,omitempty"`
}

func GenerateOTP() string {
	// Create a new random number generator
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	otp := rng.Intn(1000000)        // Generate a number between 0 and 999999
	return fmt.Sprintf("%06d", otp) // Format it as a six-digit string
}

// func GetApprovalUsers(creatorUserId int, requestID int) []ApprovalUser {
// 	dbTx := database.PrimeDB
// 	var config Config
// 	configFile, err := os.Open("config/approval_config.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer configFile.Close()
// 	decoder := json.NewDecoder(configFile)
// 	if err := decoder.Decode(&config); err != nil {
// 		panic(err)
// 	}

// 	var userIdArr []ApprovalUser
// 	for _, item := range config.TravelRequestConfig {

// 		if item.Type != "travel_desk" {
// 			var userDetailObj services.UserDetailObj
// 			var userObj services.UserObjResult
// 			var chkExistObj services.TravelApprovalUsersResult
// 			dbTx.Table("user_details as ud").Select(item.Type+" as approval_user").Where("ud.user_relat_id = ?", creatorUserId).First(&userDetailObj)
// 			dbTx.Table("users").Select("users.*").Where("email = ?", userDetailObj.ApprovalUser).First(&userObj)
// 			dbTx.Table("travel_request_approval_users as trap").Select("trap.*").Where("trap.user_id = ?", userObj.UserId).Where("trap.request_id = ?", requestID).First(&chkExistObj)

// 			if chkExistObj == (services.TravelApprovalUsersResult{}) {
// 				// userIdArr = append(userIdArr, userObj.UserId)
// 				userIdArr = append(userIdArr, ApprovalUser{
// 					UserID: userObj.UserId,
// 					Type:   item.Type, // Replace with the actual extra parameter value
// 				})
// 				fmt.Println("userIdArr", userIdArr)
// 				return userIdArr
// 			}

// 		} else {
// 			mainStatusUpdateResult := dbTx.Model(&model.TravelRequests{}).
// 				Where("id = ?", requestID).
// 				Update("status", "approved")
// 			if mainStatusUpdateResult.Error != nil {
// 				fmt.Println("Error updating IsAction:", mainStatusUpdateResult.Error)
// 			}
// 			var userIdArr []ApprovalUser
// 			fmt.Println(item.ID)
// 			travelDeskUser := []int{1, 2, 3}
// 			for _, approvalUserID := range travelDeskUser {
// 				fmt.Println("approvalUserID", approvalUserID)
// 				userIdArr = append(userIdArr, ApprovalUser{
// 					UserID: approvalUserID,
// 					Type:   item.Type, // Replace with the actual extra parameter value
// 				})
// 			}
// 			fmt.Println("userIdArr", userIdArr)
// 			return userIdArr
// 		}
// 	}
// 	return userIdArr
// }

// func GetTravelDeskUsers() []int {
// 	var config Config
// 	configFile, err := os.Open("config/approval_config.json")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer configFile.Close()
// 	decoder := json.NewDecoder(configFile)
// 	if err := decoder.Decode(&config); err != nil {
// 		panic(err)
// 	}

// 	var travelDeskIDs []int
// 	for _, config := range config.TravelRequestConfig {
// 		travelDeskIDs = append(travelDeskIDs, config.ID...)
// 	}

// 	return travelDeskIDs
// }

// func TitleCase(s string) string {
// 	return strings.Title(s)
// }

// func IsValidPhoneNo(PhoneNo string) bool {
// 	phoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
// 	return phoneRegex.MatchString(PhoneNo)
// }

// // ParseDateRange parses the start and end dates from a custom format and returns them in the database format.
// func ParseDateRange(startDateStr, endDateStr string) (string, string, error) {
// 	// Define the input date format and output date format
// 	inputFormat := "2-01-2006"
// 	outputFormat := "2006-01-02 15:04:05"

// 	// Parse the start date
// 	startDate, err := time.Parse(inputFormat, startDateStr)
// 	if err != nil {
// 		return "", "", errors.New("failed to parse start date: " + err.Error())
// 	}

// 	// Parse the end date
// 	endDate, err := time.Parse(inputFormat, endDateStr)
// 	if err != nil {
// 		return "", "", errors.New("failed to parse end date: " + err.Error())
// 	}

// 	// Set start date time to 00:00:00
// 	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

// 	// Set end date time to 23:59:59
// 	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

// 	// Format the dates for the database
// 	return startDate.Format(outputFormat), endDate.Format(outputFormat), nil
// }

// func GenerateSlug(name string) string {
// 	slug := strings.ToLower(name)
// 	slug = regexp.MustCompile(`[^\w\s-]`).ReplaceAllString(slug, "")
// 	slug = regexp.MustCompile(`[\s-]+`).ReplaceAllString(slug, "-")
// 	slug = strings.Trim(slug, "-")
// 	return slug
// }

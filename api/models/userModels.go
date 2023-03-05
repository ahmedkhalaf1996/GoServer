package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ManualLoginModel struct {
	Email    string `json:"email" bson:"email" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required,min=5"`
}

type SocialAuthModel struct {
	Email     string `json:"email" bson:"email" validate:"required"`
	FirstName string `json:"firstname" bson:"firstname" validate:"required"`
	LastName  string `json:"lastname" bson:"lastname" validate:"required"`
}

type PhoneAuthModel struct {
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber" validate:"required"`
	FirstName   string `json:"firstname" bson:"firstname"`
	LastName    string `json:"lastname" bson:"lastname"`
	Password    string `json:"password" bson:"password" validate:"required,min=5"`
}

type UserModel struct {
	ID                    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IsAccountVerified     bool               `json:"isAccountVerified" bson:"isAccountVerified"`
	IsPhoneNumberVerified bool               `json:"isPhoneNumberVerified" bson:"isPhoneNumberVerified"`
	FirstName             string             `json:"firstname" bson:"firstname" validate:"required"`
	LastName              string             `json:"lastname" bson:"lastname" validate:"required"`
	DateOfBirth           string             `json:"dateOfBirth" bson:"dateOfBirth" validate:"required"`
	Email                 string             `json:"email" bson:"email" validate:"required"`
	Password              string             `json:"password" bson:"password" validate:"required,min=5"`
	PhoneNumber           string             `json:"phoneNumber" bson:"phoneNumber"`
	StatusOfuser          string             `json:"statusOfUser" bson:"statusOfUser"` // Dancer Instructor DancerAndInstructor Wondering
	UserProfilePhoto      string             `json:"userProfilePhoto" bson:"userProfilePhoto"`
	UserMediaPhoto        []string           `json:"userMediaPhoto" bson:"userMediaPhoto"` // base64
	UserMediaVideo        []string           `json:"userMediaVideo" bson:"userMediaVideo"` // base64
	UserLanguages         []string           `json:"userLanguages" bson:"userLanguages"`
	UserHeight            string             `json:"userHeight" bson:"userHeight"`
	UserDances            []string           `json:"userDances" bson:"userDances"`             // country city
	UserLocation          []string           `json:"userLocation" bson:"userLocation"`         // country city
	LocationDetails       []string           `json:"locationDetails" bson:"locationDetails"`   // country city
	Gender                string             `json:"gender" bson:"gender"`                     // woman man
	Orientation           string             `json:"orientation" bson:"orientation"`           // bisexual Asexual Heterosexual straight
	UserRelationship      string             `json:"userRelationship" bson:"userRelationship"` // marred single complecated
	UserLookFor           string             `json:"userLookFor" bson:"userLookFor"`           // partners friends team
	UserRole              string             `json:"userRole" bson:"userRole"`                 // lider flexible follower
	UserHobbyes           []string           `json:"userHobbyes" bson:"userHobbyes"`           // lider flexible follower
	UserZodlac            []string           `json:"userZodlac" bson:"userZodlac"`
	School                string             `json:"school" bson:"school"`
	PlaceOfJob            string             `json:"placeOfJob" bson:"placeOfJob"`
	JobTitle              string             `json:"jobTitle" bson:"jobTitle"`
	UserPets              []string           `json:"userPets" bson:"userPets"`                     // dog cat repbile bird fish pet-free
	IsUserSmoking         bool               `json:"isUserSmoking" bson:"isUserSmoking"`           // frequently Socially never
	PhyslcalAttraction    []string           `json:"physlcalAttraction" bson:"physlcalAttraction"` // eyes breasts ears biceps feet hair hands lips
	UserTurnON            []string           `json:"userTurnON" bson:"userTurnON"`                 // eyes breasts ears biceps feet hair hands lips
	UserStyle             []string           `json:"userStyle" bson:"userStyle"`                   // eyes breasts ears biceps feet hair hands lips
	UserBody              []string           `json:"userBody" bson:"userBody"`                     // eyes breasts ears biceps feet hair hands lips
	UserCurrentMood       string             `json:"userCurrentMood" bson:"userCurrentMood"`       // Amused or any
	UserSesson            []string           `json:"userSesson" bson:"userSesson"`                 // Amused or any
	UserMovies            []string           `json:"userMovies" bson:"userMovies"`
	UserDrink             []string           `json:"userDrink" bson:"userDrink"`
}

package routes

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse
import (
	"main/api/controllers"
	"main/api/middleware"
	"main/api/validation"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("Private/GetAllUsers", controllers.GetAllUsers)

	// auth Start
	// manual signup
	app.Post("api/user/manual-signup",
		validation.ValidateUser,
		controllers.ManualRegister)
	// manual signin
	app.Post("api/user/manual-signin",
		validation.ValidateUserManualLogin,
		controllers.ManualLogin)

	// ------ Social Auth --------------//
	// Social login & Register google
	app.Post("api/user/Auth-social-Google",
		validation.ValidateUserSocialAuth,
		controllers.AuthWithGoogle)

	// Social login & Register facebook
	app.Post("api/user/Auth-social-Facebook",
		validation.ValidateUserSocialAuth,
		controllers.AuthWithFaceBook)

	// Social login & Register apple
	app.Post("api/user/Auth-social-Apple",
		validation.ValidateUserSocialAuth,
		controllers.AuthWithApple)
	// --------------Phone-------------------//
	app.Get("api/user/GetPhoneCountriesCodesList", controllers.GetPhoneCountriesCodesList)
	// Auth login & Register With Phone Number
	app.Post("api/user/Auth-Phone-Number",
		validation.ValidateUserPhoneAuth,
		controllers.AuthWithPhoneNumber)

	// get dances list
	app.Get("api/GetDaces", controllers.GetDancesList)

	// Update User Data Any Time or After Siginup
	// we should have the token bearar and send it as Authorization header & user ID
	app.Patch("api/user/UpdateInfo/:id",
		middleware.AuthMiddleware,
		controllers.UpdateUserInfo)
	// update user sug like or hate
	app.Patch("api/user/UpdateSug/:mid/LikeOrHate/:nid",
		controllers.UpUserSug)
	// suggested partner
	app.Get("api/user/Suggested/:id",
		controllers.SuggestedPartner)

	// Message

	// get rooms ides for chat service
	app.Get("Private/GetAllPrivateRoomsID", controllers.GetAllPrivateRoomsID)

	// Get User Message Room
	app.Get("api/message/:FuserID/:SuserID", controllers.GetPrivateRoomeID)

	// SendMessageToPrivateRoom
	app.Post("api/message", controllers.SendMessageToPrivateRoom)

	// GetMessageByNumbers // get the copy of the user
	app.Post("api/message/GetMessages/:roomId/:userId", controllers.GetMessageByNumbers)

	// RemoveUserCopyOfChat/:roomId/:userId
	app.Delete("api/message/RemoveUserCopyOfChat/:roomId/:userId", controllers.RemoveUserCopyOfChat)

	// PrivateFunGetRoomData
	app.Get("Private/:RoomID", controllers.PrivateFunGetRoomData)

	// ---------------------------- Group ----------------------- //

	// GetGroupRoomID

	// get Group rooms ides for chat service
	app.Get("Gorup/GetAllGroupRoomsID", controllers.GetAllGroupRoomsID)

	app.Post("api/GroupMessage/:GroupName/:CreatorID", controllers.GetGroupRoomID)

	// SendMessageToGroupChat
	app.Post("api/SendMessageToGroup/:GroupID/:SenerId", controllers.SendMessageToGroupChat)

	// GetGroupMessageByNumbers
	app.Post("api/GetGroupMessageByNumbers/:roomId/:userId", controllers.GetGroupMessageByNumbers)

	// AddNewUserToGroup || should notify the user
	app.Post("api/AddNewUserToGroup/:roomId/:JoinUserId", middleware.AuthMiddleware, controllers.AddNewUserToGroup)

	// RemoveMemberFromchatGroup
	app.Delete("api/RemoveMemberFromchatGroup/:roomId/:JoinUserId", middleware.AuthMiddleware, controllers.RemoveMemberFromchatGroup)

	// ------ Story ------- //

	app.Post("api/AddStory/:userId", controllers.AddStory) // Add New Story

	app.Get("api/GetStoryes/:userId", controllers.GetStory) // return all User Storys

	app.Delete("api/RemoveStory/:userId/:StoryNumber", controllers.RemoveStory) // return Remove One Story

	// ----- Events ---------- //

	// create event
	app.Post("api/CreateEvent/:userId", controllers.CreateEvent)
	// Get event by id
	app.Get("api/GetEvent/:eventId", controllers.GetEventByid)
	// Get events created by user
	app.Get("api/GetEventsCreatedByUser/:userId", controllers.GetEventsCreatedByUser)
	// Get events that user intested in
	app.Get("api/GetUserNearEvents/:userId", controllers.GetUserNearEvents)
	// invite User To an Event
	app.Get("api/inviteUserToEvent/:eventId/:SuserId/:RuserId", controllers.InviteUserToEvent)

	// UserGoingOrNotToEvent
	app.Patch("api/UserGoingOrNotToEvent/:eventId/:userId", controllers.UserGoingOrNotToEvent)

	// CheckIfUserGoingOrNotToEvent
	app.Get("api/CheckIfUserGoingOrNotToEvent/:eventId/:userId", controllers.CheckIfUserGoingOrNotToEvent)
	// add Comment to Event Page
	app.Post("api/AddCommentToEvent/:eventId", controllers.AddCommentToEventPage)

	// update event deatils Need To Auth
	app.Patch("api/UpdateEventDeatils/:eventId/:userId",
		middleware.AuthMiddleware,
		controllers.UpdateEventDeatils)
	//GetLeftDaysToStartTheEvent
	app.Get("api/GetLeftDaysToStartTheEvent/:eventId", controllers.GetLeftDaysToStartTheEvent)

	// create and send veryfied account Mail
	app.Post("api/ReSendVeryficationCodeToMail", controllers.ReSendVeryficationCodeToMail)

	// crete and send veryfied account Phone Number
	app.Post("api/ReSendVeryficationCodeToPhoneNumber", controllers.ReSendVeryficationCodeToPhone)

	// VeryfiedProfileAccount
	app.Post("api/VeryfiedProfileAccount/:userId", controllers.VeryfiedProfileAccount)

	// -----------------------------------
	// AddNewEmailToUserProfile
	app.Post("api/AddNewEmailToUserProfile/:userId", controllers.AddNewEmailToUserProfile)

	// VeryfiedAddedEmailToProfile
	app.Post("api/VeryfiedAddedEmailToProfile/:userId", controllers.VeryfiedAddedEmailToProfile)

	// AddNewPhoneNumberToUserProfile
	app.Post("api/AddNewPhoneNumberToUserProfile/:userId", controllers.AddNewPhoneNumberToUserProfile)

	// VeryfiedAddedPhoneToProfile
	app.Post("api/VeryfiedAddedPhoneToProfile/:userId", controllers.VeryfiedAddedPhoneToProfile)

	// -------- User Requestes ------------ //
	// RequestToChat/:SuserId/:RuserId
	app.Get("api/RequestToChat/:SuserId/:RuserId", controllers.RequestToChat)

	// SendLovetoUser/:SuserId/:RuserId
	app.Get("api/SendLovetoUser/:SuserId/:RuserId", controllers.SendLovetoUser)

	// SendBuzztoUser/:SuserId/:RuserId
	app.Get("api/SendBuzztoUser/:SuserId/:RuserId", controllers.SendBuzzToUser)

	// SendStartoUser/:SuserId/:RuserId
	app.Get("api/SendStartoUser/:SuserId/:RuserId", controllers.SendStarToUser)

	// api/AcceptRefuaseRequests/:ReqId/:userId
	app.Post("api/AcceptRefuaseRequests/:ReqId/:userId", controllers.AcceptRefuaseRequests)

	// ---------------Blockeded User Part--------------

	//  AddUserToBlockingList/:mainUID/secUID
	app.Post("api/AddUserToBlockingList/:mainUID/:secUID", controllers.AddUserToBlockingList)

	// api/UnBlockUser/:mainUID/:secUID
	app.Post("api/UnBlockUser/:mainUID/:secUID", controllers.UnBlockUser)

	// api/Getblockeduserslist/:mainUID
	app.Get("/api/Getblockeduserslist/:mainUID", controllers.Getblockeduserslist)

	// api/IsUserBlockedOrNot/:mainUID/:secUID
	app.Get("/api/IsUserBlockedOrNot/:mainUID/:secUID", controllers.IsUserBlockedOrNot)

}

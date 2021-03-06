package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"projekat/backend/userService/dto"
	"projekat/backend/userService/model"
	"projekat/backend/userService/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegularUserService struct {
	RegularUserRepository *repository.RegularUserRepository
}

func (service *RegularUserService) Register(regularUserRegistrationDto dto.RegularUserRegistrationDTO) error {
	fmt.Println("Creating regular user")

	if service.RegularUserRepository.ExistByUsername(regularUserRegistrationDto.Username) {
		return fmt.Errorf("username is already taken")
	}

	var regularUser = createRegularUserFromRegularUserRegistrationDTO(&regularUserRegistrationDto)
	createdUserId, err := service.RegularUserRepository.Register(regularUser)
	if err != nil {
		return err
	}
	err2 := service.registerUserInAuthenticationService(regularUserRegistrationDto, createdUserId)
	if err2 != nil {
		return err2
	}
	return nil
}

func (service *RegularUserService) RegisterAgent(regularUserRegistrationDto dto.RegularUserRegistrationDTO) error {
	fmt.Println("Creating agent")

	if service.RegularUserRepository.ExistByUsername(regularUserRegistrationDto.Username) {
		return fmt.Errorf("username is already taken")
	}

	var regularUser = createRegularUserFromRegularUserRegistrationDTO(&regularUserRegistrationDto)
	createdUserId, err := service.RegularUserRepository.Register(regularUser)
	if err != nil {
		return err
	}
	err2 := service.registerUserInAuthenticationService(regularUserRegistrationDto, createdUserId)
	if err2 != nil {
		return err2
	}
	return nil
}

func (service *RegularUserService) registerUserInAuthenticationService(regularUserRegistrationDto dto.RegularUserRegistrationDTO, createdUserId string) error {
	postBody, _ := json.Marshal(map[string]string{
		"userId":   createdUserId,
		"email":    regularUserRegistrationDto.Email,
		"password": regularUserRegistrationDto.Password,
		"username": regularUserRegistrationDto.Username,
		"name":     regularUserRegistrationDto.Name,
		"surname":  regularUserRegistrationDto.Surname,
	})
	requestUrl := fmt.Sprintf("http://%s:%s/register", os.Getenv("AUTHENTICATION_SERVICE_DOMAIN"), os.Getenv("AUTHENTICATION_SERVICE_PORT"))
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func (service *RegularUserService) UpdatePersonalInformations(regularUserUpdateDto dto.RegularUserUpdateDTO) error {
	fmt.Println("Updating regular user")

	if service.RegularUserRepository.ExistByUsername(regularUserUpdateDto.Username) {
		id, _ := primitive.ObjectIDFromHex(regularUserUpdateDto.Id)
		if service.RegularUserRepository.UsernameChanged(regularUserUpdateDto.Username, id) {
			return fmt.Errorf("username is already taken")
		}
	}
	id := regularUserUpdateDto.Id
	var regularUser = createRegularUserFromRegularUserUpdateDTO(&regularUserUpdateDto)
	err := service.RegularUserRepository.UpdatePersonalInformations(regularUser)
	if err != nil {
		return err
	}
	err2 := service.updateUserInAuthenticationService(regularUserUpdateDto, id)
	if err2 != nil {
		return err2
	}
	return nil
}

func (service *RegularUserService) updateUserInAuthenticationService(regularUserUpdateDto dto.RegularUserUpdateDTO, createdUserId string) error {
	postBody, _ := json.Marshal(map[string]string{
		"_id":      createdUserId,
		"email":    regularUserUpdateDto.Email,
		"username": regularUserUpdateDto.Username,
		"name":     regularUserUpdateDto.Name,
		"surname":  regularUserUpdateDto.Surname,
	})
	requestUrl := fmt.Sprintf("http://%s:%s/update", os.Getenv("AUTHENTICATION_SERVICE_DOMAIN"), os.Getenv("AUTHENTICATION_SERVICE_PORT"))
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func (service *RegularUserService) DeleteRegularUser(deleteUserDto dto.DeleteUserDTO) error {
	id, err := primitive.ObjectIDFromHex(deleteUserDto.Id)
	if err != nil {
		return err
	}
	err1 := service.RegularUserRepository.DeleteRegularUser(id)
	if err1 != nil {
		return err1
	}
	err2 := service.deleteUserInAuthenticationService(id.Hex())
	if err2 != nil {
		return err2
	}
	err3 := service.deleteUserDataInMediaContentService(id.Hex())
	if err3 != nil {
		return err3
	}
	err4 := service.deleteUserDataInFollowService(id.Hex())
	if err4 != nil {
		return err4
	}
	err5 := service.deleteUserDataInStoryService(id.Hex())
	if err5 != nil {
		return err4
	}
	return nil
}

func (service *RegularUserService) deleteUserInAuthenticationService(id string) error {
	postBody, _ := json.Marshal(map[string]string{
		"userId": id,
	})
	requestUrl := fmt.Sprintf("http://%s:%s/delete-user", os.Getenv("AUTHENTICATION_SERVICE_DOMAIN"), os.Getenv("AUTHENTICATION_SERVICE_PORT"))
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func (service *RegularUserService) deleteUserDataInMediaContentService(id string) error {
	postBody, _ := json.Marshal(map[string]string{
		"userId": id,
	})
	requestUrl := fmt.Sprintf("http://%s:%s/delete-user-media-content", os.Getenv("MEDIA_CONTENT_SERVICE_DOMAIN"), os.Getenv("MEDIA_CONTENT_SERVICE_PORT"))
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func (service *RegularUserService) deleteUserDataInStoryService(id string) error {
	postBody, _ := json.Marshal(map[string]string{
		"userId": id,
	})
	requestUrl := fmt.Sprintf("http://%s:%s/delete-user-stories", os.Getenv("MEDIA_CONTENT_SERVICE_DOMAIN"), os.Getenv("MEDIA_CONTENT_SERVICE_PORT"))
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func (service *RegularUserService) deleteUserDataInFollowService(id string) error {
	postBody, _ := json.Marshal(map[string]string{
		"userId": id,
	})
	requestUrl := fmt.Sprintf("http://%s:%s/delete-user", os.Getenv("FOLLOW_SERVICE_DOMAIN"), os.Getenv("FOLLOW_SERVICE_PORT"))
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func (service *RegularUserService) UpdateProfilePrivacy(profilePrivacyDto dto.ProfilePrivacyDTO) error {
	fmt.Println("Updating regular user")

	var regularUser = createRegularUserFromProfilePrivacyDTO(&profilePrivacyDto)
	err := service.RegularUserRepository.UpdateProfilePrivacy(regularUser)
	if err != nil {
		return err
	}

	postBody, _ := json.Marshal(map[string]string{
		"_id":         profilePrivacyDto.Id,
		"privacyType": string(profilePrivacyDto.PrivacyType),
	})
	err2 := service.updatePostsPrivacy(postBody)
	err3 := service.updateStoriesPrivacy(postBody)
	if err2 != nil || err3 != nil {
		return err2
	}

	return nil
}

func (service *RegularUserService) updatePostsPrivacy(postBody []byte) error {
	requestUrl := fmt.Sprintf("http://%s:%s/update-posts-privacy", os.Getenv("MEDIA_CONTENT_SERVICE_DOMAIN"), os.Getenv("MEDIA_CONTENT_SERVICE_PORT"))
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func (service *RegularUserService) updateStoriesPrivacy(postBody []byte) error {
	requestUrl := fmt.Sprintf("http://%s:%s/update-stories-privacy", os.Getenv("MEDIA_CONTENT_SERVICE_DOMAIN"), os.Getenv("MEDIA_CONTENT_SERVICE_PORT"))
	resp, err := http.Post(requestUrl, "application/json", bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(resp.StatusCode)
	return nil
}

func (service *RegularUserService) FindUserById(userId primitive.ObjectID) (*model.RegularUser, error) {
	fmt.Print("Searching for logged user...")
	regularUser, err := service.RegularUserRepository.FindUserById(userId)
	if err != nil {
		return nil, err
	}
	return regularUser, err
}

func (service *RegularUserService) CreateRegularUserPostDTOByUsername(username string) (*dto.RegularUserPostDTO, error) {
	regularUser, err := service.RegularUserRepository.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	regularUserPostDto := createRegularUserPostDTOFromRegularUser(regularUser)
	return regularUserPostDto, nil
}

func (service *RegularUserService) FindRegularUserByUsername(username string) (*dto.RegularUserProfileDataDTO, error) {
	regularUser, err := service.RegularUserRepository.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	regularUserPostDto := createRegularUserProfileDataDto(regularUser)
	return regularUserPostDto, nil
}

func (service *RegularUserService) FindRegularUserLikedAndDislikedPosts(username string) (*dto.UserLikedAndDislikedDTO, error) {
	regularUser, err := service.RegularUserRepository.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	userLikedAndDislikedDTO := createRegularUserLikedAndDislikedDTO(regularUser)
	return userLikedAndDislikedDTO, nil
}

func (service *RegularUserService) GetUserSearchResults(searchInput string) ([]model.RegularUser, error) {
	searchPublicRegularUser, err := service.RegularUserRepository.GetAllRegularUsers()
	if err != nil {
		return nil, err
	}
	searchPublicRegularUserModel := CreateUserFromDocuments(searchPublicRegularUser)
	searchPublicRegularUserResults := service.RegularUserRepository.GetUserSearchResults(searchInput, searchPublicRegularUserModel)

	return searchPublicRegularUserResults, nil
}

func (service *RegularUserService) GetAllPublicRegularUsers() ([]dto.RegularUserDTO, error) {
	allRegularUsers, err := service.RegularUserRepository.GetAllPublicRegularUsers()
	if err != nil {
		return nil, err
	}

	allRegularUsersModel := CreateUserFromDocuments(allRegularUsers)

	allRegularUsersDto := createRegularUserDtoFromRegularUser(allRegularUsersModel)
	return allRegularUsersDto, nil
}

func (service *RegularUserService) GetAllRegularUsers() ([]dto.RegularUserDTO, error) {
	allRegularUsers, err := service.RegularUserRepository.GetAllRegularUsers()
	if err != nil {
		return nil, err
	}

	allRegularUsersModel := CreateUserFromDocuments(allRegularUsers)

	allRegularUsersDto := createRegularUserDtoFromRegularUser(allRegularUsersModel)
	return allRegularUsersDto, nil
}

/*
func (service *RegularUserService) GetAllAgentRequests() ([]dto.RegularUserDTO, error){
	allRegularUsers,err := service.RegularUserRepository.GetAllAgentRequests()
	if err != nil {
		return nil, err
	}

	allRegularUsersModel := CreateUserFromDocuments(allRegularUsers)

	allRegularUsersDto := createRegularUserDtoFromRegularUser(allRegularUsersModel)
	return allRegularUsersDto,nil
}*/

func CreateUserFromDocuments(UserDocuments []bson.D) []model.RegularUser {
	var users []model.RegularUser
	for i := 0; i < len(UserDocuments); i++ {
		var user model.RegularUser
		bsonBytes, _ := bson.Marshal(UserDocuments[i])
		_ = bson.Unmarshal(bsonBytes, &user)
		users = append(users, user)
	}
	return users
}
func (service *RegularUserService) FindUsersByIds(usersIds []string) (*[]dto.UserFollowDTO, error) {
	var users []model.RegularUser
	for i := 0; i < len(usersIds); i++ {
		id, _ := primitive.ObjectIDFromHex(usersIds[i])
		regularUser, err := service.RegularUserRepository.FindUserById(id)
		if err != nil {
			return nil, err
		}
		users = append(users, *regularUser)
	}

	userFollowDTOs := createUserFollowDTOsFromRegularUsers(users)
	return userFollowDTOs, nil
}

func createRegularUserPostDTOFromRegularUser(regularUser *model.RegularUser) *dto.RegularUserPostDTO {
	var regularUserPostDto dto.RegularUserPostDTO
	regularUserPostDto.Id = regularUser.Id.Hex()
	regularUserPostDto.PrivacyType = &regularUser.ProfilePrivacy.PrivacyType

	return &regularUserPostDto
}

func createRegularUserFromRegularUserRegistrationDTO(regularUserDto *dto.RegularUserRegistrationDTO) *model.RegularUser {
	profilePrivacy := model.ProfilePrivacy{
		PrivacyType:        model.PrivacyType(0),
		AllMessageRequests: true,
		TagsAllowed:        true,
	}
	var regularUser model.RegularUser
	regularUser.Name = regularUserDto.Name
	regularUser.Surname = regularUserDto.Surname
	regularUser.Username = regularUserDto.Username
	regularUser.Password = regularUserDto.Password
	regularUser.Email = regularUserDto.Email
	regularUser.PhoneNumber = regularUserDto.PhoneNumber
	regularUser.BirthDate = regularUserDto.BirthDate
	regularUser.Biography = regularUserDto.Biography
	regularUser.WebSite = regularUserDto.WebSite
	regularUser.ProfilePrivacy = profilePrivacy
	regularUser.IsDisabled = false
	regularUser.UserRole = model.UserRole(0)
	regularUser.UserType = model.UserType(0)
	regularUser.Gender = regularUserDto.Gender

	return &regularUser
}

func createAgentFromRegularUserRegistrationDTO(regularUserDto *dto.RegularUserRegistrationDTO) *model.Agent {
	profilePrivacy := model.ProfilePrivacy{
		PrivacyType:        model.PrivacyType(0),
		AllMessageRequests: true,
		TagsAllowed:        true,
	}
	var agent model.Agent
	agent.Name = regularUserDto.Name
	agent.Surname = regularUserDto.Surname
	agent.Username = regularUserDto.Username
	agent.Password = regularUserDto.Password
	agent.Email = regularUserDto.Email
	agent.PhoneNumber = regularUserDto.PhoneNumber
	agent.BirthDate = regularUserDto.BirthDate
	agent.Biography = regularUserDto.Biography
	agent.WebSite = regularUserDto.WebSite
	agent.ProfilePrivacy = profilePrivacy
	agent.IsDisabled = false
	agent.UserRole = model.UserRole(2)
	agent.Gender = regularUserDto.Gender
	agent.Verified = false

	return &agent
}

func createRegularUserFromRegularUserUpdateDTO(userUpdateDto *dto.RegularUserUpdateDTO) *model.RegularUser {
	id, _ := primitive.ObjectIDFromHex(userUpdateDto.Id)
	var regularUser model.RegularUser
	regularUser.Id = id
	regularUser.Name = userUpdateDto.Name
	regularUser.Surname = userUpdateDto.Surname
	regularUser.Username = userUpdateDto.Username
	regularUser.Email = userUpdateDto.Email
	regularUser.PhoneNumber = userUpdateDto.PhoneNumber
	regularUser.Gender = userUpdateDto.Gender
	regularUser.BirthDate = userUpdateDto.BirthDate
	regularUser.Biography = userUpdateDto.Biography
	regularUser.WebSite = userUpdateDto.WebSite

	return &regularUser
}

func createRegularUserFromProfilePrivacyDTO(profilePrivacyDto *dto.ProfilePrivacyDTO) *model.RegularUser {
	id, _ := primitive.ObjectIDFromHex(profilePrivacyDto.Id)
	var regularUser model.RegularUser
	regularUser.Id = id
	regularUser.ProfilePrivacy.PrivacyType = profilePrivacyDto.PrivacyType
	regularUser.ProfilePrivacy.AllMessageRequests = profilePrivacyDto.AllMessagesRequests
	regularUser.ProfilePrivacy.TagsAllowed = profilePrivacyDto.TagsAllowed

	return &regularUser
}

func createRegularUserDtoFromRegularUser(allRegularUsers []model.RegularUser) []dto.RegularUserDTO {

	var regularUser []dto.RegularUserDTO
	for i := 0; i < len(allRegularUsers); i++ {
		var regularUserIteration dto.RegularUserDTO
		regularUserIteration.Id = allRegularUsers[i].Id
		regularUserIteration.Username = allRegularUsers[i].Username
		regularUserIteration.Name = allRegularUsers[i].Name
		regularUserIteration.Surname = allRegularUsers[i].Surname
		regularUser = append(regularUser, regularUserIteration)
	}
	return regularUser
}

func createRegularUserProfileDataDto(regularUser *model.RegularUser) *dto.RegularUserProfileDataDTO {
	var regularUserProfileDataDto dto.RegularUserProfileDataDTO

	regularUserProfileDataDto.Id = regularUser.Id
	regularUserProfileDataDto.Name = regularUser.Name
	regularUserProfileDataDto.Surname = regularUser.Surname
	regularUserProfileDataDto.Username = regularUser.Username
	regularUserProfileDataDto.Biography = regularUser.Biography
	regularUserProfileDataDto.WebSite = regularUser.WebSite
	regularUserProfileDataDto.ProfilePrivacy = regularUser.ProfilePrivacy

	return &regularUserProfileDataDto
}
func createUserFollowDTOsFromRegularUsers(regularUsers []model.RegularUser) *[]dto.UserFollowDTO {
	var userFollowDTOs []dto.UserFollowDTO
	for i := 0; i < len(regularUsers); i++ {
		var userFollowDto dto.UserFollowDTO
		userFollowDto.Username = regularUsers[i].Username
		userFollowDto.UserId = regularUsers[i].Id.Hex()
		userFollowDTOs = append(userFollowDTOs, userFollowDto)
	}

	return &userFollowDTOs
}

func createRegularUserLikedAndDislikedDTO(regularUser *model.RegularUser) *dto.UserLikedAndDislikedDTO {
	var userLikedAndDislikedDTO dto.UserLikedAndDislikedDTO

	userLikedAndDislikedDTO.LikedPostsIds = regularUser.LikedPosts
	userLikedAndDislikedDTO.DislikedPostsIds = regularUser.DislikedPosts

	return &userLikedAndDislikedDTO
}

func (service *RegularUserService) UpdateLikedPosts(postLikeDTO dto.UpdatePostLikeAndDislikeDTO) error {
	fmt.Println("Updating regular user liked posts...")

	regularUser, err := service.RegularUserRepository.FindUserByUsername(postLikeDTO.Username)
	if err != nil {
		return err
	}
	if postLikeDTO.IsAdd == "yes" {
		appendedLikes := append(regularUser.LikedPosts, postLikeDTO.PostId)
		regularUser.LikedPosts = appendedLikes
	} else {
		removedLikes := removeFromSlice(regularUser.LikedPosts, postLikeDTO.PostId)
		regularUser.LikedPosts = removedLikes
	}
	err = service.RegularUserRepository.UpdatePersonalInformations(regularUser)
	if err != nil {
		return err
	}

	return nil
}

func (service *RegularUserService) UpdateDislikedPosts(postLikeDTO dto.UpdatePostLikeAndDislikeDTO) error {
	fmt.Println("Updating regular user disliked posts...")

	regularUser, err := service.RegularUserRepository.FindUserByUsername(postLikeDTO.Username)
	if err != nil {
		return err
	}
	if postLikeDTO.IsAdd == "yes" {
		appendedDislikes := append(regularUser.DislikedPosts, postLikeDTO.PostId)
		regularUser.DislikedPosts = appendedDislikes
	} else {
		removedDislikes := removeFromSlice(regularUser.DislikedPosts, postLikeDTO.PostId)
		regularUser.DislikedPosts = removedDislikes
	}
	err = service.RegularUserRepository.UpdatePersonalInformations(regularUser)
	if err != nil {
		return err
	}

	return nil
}

func (service *RegularUserService) SavePost(postSaveDTO dto.PostSaveDTO) error {
	fmt.Println("Updating regular user saved posts...")

	regularUser, err := service.RegularUserRepository.FindUserByUsername(postSaveDTO.Username)
	if err != nil {
		return err
	}
	if postSaveDTO.IsAdd == "yes" {
		var savedPost model.SavedPost
		savedPost.CollectionName = "allPosts"
		savedPost.PostId = postSaveDTO.PostId
		appendedSaved := append(regularUser.SavedPosts, savedPost)
		regularUser.SavedPosts = appendedSaved
	}
	err = service.RegularUserRepository.UpdatePersonalInformations(regularUser)
	if err != nil {
		return err
	}

	return nil
}

func (service *RegularUserService) FindRegularUserSavedPosts(username string) ([]model.SavedPost, error) {
	regularUser, err := service.RegularUserRepository.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return regularUser.SavedPosts, nil
}

func (service *RegularUserService) VerifyUser(userVerificationDto dto.UserVerificationDTO) error {
	fmt.Println("Verifying user ...")

	id, _ := primitive.ObjectIDFromHex(userVerificationDto.UserId)
	user, err1 := service.RegularUserRepository.FindUserById(id)
	if err1 != nil {
		return err1
	}
	user.UserType = userVerificationDto.VerificationType
	err2 := service.RegularUserRepository.UpdateUserType(user)
	if err2 != nil {
		return err2
	}
	return nil
}

func removeFromSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

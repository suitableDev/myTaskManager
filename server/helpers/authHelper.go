package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(context *gin.Context, role string) (err error) {
	userType := context.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("unauthorised to acces this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(context *gin.Context, UserID string) (err error) {
	userType := context.GetString("user_type")
	uid := context.GetString("uid")

	if userType == "USER" && uid != UserID {
		err = errors.New("unauthorised to access this resource")
		if err != nil {
			return err
		}
	}

	err = CheckUserType(context, userType)
	return err
}

package apiThird

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRTCInvitationInfo(c *gin.Context) {
	var (
		req  api.GetRTCInvitationInfoReq
		resp api.GetRTCInvitationInfoResp
	)
	if err := c.Bind(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	var ok bool
	var errInfo string
	ok, _, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	var err error
	invitationInfo, err := db.DB.GetSignalInfoFromCache(req.ClientMsgID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetSignalInfoFromCache", err.Error(), req)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.Data.OpUserID = invitationInfo.OpUserID
	resp.Data.Invitation.RoomID = invitationInfo.Invitation.RoomID
	resp.Data.Invitation.SessionType = invitationInfo.Invitation.SessionType
	resp.Data.Invitation.GroupID = invitationInfo.Invitation.GroupID
	resp.Data.Invitation.InviterUserID = invitationInfo.Invitation.InviterUserID
	resp.Data.Invitation.InviteeUserIDList = invitationInfo.Invitation.InviteeUserIDList
	resp.Data.Invitation.MediaType = invitationInfo.Invitation.MediaType
	resp.Data.Invitation.Timeout = invitationInfo.Invitation.Timeout
	c.JSON(http.StatusOK, resp)
}
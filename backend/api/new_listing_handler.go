package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AustrianDataLAB/GeWoScout/backend/cosmos"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/go-chi/render"
)

func (h *Handler) HandleNewListingResult(w http.ResponseWriter, r *http.Request) {
	injectedData := models.QueueBindingInput{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&injectedData); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NewListingHandler | Failed to read invoke request body: %s", err.Error())},
		})
		return
	}

	msgId, err := strconv.Unquote(injectedData.Metadata.Id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NewListingHandler %s | Failed to unquote message ID: %s", injectedData.Metadata.Id, err.Error())},
		})
		return
	}

	msgPlain, err := strconv.Unquote(injectedData.Data.Msg)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NewListingHandler %s | Failed to unquote message: %s", msgId, err.Error())},
		})
		return
	}

	logs := make([]string, 0)

	newListing := models.Listing{}
	err = json.Unmarshal([]byte(msgPlain), &newListing)
	if err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NewListingHandler %s | Failed to unmarshal message: %s", msgId, err.Error())},
		})
		return
	}

	container := h.GetNotificationSettingsByCityContainerClient()

	getCtx, getCancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer getCancel()
	mails, err := cosmos.GetUsersMatchingWithListing(getCtx, container, &newListing)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NewListingHandler %s | Failed to get users: %s", msgId, err.Error())},
		})
		return
	}

	if len(mails) == 0 {
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NewListingHandler %s | No users found matching the listing %s", msgId, newListing.ID)},
		})
		return
	}

	logs = append(logs, fmt.Sprintf("NewListingHandler %s | Found %d users matching the listing %s", msgId, len(mails), newListing.ID))

	fmt.Println(mails)

	notificationData := models.NotificationData{
		Emails:  mails,
		Listing: newListing,
	}

	// For each non-existant ID, insert the listing and create a queue message
	invokeResponse := models.InvokeResponse{
		Logs: logs,
		Outputs: map[string]interface{}{
			"msgOut": []models.NotificationData{
				notificationData,
			},
		},
	}
	render.JSON(w, r, invokeResponse)
}

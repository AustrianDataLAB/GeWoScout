package api

import (
	"encoding/json"
	"fmt"
	"github.com/AustrianDataLAB/GeWoScout/backend/models"
	"github.com/AustrianDataLAB/GeWoScout/backend/notification"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

func (h *Handler) HandleNotification(w http.ResponseWriter, r *http.Request) {
	injectedData := models.QueueBindingInput{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&injectedData); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NotificationHandler | Failed to read invoke request body: %s", err.Error())},
		})
		return
	}

	msgId, err := strconv.Unquote(injectedData.Metadata.Id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NotificationHandler %s | Failed to unquote message ID: %s", injectedData.Metadata.Id, err.Error())},
		})
		return
	}

	msgPlain, err := strconv.Unquote(injectedData.Data.Msg)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NotificationHandler %s | Failed to unquote message: %s", msgId, err.Error())},
		})
		return
	}

	logs := make([]string, 0)

	newListing := models.NotificationData{}
	err = json.Unmarshal([]byte(msgPlain), &newListing)
	if err != nil {
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NotificationHandler %s | Failed to unmarshal message: %s", msgId, err.Error())},
		})
		return
	}

	err = notification.SendNotification(newListing.Listing, newListing.Emails)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.InvokeResponse{
			Logs: []string{fmt.Sprintf("NotificationHandler %s | Failed to send notification: %s", msgId, err.Error())},
		})
		return
	}

	logs = append(logs, fmt.Sprintf("NotificationHandler %s | Successfully notified %d users about listing %s", msgId, len(newListing.Emails), newListing.Listing.ID))

	invokeResponse := models.InvokeResponse{
		Logs:    logs,
		Outputs: map[string]interface{}{},
	}
	render.JSON(w, r, invokeResponse)
}

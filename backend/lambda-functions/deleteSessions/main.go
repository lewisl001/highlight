package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/highlight-run/highlight/backend/lambda-functions/deleteSessions/handlers"
	"github.com/highlight-run/highlight/backend/lambda-functions/deleteSessions/utils"
	"github.com/highlight-run/highlight/backend/util"
	"github.com/highlight/highlight/sdk/highlight-go"
	hlog "github.com/highlight/highlight/sdk/highlight-go/log"
)

// Meant for local invocation for testing the lambda handler stack
func main() {
	if !util.IsDevOrTestEnv() {
		return
	}

	highlight.SetProjectID("1jdkoe52")
	highlight.Start()
	defer highlight.Stop()
	hlog.Init()

	h := handlers.NewHandlers()
	input := utils.QuerySessionsInput{
		ProjectId:    1,
		Email:        "zane@highlight.io",
		FirstName:    "Zane",
		Query:        "{\"bool\":{\"must\":[{\"bool\":{\"should\":[{\"term\":{\"processed\":\"true\"}}]}},{\"bool\":{\"should\":[{\"range\":{\"created_at\":{\"gte\":\"2022-07-15T23:00:25.525Z\",\"lte\":\"2022-08-01T23:00:25.525Z\"}}}]}}]}}",
		SessionCount: 256,
		DryRun:       true,
	}
	ctx := context.TODO()
	out, err := h.GetSessionIdsByQuery(ctx, input)
	if err != nil {
		log.WithContext(ctx).Fatal(err)
	}

	if len(out) == 0 {
		log.WithContext(ctx).Fatal("GetSessionIdsByQuery returned an empty list")
	}

	if _, err := h.DeleteSessionBatchFromPostgres(ctx, out[0]); err != nil {
		log.WithContext(ctx).Fatal(err)
	}
	if _, err := h.DeleteSessionBatchFromS3(ctx, out[0]); err != nil {
		log.WithContext(ctx).Fatal(err)
	}
	if _, err := h.DeleteSessionBatchFromOpenSearch(ctx, out[0]); err != nil {
		log.WithContext(ctx).Fatal(err)
	}

	if err := h.SendEmail(ctx, input); err != nil {
		log.WithContext(ctx).Fatal(err)
	}
}

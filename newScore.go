// Package p contains a Firestore Cloud Function.
package p

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/functions/metadata"
	"google.golang.org/api/iterator"
)

// These variables are used for logging and are automatically
// set by the Cloud Functions runtime.
var (
	projectID = os.Getenv("GCLOUD_PROJECT")
)

// FirestoreEvent is the payload of a Firestore event.
// Please refer to the docs for additional information
// regarding Firestore events.
type FirestoreEvent struct {
	OldValue FirestoreValue `json:"oldValue"`
	Value    FirestoreValue `json:"value"`
}

func getDocumentID(ctx context.Context) (string, error) {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return "", err
	}
	parts := strings.Split(meta.Resource.Name, "/")
	if len(parts) == 0 {
		return "", errors.New("Error getting ID from context")
	}
	return parts[len(parts)-1], nil
}

func getCreateTime(ctx context.Context) (time.Time, error) {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return meta.Timestamp, nil
}

// OnNewScore is triggered by a change to a Firestore document.
func OnNewScore(ctx context.Context, e FirestoreEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)

	// Get data from Score
	score, err := NewFromFirestoreValue(e.Value)
	if err != nil {
		return err
	}
	log.Printf("%+v", score)

	client, err := firestore.NewClient(ctx, projectID)
	defer client.Close()

	// Compute the rank, based on the Points collection

	limit := 100
	var maxRank int64 = 300
	iter := client.Collection("Points").Where("score", ">=", score.Points).Limit(limit).Documents(ctx)
	var rank int64 = 1
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		score := doc.Data()
		count, ok := score["count"].(int64)
		if !ok {
			return errors.New("Error extracting count from Points")
		}
		rank += count
		if rank > maxRank {
			rank = 0
			break
		}
	}

	log.Printf("RANK: %d", rank)

	docID, err := getDocumentID(ctx)
	if err != nil {
		return nil
	}

	// Add rank value into the Score document
	client.Collection("Score").Doc(docID).Set(ctx, map[string]interface{}{
		"rank": rank,
	}, firestore.MergeAll)

	// Add this score to the Points collection
	createTime, err := getCreateTime(ctx)
	if err != nil {
		return nil
	}
	addToPoints(ctx, client, docID, createTime, score)

	return nil
}

func addToPoints(ctx context.Context, client *firestore.Client, scoresId string, createTime time.Time, score *Score) error {

	pointsID := fmt.Sprintf("%d", score.Points)
	snap, err := client.Collection("Points").Doc(pointsID).Get(ctx)

	// scoresList is the list of scores for these points
	var scoresList []interface{}

	if err != nil {
		// document not found, create an empty list of scores
		scoresList = []interface{}{}
	} else {
		// document found, get previous scores
		data := snap.Data()
		var ok bool
		scoresList, ok = data["scoresList"].([]interface{})
		if !ok {
			return errors.New(fmt.Sprintf("Error extracting scores for %d points", score.Points))
		}
	}

	// Add this score to the list
	scoresList = append(scoresList, map[string]interface{}{
		"uid":     score.Uid,
		"name":    score.Name,
		"country": score.Country,
		"details": score.Details,
		"time":    createTime,
	})

	// Write new content
	client.Collection("Points").Doc(pointsID).Set(ctx, map[string]interface{}{
		"score":      score.Points,
		"scoresList": scoresList,
		"count":      len(scoresList),
	}, firestore.MergeAll)

	return nil
}

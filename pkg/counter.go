package pkg

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type counterHandler struct {
	Active            bool
	mongodbClient     *mongo.Client
	mongodbDatabase   string
	mongodbCollection string
	mongodbDocumentID string
	counterID         primitive.ObjectID
	collection        *mongo.Collection
}

// Counter storage model
type Counter struct {
	Count int `bson:"count"`
}

var initOnce sync.Once

// Content returns the counter document from the backing store
func (h *counterHandler) Content(reqCtx context.Context) Counter {
	initOnce.Do(func() {
		h.collection = h.mongodbClient.Database(h.mongodbDatabase).Collection(h.mongodbCollection)
		h.counterID, _ = primitive.ObjectIDFromHex(h.mongodbDocumentID)
	})

	h.incrementCounter(reqCtx)

	ctx, cancel := context.WithTimeout(reqCtx, 3*time.Second)
	defer cancel()

	// get record from MongoDB
	result := Counter{}
	err := h.collection.FindOne(ctx, bson.M{"_id": h.counterID}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

// incrementCounter is used to update the visit counter on each request to /counter
func (h *counterHandler) incrementCounter(reqCtx context.Context) {
	ctx, cancel := context.WithTimeout(reqCtx, 3*time.Second)
	defer cancel()

	result, err := h.collection.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: h.counterID}},
		bson.D{{Key: "$inc", Value: bson.D{{Key: "count", Value: 1}}}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Updated %v Documents!\n", result.ModifiedCount)
}

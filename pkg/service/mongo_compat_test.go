package service

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBridgeObjectIDRoundTrip(t *testing.T) {
	oid := primitive.NewObjectID()

	gotUUID := bridgeObjectIDToUUID(oid)
	gotOID, err := bridgeUUIDToObjectID(gotUUID)
	if err != nil {
		t.Fatalf("bridgeUUIDToObjectID returned error: %v", err)
	}
	if gotOID != oid {
		t.Fatalf("round trip mismatch: got %s want %s", gotOID.Hex(), oid.Hex())
	}
}

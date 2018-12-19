package main

import (
	"github.com/globalsign/mgo"

	pb "github.com/aaronflower/shippy-service-consignment/proto/consignment"
)

const (
	dbName                = "shippy"
	consignmentCollection = "consignments"
)

// Repository defines Repository interface.
type Repository interface {
	Create(*pb.Consignment) error
	GetAll() ([]*pb.Consignment, error)
	Close()
}

// ConsignmentRepository - Dummy repository, this simulates the use oa a datastroe
// of some kind. We'll replace this with a real implementation later on.
type ConsignmentRepository struct {
	session *mgo.Session
}

// Create a new consignment
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) error {
	return repo.collection().Insert(consignment)
}

// GetAll consignments
func (repo *ConsignmentRepository) GetAll() ([]*pb.Consignment, error) {
	var consignments []*pb.Consignment
	err := repo.collection().Find(nil).All(&consignments)
	return consignments, err
}

// Close closes the database session after each query has ran.
func (repo *ConsignmentRepository) Close() {
	repo.session.Close()
}

func (repo *ConsignmentRepository) collection() *mgo.Collection {
	return repo.session.DB(dbName).C(consignmentCollection)
}

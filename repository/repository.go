package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ServiceCoverage struct {
	ServiceName string  `json:"service_name"`
	Coverage    float32 `json:"coverage"`
}

type Repository interface {
	UpdateServiceCoverage(ServiceCoverage) error
	ListServiceCoverage() []ServiceCoverage
}

type DynamodbRepository struct {
	TableName string
	CL        dynamodb.Client
}

func (dr *DynamodbRepository) UpdateServiceCoverage(sc ServiceCoverage) (err error) {
	item, _ := attributevalue.MarshalMap(sc)
	_, err = dr.cl.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: &dr.tableName})
	return
}

func (dr *DynamodbRepository) ListServiceCoverage() (sc []ServiceCoverage) {
	scan, _ := dr.cl.Scan(context.TODO(), &dynamodb.ScanInput{TableName: &dr.tableName})
	attributevalue.UnmarshalListOfMaps(scan.Items, sc)
	return
}

type LocalRepository struct {
	Services []ServiceCoverage
}

func (lr *LocalRepository) UpdateServiceCoverage(sc ServiceCoverage) (err error) {
	for i, v := range lr.Services {
		if v.ServiceName == sc.ServiceName {
			lr.Services[i].Coverage = sc.Coverage
			return
		}
	}
	lr.Services = append(lr.Services, sc)
	return
}

func (lr *LocalRepository) ListServiceCoverage() (sc []ServiceCoverage) {
	sc = lr.Services
	return
}

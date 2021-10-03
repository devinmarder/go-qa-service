package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ServiceCoverage struct {
	ServiceName string  `json:"service_name"`
	Coverage    float32 `json:"coverage"`
}

type Repository interface {
	UpdateServiceCoverage(serviceName string, coverage float32) error
	ListServiceCoverage() []ServiceCoverage
}

type DynamodbRepository struct {
	tableName string
	cl        dynamodb.Client
}

func (dr *DynamodbRepository) UpdateServiceCoverage(sn string, cov float32) (err error) {
	_, err = dr.cl.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"service_name": &types.AttributeValueMemberS{Value: sn},
			"coverage":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", cov)}},
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

func (lr *LocalRepository) UpdateServiceCoverage(sn string, cov float32) (err error) {
	for i, v := range lr.Services {
		if v.ServiceName == sn {
			lr.Services[i].Coverage = cov
			return
		}
	}
	newServiceCoverage := ServiceCoverage{ServiceName: sn, Coverage: cov}
	lr.Services = append(lr.Services, newServiceCoverage)
	return
}

func (lr *LocalRepository) ListServiceCoverage() (sc []ServiceCoverage) {
	sc = lr.Services
	return
}

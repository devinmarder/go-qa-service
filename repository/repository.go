package repository

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ServiceCoverage struct {
	ServiceName string  `json:"service_name" dynamodbav:"service_name"`
	Coverage    float32 `json:"coverage" dynamodbav:"coverage"`
}

type Repository interface {
	UpdateServiceCoverage(ServiceCoverage) error
	ListServiceCoverage() ([]ServiceCoverage, error)
}

func ConfigureRepository(repo *Repository) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	cl := dynamodb.NewFromConfig(cfg)
	*repo = &DynamodbRepository{Cl: cl, TableName: os.Getenv("QA_TABLE_NAME")}
}

type DynamodbRepository struct {
	TableName string
	Cl        *dynamodb.Client
}

func (dr *DynamodbRepository) UpdateServiceCoverage(sc ServiceCoverage) (err error) {
	item, _ := attributevalue.MarshalMap(sc)
	_, err = dr.Cl.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: &dr.TableName})
	return
}

func (dr *DynamodbRepository) ListServiceCoverage() (sc []ServiceCoverage, err error) {
	scan, err := dr.Cl.Scan(context.TODO(), &dynamodb.ScanInput{TableName: &dr.TableName})
	attributevalue.UnmarshalListOfMaps(scan.Items, &sc)
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

func (lr *LocalRepository) ListServiceCoverage() (sc []ServiceCoverage, err error) {
	sc = lr.Services
	return
}

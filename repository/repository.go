package repository

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//ServiceCoverage is the schema used for the repository items.
type ServiceCoverage struct {
	ServiceName string  `json:"service_name" dynamodbav:"service_name"`
	Coverage    float32 `json:"coverage" dynamodbav:"coverage"`
}

//Repository is an interface used for updating and listing the stored items.
type Repository interface {
	UpdateServiceCoverage(ServiceCoverage) error
	ListServiceCoverage() ([]ServiceCoverage, error)
}

//ConfigureRepository creates the aws-sdk client and attaches it to the provided repository.
func ConfigureRepository(repo *Repository) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	cl := dynamodb.NewFromConfig(cfg)
	*repo = &DynamodbRepository{Cl: cl, TableName: os.Getenv("QA_TABLE_NAME")}
}

//DynamodbRepository is the repository used for interacting with a DynamodbTable.
type DynamodbRepository struct {
	TableName string
	Cl        *dynamodb.Client
}

//UpdateServiceCoverage updates the repository with the provided ServiceCoverage.
//If the service already exists, then the coverage is overwritten, otherwise a new item is created.
func (dr *DynamodbRepository) UpdateServiceCoverage(sc ServiceCoverage) (err error) {
	item, _ := attributevalue.MarshalMap(sc)
	_, err = dr.Cl.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: &dr.TableName})
	return
}

//ListServiceCoverage returns a list of all services on the Repo.
func (dr *DynamodbRepository) ListServiceCoverage() (sc []ServiceCoverage, err error) {
	scan, err := dr.Cl.Scan(context.TODO(), &dynamodb.ScanInput{TableName: &dr.TableName})
	attributevalue.UnmarshalListOfMaps(scan.Items, &sc)
	return
}

//LocalRepository stores items in memory and can be used for debugging.
type LocalRepository struct {
	Services []ServiceCoverage
}

//UpdateServiceCoverage updates the repository with the provided ServiceCoverage.
//If the service already exists, then the coverage is overwritten, otherwise a new item is created.
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

//ListServiceCoverage returns a list of all services on the Repo.
func (lr *LocalRepository) ListServiceCoverage() (sc []ServiceCoverage, err error) {
	sc = lr.Services
	return
}

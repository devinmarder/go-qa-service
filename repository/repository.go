package repository

type ServiceCoverage struct {
	ServiceName string `json:"service_name"`
	Coverage    int    `json:"coverage"`
}

type Repository interface {
	UpdateServiceCoverage(serviceName string, coverage int)
	ListServiceCoverage() []ServiceCoverage
}

type LocalRepository struct {
	Services []ServiceCoverage
}

func (lr *LocalRepository) UpdateServiceCoverage(sn string, cov int) {
	for i, v := range lr.Services {
		if v.ServiceName == sn {
			lr.Services[i].Coverage = cov
			return
		}
	}
	newServiceCoverage := ServiceCoverage{ServiceName: sn, Coverage: cov}
	lr.Services = append(lr.Services, newServiceCoverage)
}

func (lr *LocalRepository) ListServiceCoverage() []ServiceCoverage {
	return lr.Services
}

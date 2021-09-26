package repository

type ServiceCoverage struct {
	ServiceName string  `json:"service_name"`
	Coverage    float32 `json:"coverage"`
}

type Repository interface {
	UpdateServiceCoverage(serviceName string, coverage float32)
	ListServiceCoverage() []ServiceCoverage
}

type LocalRepository struct {
	Services []ServiceCoverage
}

func (lr *LocalRepository) UpdateServiceCoverage(sn string, cov float32) {
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

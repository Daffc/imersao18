package service

import "fmt"

type PartnerFactory interface {
	CreatePartner(partnerId int) (Partner, error)
}

type DefaultPartnerFactory struct {
	partnerBaseURLs map[int]string
}

func NewPartnerfactory(partnerBaseURLs map[int]string) PartnerFactory {
	return &DefaultPartnerFactory{partnerBaseURLs: partnerBaseURLs}
}

func (f *DefaultPartnerFactory) CreatePartner(partnerId int) (Partner, error) {
	baseURL, ok := f.partnerBaseURLs[partnerId]
	if !ok {
		return nil, fmt.Errorf("partner with Id %d not found", partnerId)
	}

	switch partnerId {
	case 1:
		return &Partner1{BaseURL: baseURL}, nil
	case 2:
		return &Partner2{BaseURL: baseURL}, nil
	default:
		return nil, fmt.Errorf("partner with Id %d not found", partnerId)
	}
}

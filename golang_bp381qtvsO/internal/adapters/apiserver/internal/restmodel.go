package internal

import (
	"errors"
	"example_consumer/internal/core/model"
	"fmt"
	"github.com/samber/lo"
)

type ContactToSaveRest struct {
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Phones    []PhoneRest `json:"phones"`
}

type PhoneRest struct {
	PhoneType   string `json:"phone_type"`
	PhoneNumber string `json:"phone_number"`
}

type ContactRest struct {
	ID        string      `json:"id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Phones    []PhoneRest `json:"phones"`
}

func (r *ContactToSaveRest) toModel() (*model.ContactToSave, error) {
	if r.FirstName == "" {
		return nil, errors.New("first_name must not be empty")
	}
	if r.LastName == "" {
		return nil, errors.New("last_name must not be empty")
	}
	phones := make([]*model.ContactPhoneToSave, len(r.Phones))
	for i, phone := range r.Phones {
		phoneModel, err := phone.toContactPhoneToSaveModel()
		if err != nil {
			return nil, err
		}
		phones[i] = phoneModel
	}
	return &model.ContactToSave{
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Phones:    phones,
	}, nil
}

func (r PhoneRest) toContactPhoneToSaveModel() (*model.ContactPhoneToSave, error) {
	phoneType, err := phoneTypeRestToModel(r.PhoneType)
	if err != nil {
		return nil, err
	}
	return &model.ContactPhoneToSave{
		PhoneType:   *phoneType,
		PhoneNumber: r.PhoneNumber,
	}, nil
}

func phoneTypeRestToModel(phoneType string) (*model.ContactPhoneType, error) {
	switch phoneType {
	case "mobile":
		return lo.ToPtr(model.ContactPhoneTypeMobile), nil
	case "home":
		return lo.ToPtr(model.ContactPhoneTypeHome), nil
	case "work":
		return lo.ToPtr(model.ContactPhoneTypeWork), nil
	default:
		return nil, fmt.Errorf("unknown contact phone type: %s", phoneType)
	}
}

func phoneTypeModelToRest(phoneType model.ContactPhoneType) string {
	switch phoneType {
	case model.ContactPhoneTypeMobile:
		return "mobile"
	case model.ContactPhoneTypeHome:
		return "home"
	case model.ContactPhoneTypeWork:
		return "work"
	default:
		panic(fmt.Sprintf("Unexpected phone type model: %s", phoneType))
	}
}

func contactModelToRest(m *model.Contact) *ContactRest {
	phones := lo.Map(m.Phones, func(item *model.ContactPhone, _ int) PhoneRest {
		return PhoneRest{
			PhoneType:   phoneTypeModelToRest(item.PhoneType),
			PhoneNumber: item.PhoneNumber,
		}

	})
	return &ContactRest{
		ID:        m.ID,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Phones:    phones,
	}
}

type VersionRest struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Build   string `json:"build"`
}

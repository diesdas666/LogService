package outport

import (
	"context"
	"example_consumer/internal/core/model"
)

type AddrBook interface {
	LoadAllContacts(ctx context.Context) ([]*model.Contact, error)
	LoadContactByID(ctx context.Context, ID string) (*model.Contact, error)
	AddContact(ctx context.Context, c *model.ContactToSave) (*model.Contact, error)
	UpdateContact(ctx context.Context, ID string, c *model.ContactToSave) (*model.Contact, error)
	DeleteContact(ctx context.Context, ID string) (found bool, err error)
}

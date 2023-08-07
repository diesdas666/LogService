package usecase

import (
	"context"
	"example_consumer/internal/core/app"
	"example_consumer/internal/core/model"
)

func (uc *UseCases) LoadAddrBookContacts(
	ctx context.Context,
) ([]*model.Contact, error) {
	app.Logger(ctx).Debug("Load all address book contacts")
	contacts, err := uc.AddrBook.LoadAllContacts(ctx)
	if err != nil {
		app.Logger(ctx).Errorf("Loading all address book contacts failed with error: %v", err)
		return nil, err
	}
	app.Logger(ctx).Debugf("Loaded address book contacts: %v", contacts)
	return contacts, nil
}

func (uc *UseCases) LoadAddrBookContactByID(
	ctx context.Context,
	ID string,
) (*model.Contact, error) {
	app.Logger(ctx).Debugf("Load address book contact by id=%s", ID)
	contact, err := uc.AddrBook.LoadContactByID(ctx, ID)
	if err != nil {
		app.Logger(ctx).Errorf("Loading address book contact by id=%s failed: %v", ID, err)
		return nil, err
	}
	if contact == nil {
		app.Logger(ctx).Infof("No address book contact found with id=%s", ID)
	} else {
		app.Logger(ctx).Debugf("Loaded address book contact: %v", contact)
	}
	return contact, nil
}

func (uc *UseCases) AddAddrBookContact(
	ctx context.Context,
	contact *model.ContactToSave,
) (*model.Contact, error) {
	app.Logger(ctx).Debugf("Add address book contact: %v", contact)
	newContact, err := uc.AddrBook.AddContact(ctx, contact)
	if err != nil {
		app.Logger(ctx).Errorf("Adding new address book contact failed with error: %v", err)
		return nil, err
	}
	app.Logger(ctx).Debugf("Added address book contact: %v", newContact)
	return newContact, nil
}

func (uc *UseCases) UpdateAddrBookContact(
	ctx context.Context,
	ID string,
	contact *model.ContactToSave,
) (updatedContact *model.Contact, found bool, err error) {
	app.Logger(ctx).Debugf("Update address book contact by id=%s with value: %v", ID, contact)
	updatedContact, err = uc.AddrBook.UpdateContact(ctx, ID, contact)
	if err != nil {
		app.Logger(ctx).Errorf("Update address book contact by id=%s failed with error: %v", ID, err)
		return nil, false, err
	}
	if updatedContact == nil {
		app.Logger(ctx).Infof("Attempt to update non-existing contact by id=%s", ID)
		return nil, false, nil
	}
	app.Logger(ctx).Debugf("Updated address book contact by ID=%s", ID)
	return updatedContact, true, nil
}

func (uc *UseCases) DeleteAddrBookContact(
	ctx context.Context,
	ID string,
) (found bool, err error) {
	app.Logger(ctx).Debugf("Delete address book contact by id=%s", ID)
	found, err = uc.AddrBook.DeleteContact(ctx, ID)
	if err != nil {
		app.Logger(ctx).Errorf("Deleting address book contact by id=%s failed with error:", ID)
		return
	}
	if found {
		app.Logger(ctx).Debugf("Deleted address book contact by id=%s", ID)
	} else {
		app.Logger(ctx).Infof("Attempt to delete non-existing contact by id=%s", ID)
	}
	return
}

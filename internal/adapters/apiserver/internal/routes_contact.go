package internal

import (
	"example_consumer/internal/core/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
)

const ContactsRoutePath = "/api/contacts"

func CreateContact(uc *usecase.UseCases) func(echo.Context) error {
	return func(c echo.Context) error {
		req := new(ContactToSaveRest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, NewBadRequestErrResponse(err))
		}
		contactToSave, err := req.toModel()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, NewBadRequestErrResponse(err))
		}
		contact, err := uc.AddAddrBookContact(c.Request().Context(), contactToSave)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, NewInternalServerErrResponse(err))
		}
		contactRest := contactModelToRest(contact)
		return c.JSON(http.StatusCreated, contactRest)
	}
}

func UpdateContact(uc *usecase.UseCases) func(echo.Context) error {
	return func(c echo.Context) error {
		ID := c.Param("id")
		req := new(ContactToSaveRest)
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, NewBadRequestErrResponse(err))
		}
		contactToSave, err := req.toModel()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, NewBadRequestErrResponse(err))
		}
		contact, found, err := uc.UpdateAddrBookContact(c.Request().Context(), ID, contactToSave)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, NewInternalServerErrResponse(err))
		}
		if !found {
			return echo.NewHTTPError(http.StatusNotFound, NotFoundErrResponse)
		}
		contactRest := contactModelToRest(contact)
		return c.JSON(http.StatusOK, contactRest)
	}
}

func ListAllContacts(uc *usecase.UseCases) func(echo.Context) error {
	return func(c echo.Context) error {
		contacts, err := uc.LoadAddrBookContacts(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, NewInternalServerErrResponse(err))
		}
		contactRestList := make([]*ContactRest, len(contacts))
		for i, contact := range contacts {
			contactRestList[i] = contactModelToRest(contact)
		}
		return c.JSON(http.StatusOK, contactRestList)
	}
}

func GetContact(uc *usecase.UseCases) func(echo.Context) error {
	return func(c echo.Context) error {
		ID := c.Param("id")
		contact, err := uc.LoadAddrBookContactByID(c.Request().Context(), ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, NewInternalServerErrResponse(err))
		}
		if contact == nil {
			return echo.NewHTTPError(http.StatusNotFound, NotFoundErrResponse)
		}
		contactRest := contactModelToRest(contact)
		return c.JSON(http.StatusOK, contactRest)
	}
}

func DeleteContact(uc *usecase.UseCases) func(echo.Context) error {
	return func(c echo.Context) error {
		ID := c.Param("id")
		found, err := uc.DeleteAddrBookContact(c.Request().Context(), ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, NewInternalServerErrResponse(err))
		}
		if !found {
			return echo.NewHTTPError(http.StatusNotFound, NotFoundErrResponse)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

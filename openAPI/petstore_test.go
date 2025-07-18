package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"open-api/api"
	"testing"

	"github.com/oapi-codegen/testutil"
	"github.com/stretchr/testify/assert"
)

func doGet(t *testing.T, handler http.Handler, url string) *httptest.ResponseRecorder {
	response := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, handler)
	return response.Recorder
}

func TestPetStore(t *testing.T) {
	var err error

	store := api.NewPetStore()
	ginPetServer := NewGinPetServer(store, "8080")
	r := ginPetServer.Handler

	t.Run("add pet", func(t *testing.T) {
		tag := "TahOfSport"
		newPet := api.NewPet{
			Name: "Spot",
			Tag:  &tag,
		}

		rr := testutil.NewRequest().Post("/pets").WithJsonBody(newPet).GoWithHTTPHandler(t, r).Recorder

		assert.Equal(t, http.StatusCreated, rr.Code)

		var resultPet api.Pet

		err = json.NewDecoder(rr.Body).Decode(&resultPet)
		assert.NoError(t, err, "error unmarshaling response")
		assert.Equal(t, newPet.Name, resultPet.Name)
		assert.Equal(t, *newPet.Tag, *resultPet.Tag)
		assert.Equal(t, resultPet.Id, int64(100))

	})

	t.Run("Find pet by ID", func(t *testing.T) {
		tag := "TahOfSport"
		pet := api.Pet{
			Id:   101,
			Name: "Spot",
			Tag:  &tag,
		}

		store.Pets[pet.Id] = pet

		// rr := doGet()

		url := fmt.Sprintf("/pets/%d", pet.Id)

		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusOK, rr.Code)

		var resultPet api.Pet

		err = json.NewDecoder(rr.Body).Decode(&resultPet)
		assert.NoError(t, err, "error unmarshaling response")
		assert.Equal(t, resultPet, pet)
	})

	t.Run("Pet not found", func(t *testing.T) {

		url := fmt.Sprintf("/pets/%d", 1)

		rr := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, r).Recorder
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("List all pets", func(t *testing.T) {
		store.Pets = map[int64]api.Pet{1: {}, 2: {}}

		// Now, list all pets, we should have two
		rr := doGet(t, r, "/pets")
		assert.Equal(t, http.StatusOK, rr.Code)

		var petList []api.Pet
		err = json.NewDecoder(rr.Body).Decode(&petList)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, 2, len(petList))
	})

	t.Run("Filter pets by tag", func(t *testing.T) {
		store.Pets = map[int64]api.Pet{1: {}, 2: {}}

		// Filter pets by non existent tag, we should have 0
		rr := doGet(t, r, "/pets?tags=NotExists")
		assert.Equal(t, http.StatusOK, rr.Code)

		var petList []api.Pet
		err = json.NewDecoder(rr.Body).Decode(&petList)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, 0, len(petList))
	})

	t.Run("Filter pets by tag", func(t *testing.T) {
		tag := "TagOfFido"

		store.Pets = map[int64]api.Pet{
			1: {
				Tag: &tag,
			},
			2: {},
		}

		// Filter pets by tag, we should have 1
		rr := doGet(t, r, "/pets?tags=TagOfFido")
		assert.Equal(t, http.StatusOK, rr.Code)

		var petList []api.Pet
		err = json.NewDecoder(rr.Body).Decode(&petList)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, 1, len(petList))
	})
	t.Run("Delete pets", func(t *testing.T) {
		store.Pets = map[int64]api.Pet{1: {}, 2: {}}

		// Let's delete non-existent pet
		rr := testutil.NewRequest().Delete("/pets/7").GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNotFound, rr.Code)

		var petError api.Error
		err = json.NewDecoder(rr.Body).Decode(&petError)
		assert.NoError(t, err, "error unmarshaling PetError")
		assert.Equal(t, int32(http.StatusNotFound), petError.Code)

		// Now, delete both real pets
		rr = testutil.NewRequest().Delete("/pets/1").GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNoContent, rr.Code)

		rr = testutil.NewRequest().Delete("/pets/2").GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusNoContent, rr.Code)

		// Should have no pets left.
		var petList []api.Pet
		rr = doGet(t, r, "/pets")
		assert.Equal(t, http.StatusOK, rr.Code)
		err = json.NewDecoder(rr.Body).Decode(&petList)
		assert.NoError(t, err, "error getting response", err)
		assert.Equal(t, 0, len(petList))
	})
}

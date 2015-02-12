package databases

import (
	"testing"

	os "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

const instanceID = "{instanceID}"

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	os.HandleCreateDBSuccessfully(t, instanceID)

	opts := os.BatchCreateOpts{
		os.CreateOpts{Name: "testingdb", CharSet: "utf8", Collate: "utf8_general_ci"},
		os.CreateOpts{Name: "sampledb"},
	}

	res := Create(fake.ServiceClient(), instanceID, opts)
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	os.HandleListDBsSuccessfully(t, instanceID)

	expectedDBs := []os.Database{
		os.Database{Name: "anotherexampledb"},
		os.Database{Name: "exampledb"},
		os.Database{Name: "nextround"},
		os.Database{Name: "sampledb"},
		os.Database{Name: "testingdb"},
	}

	pages := 0
	err := List(fake.ServiceClient(), instanceID).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := os.ExtractDBs(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, expectedDBs, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	os.HandleDeleteDBSuccessfully(t, instanceID, "{dbName}")

	err := os.Delete(fake.ServiceClient(), instanceID, "{dbName}").ExtractErr()
	th.AssertNoErr(t, err)
}

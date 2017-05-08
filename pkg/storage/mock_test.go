package storage

import (
	"github.com/sapcc/hermes/pkg/data"
	"runtime/debug"
	"testing"
)

func Test_MockStorage_EventDetail(t *testing.T) {
	eventDetail, error := Mock().GetEvent("some id")

	if error != nil {
		debug.PrintStack()
		t.FailNow()
	}

	if eventDetail.ID != "d5eed458-6666-58ec-ad06-8d3cf6bafca1" {
		debug.PrintStack()
		t.FailNow()
	}
	if eventDetail.Type != "identity.project.deleted" {
		debug.PrintStack()
		t.FailNow()
	}
	if eventDetail.Time != "2017-05-02T12:02:46.726056+0000" {
		debug.PrintStack()
		t.FailNow()
	}

}

func Test_MockStorage_Events(t *testing.T) {
	eventsList, total, error := Mock().GetEvents(data.Filter{})

	if error != nil {
		debug.PrintStack()
		t.FailNow()
	}

	if total != 24 {
		debug.PrintStack()
		t.FailNow()
	}

	if len(eventsList) != 3 {
		debug.PrintStack()
		t.FailNow()
	}

	if eventsList[0].Type != "identity.project.deleted" {
		debug.PrintStack()
		t.FailNow()
	}

	if eventsList[1].ID != "095056c9-4cbb-5200-af70-0977dbcf5000" {
		debug.PrintStack()
		t.FailNow()
	}

	if eventsList[2].Time != "2017-05-02T11:45:44.755215+0000" {
		debug.PrintStack()
		t.FailNow()
	}

}

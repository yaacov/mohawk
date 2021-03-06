package alerts

import (
	"testing"
	"time"

	"github.com/MohawkTSDB/mohawk/src/storage/memory"
)

func initTestEnv() (*memory.Storage, []*Alert, AlertRules) {
	// Testing with memory backend.
	b := &memory.Storage{}
	b.Open(nil)

	// creating some alerts.
	l := []*Alert{
		{
			ID:                "cpu usage too high",
			Tenant:            "_ops",
			Metrics:           []string{"cpu_usage"},
			AlertIfHigherThan: createFloatPtr(0.9),
		},
		{
			ID:               "free memory too low ",
			Tenant:           "_ops",
			Metrics:          []string{"free_memory"},
			AlertIfLowerThan: createFloatPtr(2000),
		},
		{
			ID:                "free memory in between ",
			Tenant:            "_ops",
			Metrics:           []string{"free_memory"},
			AlertIfLowerThan:  createFloatPtr(1000),
			AlertIfHigherThan: createFloatPtr(9000),
		},
		{
			ID:                "free memory in too high ",
			Tenant:            "_ops",
			Metrics:           []string{"free_memory"},
			AlertIfHigherThan: createFloatPtr(4000),
		},
	}

	// Create an alerts object with memory backend.
	alerts := AlertRules{
		Alerts:  l,
		Storage: b,
		Verbose: true,
	}
	alerts.Init()

	return b, l, alerts
}

func TestAlertsInit0(test *testing.T) {
	_, l, _ := initTestEnv()

	// check that init set types
	if l[0].Type != higherThan || l[1].Type != lowerThan || l[2].Type != outside || l[3].Type != higherThan {
		test.Error("Fail test 0")
	}
}

func TestAlertsInit1(test *testing.T) {
	b, l, alerts := initTestEnv()

	// Create some fake data
	// Firing alert 1
	t := int64(time.Now().UTC().Unix()*1000) - int64(2*60*1000)
	v := float64(1500)
	b.PostRawData("_ops", "free_memory", t, v)

	// run alerts worker in separate thread and push results to a channel:
	alerts.checkAlerts()

	// only alert 1 should fire!
	if l[0].State || !l[1].State || l[2].State || l[3].State {
		test.Error("Fail test 1")
	}
}

func TestAlertsInit2(test *testing.T) {
	b, l, alerts := initTestEnv()

	// Create some more fake data
	// firing alerts 1 and 2
	t := int64(time.Now().UTC().Unix()*1000) - int64(2*60*1000)
	v := float64(500)
	b.PostRawData("_ops", "free_memory", t, v)

	// run alerts worker in separate thread and push results to a channel:
	alerts.checkAlerts()

	// only alerts 1 and 2 should fire!
	if l[0].State || !l[1].State || !l[2].State || l[3].State {
		test.Error("Fail test 2")
	}
}

func TestAlertsInit3(test *testing.T) {
	b, l, alerts := initTestEnv()

	// Create some more fake data
	// firing none
	t := int64(time.Now().UTC().Unix()*1000) - int64(2*60*1000)
	v := float64(2500)
	b.PostRawData("_ops", "free_memory", t, v)

	// run alerts worker in separate thread and push results to a channel:
	alerts.checkAlerts()

	// no alerts should fire!
	if l[0].State || l[1].State || l[2].State || l[3].State {
		test.Error("Fail test 3")
	}
}

func TestAlertsInit4(test *testing.T) {
	b, l, alerts := initTestEnv()

	// Create some more fake data
	// firing alert 3
	t := int64(time.Now().UTC().Unix()*1000) - int64(2*60*1000)
	v := float64(5000)
	b.PostRawData("_ops", "free_memory", t, v)

	// run alerts worker in separate thread and push results to a channel:
	alerts.checkAlerts()

	// alert 3 should fire
	if l[0].State || l[1].State || l[2].State || !l[3].State {
		test.Error("Fail test 3")
	}
}

func createFloatPtr(v float64) *float64 {
	return &v
}

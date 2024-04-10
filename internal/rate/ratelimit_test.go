package rate

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
	"time"
)

func testStringGatherer(i int) string {
	time.Sleep(time.Duration(time.Millisecond * 100))
	return fmt.Sprintf("result-%d", i)
}

type thing struct {
	id  string
	age int
}

func testThingGatherer(_ int) thing {
	ulid := fmt.Sprintf("%d", rand.IntN(5000))
	return thing{id: ulid, age: rand.IntN(100)}
}

func TestLimiterStrings(t *testing.T) {
	type testCase struct {
		limit        int
		numberOfJobs int
	}
	for name, tc := range map[string]testCase{
		"limit_4_and_20_jobs":  {limit: 4, numberOfJobs: 20},
		"limit_16_and_30_jobs": {limit: 16, numberOfJobs: 30},
		"limit_6_and_6_jobs":   {limit: 16, numberOfJobs: 30},
		"limit_8_and_100_jobs": {limit: 8, numberOfJobs: 10},
		"limit_20_and_30_jobs": {limit: 20, numberOfJobs: 59},
	} {
		t.Run(name, func(t *testing.T) {

			limiter := New(tc.limit, tc.numberOfJobs, testStringGatherer)
			results := limiter.Spawn().Run()

			if len(results) != tc.numberOfJobs {
				t.Errorf("expected 20 but got %d\n", len(results))
			}
			for r := range results {
				expectedSubstring := "result-"
				if !strings.Contains(r, expectedSubstring) {
					t.Errorf("got %s but expected %s\n", r, expectedSubstring)
				}
			}
		})
	}
}

func TestLimiterThings(t *testing.T) {
	type testCase struct {
		limit        int
		numberOfJobs int
	}
	for name, tc := range map[string]testCase{
		"limit_4_and_20_jobs":  {limit: 4, numberOfJobs: 20},
		"limit_16_and_30_jobs": {limit: 16, numberOfJobs: 30},
		"limit_6_and_6_jobs":   {limit: 16, numberOfJobs: 30},
		"limit_8_and_100_jobs": {limit: 8, numberOfJobs: 10},
		"limit_20_and_30_jobs": {limit: 20, numberOfJobs: 59},
	} {
		t.Run(name, func(t *testing.T) {
			limiter := New(tc.limit, tc.numberOfJobs, testThingGatherer)
			results := limiter.Spawn().Run()

			if len(results) != tc.numberOfJobs {
				t.Errorf("expected 20 but got %d\n", len(results))
			}
			for r := range results {
				if r.id == "" {
					t.Errorf("expected non-empty ulid string")
				}
				if r.age > 100 {
					t.Errorf("expected age <= 100 but got %d\n", r.age)
				}
			}
		})
	}
}

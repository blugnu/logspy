package logspy

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Ensure there is no timestamp prefix or other formatting
	log.SetFlags(0)
	// Redirect to the LogSpy sink
	log.SetOutput(Sink())

	// Run tests and exit with test result
	os.Exit(m.Run())
}

func TestThatContainsFindsStringsAtAnyPositionInLog(t *testing.T) {
	// ARRANGE
	Reset()
	log.Println("test output")

	// ACT
	result := Strings()

	// ASSERT
	t.Run("expected number of entries", func(t *testing.T) {
		wanted := 1
		got := len(result)
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})

	t.Run("expected entries", func(t *testing.T) {
		wanted := "test output"
		got := result[0]
		if wanted != got {
			t.Errorf("\nwanted %#v\ngot    %#v", wanted, got)
		}
	})
}

func TestThatNumEntriesReturnsTheNumberOfLogEntries(t *testing.T) {
	// ARRANGE
	Reset()
	log.Println("output 1")
	log.Println("output 2")
	log.Println("output 3")

	// ACT
	result := Strings()

	// ASSERT
	wanted := 3
	got := len(result)
	if wanted != got {
		t.Errorf("Expected %d log entries, got %d", wanted, got)
	}
}

func TestThatNumObjectsReturnsTheCorrectResultForCompactJson(t *testing.T) {
	// ARRANGE
	Reset()
	log.Println("{\"entry\": 1}")
	log.Println("{\"entry\": 2}")
	log.Println("{\"entry\": 3}")

	// ACT
	result, err := JsonObjects()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// ASSERT
	wanted := 3
	got := len(result)
	if wanted != got {
		t.Errorf("Expected %d log objects, got %d", wanted, got)
	}
}

func TestThatNumObjectsReturnsTheCorrectResultForPrettyJson(t *testing.T) {
	// ARRANGE
	log.SetOutput(Sink())
	Reset()

	// ACT
	log.Println("{\n  \"entry\": 1\n}")
	log.Println("{\n  \"entry\": 2\n}")
	result, err := JsonObjects()

	// ASSERT

	// Number of raw entries (lines)
	wanted := 6
	got := len(Strings())
	if wanted != got {
		t.Errorf("Expected %d log objects, got %d", wanted, got)
	}

	// Number of formatted Json entries
	wanted = 2
	got = len(result)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if wanted != got {
		t.Errorf("Expected %d log objects, got %d", wanted, got)
	}
}

func TestThatStringReturnsTheCapturedLogEntries(t *testing.T) {
	// ARRANGE
	Reset()

	// ACT
	log.Println("output 1")
	log.Println("output 2")
	log.Println("output 3")

	// ASSERT
	wanted := "output 1\noutput 2\noutput 3\n"
	got := Content()
	if wanted != got {
		t.Errorf("Expected the log to contain %q, got %q", wanted, got)
	}
}

func Test_JsonObjectsReturnsObjectsUpToNonJson(t *testing.T) {
	// ARRANGE
	Reset()

	// ACT
	log.Println("{\"field\": \"value\"}")
	log.Println("{\"field\": \"value\"}")
	log.Println("THIS IS NOT JSON")

	// ASSERT
	n, err := JsonObjects()
	t.Run("returns the number of valid Json objects", func(t *testing.T) {
		wanted := 2
		got := len(n)
		if wanted != 2 {
			t.Errorf("wanted %d, got %d", wanted, got)
		}
	})
	t.Run("returns error", func(t *testing.T) {
		if err == nil {
			t.Error("wanted an error, got nil")
		}
		wanted := "invalid character 'T' looking for beginning of value"
		got := err.Error()
		if wanted != got {
			t.Errorf("wanted %q, got %q", wanted, got)

		}
	})
}

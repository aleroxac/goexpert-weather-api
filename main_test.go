package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ---------- tests
func TestCEPValidation(t *testing.T) {
	t.Run("Invalid CEP", func(t *testing.T) {
		cep := &CEP{
			CEP: "1234567",
		}
		assert.False(t, cep.Validate())
	})

	t.Run("Valid CEP", func(t *testing.T) {
		cep := &CEP{
			CEP: "12345678",
		}
		assert.True(t, cep.Validate())
	})
}

func TestAPIResponseCodes(t *testing.T) {
	go main()

	t.Run("status_200", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8080/weather?cep=13330250", nil)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, http.StatusOK)

		res_bytes, _ := io.ReadAll(resp.Body)
		var weather Weather
		err = json.Unmarshal(res_bytes, &weather)
		assert.NoError(t, err)
	})

	t.Run("status_422", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8080/weather?cep=1333025", nil)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, http.StatusUnprocessableEntity)

		res_bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		res_string := strings.TrimSpace(string(res_bytes))
		assert.Equal(t, res_string, "invalid zipcode")
	})

	t.Run("status_404", func(t *testing.T) {
		req, err := http.NewRequest("GET", "http://localhost:8080/weather?cep=12345678", nil)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		assert.Equal(t, resp.StatusCode, http.StatusNotFound)

		res_bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		res_string := strings.TrimSpace(string(res_bytes))
		assert.Equal(t, res_string, "can not find zipcode")
	})
}

// ---------- benchmarks
func BenchmarkCEPValidation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cep := &CEP{
			CEP: "1234567",
		}
		cep.Validate()
	}
}

func BenchmarkAPIResponseCodes(b *testing.B) {
	go main()

	for i := 0; i < b.N; i++ {
		b.Run("status_200", func(b *testing.B) {
			req, err := http.NewRequest("GET", "http://localhost:8080/weather?cep=13330250", nil)
			if err != nil {
				b.Fatal(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			defer resp.Body.Close()
		})

		b.Run("status_422", func(b *testing.B) {
			req, err := http.NewRequest("GET", "http://localhost:8080/weather?cep=1333025", nil)
			if err != nil {
				b.Fatal(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			defer resp.Body.Close()
		})

		b.Run("status_404", func(b *testing.B) {
			req, err := http.NewRequest("GET", "http://localhost:8080/weather?cep=12345678", nil)
			if err != nil {
				b.Fatal(err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				b.Fatal(err)
			}
			defer resp.Body.Close()
		})
	}
}

// ---------- fuzzing
func FuzzCEPValidation(f *testing.F) {
	invalid_seed := []string{
		"12345678",
		"00000000",
		"00000001",
	}
	for _, cep := range invalid_seed {
		f.Add(cep)
	}
	f.Fuzz(func(t *testing.T, cep string) {
		c := &CEP{
			CEP: cep,
		}
		assert.True(t, c.Validate())
	})
}

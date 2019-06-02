package mailer

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/adrianpk/poslan/internal/config"
)

// Unit tests only: go test -v -short
// Integration test only: go test -run Integration

var (
	protocol       = "http"
	host           = "localhost"
	port           = 8080
	serverInstance *httptest.Server
	signinURL      string
	signoutURL     string
	sendURL        string
	user1          = "5958b185-8150-4aae-b53f-0c44771ddec5"
	user2          = "3c05e701-b495-4443-b454-2c37e2ecccdf"
)

func init() {
	signinURL = fmt.Sprintf("%s://%s:%d/signin", protocol, host, port)
	signoutURL = fmt.Sprintf("%s://%s:%d/signout", protocol, host, port)
	sendURL = fmt.Sprintf("%s://%s:%d/send", protocol, host, port)
}

func TestMain(m *testing.M) {
	setup()
	e := m.Run()
	teardown()
	os.Exit(e)
}

// TestSomething is a base unit test reference.
func TestSomething(t *testing.T) {
	t.Parallel()
	t.Skip("Skiping unit test at the moment.")
}

// TestSendIntegration is a base unit test reference.
func TestSendIntegration(t *testing.T) {
	// handlers := &MyHandler{}
	// server := httptest.NewServer(handlers)
	// defer server.Close()

	if testing.Short() {
		t.Skip("Skiping integration test.")
	}

	t.Log("TestSignup started.")
	emailJSON := `
	{
		"data": {
			"to": "sendmailtest@sharklasers.com",
			"cc": "sendmailtest@sharklasers.com",
			"bcc": "sendmailtest@sharklasers.com",
			"subject": "Subject",
			"body": "Body text."
		}
	}
	`
	reader := strings.NewReader(emailJSON)
	req, _ := http.NewRequest("POST", sendURL, reader)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Logf("[ERROR] TestSendIntegration error: %s", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected: 200 - StatusOK | Received: Status: %d - %s", res.StatusCode, http.StatusText(res.StatusCode))
	}
}

func setup() {
	cfg, err := configForTest()
	if err != nil {
		log.Println("[ERROR] Cannot run test: Invalid configuraration.")
		os.Exit(1)
	}

	errchan := make(chan error)

	go TestingRun(cfg, errchan)

	select {
	case err = <-errchan:
		msg := fmt.Sprintf("test setup cannot be completed completed: %s", err.Error())
		log.Printf("[ERROR] %s", msg)
		os.Exit(1)

	default:
		log.Println("Setup completed.")
	}
}

func teardown() {
	log.Println("Teardown completed")
}

// TODO: Implement a custom leader only used for test.
// Something like the one already implemented in
// internal/config/loadFromFile(filePath string)
// Although it is probably easy to edit those values stright here
func configForTest() (*config.Config, error) {
	// Uncomment and edit if you need to add custom values.
	// os.Setenv("KEY1", "VAL1")
	// os.Setenv("KEY2", "VAL2")
	// os.Setenv("KEY2", "VAL2")

	// App
	app := config.AppConfig{
		ServerPort: port,
		LogLevel:   config.LogLevel.Debug,
	}

	provider1 := config.ProviderConfig{
		Name:     "amazon",
		Type:     "amazon-ses",
		Enabled:  true,
		Priority: 1,
		IDKey:    config.GetEnvOrDef("PROVIDER_ID_KEY_1", ""),
		APIKey:   config.GetEnvOrDef("PROVIDER_API_KEY_1", ""),
	}

	provider2 := config.ProviderConfig{
		Name:     "sendgrid",
		Type:     "sendgrid",
		Enabled:  true,
		Priority: 1,
		IDKey:    config.GetEnvOrDef("PROVIDER_ID_KEY_2", ""),
		APIKey:   config.GetEnvOrDef("PROVIDER_API_KEY_2", ""),
	}

	mailers := config.MailerConfig{
		Providers: []config.ProviderConfig{
			provider1,
			provider2,
		},
	}

	cfg := &config.Config{
		App:    app,
		Mailer: mailers,
	}

	// fmt.Printf("\n%+v\n", cfg)

	return cfg, nil
}

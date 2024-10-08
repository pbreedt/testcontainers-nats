package sub

import (
	"context"
	"fmt"
	"github/pbreedt/testcontainers/nats/pub"
	"github/pbreedt/testcontainers/nats/sub"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type natsContainer struct {
	testcontainers.Container
	URL string
}

func runNatsContainer(ctx context.Context) (*natsContainer, error) {

	// container build request
	req := testcontainers.ContainerRequest{
		Image:        "nats:latest",
		ExposedPorts: []string{"4222/tcp", "8222/tcp"},
		WaitingFor:   wait.ForLog("Listening for client connections"),
	}

	// building a generic container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	// mapped port
	mappedPort, err := container.MappedPort(ctx, "4222")
	if err != nil {
		return nil, err
	}

	// getting the container ip
	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	// generating the nats cluster uri
	uri := fmt.Sprintf("nats://%s:%s", hostIP, mappedPort.Port())

	return &natsContainer{
		Container: container,
		URL:       uri,
	}, nil
}

func TestSubscriber(t *testing.T) {
	ctx := context.Background()

	expect := 10

	// get common NATS container
	container, err := runNatsContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// start subscriber
	rCount := 0
	rCh := make(chan int)
	go sub.GoSubscribe(container.URL, expect, rCh)

	// start publisher
	pCount := 0
	pCh := make(chan int)
	go pub.GoPublish(container.URL, expect, pCh)

	// read 2 results (we don't care about the order)
	for i := 0; i <= 1; i++ {
		select {
		case pCount = <-pCh:
			t.Logf("published %d messages", pCount)
		case rCount = <-rCh:
			t.Logf("received %d messages", rCount)
		}
	}

	if rCount < expect {
		t.Errorf("received %d messages, expected %d", rCount, expect)
	}
	if pCount < expect {
		t.Errorf("published %d messages, expected %d", pCount, expect)
	}

	// cleaning container after test is complete
	t.Cleanup(func() {
		t.Log("terminating container")

		if er := container.Terminate(ctx); er != nil {
			t.Errorf("failed to terminate container: :%v", er)
		}
	})
}

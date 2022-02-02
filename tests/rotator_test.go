package integrationtests

import (
	"banners-rotator/internal/rmq"
	gw "banners-rotator/internal/server/bannersrotatorpb"
	"context"
	"os"
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRotator_RMQ(t *testing.T) {
	t.Run("publish event", func(t *testing.T) {
		conn, err := amqp.Dial(getRMQConnectionString())
		require.NoError(t, err)

		p := rmq.NewRMQProducer("rotator-test", conn)
		err = p.Connect()
		require.NoError(t, err)

		err = p.Publish(rmq.QMessage{
			Type:     "view",
			SlotID:   1,
			BannerID: 1,
			GroupID:  1,
			Date:     time.Now().Unix(),
		})
		require.NoError(t, err)
	})
}

func TestRotator_CreateSlot(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	t.Run("create slot", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := client.CreateSlot(ctx, &gw.Slot{Description: ""})
		require.Error(t, err)

		slot, err := client.CreateSlot(ctx, &gw.Slot{Description: "test desc"})
		require.NoError(t, err)
		require.Equal(t, "test desc", slot.Description)
	})
}

func TestRotator_CreateBanner(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	t.Run("create banner", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := client.CreateBanner(ctx, &gw.Banner{Description: ""})
		require.Error(t, err)

		banner, err := client.CreateBanner(ctx, &gw.Banner{Description: "test desc"})
		require.NoError(t, err)
		require.Equal(t, "test desc", banner.Description)
	})
}

func TestRotator_CreateGroup(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	t.Run("create group", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := client.CreateGroup(ctx, &gw.Group{Description: ""})
		require.Error(t, err)

		group, err := client.CreateGroup(ctx, &gw.Group{Description: "test desc"})
		require.NoError(t, err)
		require.Equal(t, "test desc", group.Description)
	})
}

func TestRotator_CreateRotation(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	t.Run("create rotation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := client.CreateRotation(ctx, &gw.Rotation{BannerId: -1, SlotId: -1})
		require.Error(t, err)

		slot, err := client.CreateSlot(ctx, &gw.Slot{Description: "test desc"})
		require.NoError(t, err)
		banner, err := client.CreateBanner(ctx, &gw.Banner{Description: "test desc"})
		require.NoError(t, err)

		msg, err := client.CreateRotation(ctx, &gw.Rotation{BannerId: banner.Id, SlotId: slot.Id})
		require.NoError(t, err)
		require.Equal(t, "Rotation was created", msg.Message)
	})
}

func TestRotator_DeleteRotation(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	t.Run("delete rotation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := client.DeleteRotation(ctx, &gw.Rotation{BannerId: -1, SlotId: -1})
		require.Error(t, err)

		slot, err := client.CreateSlot(ctx, &gw.Slot{Description: "test desc"})
		require.NoError(t, err)
		banner, err := client.CreateBanner(ctx, &gw.Banner{Description: "test desc"})
		require.NoError(t, err)
		_, err = client.CreateRotation(ctx, &gw.Rotation{BannerId: banner.Id, SlotId: slot.Id})
		require.NoError(t, err)

		msg, err := client.DeleteRotation(ctx, &gw.Rotation{BannerId: banner.Id, SlotId: slot.Id})
		require.NoError(t, err)
		require.Equal(t, "Rotation was deleted", msg.Message)

	})
}

func TestRotator_CreateClickEvent(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	t.Run("create click event", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		slot, err := client.CreateSlot(ctx, &gw.Slot{Description: "test desc"})
		require.NoError(t, err)
		banner, err := client.CreateBanner(ctx, &gw.Banner{Description: "test desc"})
		require.NoError(t, err)
		group, err := client.CreateGroup(ctx, &gw.Group{Description: "test desc"})
		require.NoError(t, err)

		_, err = client.CreateClickEvent(ctx, &gw.ClickEvent{BannerId: banner.Id, SlotId: slot.Id, GroupId: -1})
		require.Error(t, err)

		msg, err := client.CreateClickEvent(ctx, &gw.ClickEvent{BannerId: banner.Id, SlotId: slot.Id, GroupId: group.Id})
		require.NoError(t, err)
		require.Equal(t, "Click event was registered", msg.Message)
	})
}

func TestRotator_BannerForSlot(t *testing.T) {
	client, err := getClient()
	require.NoError(t, err)

	t.Run("create click event", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		slot, err := client.CreateSlot(ctx, &gw.Slot{Description: "test desc"})
		require.NoError(t, err)
		banner, err := client.CreateBanner(ctx, &gw.Banner{Description: "test desc"})
		require.NoError(t, err)
		group, err := client.CreateGroup(ctx, &gw.Group{Description: "test desc"})
		require.NoError(t, err)

		_, err = client.BannerForSlot(ctx, &gw.SlotRequest{SlotId: slot.Id, GroupId: -1})
		require.Error(t, err)

		_, err = client.BannerForSlot(ctx, &gw.SlotRequest{SlotId: slot.Id, GroupId: group.Id})
		require.Error(t, err)

		_, err = client.CreateRotation(ctx, &gw.Rotation{BannerId: banner.Id, SlotId: slot.Id})
		require.NoError(t, err)

		result, err := client.BannerForSlot(ctx, &gw.SlotRequest{SlotId: slot.Id, GroupId: group.Id})
		require.NoError(t, err)
		require.Equal(t, banner.Id, result.Id)
	})
}

func getRotatorConnectionString() string {
	const connectionString = "localhost:8080"
	cs := os.Getenv("TESTS_ROTATOR_DSN")
	if cs == "" {
		return connectionString
	}

	return cs
}

func getClient() (gw.BannersRotatorClient, error) {
	conn, err := grpc.Dial(
		getRotatorConnectionString(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := gw.NewBannersRotatorClient(conn)

	return client, nil
}

func getRMQConnectionString() string {
	const connectionString = "amqp://guest:guest@localhost:5672/"
	cs := os.Getenv("TESTS_RMQ_DSN")
	if cs == "" {
		return connectionString
	}

	return cs
}

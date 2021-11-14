package command

import (
	"testing"

	"github.com/BKrajancic/boby/m/v2/src/service"
	"github.com/BKrajancic/boby/m/v2/src/service/demoservice"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func TestRenderText(t *testing.T) {
	font, _ := truetype.Parse(goregular.TTF)
	face := truetype.NewFace(font, &truetype.Options{
		Size: 20,
		DPI:  72,
	})

	renderText(face, "hello world")
}

func DontTestMsgSend(t *testing.T) {
	demoSender := demoservice.DemoSender{
		ServiceID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	RenderText(
		testConversation,
		testSender,
		[]interface{}{"hi"},
		nil,
		demoSender.SendMessage,
	)
}

func TestMsgEmpty(t *testing.T) {
	demoSender := demoservice.DemoSender{
		ServiceID: "0",
	}
	testSender := service.User{Name: "Test_User", ServiceID: demoSender.ID()}
	testConversation := service.Conversation{
		ServiceID:      demoSender.ID(),
		ConversationID: "0",
	}

	RenderText(
		testConversation,
		testSender,
		[]interface{}{},
		nil,
		demoSender.SendMessage,
	)
}

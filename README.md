# golang-slack-webhook
Go Lang library to send messages to Slack via Incoming Webhooks.
efficiently converts data to JSON format.

## usage

package main

import "github.com/shubhambhuyan/golang-slack-webhook"
import "fmt"

func main() {

    //slack incoming webhook url
    webhookUrl := "https://hooks.slack.com/services/xxxxxxx/xxxxxxx"
    
    //message attachment
    attachment1 := slack.Attachment {
      Text:"Fly to github!!!",
      Fallback:"You were unable to choose an option",
      CallbackID:"wopr_game",Color:"#3AA3E3",
      Footer:"shubham",
      FooterIcon:"https://cdn.pixabay.com/photo/2014/03/24/13/47/tomcat-294361_960_720.png",
      Ts:778339990,
      Actions: []slack.AttachmentAction{
        {
            Name: "click",
            Text: "Click",
            Type: "button",
            URL: "http://https://github.com/shubhambhuyan/golang-slack-webhook",
            Style: "primary",
          },
          {
                Name: "cancel",
                Text: "Cancel",
                Type: "button",
                Style: "danger",
                Confirm: []slack.ConfirmationField{
                  {
                    Title: "Are you sure?",
                    Text: "prefer to cancel ",
                    OkText: "Yes",
                    DismissText: "No",
                  },
                },
          },
       },
    }

    //message to be sent is called payload
    payload := slack.Payload {
      Text: "visit github?",
      Attachments: []slack.Attachment{attachment1},
    }
    
    //send message via webhook
    err := slack.Send(webhookUrl, "", payload)
    if len(err) > 0 {
      fmt.Printf("error: %s\n", err)
    }
    
}

package zendesk

// https://developer.zendesk.com/rest_api/docs/support/triggers#via-types

type ViaType int

const (
	// ViaWebForm : Web form
	ViaWebForm ViaType = 0
	// ViaMail : Email
	ViaMail ViaType = 4
	// ViaChat : Chat
	ViaChat ViaType = 29
	// ViaTwitter : Twitter
	ViaTwitter ViaType = 30
	// ViaTwitterDM : Twitter DM
	ViaTwitterDM ViaType = 26
	// ViaTwitterFavorite : Twitter like
	ViaTwitterFavorite ViaType = 23
	// ViaVoicemail : Voicemail
	ViaVoicemail ViaType = 33
	// ViaPhoneCallInbound : Phone call (incoming)
	ViaPhoneCallInbound ViaType = 34
	// ViaPhoneCallOutbound : Phone call (outbound)
	ViaPhoneCallOutbound ViaType = 35
	// ViaAPIVoicemail : CTI voicemail
	ViaAPIVoicemail ViaType = 44
	// ViaAPIPhoneCallInbound : CTI phone call (inbound)
	ViaAPIPhoneCallInbound ViaType = 45
	// ViaAPIPhoneCallOutbound : CTI phone call (outbound)
	ViaAPIPhoneCallOutbound ViaType = 46
	// ViaSMS : SMS
	ViaSMS ViaType = 57
	// ViaGetSatisfaction : Get Satisfaction
	ViaGetSatisfaction ViaType = 16
	// ViaWebWidget : Web Widget
	ViaWebWidget ViaType = 48
	// ViaMobileSDK : Mobile SDK
	ViaMobileSDK ViaType = 49
	// ViaMobile : Mobile
	ViaMobile ViaType = 56
	// ViaHelpCenter : Help Center post
	ViaHelpCenter ViaType = 50
	// ViaWebService : Web service (API)
	ViaWebService ViaType = 5
	// ViaRule : Trigger, automation
	ViaRule ViaType = 8
	// ViaClosedTicket : Closed ticket
	ViaClosedTicket ViaType = 27
	// ViaTicketSharing : Ticket Sharing
	ViaTicketSharing ViaType = 31
	// ViaFacebookPost : Facebook post
	ViaFacebookPost ViaType = 38
	// ViaFacebookMessage : Facebook private message
	ViaFacebookMessage ViaType = 41
	// ViaSatisfactionPrediction : Satisfaction prediction
	ViaSatisfactionPrediction ViaType = 54
	// ViaAnyChannel : Channel framework
	ViaAnyChannel ViaType = 55
)

var viaTypeMapText = map[ViaType]string{
	ViaWebForm:                "web_form",
	ViaMail:                   "mail",
	ViaChat:                   "chat",
	ViaTwitter:                "twitter",
	ViaTwitterDM:              "twitter_dm",
	ViaTwitterFavorite:        "twitter_favorite",
	ViaVoicemail:              "voicemail",
	ViaPhoneCallInbound:       "phone_call_inbound",
	ViaPhoneCallOutbound:      "phone_call_outbound",
	ViaAPIVoicemail:           "api_voicemail",
	ViaAPIPhoneCallInbound:    "api_phone_call_inbound",
	ViaAPIPhoneCallOutbound:   "api_phone_call_outbound",
	ViaSMS:                    "sms",
	ViaGetSatisfaction:        "get_satisfaction",
	ViaWebWidget:              "web_widget",
	ViaMobileSDK:              "mobile_sdk",
	ViaMobile:                 "mobile",
	ViaHelpCenter:             "helpcenter",
	ViaWebService:             "web_service",
	ViaRule:                   "rule",
	ViaClosedTicket:           "closed_ticket",
	ViaTicketSharing:          "ticket_sharing",
	ViaFacebookPost:           "facebook_post",
	ViaFacebookMessage:        "facebook_message",
	ViaSatisfactionPrediction: "satisfaction_prediction",
	ViaAnyChannel:             "any_channel",
}

var ViaTypeMap = map[string]ViaType{
	"web_form":                ViaWebForm,
	"mail":                    ViaMail,
	"chat":                    ViaChat,
	"twitter":                 ViaTwitter,
	"twitter_dm":              ViaTwitterDM,
	"twitter_favorite":        ViaTwitterFavorite,
	"voicemail":               ViaVoicemail,
	"phone_call_inbound":      ViaPhoneCallInbound,
	"phone_call_outbound":     ViaPhoneCallOutbound,
	"api_voicemail":           ViaAPIVoicemail,
	"api_phone_call_inbound":  ViaAPIPhoneCallInbound,
	"api_phone_call_outbound": ViaAPIPhoneCallOutbound,
	"sms":                     ViaSMS,
	"get_satisfaction":        ViaGetSatisfaction,
	"web_widget":              ViaWebWidget,
	"mobile_sdk":              ViaMobileSDK,
	"mobile":                  ViaMobile,
	"helpcenter":              ViaHelpCenter,
	"web_service":             ViaWebService,
	"rule":                    ViaRule,
	"closed_ticket":           ViaClosedTicket,
	"ticket_sharing":          ViaTicketSharing,
	"facebook_post":           ViaFacebookPost,
	"facebook_message":        ViaFacebookMessage,
	"satisfaction_prediction": ViaSatisfactionPrediction,
	"any_channel":             ViaAnyChannel,
}

// ViaTypeText takes via_id and returns via_type
func ViaTypeText(viaID ViaType) string {
	return viaTypeMapText[viaID]
}

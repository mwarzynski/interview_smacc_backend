package mail

import (
	"testing"
)

func TestMessageValidate(t *testing.T) {
	type fields struct {
		From    string
		To      string
		Subject string
		Text    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "empty fields",
			wantErr: true,
		},
		{
			name: "valid fields",
			fields: fields{
				From:    "a@smacc.com",
				To:      "b@smacc.com",
				Subject: "it's a valid subject",
				Text:    "some text",
			},
			wantErr: false,
		},
		{
			name: "invalid from",
			fields: fields{
				From:    "a",
				To:      "b@smacc.com",
				Subject: "from is invalid",
			},
			wantErr: true,
		},
		{
			name: "subject is required",
			fields: fields{
				From: "a@smacc.com",
				To:   "b@smacc.com",
			},
			wantErr: true,
		},
		{
			name: "text field is not required",
			fields: fields{
				From:    "a@smacc.com",
				To:      "b@smacc.com",
				Subject: "it's a valid subject",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				From:    tt.fields.From,
				To:      tt.fields.To,
				Subject: tt.fields.Subject,
				Text:    tt.fields.Text,
			}
			if err := m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Message.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

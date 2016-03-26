package steam

import (
	"encoding/json"
	"errors"
	"net/url"
)

var (
	ErrInvalidPhoneNumber = errors.New("invalid phone number specified")
)

type PhoneAPIResponse struct {
	State     string `json:"state"`
	ErrorText string `json:"errorText"`
}

func (community *Community) CopyCookiesToSteamStore() {
	commu, _ := url.Parse("https://steamcommunity.com")
	store, _ := url.Parse("https://store.steampowered.com")

	community.client.Jar.SetCookies(store, community.client.Jar.Cookies(commu))
}

func (community *Community) ValidatePhoneNumber(number string) error {
	resp, err := community.client.Get("https://store.steampowered.com/phone/validate?phoneNumber=" + url.QueryEscape(number))
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	type Response struct {
		Success bool `json:"success"`
	}

	var response Response
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if !response.Success {
		return ErrInvalidPhoneNumber
	}

	return nil
}

func (community *Community) AddPhoneNumber(number string) error {
	resp, err := community.client.Get("https://store.steampowered.com/phone/add_ajaxop?" + url.Values{
		"op":        {"get_phone_number"},
		"input":     {number},
		"sessionID": {community.sessionID},
		"confirmed": {"0"},
	}.Encode())
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	var response PhoneAPIResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response.State != "get_sms_code" {
		return errors.New(response.ErrorText)
	}

	return nil
}

func (community *Community) VerifyPhoneNumber(code string) error {
	resp, err := community.client.Get("https://store.steampowered.com/phone/add_ajaxop?" + url.Values{
		"op":        {"get_sms_code"},
		"input":     {code},
		"sessionID": {community.sessionID},
		"confirmed": {"0"},
	}.Encode())
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	var response PhoneAPIResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	if response.State != "done" {
		return errors.New(response.ErrorText)
	}

	return nil
}
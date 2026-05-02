package fbd

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

const signInURL string = "https://www.football-bet-data.com/signin/"
const sessionCookieName string = "ASP.NET_SessionId"

func (c *Client) SignIn(email, password string) error {
	if err := c.ensureReady(); err != nil {
		return err
	}

	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return errors.New("email and password are required")
	}

	slog.Info("Signing in to FBD")

	fields, err := c.getWebFormsFields(signInURL)
	if err != nil {
		return err
	}

	fields.Set("ctl00$ContentPlaceHolder2$unameTextBox", email)
	fields.Set("ctl00$ContentPlaceHolder2$pwordTextBox", password)
	fields.Set("ctl00$ContentPlaceHolder2$submitButton", "Submit")

	if _, err := c.httpClient.PostForm(signInURL, fields); err != nil {
		return fmt.Errorf("post sign-in form: %w", err)
	}

	hasSessionCookie := false
	for _, cookie := range c.httpClient.Jar.Cookies(c.baseURL) {
		if cookie.Name == sessionCookieName && cookie.Value != "" {
			hasSessionCookie = true
			break
		}
	}
	if !hasSessionCookie {
		return fmt.Errorf("signin failed: %s cookie missing or empty, probably invalid credentials", sessionCookieName)
	}
	c.enforceSession = true

	return nil
}

func (c *Client) SignOut() error {
	if err := c.ensureReady(); err != nil {
		return err
	}

	fields, err := c.getWebFormsFields(homeURL)
	if err != nil {
		return err
	}

	slog.Info("Signing out of FBD")

	fields.Set("ctl00$logoutButton2", "Sign Out")

	if _, err := c.httpClient.PostForm(homeURL, fields); err != nil {
		return fmt.Errorf("post sign-out form: %w", err)
	}

	c.enforceSession = false

	return nil
}

package resolver

import (
	"context"
	"errors"

	"github.com/aripalo/vegas-credentials/internal/multinput"
	"github.com/aripalo/vegas-credentials/internal/prompt"
)

var guiPrompt = prompt.Dialog

func GUI(ctx context.Context) (*multinput.Result, error) {
	value, err := guiPrompt(ctx, "Multifactor Authentication", "Enter TOTP MFA Token Code:")
	return &multinput.Result{Value: value, ResolverID: ResolverIdGuiDialog}, err
}

func ConfigureGUI(enabled bool) multinput.InputResolver {
	if enabled {
		return GUI
	}
	// TODO fix this
	// To avoid nil pointer reference, return just a resolver that resolves
	// into an emtpy value with an error
	return func(ctx context.Context) (*multinput.Result, error) {
		return &multinput.Result{Value: "", ResolverID: ResolverIdGuiDialog}, errors.New("gui disabled")
	}
}

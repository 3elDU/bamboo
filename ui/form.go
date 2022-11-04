/*
Interactive form with multiple input propmts.
Useful for different dialogs
*/

package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type FormPrompt struct {
	Title string

	input InputView
}

// a simple wrapper around basic components, to simplify creation of prompts
type formComponent struct {
	baseView

	formData chan []string

	view View
	// slice with references to the input components themselves
	prompts       []FormPrompt
	composeButton ButtonView
}

// When compose button is pressed, form data will be sent through `formData` channel,
// in the proper order
func Form(submitButtonTitle string, formData chan []string, prompts ...FormPrompt) *formComponent {
	form := &formComponent{
		baseView: newBaseView(),
		formData: formData,
		prompts:  make([]FormPrompt, len(prompts)),
	}

	// assemble a view with the prompts
	promptsStack := Stack(StackOptions{Direction: VerticalStack, Spacing: 2.0})

	// assemble a view for each prompt
	for i, prompt := range prompts {
		// the first input in the form will be focused
		var focus bool
		if i == 0 {
			focus = true
		}

		// keep a reference to the input field, we'll need it
		inp := Input(func(s string) { log.Printf("%s - %s", prompt.Title, s) }, ebiten.KeyEnter, focus)

		// create a view for the prompt
		promptView := Stack(StackOptions{Direction: VerticalStack, Spacing: 1.0},
			Label(DefaultLabelOptions(), prompt.Title), inp,
		)

		// add it to the stack
		promptsStack.AddChild(promptView)

		form.prompts[i] = FormPrompt{
			Title: prompt.Title,
			input: inp,
		}
	}

	composeButton := Button(
		// Leave the handler empty, because we'll check for the button press manually
		func() {},
		Label(DefaultLabelOptions(), submitButtonTitle),
	)
	promptsStack.AddChild(composeButton)

	form.composeButton = composeButton
	form.view = promptsStack
	return form
}

func (f *formComponent) MaxSize() (float64, float64) {
	return f.view.MaxSize()
}
func (f *formComponent) ComputedSize() (float64, float64) {
	return f.view.ComputedSize()
}
func (f *formComponent) CapacityForChild(_ View) (float64, float64) {
	return 0, 0
}
func (f *formComponent) Children() []View {
	return []View{f.view}
}
func (f *formComponent) Update() error {
	if f.composeButton.IsPressed() {
		log.Println("formComponent.Update() - got button press")
		inputs := make([]string, len(f.prompts))
		for i, prompt := range f.prompts {
			inputs[i] = prompt.input.Input()
		}
		log.Printf("formComponent.Update() - [%v] [%v]", len(inputs), inputs)

		if len(inputs) != len(f.prompts) {
			log.Panicf("formComponent.Update() - len(inputs) != len(prompts)")
		}

		// make sure to send data in the separate goroutine, so we won't get a deadlock
		go func() {
			f.formData <- inputs
		}()
	}

	return f.view.Update()
}
func (f *formComponent) Draw(screen *ebiten.Image, x, y float64) error {
	return f.view.Draw(screen, x, y)
}

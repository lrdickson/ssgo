package views

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"

	"github.com/lrdickson/ssgo/internal/viewmodels"
)

type variableInfo struct {
	code   *string
	name   binding.String
	output binding.String
}

func checkErrFatal(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func getVariable(variables binding.UntypedList, id widget.ListItemID) variableInfo {
	variablesInterface, err := variables.Get()
	checkErrFatal("Failed to get variable interface array:", err)
	return variablesInterface[id].(variableInfo)
}

func NewMainView() *fyne.Container {

	// Create the editor
	variableEditor := widget.NewMultiLineEntry()
	variableEditor.SetPlaceHolder("Enter text...")

	// Display the output
	variables := binding.NewUntypedList()
	variableList := widget.NewListWithData(
		variables,
		func() fyne.CanvasObject {
			name := container.NewBorder(
				// The width of the variable pane can be controlled by the length of this label
				widget.NewLabel("AAAAAAAAAAAAAAAAAAAAAAA"),
				nil, nil, nil,
				widget.NewEntry())
			output := widget.NewLabel("Output")
			return container.NewBorder(name, nil, nil, nil, output)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			// Get the variable
			v, err := item.(binding.Untyped).Get()
			checkErrFatal("Failed to get variable data:", err)
			variable := v.(variableInfo)

			// Set the output
			output := obj.(*fyne.Container).Objects[0].(*widget.Label)
			output.Bind(variable.output)
			output.Refresh()

			// Set the name
			name := obj.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Label)
			name.Bind(variable.name)
			nameEntry := obj.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Entry)
			nameEntry.Bind(variable.name)
			name.Refresh()
		})

	// Create a new variable
	newVariableButton := widget.NewButton("New", func() {
		code := ""
		name := binding.NewString()
		name.Set("NewVariable")
		output := binding.NewString()
		output.Set("")
		newVariable := variableInfo{&code, name, output}
		variables.Append(newVariable)
		variableList.Refresh()
	})

	// Edit the code of the selected variable
	var selectedVariableId widget.ListItemID
	isVariableSelected := false
	variableList.OnSelected = func(id widget.ListItemID) {
		// Assign the code to the editor
		selectedVariableCode := binding.BindString(getVariable(variables, id).code)
		variableEditor.Bind(selectedVariableCode)

		// Set the selected variable
		selectedVariableId = id
		isVariableSelected = true
	}

	// Run variable code button
	mainViewModel := viewmodels.NewMainViewModel()
	runButton := widget.NewButton("Run", func() {
		if !isVariableSelected {
			return
		}
		variable := getVariable(variables, selectedVariableId)
		mainViewModel.EditorCode = *variable.code
		variable.output.Set(mainViewModel.RunCode())
	})

	// Put everything together
	content := container.NewBorder(
		nil, nil,
		container.NewBorder(nil, newVariableButton, nil, nil, variableList),
		nil,
		container.NewBorder(nil, runButton, nil, nil, variableEditor))

	return content
}

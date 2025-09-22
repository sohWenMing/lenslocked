package views

type inputHTMLAttribs struct {
	Name         string
	Id           string
	InputType    string
	PlaceHolder  string
	AutoComplete string
	LabelText    string
	IsRequired   bool
	Value        string
	IsDisabled   bool
}

func (i *inputHTMLAttribs) SetName(input string) {
	i.Name = input
}
func (i *inputHTMLAttribs) SetId(input string) {
	i.Id = input
}
func (i *inputHTMLAttribs) SetInputType(input string) {
	i.InputType = input
}
func (i *inputHTMLAttribs) SetPlaceHolder(input string) {
	i.PlaceHolder = input
}
func (i *inputHTMLAttribs) SetAutoComplete(input string) {
	i.AutoComplete = input
}
func (i *inputHTMLAttribs) SetLabelText(input string) {
	i.LabelText = input
}
func (i *inputHTMLAttribs) SetIsRequired() {
	i.IsRequired = true
}
func (i *inputHTMLAttribs) SetValue(input string) {
	i.Value = input
}
func (i *inputHTMLAttribs) SetIsDisabled() {
	i.IsDisabled = true
}

type checkBoxHTMLAttribs struct {
	Id        string
	Name      string
	LabelText string
	IsChecked bool
}

type buttonAttribs struct {
	ButtonText string
}

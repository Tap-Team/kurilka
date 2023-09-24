package message

type Button []Action

func NewButton(actions ...Action) Button {
	return Button(actions)
}

type Keyboard struct {
	Inline  bool     `json:"inline"`
	Buttons []Button `json:"buttons"`
}

type KeyboardBuilder struct {
	keyboard Keyboard
}

func NewKeyboardBuilder() *KeyboardBuilder {
	return &KeyboardBuilder{}
}

func (k *KeyboardBuilder) Build() Keyboard {
	return k.keyboard
}

func (k *KeyboardBuilder) SetInline(inline bool) *KeyboardBuilder {
	k.keyboard.Inline = true
	return k
}

func (k *KeyboardBuilder) AddButtons(b ...Button) *KeyboardBuilder {
	k.keyboard.Buttons = append(k.keyboard.Buttons, b...)
	return k
}

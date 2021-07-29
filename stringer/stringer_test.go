package stringer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReverse(t *testing.T) {
	assert.Equal(t, Reverse(""), "")
	assert.Equal(t, Reverse("X"), "X")
	assert.Equal(t, Reverse("😎⚽"), "⚽😎")
	assert.Equal(t, Reverse("This `\xc5` is an invalid UTF8 character"), "retcarahc 8FTU dilavni na si `�` sihT")
	assert.Equal(t, Reverse("The quick bròwn 狐 jumped over the lazy 犬"), "犬 yzal eht revo depmuj 狐 nwòrb kciuq ehT")
	assert.Equal(t, Reverse("رائد شوملي"), "يلموش دئار")
}

func TestReverse2(t *testing.T) {
	assert.Equal(t, Reverse2(""), "")
	assert.Equal(t, Reverse2("X"), "X")
	assert.Equal(t, Reverse2("b\u0301"), "b\u0301")
	assert.Equal(t, Reverse2("😎⚽"), "⚽😎")
	assert.Equal(t, Reverse2("Les Mise\u0301rables"), "selbare\u0301siM seL")
	assert.Equal(t, Reverse2("ab\u0301cde"), "edcb\u0301a")
	assert.Equal(t, Reverse2("This `\xc5` is an invalid UTF8 character"), "retcarahc 8FTU dilavni na si `�` sihT")
	assert.Equal(t, Reverse2("The quick bròwn 狐 jumped over the lazy 犬"), "犬 yzal eht revo depmuj 狐 nwòrb kciuq ehT")
	assert.Equal(t, Reverse2("رائد شوملي"), "يلموش دئار")
}
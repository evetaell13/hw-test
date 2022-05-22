package hw03frequencyanalysis

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

var text2 = `Lorem ipsum dolor sit amet, consectetur adipiscing elit,
	sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
	Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut 
	aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate 
	velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non 
	proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
			expected2 := []string{
				"in",
				"ut",
				"dolor",
				"dolore",
				"ad",
				"adipiscing",
				"aliqua",
				"aliquip",
				"amet",
				"anim",
			}
			require.Equal(t, expected2, Top10(text2))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}

func TestTop10OneString(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{input: "cat dog  Dog Cat ", expected: []string{"cat", "dog"}},
		{input: "cat, dog!  Dog!!! Cat ", expected: []string{"cat", "dog"}},
		{input: "cat  my  dog ,   Dog Cat ", expected: []string{"cat", "dog", "my"}},
		{input: "ногу нога Нога ноги! ногу ноги, ножка, Нога!", expected: []string{"нога", "ноги", "ногу", "ножка"}},
		{input: " once - ", expected: []string{"once"}},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result := Top10(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func Test15Words133Times(t *testing.T) {
	ts := []string{}
	var ch rune = 97
	j, i := 0, 0
	for j <= 15 {
		for i < 133 {
			ts = append(ts, string(ch))
			i++
		}
		ch++
		j++
		i = 0
	}
	ts = append(ts, "doom", "morrowind", "iddqd", "idkfa", "arcanum")
	testString := strings.Join(ts, " ")
	t.Run("test 133", func(t *testing.T) {
		expected := []string{
			"a",
			"b",
			"c",
			"d",
			"e",
			"f",
			"g",
			"h",
			"i",
			"j",
		}
		require.Equal(t, expected, Top10(testString))
	})
}

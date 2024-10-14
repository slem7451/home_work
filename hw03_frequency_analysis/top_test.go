package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	tests := []struct {
		text     string
		expected []string
		test     string
	}{
		{
			text: `Как видите, он  спускается  по  лестнице  вслед  за  своим
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
		В этот вечер...`,
			expected: []string{
				"а",
				"он",
				"и",
				"ты",
				"что",
				"в",
				"его",
				"если",
				"кристофер",
				"не",
			},
			test: "original",
		},
		{
			text: `cat and dog, one dog,two cats and one man
					нога! нога нога,  'нога'  Нога
					какой-то какойто -
					dog,cat dog...cat dogcat - ------------ abcde --
					-------
					-

					'cat' "dog" "dog' cat+dog=catdog+ cat* dog\ cat= cat- 

					cat and dog, one dog,two cats and one man`,
			expected: []string{
				"cat",
				"dog",
				"нога",
				"and",
				"one",
				"cats",
				"dog,two",
				"man",
				"--",
				"-------",
			},
			test: "ru+en with punctuations",
		},
		{
			text:     `а б в г д е ж з и й к л м н о п р с т у ф х ц ч ш щ ь ы ъ э ю я`,
			expected: []string{"а", "б", "в", "г", "д", "е", "ж", "з", "и", "й"},
			test:     "ru alphabet",
		},
		{
			text:     `😎 🤪😷 😇`,
			expected: []string{"😇", "😎", "🤪😷"},
			test:     "emoji",
		},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			require.Equal(t, test.expected, Top10(test.text))
		})
	}
}

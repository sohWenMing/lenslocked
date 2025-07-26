package models

type Q2A struct {
	Question, Answer string
}

var QuestionsToAnswers = []Q2A{
	{
		"What are your opening hours?",
		"We are never open, go away",
	},
	{
		"What is your current SLA",
		"We don't owe you anything, go away",
	},
	{
		"How do I renew my membership?",
		"You can't, you're a terrible customer",
	},
}

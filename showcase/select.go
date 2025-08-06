package showcase

import (
	"github.com/fbv/go-templ-ui/view/ui"
)

// GetSampleCountries returns a list of sample countries for select examples
func GetSampleCountries() []ui.SelectItem {
	return []ui.SelectItem{
		{Name: "United States", Value: "us"},
		{Name: "Canada", Value: "ca"},
		{Name: "United Kingdom", Value: "uk"},
		{Name: "Germany", Value: "de"},
		{Name: "France", Value: "fr"},
		{Name: "Spain", Value: "es"},
		{Name: "Italy", Value: "it"},
		{Name: "Japan", Value: "jp"},
		{Name: "Australia", Value: "au"},
		{Name: "Brazil", Value: "br"},
	}
}

// GetSampleLanguages returns a list of sample languages for select examples
func GetSampleLanguages() []ui.SelectItem {
	return []ui.SelectItem{
		{Name: "English", Value: "en"},
		{Name: "Spanish", Value: "es"},
		{Name: "French", Value: "fr"},
		{Name: "German", Value: "de"},
		{Name: "Italian", Value: "it"},
		{Name: "Portuguese", Value: "pt"},
		{Name: "Russian", Value: "ru"},
		{Name: "Chinese", Value: "zh"},
		{Name: "Japanese", Value: "ja"},
		{Name: "Korean", Value: "ko"},
	}
}

// GetSampleColors returns a list of sample colors for select examples
func GetSampleColors() []ui.SelectItem {
	return []ui.SelectItem{
		{Name: "Red", Value: "red"},
		{Name: "Green", Value: "green"},
		{Name: "Blue", Value: "blue"},
		{Name: "Yellow", Value: "yellow"},
		{Name: "Purple", Value: "purple"},
		{Name: "Orange", Value: "orange"},
		{Name: "Pink", Value: "pink"},
		{Name: "Brown", Value: "brown"},
		{Name: "Black", Value: "black"},
		{Name: "White", Value: "white"},
	}
}

// GetSampleCategories returns a list of sample categories for select examples
func GetSampleCategories() []ui.SelectItem {
	return []ui.SelectItem{
		{Name: "Technology", Value: "tech"},
		{Name: "Sports", Value: "sports"},
		{Name: "Entertainment", Value: "entertainment"},
		{Name: "Business", Value: "business"},
		{Name: "Health", Value: "health"},
		{Name: "Education", Value: "education"},
		{Name: "Travel", Value: "travel"},
		{Name: "Food", Value: "food"},
		{Name: "Fashion", Value: "fashion"},
		{Name: "Science", Value: "science"},
	}
}

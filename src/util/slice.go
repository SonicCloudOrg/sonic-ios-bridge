package util

func Prepend[E any](arr []E, el ...E) []E { // https://stackoverflow.com/a/27169176/12857692
	return append(el, arr...)
}

package value_objects

type ValidaterObject[T any] interface {
	IsValid() bool
	ToString() string
	Equals(value ValidaterObject[T]) bool
	GetValue() T
}

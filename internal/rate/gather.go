package rate

type Gatherer[T any] func(int) T

func Gather[T any](number int, fn Gatherer[T]) []T {
	ch := make(chan T)
	defer close(ch)
	for i := 0; i < number; i++ {
		go func() {
			data := fn(i)
			ch <- data
		}()
	}
	var results []T
	for i := 0; i < number; i++ {
		response := <-ch
		results = append(results, response)
	}
	return results
}

type SomeData struct {
	Id   int
	Name string
}

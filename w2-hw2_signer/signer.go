package main

func ExecutePipeline(jobs ...job) {

}

func SingleHash(in chan interface{}, out chan interface{}) {
	go func() {
		for raw := range in {
			data := raw.(string)
			out <- DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
		}
	}()
}

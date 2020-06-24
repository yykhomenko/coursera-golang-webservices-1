package main

func ExecutePipeline(jobs ...job) {

}

func SingleHash(data string) string {

	data = "string"
	r := DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))

	return r
}

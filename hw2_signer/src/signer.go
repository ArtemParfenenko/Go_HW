package main

import (
	"entity"
	"strconv"
	"sync"
	"util"
)

func ExecutePipeline(jobs []job) {
	var previousOutAndCurrentIn chan interface{}
	waitGroup := &sync.WaitGroup{}
	for _, job := range jobs {
		out := make(chan interface{})
		waitGroup.Add(1)
		go executePipeline(job, previousOutAndCurrentIn, out, waitGroup)
		previousOutAndCurrentIn = out
	}
	waitGroup.Wait()
}

func executePipeline(job job, in chan interface{}, out chan interface{}, waitGroup *sync.WaitGroup) {
	defer close(out)
	defer waitGroup.Done()
	job(in, out)
}

func CalculateSingleHash(in, out chan interface{}) {
	waitGroup := &sync.WaitGroup{}
	for val := range in {
		waitGroup.Add(1)
		go func(value interface{}) {
			defer waitGroup.Done()
			out <- calculateSingleHash(value, util.MutexExecutor)
		}(val)
	}
	waitGroup.Wait()
}

func calculateSingleHash(val interface{}, syncExecutor util.SyncExecutor) string {
	decimal, ok := val.(int)
	if !ok {
		panic("The input data is not int !")
	}
	str := strconv.Itoa(decimal)
	var md5Value string

	syncExecutor(func() {
		md5Value = DataSignerMd5(str)
	})

	crc32ValueHolder := new(entity.Crc32ValueHolder)
	crc32ValueHolder.SetValue(str)
	crc32md5ValueHolder := new(entity.Crc32ValueHolder)
	crc32md5ValueHolder.SetValue(md5Value)

	waitGroup := &sync.WaitGroup{}
	for _, crc32Value := range []*entity.Crc32ValueHolder{crc32ValueHolder, crc32md5ValueHolder} {
		waitGroup.Add(1)
		go func(value *entity.Crc32ValueHolder) {
			defer waitGroup.Done()
			value.SetCrc32Value(DataSignerCrc32(value.GetValue()))
		}(crc32Value)
	}
	waitGroup.Wait()

	return crc32ValueHolder.GetCrc32Value() + "~" + crc32md5ValueHolder.GetCrc32Value()
}

func CalculateMultiHash(in, out chan interface{}) {
	waitGroup := &sync.WaitGroup{}
	for val := range in {
		waitGroup.Add(1)
		go func(value interface{}) {
			defer waitGroup.Done()
			out <- calculateMultiHash(value)
		}(val)
	}
	waitGroup.Wait()
}

func calculateMultiHash(val interface{}) string {
	str, ok := val.(string)
	if !ok {
		panic("The input data is not string !")
	}

	multiHashWaitGroup := &sync.WaitGroup{}
	multiHash := make(entity.MultiHash)
	for _, th := range []string{"0", "1", "2", "3", "4", "5"} {
		multiHashWaitGroup.Add(1)
		go func(waitGroup *sync.WaitGroup, partId string) {
			defer waitGroup.Done()
			multiHash.AppendPart(partId, calculateMultiHashPart(str, partId))
		}(multiHashWaitGroup, th)
	}
	multiHashWaitGroup.Wait()
	return multiHash.GetMultiHash()
}

func calculateMultiHashPart(str, th string) string {
	return DataSignerCrc32(th + str)
}

func CombineResults(in, out chan interface{}) {
	var resultSlice = new(entity.SortableSlice)
	for val := range in {
		str, ok := val.(string)
		if !ok {
			panic("The input data is not string !")
		}
		resultSlice.Add(str)
	}
	resultSlice.Sort()
	out <- resultSlice.String()
}

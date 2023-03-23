package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"
)

const DataFile = "loremipsum.txt"

// Return the word frequencies of the text argument.
func WordCount(text string) map[string]int {
	freqs := make(map[string]int)
	words := strings.Fields(text)

	numOfRoutines := 12
	size := len(words)

	wg := new(sync.WaitGroup)
	answerChan := make(chan map[string]int) //The channel where the routines will give their own freqs

	for i := 0; i < size; i += size / numOfRoutines { //Create different chunks of the text file. The amount is proportional to the number of routines.
		j := i + len(words)/numOfRoutines
		if j > size { //Check so that the upper bound is not over the length of the whole textfile.
			j = size //If so just make it the the length of the whole text file.
		}
		chunk := words[i:j] //create a chunk
		wg.Add(1)

		go func() {
			defer wg.Done()
			localFreqs := make(map[string]int) //Make a local frequency map
			for _, word := range chunk {       //Itterate the chunk and do the processing
				word = strings.ToLower(word)
				word = strings.TrimFunc(word, func(r rune) bool { //Remove all the non letter characters form the word
					return !unicode.IsLetter(r) //Handy dandy function
				})
				localFreqs[word]++
				//Send over the local frequency map
			}
			answerChan <- localFreqs
		}()
	}
	//Close the channel when the routines are done.
	go func() {
		wg.Wait()
		close(answerChan)
	}()
	//Receive and handle the routines.
	for {
		freqSlice, ok := <-answerChan //Recieve the local freqs form the routines
		if !ok {                      //Check if channel is open, if closed break the loop
			break
		}
		for word, freq := range freqSlice { //Itterate over all the local freq and update the main one.
			freqs[word] += freq
		}
	}

	return freqs
}

// Benchmark how long it takes to count word frequencies in text numRuns times.
//
// Return the total time elapsed.
func benchmark(text string, numRuns int) int64 {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		WordCount(text)
	}
	runtimeMillis := time.Since(start).Nanoseconds() / 1e6

	return runtimeMillis
}

// Print the results of a benchmark
func printResults(runtimeMillis int64, numRuns int) {
	fmt.Printf("amount of runs: %d\n", numRuns)
	fmt.Printf("total time: %d ms\n", runtimeMillis)
	average := float64(runtimeMillis) / float64(numRuns)
	fmt.Printf("average time/run: %.2f ms\n", average)
}

func main() {
	// read in DataFile as a string called data
	data, _ := os.ReadFile(DataFile)

	fmt.Printf("%#v", WordCount(string(data)))

	numRuns := 100
	runtimeMillis := benchmark(string(data), numRuns)
	printResults(runtimeMillis, numRuns)
}

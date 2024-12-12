package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
)

func main() {
	//read data
	irisMatrix := [][]string{}
	iris, err := os.Open("../datasets/iris.csv")
	if err != nil {
		panic(err)
	}
	defer iris.Close()

	reader := csv.NewReader(iris)
	reader.Comma = ','
	reader.LazyQuotes = true
	// Skip the header row
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

	// Read the rest of the data
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		irisMatrix = append(irisMatrix, record)
	}

	//split data into explaining and explained variables
	X := [][]float64{}
	Y := []string{}
	for _, data := range irisMatrix {

		//convert str slice data into float slice data
		temp := []float64{}
		for _, i := range data[:4] {
			parsedValue, err := strconv.ParseFloat(i, 64)
			if err != nil {
				panic(err)
			}
			temp = append(temp, parsedValue)
		}
		//explaining variables
		X = append(X, temp)

		//explained variables
		Y = append(Y, data[4])

	}

	//split data into training and test
	var (
		trainX [][]float64
		trainY []string
		testX  [][]float64
		testY  []string
	)
	for i, _ := range X {
		if i%2 == 0 {
			trainX = append(trainX, X[i])
			trainY = append(trainY, Y[i])
		} else {
			testX = append(testX, X[i])
			testY = append(testY, Y[i])
		}
	}

	//training
	knn := KNN{}
	knn.k = 5
	knn.fit(trainX, trainY)
	predicted := knn.predict(testX)

	// //check accuracy
	// correct := 0
	// for i, _ := range predicted {
	// 	if predicted[i] == testY[i] {
	// 		correct += 1
	// 	}
	// }
	// fmt.Println(correct)
	// fmt.Println(len(predicted))
	// fmt.Println(float64(correct) / float64(len(predicted)))

	// generate confusion matrix and calculate metrics
	confusionMatrix := make(map[string]map[string]int)
	for i := 0; i < len(predicted); i++ {
		actual := testY[i]
		predictedLabel := predicted[i]
		if _, exists := confusionMatrix[actual]; !exists {
			confusionMatrix[actual] = make(map[string]int)
		}
		confusionMatrix[actual][predictedLabel]++
	}

	// calculate and print precision, recall, F1 score, and overall accuracy
	classMetrics := calculateMetrics(confusionMatrix)
	printMetrics(classMetrics, len(testY), confusionMatrix)

}

// calculate metrics for each class
func calculateMetrics(confusionMatrix map[string]map[string]int) map[string]map[string]float64 {
	metrics := make(map[string]map[string]float64)
	for actual, predictedCounts := range confusionMatrix {
		truePositives := predictedCounts[actual]
		falsePositives := 0
		falseNegatives := 0
		trueNegatives := 0

		for predictedClass, count := range predictedCounts {
			if predictedClass == actual {
				continue
			}
			falsePositives += count
			falseNegatives += confusionMatrix[predictedClass][actual]
		}

		// calculate true negatives (all other counts)
		for otherActual, otherPredictedCounts := range confusionMatrix {
			if otherActual == actual {
				continue
			}
			for _, count := range otherPredictedCounts {
				trueNegatives += count
			}
		}

		precision := float64(truePositives) / float64(truePositives+falsePositives)
		recall := float64(truePositives) / float64(truePositives+falseNegatives)
		f1Score := 2 * (precision * recall) / (precision + recall)

		metrics[actual] = map[string]float64{
			"True Positives":  float64(truePositives),
			"False Positives": float64(falsePositives),
			"True Negatives":  float64(trueNegatives),
			"Precision":       precision,
			"Recall":          recall,
			"F1 Score":        f1Score,
		}
	}
	return metrics
}

func printMetrics(metrics map[string]map[string]float64, total int, confusionMatrix map[string]map[string]int) {
	fmt.Println("Reference Class    True Positives  False Positives  True Negatives  Precision   Recall   F1 Score")
	fmt.Println("---------------    --------------  ---------------  --------------  ---------   ------   --------")
	totalCorrect := 0
	for class, values := range metrics {
		truePositives := values["True Positives"]
		falsePositives := values["False Positives"]
		trueNegatives := values["True Negatives"]
		precision := values["Precision"]
		recall := values["Recall"]
		f1Score := values["F1 Score"]

		fmt.Printf("%-15s    %-14.0f  %-15.0f  %-15.0f  %-9.4f  %-8.4f  %-8.4f\n",
			class, truePositives, falsePositives, trueNegatives, precision, recall, f1Score)

		totalCorrect += int(truePositives)
	}

	overallAccuracy := float64(totalCorrect) / float64(total)
	fmt.Printf("\nOverall accuracy: %.4f\n", overallAccuracy)
}

// calculate euclidean distance betwee two slices
func Dist(source, dest []float64) float64 {
	val := 0.0
	for i, _ := range source {
		val += math.Pow(source[i]-dest[i], 2)
	}
	return math.Sqrt(val)
}

// argument sort
type Slice struct {
	sort.Interface
	idx []int
}

func (s Slice) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}

func NewSlice(n sort.Interface) *Slice {
	s := &Slice{Interface: n, idx: make([]int, n.Len())}
	for i := range s.idx {
		s.idx[i] = i
	}
	return s
}

func NewFloat64Slice(n []float64) *Slice { return NewSlice(sort.Float64Slice(n)) }

// map sort
type Entry struct {
	name  string
	value int
}
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return l[i].name < l[j].name
	} else {
		return l[i].value > l[j].value
	}
}

// count item frequence in slice
func Counter(target []string) map[string]int {
	counter := map[string]int{}
	for _, elem := range target {
		counter[elem] += 1
	}
	return counter
}

type KNN struct {
	k      int
	data   [][]float64
	labels []string
}

func (knn *KNN) fit(X [][]float64, Y []string) {
	//read data
	knn.data = X
	knn.labels = Y
}

func (knn *KNN) predict(X [][]float64) []string {

	predictedLabel := []string{}
	for _, source := range X {
		var (
			distList   []float64
			nearLabels []string
		)
		//calculate distance between predict target data and surpervised data
		for _, dest := range knn.data {
			distList = append(distList, Dist(source, dest))
		}
		//take top k nearest item's index
		s := NewFloat64Slice(distList)
		sort.Sort(s)
		targetIndex := s.idx[:knn.k]

		//get the index's label
		for _, ind := range targetIndex {
			nearLabels = append(nearLabels, knn.labels[ind])
		}

		//get label frequency
		labelFreq := Counter(nearLabels)

		//the most frequent label is the predict target label
		a := List{}
		for k, v := range labelFreq {
			e := Entry{k, v}
			a = append(a, e)
		}
		sort.Sort(a)
		predictedLabel = append(predictedLabel, a[0].name)
	}
	return predictedLabel

}

package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/keiai0/go-demo/demo/algorithm"
)

func AlgorithmDemoHandler(w http.ResponseWriter, r *http.Request) {
	// 線形探索
	lss := generateRandomArray(1, 5)
	lst := generateRandomInt(1, 5)
	if _, err := w.Write([]byte(fmt.Sprintf("Array: %v\n", lss))); err != nil {
		return
	}
	if _, err := w.Write([]byte(fmt.Sprintf("Search target: %d\n", lst))); err != nil {
		return
	}

	lsr := algorithm.LinearSearch(lss, lst)
	if lsr != -1 {
		if _, err := w.Write([]byte(fmt.Sprintf("二分探索(配列内の値の重複なし)：ターゲット %d はインデックス %d にあります\n", lst, lsr))); err != nil {
			return
		}
	} else {
		if _, err := w.Write([]byte("二分探索(配列内の値の重複なし)：該当する値は存在しません。")); err != nil {
			return
		}
	}
	if _, err := w.Write([]byte("=========================\n")); err != nil {
		return
	}

	// バブルソート(昇順)
	bass := generateRandomArray(1, 50)
	if _, err := w.Write([]byte(fmt.Sprintf("Array: %v\n", bass))); err != nil {
		return
	}
	bssr := algorithm.BubbleAscSort(bass)
	if _, err := w.Write([]byte(fmt.Sprintf("バブルソート(昇順): %d\n", bssr))); err != nil {
		return
	}
	if _, err := w.Write([]byte("=========================\n")); err != nil {
		return
	}

	// バブルソート(降順)
	bdss := generateRandomArray(1, 50)
	if _, err := w.Write([]byte(fmt.Sprintf("Array: %v\n", bdss))); err != nil {
		return
	}
	bdsr := algorithm.BubbleDescSort(bdss)
	if _, err := w.Write([]byte(fmt.Sprintf("バブルソート(降順): %d\n", bdsr))); err != nil {
		return
	}
	if _, err := w.Write([]byte("=========================\n")); err != nil {
		return
	}

	// 二分探索(配列内の値の重複なし)
	bsAry := []int{1, 2, 3, 4, 5, 6}
	bsTarget := generateRandomInt(1, 5)
	if _, err := w.Write([]byte(fmt.Sprintf("Array: %v\n", bsAry))); err != nil {
		return
	}
	if _, err := w.Write([]byte(fmt.Sprintf("Search target: %d\n", bsTarget))); err != nil {
		return
	}

	bsRes := algorithm.BinarySearch(bsAry, bsTarget)
	if bsRes != -1 {
		if _, err := w.Write([]byte(fmt.Sprintf("二分探索(配列内の値の重複なし)：ターゲット %d はインデックス %d にあります\n", bsTarget, bsRes))); err != nil {
			return
		}
	} else {
		if _, err := w.Write([]byte("二分探索(配列内の値の重複なし)：該当する値は存在しません。")); err != nil {
			return
		}
	}
	if _, err := w.Write([]byte("=========================\n")); err != nil {
		return
	}

	// 挿入ソート
	isAry := generateRandomArray(1, 10)
	if _, err := w.Write([]byte(fmt.Sprintf("挿入ソート - ソート前: %v\n", isAry))); err != nil {
		return
	}
	isSorted := algorithm.InsertionSort(isAry)
	if _, err := w.Write([]byte(fmt.Sprintf("挿入ソート - ソート後: %v\n", isSorted))); err != nil {
		return
	}
	if _, err := w.Write([]byte("=========================\n")); err != nil {
		return
	}
}

func generateRandomArray(min, max int) []int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	randomArray := make([]int, max)
	for i := 0; i < max; i++ {
		randomArray[i] = r.Intn(max-min+1) + min
	}
	return randomArray
}

func generateRandomInt(min, max int) int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(max-min+1) + min
}

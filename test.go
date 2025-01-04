package main

func test() []int {

	nums := []int{1, 2, 5, 6}
	target := 3

	nums_map := map[int]int{}

	for index, n := range nums {
		nums_map[target-n] = index

		if  nums[nums_map[index]] +n == target {
			return []int{index, nums_map[index]}
		}

	}

}

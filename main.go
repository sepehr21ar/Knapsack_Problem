package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Item structure
type Item struct {
	index  int // starts from 1
	value  int
	weight int
}

// --------------------- Fractional Knapsack - Greedy Approach ---------------------
type GreedyItem struct {
	item  Item
	ratio float64
}

func fractionalKnapsack(items []Item, capacity int) (float64, string) {
	n := len(items)
	gItems := make([]GreedyItem, n)

	for i := 0; i < n; i++ {
		gItems[i] = GreedyItem{
			item:  items[i],
			ratio: float64(items[i].value) / float64(items[i].weight),
		}
	}

	// Sort in descending order by value/weight ratio
	sort.Slice(gItems, func(i, j int) bool {
		return gItems[i].ratio > gItems[j].ratio
	})

	totalValue := 0.0
	remaining := float64(capacity)
	description := []string{}

	for _, gi := range gItems {
		if remaining <= 0 {
			break
		}
		if float64(gi.item.weight) <= remaining {
			// Take the whole item
			totalValue += float64(gi.item.value)
			remaining -= float64(gi.item.weight)
			description = append(description, fmt.Sprintf("Item %d fully", gi.item.index))
		} else {
			// Take a fraction of the item
			fraction := remaining / float64(gi.item.weight)
			totalValue += fraction * float64(gi.item.value)
			description = append(description, fmt.Sprintf("%.2f of Item %d", fraction, gi.item.index))
			remaining = 0
		}
	}

	return totalValue, strings.Join(description, " and ")
}

// --------------------- 0/1 Knapsack - Dynamic Programming ---------------------
func zeroOneKnapsack(items []Item, capacity int) (int, int, string) {
	n := len(items)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	// Fill the DP table
	for i := 1; i <= n; i++ {
		for w := 0; w <= capacity; w++ {
			if items[i-1].weight <= w {
				dp[i][w] = max(dp[i-1][w], items[i-1].value+dp[i-1][w-items[i-1].weight])
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	maxValue := dp[n][capacity]

	// Recover selected items
	selected := []int{}
	w := capacity
	i := n
	for i > 0 && w > 0 {
		if dp[i][w] != dp[i-1][w] {
			selected = append(selected, items[i-1].index)
			w -= items[i-1].weight
		}
		i--
	}

	// Calculate total weight of selected items
	totalWeight := 0
	for _, idx := range selected {
		for _, item := range items {
			if item.index == idx {
				totalWeight += item.weight
				break
			}
		}
	}

	desc := []string{}
	for _, idx := range selected {
		desc = append(desc, fmt.Sprintf("Item %d", idx))
	}
	description := strings.Join(desc, " and ")
	if description == "" {
		description = "none"
	}

	return maxValue, totalWeight, description
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// --------------------- Main function with user input ---------------------
func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Knapsack capacity (W): ")
	scanner.Scan()
	capacity, _ := strconv.Atoi(scanner.Text())

	fmt.Print("Number of items (n): ")
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	items := make([]Item, n)
	fmt.Println("Please enter weight and value for each item on a separate line (weight first, then value):")
	for i := 0; i < n; i++ {
		scanner.Scan()
		parts := strings.Fields(scanner.Text())
		weight, _ := strconv.Atoi(parts[0])
		value, _ := strconv.Atoi(parts[1])
		items[i] = Item{index: i + 1, weight: weight, value: value}
	}

	// Display received input
	fmt.Printf("\nReceived input:\nCapacity = %d, Number of items = %d\n", capacity, n)
	for _, item := range items {
		fmt.Printf("Item %d: weight = %d, value = %d\n", item.index, item.weight, item.value)
	}
	fmt.Println()

	// Fractional Knapsack
	fracValue, fracDesc := fractionalKnapsack(items, capacity)
	fmt.Printf("Fractional Knapsack (Greedy): Total value = %.2f\nSelected: %s\n\n", fracValue, fracDesc)

	// 0/1 Knapsack
	dpValue, dpWeight, dpDesc := zeroOneKnapsack(items, capacity)
	fmt.Printf("0/1 Knapsack (Dynamic Programming): Maximum value = %d\nTotal weight used = %d\nSelected items: %s\n\n", dpValue, dpWeight, dpDesc)

	// Answers to theoretical questions
	fmt.Println("Answers to theoretical questions:")
	fmt.Println("1. Why is the greedy algorithm optimal for the fractional knapsack problem?")
	fmt.Println("   It possesses the greedy choice property and optimal substructure: selecting the highest value/weight ratio at each step always leads to an optimal solution.")
	fmt.Println("")
	fmt.Println("2. Why can the greedy algorithm not be applied to the 0/1 knapsack problem?")
	fmt.Println("   The 0/1 variant lacks the greedy choice property; a locally optimal choice may lead to a suboptimal global solution.")
	fmt.Println("")
	fmt.Println("3. What is the main difference between the two problem models?")
	fmt.Println("   In the fractional model, fractions of items are allowed; in the 0/1 model, each item must be taken entirely or not at all.")
	fmt.Println("")
	fmt.Println("4. Is the fractional knapsack solution always ≥ the 0/1 solution?")
	fmt.Println("   Yes, because the fractional version has fewer constraints and includes all feasible solutions of the 0/1 version as special cases.")
	fmt.Println("")
	fmt.Println("5. Time complexity comparison:")
	fmt.Println("   Greedy (fractional): O(n log n) due to sorting")
	fmt.Println("   Dynamic Programming (0/1): O(n × W)")
	fmt.Println("")
	fmt.Println("6. When are the solutions of both methods equal?")
	fmt.Println("   When the optimal solution does not require any fractional items, i.e., the capacity is fully utilized with whole items.")
}

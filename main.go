package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type Item struct {
	Index  int `json:"index"`
	Value  int `json:"value"`
	Weight int `json:"weight"`
}

type GreedyItem struct {
	Item  Item
	Ratio float64
}

type RequestData struct {
	Capacity int    `json:"capacity"`
	Items    []Item `json:"items"`
}

func fractionalKnapsack(items []Item, capacity int) (float64, string) {
	gItems := make([]GreedyItem, len(items))

	for i := range items {
		gItems[i] = GreedyItem{
			Item:  items[i],
			Ratio: float64(items[i].Value) / float64(items[i].Weight),
		}
	}

	sort.Slice(gItems, func(i, j int) bool {
		return gItems[i].Ratio > gItems[j].Ratio
	})

	total := 0.0
	remain := float64(capacity)
	desc := []string{}

	for _, gi := range gItems {
		if remain <= 0 {
			break
		}
		if float64(gi.Item.Weight) <= remain {
			total += float64(gi.Item.Value)
			remain -= float64(gi.Item.Weight)
			desc = append(desc, fmt.Sprintf("Item %d fully", gi.Item.Index))
		} else {
			f := remain / float64(gi.Item.Weight)
			total += f * float64(gi.Item.Value)
			desc = append(desc, fmt.Sprintf("%.2f of Item %d", f, gi.Item.Index))
			break
		}
	}

	return total, strings.Join(desc, " and ")
}

func zeroOneKnapsack(items []Item, capacity int) (int, int, string) {
	n := len(items)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for w := 0; w <= capacity; w++ {
			if items[i-1].Weight <= w {
				dp[i][w] = max(
					dp[i-1][w],
					items[i-1].Value+dp[i-1][w-items[i-1].Weight],
				)
			} else {
				dp[i][w] = dp[i-1][w]
			}
		}
	}

	selected := []int{}
	w := capacity
	for i := n; i > 0; i-- {
		if dp[i][w] != dp[i-1][w] {
			selected = append(selected, items[i-1].Index)
			w -= items[i-1].Weight
		}
	}

	totalWeight := 0
	desc := []string{}
	for _, idx := range selected {
		for _, it := range items {
			if it.Index == idx {
				totalWeight += it.Weight
				desc = append(desc, fmt.Sprintf("Item %d", idx))
			}
		}
	}

	return dp[n][capacity], totalWeight, strings.Join(desc, " and ")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ---------- Main ----------
func main() {
	r := gin.Default()

	// Serve static HTML
	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.POST("/solve", func(c *gin.Context) {
		var req RequestData
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		fVal, fDesc := fractionalKnapsack(req.Items, req.Capacity)
		zVal, zWeight, zDesc := zeroOneKnapsack(req.Items, req.Capacity)

		c.JSON(200, gin.H{
			"fractional": gin.H{
				"total_value": fVal,
				"selected":    fDesc,
			},
			"zero_one": gin.H{
				"max_value":    zVal,
				"total_weight": zWeight,
				"selected":     zDesc,
			},
		})
	})

	fmt.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}

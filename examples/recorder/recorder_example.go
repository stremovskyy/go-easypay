/*
 * MIT License
 *
 * Copyright (c) 2024 Anton Stremovskyy
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"context"
	"fmt"
	"log"

	go_easypay "github.com/stremovskyy/go-easypay"
	"github.com/stremovskyy/recorder/redis_recorder"
)

const (
	orderID       = "079cd1f0-9917-4fde-8c65-98a35ac70b0b"
	requestID     = "54c02cf3-f089-4db3-b548-d7c9bc382941"
	transactionID = "141876024"
)

func main() {
	recorder := redis_recorder.NewRedisRecorder(
		&redis_recorder.Options{
			Addr:           "localhost:6379",
			DB:             5,
			CompressionLvl: 9,
			Prefix:         "http:easypay",
		},
	)

	client := go_easypay.NewClientWithRecorder(recorder)
	ctx := context.Background()

	fmt.Println("=== Recorder Enhanced Usage Examples ===")

	// Example 1: Get exchanges by Order ID
	fmt.Println("\n2. Getting exchanges by Order ID:")
	orderExchanges, err := client.GetExchangesByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("Error getting exchanges by order ID: %v", err)
	} else {
		fmt.Printf("Found %d exchanges for order\n", len(orderExchanges))
		for i, exchange := range orderExchanges {
			fmt.Printf("  Exchange %d: RequestID=%s, Tags=%v\n", i+1, exchange.RequestID, exchange.Tags)
			fmt.Printf("    Request: %s\n", string(exchange.Request))
			fmt.Printf("    Response: %s\n", string(exchange.Response))
		}
	}

	// Example 2: Get exchanges by Transaction ID
	fmt.Println("\n3. Getting exchanges by Transaction ID:")
	txExchanges, err := client.GetExchangesByTransactionID(ctx, transactionID)
	if err != nil {
		log.Printf("Error getting exchanges by transaction ID: %v", err)
	} else {
		fmt.Printf("Found %d exchanges for transaction\n", len(txExchanges))
		for i, exchange := range txExchanges {
			fmt.Printf("  Exchange %d: RequestID=%s, Tags=%v\n", i+1, exchange.RequestID, exchange.Tags)
		}
	}
}

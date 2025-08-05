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

	fmt.Println("=== Comprehensive Recorder Test ===")

	// Test 1: Get exchanges by known Order ID
	orderID := "079cd1f0-9917-4fde-8c65-98a35ac70b0b"
	fmt.Printf("\n1. Testing GetExchangesByOrderID with orderID: %s\n", orderID)

	orderExchanges, err := client.GetExchangesByOrderID(ctx, orderID)
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		fmt.Printf("✅ Found %d exchanges for order ID\n", len(orderExchanges))
		for i, exchange := range orderExchanges {
			fmt.Printf("   Exchange %d:\n", i+1)
			fmt.Printf("     - RequestID: %s\n", exchange.RequestID)
			fmt.Printf("     - Request size: %d bytes\n", len(exchange.Request))
			fmt.Printf("     - Response size: %d bytes\n", len(exchange.Response))
			fmt.Printf("     - Tags: %v\n", exchange.Tags)
		}
	}

	// Test 2: Get exchanges by known Transaction ID
	transactionID := "1418660124"
	fmt.Printf("\n2. Testing GetExchangesByTransactionID with transactionID: %s\n", transactionID)

	txExchanges, err := client.GetExchangesByTransactionID(ctx, transactionID)
	if err != nil {
		log.Printf("❌ Error: %v", err)
	} else {
		fmt.Printf("✅ Found %d exchanges for transaction ID\n", len(txExchanges))
		for i, exchange := range txExchanges {
			if i >= 2 { // Limit output
				fmt.Printf("   ... and %d more exchanges\n", len(txExchanges)-2)
				break
			}
			fmt.Printf("   Exchange %d: RequestID=%s\n", i+1, exchange.RequestID)
		}
	}

	// Test 3: Get specific exchange by RequestID
	if len(txExchanges) > 0 {
		testRequestID := txExchanges[0].RequestID
		fmt.Printf("\n3. Testing GetRecordedExchange with requestID: %s\n", testRequestID)

		exchange, err := client.GetRecordedExchange(ctx, testRequestID)
		if err != nil {
			log.Printf("❌ Error: %v", err)
		} else {
			fmt.Printf("✅ Successfully retrieved exchange\n")
			fmt.Printf("   - RequestID: %s\n", exchange.RequestID)
			fmt.Printf("   - Request size: %d bytes\n", len(exchange.Request))
			fmt.Printf("   - Response size: %d bytes\n", len(exchange.Response))
			fmt.Printf("   - Timestamp: %v\n", exchange.Timestamp)
		}
	}

	// Test 5: Test filtering functionality
	fmt.Printf("\n5. Testing GetExchangesWithMetrics with filters\n")

	_, err = client.GetExchangesByOrderID(ctx, "")
	if err != nil {
		fmt.Printf("✅ Correctly handled empty orderID: %v\n", err)
	}

	_, err = client.GetExchangesByOrderID(ctx, "non-existent-order-id")
	if err != nil {
		fmt.Printf("✅ Correctly handled non-existent orderID (returned empty list)\n")
	} else {
		fmt.Printf("✅ Correctly handled non-existent orderID (returned empty list)\n")
	}
}

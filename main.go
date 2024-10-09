package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	//send get request to the host on srv.msk01.gigacorp.local
	//and print the response

	error_count := 0

	for {
		time.Sleep(2 * time.Second)

		if error_count >= 3 {
			fmt.Println("Unable to fetch server statistic")
			error_count = 0
		}

		resp, err := http.Get("http://srv.msk01.gigacorp.local")

		//check if http request was successful
		if err != nil {
			error_count++
			continue
		}

		//check if http responce status code is ok
		if resp.StatusCode == 200 {

			//check if http responce contents are ok
			content_type := resp.Header.Get("Content-Type")

			if content_type != "text/plain; charset=UTF-8" {
				error_count++
				continue
			}

		} else {
			//http responce status code is not ok (not 200)
			error_count++
			continue
		}

		//read responce body
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			error_count++
			continue
		}

		//split responce body by comma
		values := strings.Split(string(body), ",")

		//check if values are ok
		if len(values) != 7 {
			error_count++
			continue
		}

		average_load, err := strconv.Atoi(values[0])
		if err != nil {
			continue
		}

		max_ram, err := strconv.Atoi(values[1])
		if err != nil {
			continue
		}

		used_ram, err := strconv.Atoi(values[2])
		if err != nil {
			continue
		}

		disk_space_bytes, err := strconv.Atoi(values[3])
		if err != nil {
			continue
		}

		used_disk_space_bytes, err := strconv.Atoi(values[4])
		if err != nil {
			continue
		}

		net_throughput_bytes_per_s, err := strconv.Atoi(values[5])
		if err != nil {
			continue
		}

		net_load_bytes_per_s, err := strconv.Atoi(values[6])
		if err != nil {
			continue
		}

		if average_load > 30 {
			fmt.Println("Load Average is too high:", average_load)
		}

		if float64(used_ram/max_ram) > 0.8 {
			ram_usage_percent_str := strconv.FormatFloat(float64(used_ram/max_ram)*100, 'f', 2, 64) + "%"
			fmt.Println("Memory usage too high:", ram_usage_percent_str)
		}

		if float64(used_disk_space_bytes/disk_space_bytes) > 0.9 {
			disk_usage_percent_str := strconv.FormatInt(int64(used_disk_space_bytes-disk_space_bytes/1048576), 10)
			fmt.Println("Free disk space is too low:", disk_usage_percent_str, "Mb left")
		}

		if float64(net_load_bytes_per_s/net_throughput_bytes_per_s) > 0.9 {
			net_load_percent_str := strconv.FormatInt(int64(net_throughput_bytes_per_s-net_load_bytes_per_s/1048576), 10)
			fmt.Println("Network bandwidth usage high:", net_load_percent_str, "Mbit/s available")
		}

	}
}

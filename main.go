package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	//send get request to the host on srv.msk01.gigacorp.local
	//and print the response

	fetch_success := false

	for i := 0; i < 3; i++ {
		resp, err := http.Get("http://srv.msk01.gigacorp.local")

		//check if http request was successful
		if err != nil {
			fmt.Println("Http receive error. Trying again...")
			continue
		}

		//check if http responce status code is ok
		if resp.StatusCode == 200 {

			//check if http responce contents are ok
			content_type := resp.Header.Get("Content-Type")

			if content_type != "text/plain; charset=UTF-8" {
				continue
			}

		} else {
			//http responce status code is not ok (not 200)
			continue
		}

		//read responce body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		//split responce body by comma
		values := strings.Split(string(body), ",")

		//check if values are ok
		if len(values) != 7 {
			continue
		}

		average_load, err := strconv.Atoi(values[0])
		max_ram, err := strconv.Atoi(values[1])
		used_ram, err := strconv.Atoi(values[2])
		disk_space_bytes, err := strconv.Atoi(values[3])
		used_disk_space_bytes, err := strconv.Atoi(values[4])
		net_throughput_bytes_per_s, err := strconv.Atoi(values[5])
		net_load_bytes_per_s, err := strconv.Atoi(values[6])

		if err != nil {
			continue
		}

		if average_load > 30 {
			fmt.Println("Load Average is too high:", average_load)
		}

		if float64(used_ram/max_ram) > 0.8 {
			ram_usage_percent_str := strconv.FormatFloat(float64(used_ram/max_ram)*100, 'f', 2, 64)
			fmt.Println("Memory usage too high:", ram_usage_percent_str)
		}

		if float64(used_disk_space_bytes/disk_space_bytes) > 0.9 {
			disk_usage_percent_str := strconv.FormatFloat(float64(used_disk_space_bytes/disk_space_bytes)*100, 'f', 2, 64)
			fmt.Println("Free disk space is too low:", disk_usage_percent_str, "Mb left")
		}

		if float64(net_load_bytes_per_s/net_throughput_bytes_per_s) > 0.9 {
			net_load_percent_str := strconv.FormatFloat(float64(net_load_bytes_per_s/net_throughput_bytes_per_s)*100, 'f', 2, 64)
			fmt.Println("Network load is too high:", net_load_percent_str)
		}

		fetch_success = true

	}

	if !fetch_success {
		fmt.Println("Unable to fetch server statistic")
	}
}

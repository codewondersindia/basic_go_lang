package sanjay

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func HelloSanjay() {
	fmt.Println("sanjay: Hello Sanjay")
}

const new_line = "\n"
const carriage_return = "\r\n"
const separator = string(os.PathSeparator)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const counter_filename = "tmp/counter.txt"
const file_name = "master_logs.sql"

func HelpingToolGetCounter() string {
	const new_line = "\n"
	// Read counter
	data, err := os.ReadFile(counter_filename)
	check(err)
	str_counter := string(data)
	// Replace new line charactor
	str_counter = strings.Replace(str_counter, new_line, "", -1)
	return str_counter
}

// Reads the current value from the file
func readValue() (int, error) {
  data, err := os.ReadFile(counter_filename)
  if err != nil {
    if os.IsNotExist(err) {
      // File doesn't exist, create it with default value "0"
      return 0, writeValue(0)
    }
    return 0, err
  }
  str := string(data)

  // handle new line cleanup
  if strings.Contains(str, carriage_return) {
    str = strings.ReplaceAll(str, carriage_return, "")
  } else if strings.Contains(str, new_line) {
    str = strings.ReplaceAll(str, new_line, "")
  }

  value, err := strconv.Atoi(str)
  if err != nil {
    return 0, fmt.Errorf("error converting value from file: %w", err)
  }
  return value, nil
}

// Writes the given value to the file
func writeValue(value int) error {
	data := []byte(strconv.Itoa(value))
	err := os.WriteFile(counter_filename, data, 0644)
	return err
}

// Increments the value, updates the file, and returns the new value
func incrementValue() (int, error) {
	currentValue, err := readValue()
	if err != nil {
		return 0, err
	}
	newValue := currentValue + 1
	err = writeValue(newValue)
	return newValue, err
}

func shouldAdd(args []string) bool {
	// By default, assume not
	shouldAdd := false

	if len(args) > 1 {
		if args[1] == "-a" {
			shouldAdd = true
		}
	}

	return shouldAdd
}

func contains(arr []string, value string) bool {
	for _, item := range arr {
		if strings.Contains(item, value) {
			return true
		}
	}
	return false
}

func shouldIncrement(args []string) bool {
	// By default, assume increment
	shouldIncrement := false

	if len(args) > 1 {
		firstArg := args[1]
		if firstArg == "-i" {
			shouldIncrement = true
		}
	}

	return shouldIncrement
}

func isVerbose(args []string) bool {
	for _, arg := range args {
		if arg == "-v" || arg == "--verbose" { // Check for both "-v" and "--verbose" flags
			fmt.Println("Verbose mode enabled (found -v or --verbose flag).")
			return true
		}
	}
	return false
}

func HelpingTool() {
	verbose := isVerbose(os.Args)

	if shouldIncrement(os.Args) {
		newValue, err := incrementValue()
		if err != nil {
			fmt.Println("Error incrementing value:", err)
			return
		}
		if verbose {
			fmt.Println("New incremented value:", newValue)
		}
		return
	}

	//
	counter, err := readValue()
	check(err)
	// fmt.Println(counter)
	// return

	// get directory name that has all the table ddl log
	dir := "tmp"
	dir += separator
	dir += "table_ddl"
	dir += separator
	dir += fmt.Sprintf("%d", counter+1)

	if shouldAdd(os.Args) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		if verbose {
			fmt.Println("Directory created successfully!")
			fmt.Println(dir)
		}
		return
	}

	// get all files
	table_ddl_dir, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		fmt.Println("No files!")
		return
	}
	check(err)

	// fmt.Println(
	// 	"Found log files:",
	// 	new_line,
	// 	strings.Join([]string(table_ddl), new_line),
	// 	new_line)

	if verbose {
		fmt.Println("Found ", len(table_ddl_dir), "files")
	}
	var ignore_names = []string{file_name, "test", "test.sql"}
	var filtered_dir []string
	for _, str := range table_ddl_dir {
		if !contains(ignore_names, str.Name()) {
			filtered_dir = append(filtered_dir, str.Name())
		}
	}
	if 0 == len(filtered_dir) {
		fmt.Println("Existing...")
		return
	}

	file_path := dir
	file_path += separator
	file_path += file_name
	// file, err := os.OpenFile(file_path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	file, err := os.OpenFile(file_path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	// file, err := os.Create(master_log_file)
	check(err)
	defer file.Close() // Close the file when the function finishes

	// read all the files in memory
	// files := make(map[string][]byte)
	content := "-- this file is generated by nilesh suthar.\n\n"

	for _, entry := range table_ddl_dir {
		if entry.IsDir() {
			if verbose {
				fmt.Println("Skipping...Dir", entry.Name())
			}
			continue
		}
		if !contains(filtered_dir, entry.Name()) {
			if verbose {
				fmt.Println("Skipping...File", entry.Name())
			}
			continue
		}
		file_path := dir + separator + entry.Name()
		if verbose {
			fmt.Println("Reading", file_path)
		}
		data, err := os.ReadFile(file_path)
		check(err)
		content += "-- " + file_path + "\n"
		content += string(data) + "\n"
	}
	_, err = io.WriteString(file, content)
	check(err)

	fmt.Println("File created and content written successfully!")
}

package main

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "log"
)

// Finds current directory
func currentDir() (string) {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path
}

// Change directory function
func cd(args []string) (error) {

	// Return error if no path provided
        if len(args) < 2 {
            	return ErrNoPath
        }
        
        // Change directory and return any errors
        return os.Chdir(args[1])
}

func main() {
	// Display shell name
	fmt.Println("\n----------Welcome to CofC Shell!----------")
	fmt.Println("----------       ( CCSH )       ----------")
	fmt.Println("---------- TYPE: 'exit' to end  ----------")
	fmt.Println("----------    To run a file:    ----------")
	fmt.Println("---------- TYPE: 'read filename'----------\n")

	reader := bufio.NewReader(os.Stdin)
	
	// Loop
	for {
		// Get current directory
		var path = currentDir()

		// Create Shell Name with dir
		var shellname = "CCSH:" + path + "> "

		// Print Shell Name
		fmt.Print(shellname)
    	
        	// Get input from line
        	input, err := reader.ReadString('\n')
        	if err != nil {
            	fmt.Fprintln(os.Stderr, err)
        }
        // Handle the execution of the input
        if err = execute(input); err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    }
}

// Error message for cd
var ErrNoPath = errors.New("path required")

func execute(input string) error {
	
    // Remove newline character
    input = strings.TrimSuffix(input, "\n")

    // Separate the command and the arguments
    args := strings.Split(input, " ")
    
    // Check for multiple commands
    var multi bool
    multi = false
    
    for i := 0; i < len(args); i++ {
    	switch args[i] {
    	case "&":
    		multi = true
    	case "read":
    		// Check for file
            	file, err2 := os.Open("./" + args[(i+1)])
            	
            	// Initialize variable
            	var lines []string
            	
            	// If file not found
            	if err2 == nil {
            		
            		// Iterate through file
            		scanner := bufio.NewScanner(file) 
            		for scanner.Scan() {
            			// Add lines to array
            			lines = append(lines, scanner.Text())
            		}
            		
            		// Break if error
            		if err3 := scanner.Err(); err3 != nil {
            			log.Fatal(err3)
            		}
            		
            		// Set lines to input
    			args = lines
            	}
            	
            	multi = true
    	}
    }  
    
    // If multiple commands
    if (multi) {
    	
    	// Initialize variables
    	var newArgs []string
    	
    	// Iterate through line
    	for i := 0; i <len(args); i++ {
    		// Find command
    		if (args[i] != "&") { 
    			var x string
    			switch args[i] {
    			// Built in functions
    			case "cd":
    				// Change directory
    				err := cd(args[i:(i+2)])
    				if (err != nil) {
    					fmt.Println(err)
    				}
    				i++
    			case "exit":
    				// Exit
    				os.Exit(0)
    			default:
    				x = args[i] + ";"
    				// Append to command list
    				newArgs = append(newArgs, x)
    			}
    		}
    	}
    	
    	// Format command
    	str := strings.Join(newArgs, " ")
    	
    	// Prepare command
    	cmd := exec.Command("/bin/sh", "-c", str)
    	
    	// Set the correct output device
    	cmd.Stderr = os.Stderr
    	cmd.Stdout = os.Stdout

    	// Execute command and return any errors
    	return cmd.Run()
    	
    } else {
    
    	// Check for built-in commands.
    	switch args[0] {
    	
    	// Change directory
    	case "cd":

        	// Change directory and return any errors
        	return cd(args)
        	
        // Exit shell
    	case "exit":
    	
        	os.Exit(0)
    	}
    	
    	// Prepare the command to execute
    	cmd := exec.Command(args[0], args[1:]...)
    
    	// Set the correct output device
    	cmd.Stderr = os.Stderr
    	cmd.Stdout = os.Stdout
    	
    	// Execute command and return any errors
    	return cmd.Run()
    }
    
    return nil
}

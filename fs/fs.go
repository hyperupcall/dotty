package fs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/eankeen/globe/internal/util"
)

// CopyFile copies a file from a source to destination. If there are any errors,
// it prints the error to the screen and immediately panics
func CopyFile(srcFile string, destFile string, relFile string) {
	util.PrintDebug("srcFile: %s\n", srcFile)
	util.PrintDebug("destFile: %s\n", destFile)

	// ensure parent directory exists
	err := os.MkdirAll(path.Dir(destFile), 0755)
	if err != nil {
		util.PrintError("An error occured when trying to recurisvely create a directory at '%s'. Exiting", destFile)
		panic(err)
	}

	srcContents, err := ioutil.ReadFile(srcFile)
	if err != nil {
		util.PrintError("An error occured when trying to read the file '%s'. Exiting", srcFile)
		panic(err)
	}

	// prompt to remove preexisting file if it exists
	destFilePossiblyExists, err := FilePossiblyExists(destFile)
	if err != nil {
		fmt.Printf("Could not determine if destination file '%s' exists. It could, but we received an error when trying to determine so. Exiting\n", destFile)
		panic(err)
	}
	// since we panic if there is an error, from now on we can
	// be certain that the boolean indicates if the file exist
	fileExists := destFilePossiblyExists

	util.PrintDebug("fileExists: %v\n", fileExists)

	// only continue if we are sure the destination file does not exist. of course, there can still be races, but we'll make sure to print errors
	if fileExists {
		// if the file buffers are the same, return no need to copy
		destContents, err := ioutil.ReadFile(destFile)
		if err != nil {
			util.PrintError("An error occured when trying to read the file '%s'. Exiting", destContents)
			panic(err)
		}

		// if the files are the same, don't copy and return
		if bytes.Compare(srcContents, destContents) == 0 {
			util.PrintDebug("Skipping unchanged '%s' file\n", relFile)
			return
		}

		// file exists and are different, we ask if we should remove file
		shouldRemove := shouldRemoveExistingFile(destFile, relFile, destContents, srcContents)
		if shouldRemove == false {
			return
		}
	}

	// if we got here, it means the file DOES NOT exist or
	// the user wants to OVERWRITE the existing file
	util.PrintInfo("Copying %s to %s\n", srcFile, destFile)
	err = ioutil.WriteFile(destFile, srcContents, 0644)
	if err != nil {
		util.PrintError("There was an error trying to write to file '%s' (from original file '%s'). Exiting", destFile, srcContents)
		panic(err)
	}
}

// RemoveFile removes a file. If there are any errors in doing so, it immediately panics
func RemoveFile(destFile string) {
	err := os.Remove(destFile)
	if err != nil {
		util.PrintError("Error when trying to remove file '%s'. Exiting\n", destFile)
		panic(err)
	}
}

// FilePossiblyExists stops the program if the file possiblyExists
// If no error is returned, we can be certain that boolean has
// integrity. If there is an error, then the file _possibly_ exists
// and the boolean does _not_ have integrity
func FilePossiblyExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)

	if err != nil {
		if os.IsNotExist(err) {
			// return nil because is a known error
			// that the value of the boolean depends on
			return false, nil
		}
		return true, err
	}
	return true, nil
}
